# stromgedacht-exporter

Simple [Prometheus](https://prometheus.io) exporter for the [StromGedacht API](https://api.stromgedacht.de). Thanks to [jxsl13/stromgedacht](https://github.com/jxsl13/stromgedacht) and the StromGedacht team at TransnetBW.

## Installation
```bash
$ go install github.com/muety/stromgedacht-exporter@latest
```

## Usage
```bash
$ stromgedacht-exporter -web.listen-address=127.0.0.1:9321
```

## Example
`GET http://localhost:9321/metrics?zip=76149`

```
# HELP stromgedacht_load Current in kWh
# TYPE stromgedacht_load gauge
stromgedacht_load 5867
# HELP stromgedacht_renewable_energy Current supply of renewables in kWh
# TYPE stromgedacht_renewable_energy gauge
stromgedacht_renewable_energy 2219
# HELP stromgedacht_residual_load Current residual load in kWh
# TYPE stromgedacht_residual_load gauge
stromgedacht_residual_load 3647
# HELP stromgedacht_state_now Current state (supergreen, green, yellow, orange or red)
# TYPE stromgedacht_state_now gauge
stromgedacht_state_now 1
# HELP stromgedacht_supergreen_threshold Current threshold on kWh for supergreen state
# TYPE stromgedacht_supergreen_threshold gauge
stromgedacht_supergreen_threshold 2830

```

## License
MIT