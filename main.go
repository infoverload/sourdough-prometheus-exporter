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
	humidityDesc    = prometheus.NewDesc("bme280_humidity_percentage", "Humidity in percentage of relative humidity", nil, nil)
)

// make a collector type
// a collector is a prometheus.Collector for a service
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
		ch <- prometheus.NewInvalidMetric(temperatureDesc, err)
		return
	}
	ch <- prometheus.MustNewConstMetric(temperatureDesc, prometheus.GaugeValue, float64(temperature))

	pressure, err := c.sensorDriver.Pressure()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(pressureDesc, err)
		return
	}
	ch <- prometheus.MustNewConstMetric(pressureDesc, prometheus.GaugeValue, float64(pressure)/100)

	humidity, err := c.sensorDriver.Humidity()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(humidityDesc, err)
		return
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
	// make Prometheus client aware of our collector
	registry.MustRegister(collector)

	// set up HTTP handler for metrics and root endpoints
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			`<html>
			<head><title>BME280 Node Exporter for Sourdough Monitoring</title></head>
			<body>
			<h1>BME280 Node Exporter for Sourdough Monitoring</h1>
			<p>To see metrics, go to the following endpoint: /metrics </p>
			</body>
			</html>`))
	})

	// start listening for HTTP connections
	port := ":8080"
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
