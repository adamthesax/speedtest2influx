# speedtest2influx

Runs a [Speedtest](https://www.speedtest.net/) and logs the results to [InfluxDB](https://www.influxdata.com/).

# Usage
```
Usage:
  speedtest2influx [OPTIONS]

Application Options:
      --interval=         how often in seconds to run a speedtest and track the
                          results. If a interval is not provided, this will
                          track once then exit
      --influxdb-url=     URL to influxdb server
      --influxdb-token=   API token for influxdb server
      --influxdb-org=     influxdb org to write to
      --influxdb-bucket=  influxdb bucket to write to
      --speedtest-server= speedtest server id, if none is provided the nearest
                          server will be used

Help Options:
  -h, --help              Show this help message
```


