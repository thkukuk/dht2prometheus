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

package dht2Prometheus

import (
        "net/http"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/d2r2/go-logger"
	log "github.com/thkukuk/dht2prometheus/pkg/logger"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
  defListen = ":9420"
)

type ConfigType struct {
	Name                   string `yaml:"name"`
	Sensor                 string `yaml:"sensor"`
	Pin                    int `yaml:"gpio_pin"`
	Listen                 string `yaml:"listen"`
}

var (
	Version = "unreleased"
	Quiet   = false
	Verbose = false
	Config ConfigType
)

func RunServer() {
	if !Quiet {
		log.Infof("DHT to Prometheus Exporter (dht2prometheus) %s is starting...\n", Version)
	}

	logger.ChangePackageLogLevel("dht", logger.InfoLevel)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		<-quit
		if !Quiet {
			log.Info("Terminated via Signal. Shutting down...")
		}
		os.Exit(0)
	}()

	if len(Config.Listen) == 0 {
		Config.Listen = defListen
	}

	collector := newCollector(Config)
	prometheus.MustRegister(collector)
        http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
                // XXX ErrorLog: log,
        }))
	if !Quiet {
		log.Infof("Starting http server on %s", Config.Listen)
	}
        log.Fatal(http.ListenAndServe(Config.Listen, nil))
}
