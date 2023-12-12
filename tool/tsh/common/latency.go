package common

import (
	"context"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gravitational/trace"
	"github.com/guptarohit/asciigraph"

	"github.com/gravitational/teleport/lib/utils"
	"github.com/gravitational/teleport/lib/utils/diagnostics/latency"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit}}
}

var helpKeys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type LatencyModel struct {
	ready          bool
	client, server string
	h, w           int
	help           help.Model

	maxP, maxH    int64
	last          latency.Statistics
	proxyLatency  *utils.CircularBuffer
	targetLatency *utils.CircularBuffer
}

func NewLatencyModel(clientLabel, serverLabel string) (*LatencyModel, error) {
	clientStats, err := utils.NewCircularBuffer(50)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	serverStats, err := utils.NewCircularBuffer(50)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &LatencyModel{
		client:        clientLabel,
		server:        serverLabel,
		help:          help.New(),
		proxyLatency:  clientStats,
		targetLatency: serverStats,
	}, nil
}

func (m *LatencyModel) Init() tea.Cmd {
	m.help = help.New()
	return nil
}

func (m *LatencyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.ready = true
		m.h = msg.Height
		m.w = msg.Width
	case latency.Statistics:
		m.proxyLatency.Add(float64(msg.Client))
		m.targetLatency.Add(float64(msg.Server))

		m.maxH = max(m.maxH, msg.Server)
		m.maxP = max(m.maxP, msg.Client)
		m.last = msg
	}
	return m, nil
}

func (m *LatencyModel) View() string {
	if !m.ready {
		return ""
	}

	proxy := m.proxyLatency.Data(150)
	target := m.targetLatency.Data(150)

	if proxy == nil || target == nil {
		return ""
	}

	plotH := m.h / 3
	proxyPlot := asciigraph.Plot(
		proxy,
		asciigraph.Height(plotH),
		asciigraph.Width(m.w),
		asciigraph.SeriesColors(asciigraph.Blue),
		asciigraph.Caption(m.server),
	)

	hostPlot := asciigraph.Plot(
		target,
		asciigraph.Height(plotH),
		asciigraph.Width(m.w),
		asciigraph.SeriesColors(asciigraph.Goldenrod),
		asciigraph.Caption(m.client),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(s.Cell.GetForeground()).
		Background(s.Cell.GetBackground()).
		Bold(false)

	t := table.New(
		table.WithColumns([]table.Column{{Title: "", Width: 10}, {Title: "Last Seen", Width: 10}, {Title: "Max", Width: 10}}),
		table.WithRows([]table.Row{
			[]string{m.server, strconv.FormatInt(m.last.Client, 10) + "ms", strconv.FormatInt(m.maxP, 10) + "ms"},
			[]string{m.client, strconv.FormatInt(m.last.Server, 10) + "ms", strconv.FormatInt(m.maxH, 10) + "ms"},
		}),
		table.WithHeight(2),
		table.WithStyles(s),
	)

	return lipgloss.JoinVertical(0, "Network Latency\n",
		lipgloss.JoinVertical(0,
			lipgloss.JoinVertical(0, lipgloss.JoinVertical(0, proxyPlot+"\n", hostPlot+"\n"),
				lipgloss.PlaceHorizontal(m.w, lipgloss.Center, t.View()),
			),
			lipgloss.PlaceHorizontal(m.w, lipgloss.Right, m.help.View(helpKeys)),
		))
}

func showLatency(ctx context.Context, clientPinger, serverPinger latency.Pinger, clientLabel, serverLabel string) error {
	m, err := NewLatencyModel(clientLabel, serverLabel)
	if err != nil {
		return trace.Wrap(err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithContext(ctx), tea.WithoutSignalHandler())

	monitor, err := latency.NewMonitor(latency.MonitorConfig{
		InitialPingInterval:   time.Millisecond,
		InitialReportInterval: 500 * time.Millisecond,
		PingInterval:          time.Second,
		ReportInterval:        time.Second,
		ClientPinger:          clientPinger,
		ServerPinger:          serverPinger,
		Reporter: latency.ReporterFunc(func(ctx context.Context, stats latency.Statistics) error {
			p.Send(stats)
			return nil
		}),
	})
	if err != nil {
		return trace.Wrap(err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		monitor.Run(ctx)
	}()

	if _, err := p.Run(); err != nil {
		return trace.Wrap(err)
	}

	return nil
}
