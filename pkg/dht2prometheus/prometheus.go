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
	"os"
	"strings"

	log "github.com/thkukuk/mqtt-exporter/pkg/logger"
	"github.com/thkukuk/go-dht"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	name              string
	sensorType        dht.SensorType
	gpioPin           int
	offsetTemperature float32
	offsetHumidity    float32
	temperatureMetric *prometheus.Desc
	humidityMetric    *prometheus.Desc
}

func newCollector(config ConfigType) *Collector {
	if Verbose {
		log.Debugf("Creating prometheus collector for sensor %q...",
			config.Name)
	}

	var sensorType dht.SensorType
	if strings.ToLower(config.Sensor) == "dht22" ||
		strings.ToLower(config.Sensor) == "am2302" {
                sensorType = dht.DHT22
        } else if strings.ToLower(config.Sensor) == "dht11" {
                sensorType = dht.DHT11
	} else if strings.ToLower(config.Sensor) == "dht12" {
                sensorType = dht.DHT12
        } else {
		log.Fatalf("ERROR: Unknown sensor type %q",
			config.Sensor)
	}

	return &Collector{
		name: config.Name,
		sensorType: sensorType,
		gpioPin: config.Pin,
		offsetTemperature: config.TemperatureOffset,
		offsetHumidity: config.HumidityOffset,
		temperatureMetric: prometheus.NewDesc("dht_temperature",
			"Temperature (Celsius) measured by the sensor",
			[]string{"dht_name", "hostname"}, nil,
		),
		humidityMetric: prometheus.NewDesc("dht_humidity",
			"Humidity (percent) measured by the sensor",
			[]string{"dht_name", "hostname"}, nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.temperatureMetric
	ch <- c.humidityMetric
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {

	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(c.sensorType, c.gpioPin, false, 10)
	if err != nil {
		log.Errorf("Error reading sensor: %v", err)
		return
        }

	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("Failed to get hostname: %v", err)
	}

	temperature = temperature + c.offsetTemperature
	humidity = humidity + c.offsetHumidity

	ch <- prometheus.MustNewConstMetric(c.temperatureMetric, prometheus.CounterValue, float64(temperature), c.name, hostname)
	ch <- prometheus.MustNewConstMetric(c.humidityMetric, prometheus.CounterValue, float64(humidity), c.name, hostname)
}
