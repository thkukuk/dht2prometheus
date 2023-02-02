// Copyright 2023 Thorsten Kukuk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
        "io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	log "github.com/thkukuk/dht2prometheus/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/thkukuk/dht2prometheus/pkg/dht2prometheus"
)

var (
	configFile = "config.yaml"
)

func read_yaml_config(conffile string) (dht2Prometheus.ConfigType, error) {

        var config dht2Prometheus.ConfigType

        file, err := ioutil.ReadFile(conffile)
        if err != nil {
                return config, fmt.Errorf("Cannot read %q: %v", conffile, err)
        }
        err = yaml.Unmarshal(file, &config)
        if err != nil {
                return config, fmt.Errorf("Unmarshal error: %v", err)
        }

        return config, nil
}


func main() {
// dht2PrometheusCmd represents the dht2prometheus command
	dht2PrometheusCmd := &cobra.Command{
		Use:   "dht2prometheus",
		Short: "Exports DHTXX values as prometheus metrics",
		Long: `Starts a daemon which exports the values of the DHTXX sensors as metrics for Proemtheus.`,
		Run: runMqttExporterCmd,
		Args:  cobra.ExactArgs(0),
	}

        dht2PrometheusCmd.Version = dht2Prometheus.Version

	dht2PrometheusCmd.Flags().StringVarP(&configFile, "config", "c", configFile, "configuration file")

	dht2PrometheusCmd.Flags().BoolVarP(&dht2Prometheus.Quiet, "quiet", "q", dht2Prometheus.Quiet, "don't print any informative messages")
	dht2PrometheusCmd.Flags().BoolVarP(&dht2Prometheus.Verbose, "verbose", "v", dht2Prometheus.Verbose, "become really verbose in printing messages")

	if err := dht2PrometheusCmd.Execute(); err != nil {
                os.Exit(1)
        }
}

func runMqttExporterCmd(cmd *cobra.Command, args []string) {
	var err error

	if !dht2Prometheus.Quiet {
		log.Infof("Read yaml config %q\n", configFile)
	}
	dht2Prometheus.Config, err = read_yaml_config(configFile)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	dht2Prometheus.RunServer()
}
