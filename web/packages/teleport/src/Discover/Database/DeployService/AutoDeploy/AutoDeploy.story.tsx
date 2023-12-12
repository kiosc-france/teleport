/**
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

import React from 'react';
import { initialize, mswLoader } from 'msw-storybook-addon';
import { rest } from 'msw';

import cfg from 'teleport/config';
import { getDbMeta, TeleportProvider } from 'teleport/Discover/fixtures';

import { AutoDeploy } from './AutoDeploy';

export default {
  title: 'Teleport/Discover/Database/Deploy/Auto',
  loaders: [mswLoader],
};

initialize();

export const Init = () => {
  return (
    <TeleportProvider agentMeta={getDbMeta()}>
      <AutoDeploy />
    </TeleportProvider>
  );
};

Init.parameters = {
  msw: {
    handlers: [
      rest.post(
        cfg.getListSecurityGroupsUrl('test-integration'),
        (req, res, ctx) =>
          res(ctx.json({ securityGroups: securityGroupsResponse }))
      ),
    ],
  },
};

export const InitWithAutoEnroll = () => {
  return (
    <TeleportProvider
      agentMeta={{
        ...getDbMeta(),
        autoDiscoveryConfig: { name: '', discoveryGroup: '', aws: [] },
      }}
    >
      <AutoDeploy />
    </TeleportProvider>
  );
};
InitWithAutoEnroll.parameters = {
  msw: {
    handlers: [
      rest.post(
        cfg.getListSecurityGroupsUrl('test-integration'),
        (req, res, ctx) =>
          res(ctx.json({ securityGroups: securityGroupsResponse }))
      ),
    ],
  },
};

export const InitWithLabels = () => {
  return (
    <TeleportProvider
      agentMeta={{
        ...getDbMeta(),
        agentMatcherLabels: [
          { name: 'env', value: 'staging' },
          { name: 'os', value: 'windows' },
        ],
      }}
    >
      <AutoDeploy />
    </TeleportProvider>
  );
};

InitWithLabels.parameters = {
  msw: {
    handlers: [
      rest.post(
        cfg.getListSecurityGroupsUrl('test-integration'),
        (req, res, ctx) =>
          res(ctx.json({ securityGroups: securityGroupsResponse }))
      ),
    ],
  },
};

export const InitSecurityGroupsLoadingFailed = () => {
  return (
    <TeleportProvider agentMeta={getDbMeta()}>
      <AutoDeploy />
    </TeleportProvider>
  );
};

InitSecurityGroupsLoadingFailed.parameters = {
  msw: {
    handlers: [
      rest.post(
        cfg.getListSecurityGroupsUrl('test-integration'),
        (req, res, ctx) =>
          res(
            ctx.status(403),
            ctx.json({
              message: 'some error when trying to list security groups',
            })
          )
      ),
    ],
  },
};

export const InitSecurityGroupsLoading = () => {
  return (
    <TeleportProvider agentMeta={getDbMeta()}>
      <AutoDeploy />
    </TeleportProvider>
  );
};

InitSecurityGroupsLoading.parameters = {
  msw: {
    handlers: [
      rest.post(
        cfg.getListSecurityGroupsUrl('test-integration'),
        (req, res, ctx) => res(ctx.delay('infinite'))
      ),
    ],
  },
};

const securityGroupsResponse = [
  {
    name: 'security-group-1',
    id: 'sg-1',
    description: 'this is security group 1',
    inboundRules: [
      {
        ipProtocol: 'tcp',
        fromPort: '0',
        toPort: '0',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '443',
        toPort: '443',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '2000',
        toPort: '5000',
        cidrs: [
          { cidr: '192.168.1.0/24', description: 'Subnet Mask 255.255.255.0' },
        ],
      },
    ],
    outboundRules: [
      {
        ipProtocol: 'tcp',
        fromPort: '0',
        toPort: '0',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '22',
        toPort: '22',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '2000',
        toPort: '5000',
        cidrs: [
          { cidr: '10.0.0.0/16', description: 'Subnet Mask 255.255.0.0' },
        ],
      },
    ],
  },
  {
    name: 'security-group-2',
    id: 'sg-2',
    description: 'this is security group 2',
    inboundRules: [
      {
        ipProtocol: 'tcp',
        fromPort: '0',
        toPort: '0',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '443',
        toPort: '443',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '2000',
        toPort: '5000',
        cidrs: [
          { cidr: '192.168.1.0/24', description: 'Subnet Mask 255.255.255.0' },
        ],
      },
    ],
    outboundRules: [
      {
        ipProtocol: 'tcp',
        fromPort: '0',
        toPort: '0',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '22',
        toPort: '22',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '2000',
        toPort: '5000',
        cidrs: [
          { cidr: '10.0.0.0/16', description: 'Subnet Mask 255.255.0.0' },
        ],
      },
    ],
  },
  {
    name: 'security-group-3',
    id: 'sg-3',
    description: 'this is security group 3',
    inboundRules: [
      {
        ipProtocol: 'tcp',
        fromPort: '0',
        toPort: '0',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '443',
        toPort: '443',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '2000',
        toPort: '5000',
        cidrs: [
          { cidr: '192.168.1.0/24', description: 'Subnet Mask 255.255.255.0' },
        ],
      },
    ],
    outboundRules: [
      {
        ipProtocol: 'tcp',
        fromPort: '0',
        toPort: '0',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '22',
        toPort: '22',
        cidrs: [{ cidr: '0.0.0.0/0', description: 'Everything' }],
      },
      {
        ipProtocol: 'tcp',
        fromPort: '2000',
        toPort: '5000',
        cidrs: [
          { cidr: '10.0.0.0/16', description: 'Subnet Mask 255.255.0.0' },
        ],
      },
    ],
  },
];
