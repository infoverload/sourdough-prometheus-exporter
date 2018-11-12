## Sourdough Prometheus Exporter

This is a simple Prometheus exporter written in Go for the Bosch BME280 Sensor.  It is part of a personal project to monitor the temperature and humidity of my sourdough cultures. 

The following metrics are exported:
Temperature (celsius)
Barometric pressure (hPa)
Humidity (percentage of relative humidity)


### Build it

Raspberry Pi is being used with the sensor. If you are also using a Raspberry Pi, you can connect the sensor to the Pi according to this diagram.

BME280 | Desc    | GPIO Header Pins
------ | ------- |------------------
VIN    | 3.3V    | P1-01
GND    | Ground  | P1-06
SCL    | I2C SCL | P1-05
SDA    | I2C SDA | P1-03


### Run it

```
git clone git://github.com/infoverload/sourdough-prometheus-exporter
cd sourdough-prometheus-exporter
go get ./...
go build 
./sourdough-prometheus-exporter
```

### Start automatically on boot
If you want to automatically start the program on boot, you can use the provided systemd unit (sourdough.service).


### Get current IP address on device running it
If your Pi is connected to the WiFi and your router assigns it a new IP address periodically, you can use the script provided to get alerts to your Slack about the IP changes.  If you want to start this script on boot, you can use the provided systemd unit (gethostname.service).


### To do
- [ ] write test
- [ ] implement simple service discovery
- [ ] authorisation
