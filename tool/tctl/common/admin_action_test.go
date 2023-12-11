/*
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package common_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gravitational/trace"
	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport"
	"github.com/gravitational/teleport/api/client"
	"github.com/gravitational/teleport/api/client/proto"
	"github.com/gravitational/teleport/api/constants"
	"github.com/gravitational/teleport/api/mfa"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/auth"
	"github.com/gravitational/teleport/lib/auth/mocku2f"
	wancli "github.com/gravitational/teleport/lib/auth/webauthncli"
	wantypes "github.com/gravitational/teleport/lib/auth/webauthntypes"
	libclient "github.com/gravitational/teleport/lib/client"
	libmfa "github.com/gravitational/teleport/lib/client/mfa"
	"github.com/gravitational/teleport/lib/service/servicecfg"
	"github.com/gravitational/teleport/lib/utils"
	tctl "github.com/gravitational/teleport/tool/tctl/common"
	testserver "github.com/gravitational/teleport/tool/teleport/testenv"
	tsh "github.com/gravitational/teleport/tool/tsh/common"
)

func TestAdminActionMFA(t *testing.T) {
	s := newAdminActionTestSuite(t)

	t.Run("Users", s.testAdminActionMFA_Users)
	t.Run("Roles", s.testAdminActionMFA_Roles)
	t.Run("SAMLIdpServiceProvider", s.testAdminActionMFA_SAMLIdpServiceProvider)
}

func (s *adminActionTestSuite) testAdminActionMFA_Users(t *testing.T) {
	ctx := context.Background()

	user, err := types.NewUser("teleuser")
	require.NoError(t, err)

	createUser := func() error {
		_, err := s.authServer.CreateUser(ctx, user)
		return trace.Wrap(err)
	}

	deleteUser := func() error {
		return s.authServer.DeleteUser(ctx, "teleuser")
	}

	t.Run("UserCommands", func(t *testing.T) {
		for _, tc := range []adminActionTestCase{
			{
				command:    "users add teleuser --roles=access",
				cliCommand: &tctl.UserCommand{},
				cleanup:    deleteUser,
			}, {
				command:    "users update teleuser --set-roles=access,auditor",
				cliCommand: &tctl.UserCommand{},
				setup:      createUser,
				cleanup:    deleteUser,
			}, {
				command:    "users rm teleuser",
				cliCommand: &tctl.UserCommand{},
				setup:      createUser,
				cleanup:    deleteUser,
			},
		} {
			t.Run(tc.command, func(t *testing.T) {
				s.runTestCase(t, ctx, tc)
			})
		}
	})

	t.Run("ResourceCommands", func(t *testing.T) {
		s.testAdminActionMFA_ResourceCommand(t, ctx, resourceCommandTestCase{
			resource:       user,
			resourceCreate: createUser,
			resourceDelete: deleteUser,
		})
	})
}

func (s *adminActionTestSuite) testAdminActionMFA_Roles(t *testing.T) {
	ctx := context.Background()

	role, err := types.NewRole("telerole", types.RoleSpecV6{})
	require.NoError(t, err)

	createRole := func() error {
		_, err := s.authServer.CreateRole(ctx, role)
		return trace.Wrap(err)
	}

	getRole := func() (types.Resource, error) {
		return s.authServer.GetRole(ctx, role.GetName())
	}

	deleteRole := func() error {
		return s.authServer.DeleteRole(ctx, role.GetName())
	}

	t.Run("ResourceCommands", func(t *testing.T) {
		s.testAdminActionMFA_ResourceCommand(t, ctx, resourceCommandTestCase{
			resource:       role,
			resourceCreate: createRole,
			resourceDelete: deleteRole,
		})
	})

	t.Run("EditCommand", func(t *testing.T) {
		s.testAdminActionMFA_EditCommand(t, ctx, editCommandTestCase{
			resourceRef:    getResourceRef(role),
			resourceCreate: createRole,
			resourceGet:    getRole,
			resourceDelete: deleteRole,
		})
	})
}

func (s *adminActionTestSuite) testAdminActionMFA_SAMLIdpServiceProvider(t *testing.T) {
	ctx := context.Background()

	sp, err := types.NewSAMLIdPServiceProvider(types.Metadata{
		Name: "test-saml-app",
	}, types.SAMLIdPServiceProviderSpecV1{
		// A test entity descriptor from https://sptest.iamshowcase.com/testsp_metadata.xml.
		EntityDescriptor: `<?xml version="1.0" encoding="UTF-8"?>
		<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" xmlns:ds="http://www.w3.org/2000/09/xmldsig#" entityID="test-saml-app" validUntil="2025-12-09T09:13:31.006Z">
			 <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="true" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
					<md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
					<md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress</md:NameIDFormat>
					<md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://sptest.iamshowcase.com/acs" index="0" isDefault="true"/>
			 </md:SPSSODescriptor>
		</md:EntityDescriptor>`,
		EntityID: "test-saml-app",
	})
	require.NoError(t, err)

	CreateSAMLIdPServiceProvider := func() error {
		return s.authServer.CreateSAMLIdPServiceProvider(ctx, sp)
	}

	getSAMLIdPServiceProvider := func() (types.Resource, error) {
		return s.authServer.GetSAMLIdPServiceProvider(ctx, sp.GetName())
	}

	deleteSAMLIdPServiceProvider := func() error {
		return s.authServer.DeleteSAMLIdPServiceProvider(ctx, sp.GetName())
	}

	t.Run("ResourceCommands", func(t *testing.T) {
		s.testAdminActionMFA_ResourceCommand(t, ctx, resourceCommandTestCase{
			resource:       sp,
			resourceCreate: CreateSAMLIdPServiceProvider,
			resourceDelete: deleteSAMLIdPServiceProvider,
		})
	})

	t.Run("EditCommand", func(t *testing.T) {
		s.testAdminActionMFA_EditCommand(t, ctx, editCommandTestCase{
			resourceRef:    getResourceRef(sp),
			resourceCreate: CreateSAMLIdPServiceProvider,
			resourceGet:    getSAMLIdPServiceProvider,
			resourceDelete: deleteSAMLIdPServiceProvider,
		})
	})
}

type resourceCommandTestCase struct {
	resource       types.Resource
	resourceCreate func() error
	resourceDelete func() error
}

func (s *adminActionTestSuite) testAdminActionMFA_ResourceCommand(t *testing.T, ctx context.Context, tc resourceCommandTestCase) {
	t.Helper()

	resourceYamlPath := filepath.Join(t.TempDir(), fmt.Sprintf("%v.yaml", tc.resource.GetKind()))
	f, err := os.Create(resourceYamlPath)
	require.NoError(t, err)
	require.NoError(t, utils.WriteYAML(f, tc.resource))

	t.Run(fmt.Sprintf("create %v.yaml", tc.resource.GetKind()), func(t *testing.T) {
		s.runTestCase(t, ctx, adminActionTestCase{
			command:    fmt.Sprintf("create %v", resourceYamlPath),
			cliCommand: &tctl.ResourceCommand{},
			cleanup:    tc.resourceDelete,
		})
	})

	t.Run(fmt.Sprintf("create -f %v.yaml", tc.resource.GetKind()), func(t *testing.T) {
		s.runTestCase(t, ctx, adminActionTestCase{
			command:    fmt.Sprintf("create -f %v", resourceYamlPath),
			cliCommand: &tctl.ResourceCommand{},
			setup:      tc.resourceCreate,
			cleanup:    tc.resourceDelete,
		})
	})

	rmCommand := fmt.Sprintf("rm %v", getResourceRef(tc.resource))
	t.Run(rmCommand, func(t *testing.T) {
		s.runTestCase(t, ctx, adminActionTestCase{
			command:    rmCommand,
			cliCommand: &tctl.ResourceCommand{},
			setup:      tc.resourceCreate,
			cleanup:    tc.resourceDelete,
		})
	})
}

type editCommandTestCase struct {
	resourceRef    string
	resourceCreate func() error
	resourceGet    func() (types.Resource, error)
	resourceDelete func() error
}

func (s *adminActionTestSuite) testAdminActionMFA_EditCommand(t *testing.T, ctx context.Context, tc editCommandTestCase) {
	editCommand := fmt.Sprintf("edit %v", tc.resourceRef)
	t.Run(editCommand, func(t *testing.T) {
		s.runTestCase(t, ctx, adminActionTestCase{
			command: editCommand,
			setup:   tc.resourceCreate,
			cliCommand: &tctl.EditCommand{
				Editor: func(filename string) error {
					// Get the latest version of the resource with the correct revision ID.
					resource, err := tc.resourceGet()
					require.NoError(t, err)

					// Update the expiry so that the edit goes through.
					resource.SetExpiry(time.Now())

					f, err := os.Create(filename)
					require.NoError(t, err)
					require.NoError(t, utils.WriteYAML(f, resource))
					return nil
				},
			},
			cleanup: tc.resourceDelete,
		})
	})
}

type adminActionTestSuite struct {
	authServer *auth.Server
	// userClientWithMFA supports MFA prompt for admin actions.
	userClientWithMFA auth.ClientI
	// userClientWithMFA does not support MFA prompt for admin actions.
	userClientNoMFA auth.ClientI
}

func newAdminActionTestSuite(t *testing.T) *adminActionTestSuite {
	t.Helper()
	ctx := context.Background()

	authPref, err := types.NewAuthPreference(types.AuthPreferenceSpecV2{
		Type:         constants.Local,
		SecondFactor: constants.SecondFactorOptional,
		Webauthn: &types.Webauthn{
			RPID: "127.0.0.1",
		},
	})
	require.NoError(t, err)
	authPref.SetOrigin(types.OriginDefaults)

	process := testserver.MakeTestServer(t, testserver.WithAuthPreference(authPref))
	proxyAddr, err := process.ProxyWebAddr()
	require.NoError(t, err)
	authAddr, err := process.AuthAddr()
	require.NoError(t, err)
	authServer := process.GetAuthServer()

	// create admin role and user.
	username := "admin"
	adminRole, err := types.NewRole(username, types.RoleSpecV6{
		Allow: types.RoleConditions{
			Rules: []types.Rule{
				{
					Resources: []string{types.Wildcard},
					Verbs:     []string{types.Wildcard},
				},
			},
		},
	})
	require.NoError(t, err)
	adminRole, err = authServer.CreateRole(ctx, adminRole)
	require.NoError(t, err)

	user, err := types.NewUser(username)
	user.SetRoles([]string{adminRole.GetName()})
	require.NoError(t, err)
	_, err = authServer.CreateUser(ctx, user)
	require.NoError(t, err)

	mockWebauthnLogin := setupWebAuthn(t, authServer, username)
	mockMFAPromptConstructor := func(opts ...mfa.PromptOpt) mfa.Prompt {
		promptCfg := libmfa.NewPromptConfig(proxyAddr.String(), opts...)
		promptCfg.WebauthnLoginFunc = mockWebauthnLogin
		return libmfa.NewCLIPrompt(promptCfg, os.Stderr)
	}

	// Login as the admin user.
	tshHome := t.TempDir()
	err = tsh.Run(context.Background(), []string{
		"login",
		"--insecure",
		"--debug",
		"--user", username,
		"--proxy", proxyAddr.String(),
		"--auth", constants.PasswordlessConnector,
	},
		setHomePath(tshHome),
		setKubeConfigPath(filepath.Join(t.TempDir(), teleport.KubeConfigFile)),
		func(c *tsh.CLIConf) error {
			c.WebauthnLogin = mockWebauthnLogin
			return nil
		},
	)
	require.NoError(t, err)

	userClientNoMFA, err := auth.NewClient(client.Config{
		Addrs: []string{authAddr.String()},
		Credentials: []client.Credentials{
			client.LoadProfile(tshHome, ""),
		},
	})
	require.NoError(t, err)

	userClientWithMFA, err := auth.NewClient(client.Config{
		Addrs: []string{authAddr.String()},
		Credentials: []client.Credentials{
			client.LoadProfile(tshHome, ""),
		},
		MFAPromptConstructor: mockMFAPromptConstructor,
	})
	require.NoError(t, err)

	return &adminActionTestSuite{
		authServer:        authServer,
		userClientNoMFA:   userClientNoMFA,
		userClientWithMFA: userClientWithMFA,
	}
}

type adminActionTestCase struct {
	command    string
	cliCommand tctl.CLICommand
	setup      func() error
	cleanup    func() error
}

func (s *adminActionTestSuite) runTestCase(t *testing.T, ctx context.Context, tc adminActionTestCase) {
	t.Helper()

	t.Run("OK with MFA", func(t *testing.T) {
		err := runTestSubCase(t, ctx, s.userClientWithMFA, tc)
		require.NoError(t, err)
	})

	t.Run("NOK without MFA", func(t *testing.T) {
		err := runTestSubCase(t, ctx, s.userClientNoMFA, tc)
		require.ErrorContains(t, err, mfa.ErrAdminActionMFARequired.Message)
	})

	t.Run("OK mfa off", func(t *testing.T) {
		// turn MFA off, admin actions should not require MFA now.
		authPref := types.DefaultAuthPreference()
		authPref.SetSecondFactor(constants.SecondFactorOff)
		originalAuthPref, err := s.authServer.GetAuthPreference(ctx)
		require.NoError(t, err)

		err = runTestSubCase(t, ctx, s.userClientNoMFA, adminActionTestCase{
			command:    tc.command,
			cliCommand: tc.cliCommand,
			setup: func() error {
				if err := s.authServer.SetAuthPreference(ctx, authPref); err != nil {
					return trace.Wrap(err)
				}
				if tc.setup != nil {
					return tc.setup()
				}
				return nil
			},
			cleanup: func() error {
				if err := s.authServer.SetAuthPreference(ctx, originalAuthPref); err != nil {
					return trace.Wrap(err)
				}
				if tc.cleanup != nil {
					return tc.cleanup()
				}
				return nil
			},
		})
		require.NoError(t, err)
	})
}

func runTestSubCase(t *testing.T, ctx context.Context, client auth.ClientI, tc adminActionTestCase) error {
	t.Helper()

	if tc.setup != nil {
		require.NoError(t, tc.setup(), "unexpected error during setup")
	}
	if tc.cleanup != nil {
		t.Cleanup(func() {
			if err := tc.cleanup(); err != nil && !trace.IsNotFound(err) {
				t.Errorf("unexpected error during cleanup: %v", err)
			}
		})
	}

	app := utils.InitCLIParser("tctl", tctl.GlobalHelpString)
	cfg := servicecfg.MakeDefaultConfig()
	tc.cliCommand.Initialize(app, cfg)

	args := strings.Split(tc.command, " ")
	commandName, err := app.Parse(args)
	require.NoError(t, err)

	match, err := tc.cliCommand.TryRun(ctx, commandName, client)
	require.True(t, match)
	return err
}

func getResourceRef(r types.Resource) string {
	switch kind := r.GetKind(); kind {
	case types.KindClusterAuthPreference:
		// single resources are referred to by kind alone.
		return kind
	default:
		return fmt.Sprintf("%v/%v", r.GetKind(), r.GetName())
	}
}

func setupWebAuthn(t *testing.T, authServer *auth.Server, username string) libclient.WebauthnLoginFunc {
	t.Helper()
	ctx := context.Background()

	const origin = "https://127.0.0.1"
	device, err := mocku2f.Create()
	require.NoError(t, err)
	device.SetPasswordless()

	token, err := authServer.CreateResetPasswordToken(ctx, auth.CreateUserTokenRequest{
		Name: username,
	})
	require.NoError(t, err)

	tokenID := token.GetName()
	res, err := authServer.CreateRegisterChallenge(ctx, &proto.CreateRegisterChallengeRequest{
		TokenID:     tokenID,
		DeviceType:  proto.DeviceType_DEVICE_TYPE_WEBAUTHN,
		DeviceUsage: proto.DeviceUsage_DEVICE_USAGE_PASSWORDLESS,
	})
	require.NoError(t, err)
	cc := wantypes.CredentialCreationFromProto(res.GetWebauthn())

	userWebID := res.GetWebauthn().PublicKey.User.Id

	ccr, err := device.SignCredentialCreation(origin, cc)
	require.NoError(t, err)
	_, err = authServer.ChangeUserAuthentication(ctx, &proto.ChangeUserAuthenticationRequest{
		TokenID: tokenID,
		NewMFARegisterResponse: &proto.MFARegisterResponse{
			Response: &proto.MFARegisterResponse_Webauthn{
				Webauthn: wantypes.CredentialCreationResponseToProto(ccr),
			},
		},
	})
	require.NoError(t, err)

	return func(ctx context.Context, origin string, assertion *wantypes.CredentialAssertion, prompt wancli.LoginPrompt, opts *wancli.LoginOpts) (*proto.MFAAuthenticateResponse, string, error) {
		car, err := device.SignAssertion(origin, assertion)
		if err != nil {
			return nil, "", err
		}
		car.AssertionResponse.UserHandle = userWebID

		return &proto.MFAAuthenticateResponse{
			Response: &proto.MFAAuthenticateResponse_Webauthn{
				Webauthn: wantypes.CredentialAssertionResponseToProto(car),
			},
		}, "", nil
	}
}

func setHomePath(path string) tsh.CliOption {
	return func(cf *tsh.CLIConf) error {
		cf.HomePath = path
		return nil
	}
}

func setKubeConfigPath(path string) tsh.CliOption {
	return func(cf *tsh.CLIConf) error {
		cf.KubeConfigPath = path
		return nil
	}
}
