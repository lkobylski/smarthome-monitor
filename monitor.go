package main

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Device struct to represent registered devices and their data structure
type Device struct {
	ID         string `yaml:"id"`
	Topic      string `yaml:"topic"`
	DataFormat string `yaml:"data_format"`
}

// Monitor struct to encapsulate MQTT client and configuration
type Monitor struct {
	config  *Config
	client  mqtt.Client
	msgChan chan mqtt.Message
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	metrics map[string]map[string]prometheus.Gauge
}

// NewMonitor initializes a new Monitor instance
func NewMonitor(configFile string) (*Monitor, error) {
	config, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	msgChan := make(chan mqtt.Message, 100)

	// Connection options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientID)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("connection error: %v", token.Error())
	}

	// Initialize Prometheus metrics
	metrics := make(map[string]map[string]prometheus.Gauge)
	for _, device := range config.Devices {
		log.Info().Msgf("registred device %s on topic %s", device.ID, device.Topic)
		metrics[device.ID] = make(map[string]prometheus.Gauge)
	}

	monitor := &Monitor{
		config:  config,
		client:  client,
		msgChan: msgChan,
		ctx:     ctx,
		cancel:  cancel,
		metrics: metrics,
	}

	return monitor, nil
}

// Start begins monitoring the MQTT topics
func (m *Monitor) Start() {
	m.subscribeTopics()

	//Start Prometheus HTTP server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Info().Msgf("Start Prometheus metrics server on :2112")
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			log.Fatal().Err(err)

		}
	}()

	// Goroutine to process messages from the channel
	m.wg.Add(1)
	go m.processMessages()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		<-signalChan
		fmt.Println("Received shutdown signal")
		m.cancel()
		close(m.msgChan) // Close the channel to stop the message processing goroutine
	}()

	// Wait for context cancellation
	<-m.ctx.Done()

	// Wait for all goroutines to finish
	m.wg.Wait()

	// Disconnect client
	m.client.Disconnect(uint(m.config.DisconnectTimeout))
}

// subscribeTopics subscribes to all topics defined in the config
func (m *Monitor) subscribeTopics() {
	for _, topic := range m.config.Topics {
		topic := topic // Capture the loop variable
		topicHandler := func(client mqtt.Client, msg mqtt.Message) {
			m.msgChan <- msg // Send message to channel for processing
		}
		token := m.client.Subscribe(topic, 1, topicHandler)
		token.Wait()
		log.Info().Msgf("Subscribed to topic: %s", topic)
	}
}

// processMessages handles messages received on subscribed topics
func (m *Monitor) processMessages() {
	defer m.wg.Done()
	for msg := range m.msgChan {
		m.processMessage(msg)
	}
}

// processMessage processes a single MQTT message
func (m *Monitor) processMessage(msg mqtt.Message) {
	for _, device := range m.config.Devices {
		if msg.Topic() == device.Topic {
			log.Debug().Msgf("Processing message for device [%s] with data format [%s]: %s", device.ID, device.DataFormat, msg.Payload())

			parseAndSetMetrics(device.ID, msg.Payload(), m.metrics)
		}
	}
}
