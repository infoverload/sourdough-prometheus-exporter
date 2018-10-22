package main

import (
	"log"
	"net/http"

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

// implement Describe method to satisfy Collector interface in client_golang/prometheus/collector.go
func (c collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- temperatureDesc
	ch <- pressureDesc
	ch <- humidityDesc
}

// implement Collect method to satisfy Collector interface in client_golang/prometheus/collector.go
func (c collector) Collect(ch chan<- prometheus.Metric) {
	temperature, err := c.sensorDriver.Temperature()
	if err != nil {
		log.Printf("Error getting temperature: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(temperatureDesc, prometheus.GaugeValue, float64(temperature))

	pressure, err := c.sensorDriver.Pressure()
	if err != nil {
		log.Printf("Error getting pressure: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(pressureDesc, prometheus.GaugeValue, float64(pressure)/100)

	humidity, err := c.sensorDriver.Humidity()
	if err != nil {
		log.Printf("Error getting humidity: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(humidityDesc, prometheus.GaugeValue, float64(humidity))
}

func main() {
	rAdaptor := raspi.NewAdaptor()
	bme280 := i2c.NewBME280Driver(rAdaptor, i2c.WithBus(1), i2c.WithAddress(0x76))

	if err := bme280.Start(); err != nil {
		log.Fatalf("Error starting driver: %s", err)
	}
	log.Print("Connected to BME280! :)")

	registry := prometheus.NewRegistry()
	collector := &collector{sensorDriver: bme280}
	registry.MustRegister(collector)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	//log.Infoln("Listening on", *listenAddress)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
