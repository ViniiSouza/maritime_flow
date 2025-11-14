package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	amqp "github.com/rabbitmq/amqp091-go"
)

type VehicleType string

const (
	metricsNamespace = "vehicles_monitoring"

	vehicleTypeLabel = "vehicle_type"

	HelicopterVehicleType VehicleType = "helicopter"
	ShipVehicleType       VehicleType = "ship"
)

var (
	registry = prometheus.NewRegistry()
	pusher   = push.New(fmt.Sprintf("%s:%s", os.Getenv("PROMETHEUS_PUSHGATEWAY_HOST"), os.Getenv("PROMETHEUS_PUSHGATEWAY_PORT")), metricsNamespace).Gatherer(registry)

	latitudeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "latitude",
			Help:      "latitude coordinate of vehicle",
			Namespace: metricsNamespace,
		},
		[]string{vehicleTypeLabel},
	)

	longitudeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "longitude",
			Help:      "longitude coordinate of vehicle",
			Namespace: metricsNamespace,
		},
		[]string{vehicleTypeLabel},
	)

	motorTemperatureMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "motor_temperature",
			Help:      "temperature of vehicle motor",
			Namespace: metricsNamespace,
		},
		[]string{vehicleTypeLabel},
	)

	fuelLevelMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "fuel_level",
			Help:      "vehicle fuel level in ratio",
			Namespace: metricsNamespace,
		},
		[]string{vehicleTypeLabel},
	)

	cpuUsagePorcMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "cpu_usage_porc",
			Help:      "cpu usage of vehicle system in ratio",
			Namespace: metricsNamespace,
		},
		[]string{vehicleTypeLabel},
	)

	memUsagePorcMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "mem_usage_porc",
			Help:      "memory usage of vehicle system in ratio",
			Namespace: metricsNamespace,
		},
		[]string{vehicleTypeLabel},
	)

	memUsageBytesMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "mem_usage_bytes",
			Help:      "memory usage of vehicle system in bytes",
			Namespace: metricsNamespace,
		},
		[]string{vehicleTypeLabel},
	)
)

type Metadata struct {
	Name        string      `json:"name"`
	VehicleType VehicleType `json:"vehicle_type"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Metrics struct {
	Coordinates      Coordinates `json:"coordinates"`
	MotorTemperature float64     `json:"motor_temperature"`
	FuelLevel        float64     `json:"fuel_level"`
	CPUUsagePorc     float64     `json:"cpu_usage_porc"`
	MemUsagePorc     float64     `json:"mem_usage_porc"`
	MemUsageBytes    int         `json:"mem_usage_bytes"`
}

type Message struct {
	Metadata Metadata `json:"metadata"`
	Metrics  Metrics  `json:"metrics"`
}

func init() {
	registry.MustRegister(latitudeMetric)
	registry.MustRegister(longitudeMetric)
	registry.MustRegister(motorTemperatureMetric)
	registry.MustRegister(fuelLevelMetric)
	registry.MustRegister(cpuUsagePorcMetric)
	registry.MustRegister(memUsagePorcMetric)
	registry.MustRegister(memUsageBytesMetric)
}

func main() {
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	queue := os.Getenv("RABBITMQ_QUEUE")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}

	msgsCh, err := registerConsumer(ch, queue)
	if err != nil {
		log.Fatal(err.Error())
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

main_loop:
	for {
		select {
		case msg := <-msgsCh:
			log.Printf("[%s] received message: %s", time.Now(), string(msg.Body))
			sendMetrics(msg.Body)

		case <-c:
			fmt.Println("interrupting...")
			ch.Close()
			conn.Close()
			break main_loop
		}
	}
}

func registerConsumer(ch *amqp.Channel, queue string) (<-chan amqp.Delivery, error) {
	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"collector",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return msgs, nil
}

func sendMetrics(data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("failed to unmarshal message content: %v", err)
		return
	}

	pusher = pusher.Grouping("vehicle_name", msg.Metadata.Name)
	vehicleType := string(msg.Metadata.VehicleType)

	latitudeMetric.WithLabelValues(vehicleType).Set(msg.Metrics.Coordinates.Latitude)
	longitudeMetric.WithLabelValues(vehicleType).Set(msg.Metrics.Coordinates.Longitude)
	motorTemperatureMetric.WithLabelValues(vehicleType).Set(msg.Metrics.MotorTemperature)
	fuelLevelMetric.WithLabelValues(vehicleType).Set(msg.Metrics.FuelLevel)
	cpuUsagePorcMetric.WithLabelValues(vehicleType).Set(msg.Metrics.CPUUsagePorc)
	memUsagePorcMetric.WithLabelValues(vehicleType).Set(msg.Metrics.MemUsagePorc)
	memUsageBytesMetric.WithLabelValues(vehicleType).Set(float64(msg.Metrics.MemUsageBytes))

	if err := pusher.Add(); err != nil {
		log.Printf("failed to push metrics: %v", err)
	}
}
