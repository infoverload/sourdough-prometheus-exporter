package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

var (
	temperatureDesc = prometheus.NewDesc("bme280_temperature_celsius", "Temperature in celsius degree", nil, nil)
	pressureDesc    = prometheus.NewDesc("bme280_pressure_hpa", "Barometric pressure in hPa", nil, nil)
	humidityDesc    = prometheus.NewDesc("bme280_humidity", "Humidity in percentage of relative humidity", nil, nil)
)

type collector struct {
	sensorDriver *i2c.BME280Driver
}

// implement the Describe method to satisfy Collector interface in client_golang/prometheus/collector.go
func (c collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- temperatureDesc
	ch <- pressureDesc
	ch <- humidityDesc
}

// implement the Collect method to satisfy Collector interface in client_golang/prometheus/collector.go
func (c collector) Collect(ch chan<- prometheus.Metric) {
	temperature, err := c.sensorDriver.Temperature()
	if err != nil {
		log.Printf("error getting temperature: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(temperatureDesc, prometheus.GaugeValue, float64(temperature))

	pressure, err := c.sensorDriver.Pressure()
	if err != nil {
		log.Printf("error getting pressure: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(pressureDesc, prometheus.GaugeValue, float64(pressure)/100)

	humidity, err := c.sensorDriver.Humidity()
	if err != nil {
		log.Printf("error getting humidity: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(humidityDesc, prometheus.GaugeValue, float64(humidity))
}

func main() {
	rAdaptor := raspi.NewAdaptor()
	bme280 := i2c.NewBME280Driver(rAdaptor, i2c.WithBus(1), i2c.WithAddress(0x76))

	if err := bme280.Start(); err != nil {
		log.Fatalf("error starting driver: %s", err)
	}
	log.Print("Connected to BME280")

	registry := prometheus.NewRegistry()
	collector := &collector{sensorDriver: bme280}
	registry.MustRegister(collector)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	addr := os.Getenv("BME280_EXPORTER_ADDRESS")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting server listening on %s...", addr)
	log.Fatal(s.ListenAndServe())
}
