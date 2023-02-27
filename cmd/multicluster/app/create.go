/*
Copyright © 2023

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"github.com/electrocucaracha/multi-cluster/internal/multicluster"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
)

func NewCreateCommand(provider multicluster.DataSource) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a deployment with multiple KIND clusters",
		Long: `Create a deployment with multiple KIND clusters based on the configuration
passed as parameters.

Multicluster deployment create KIND clusters in independent bridges, that are connected
through an special container that handles the routing and the WAN emulation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			name, err := getName(flags)
			if err != nil {
				return errors.Wrap(err, "failed to retrieve the name of the multi-cluster")
			}

			configPath, err := getConfigPath(flags)
			if err != nil {
				return errors.Wrap(err, "failed to retrieve the configuration file path of the multi-cluster")
			}
			wanEmulatorImg, _ := flags.GetString("wanem")

			if err := provider.Create(name, configPath, wanEmulatorImg); err != nil {
				return errors.Wrapf(err, "failed to create %s multi-cluster", name)
			}

			return nil
		},
	}

	cmd.Flags().String(
		"name",
		cluster.DefaultName,
		"the multicluster context name",
	)

	cmd.Flags().String(
		"config",
		"./config.yml",
		"the config file with the cluster configuration",
	)

	cmd.Flags().String(
		"wanem",
		"electrocucaracha/wanem",
		"the WAN emulator docker image",
	)

	return cmd
}
