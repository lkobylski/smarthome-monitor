package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"strings"
)

func parseAndSetMetrics(deviceID string, payload []byte, metrics map[string]map[string]prometheus.Gauge) {

	var data map[string]any
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Error().Err(err).Msgf("Error parsing JSON for device [%s]: %v", deviceID, err)
		return
	}

	for key, value := range data {
		if _, exists := metrics[deviceID][key]; !exists {
			metrics[deviceID][key] = promauto.NewGauge(prometheus.GaugeOpts{
				Name: fmt.Sprintf("smarthome_%s_%s", deviceID, key),
				Help: fmt.Sprintf("Metric %s for device %s", key, deviceID),
			})
		}

		switch v := value.(type) {
		case float64:
			metrics[deviceID][key].Set(v)

		case bool:
			metricVal := 0.0
			if v {
				metricVal = 1.0
			}
			metrics[deviceID][key].Set(metricVal)
		case string:
			var metricValue float64
			switch strings.ToLower(v) {
			case "on", "open", "enabled", "positive", "lock":
				metricValue = 1.0
			case "off", "close", "closed", "disabled", "negative", "unlock":
				metricValue = 0.0
			default:
				log.Warn().Msgf("Unsupported string type for key [%s] in device [%s] with the value: %s", key, deviceID, v)
				continue
			}
			metrics[deviceID][key].Set(metricValue)
		default:
			log.Warn().Msgf("Unsupported data type for key [%s] in device [%s]", key, deviceID)
		}

	}
}
