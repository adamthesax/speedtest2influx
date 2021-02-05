package main

import (
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jessevdk/go-flags"
	"github.com/kylegrantlucas/speedtest"
	"github.com/kylegrantlucas/speedtest/http"
)

type report struct {
	download float64
	upload   float64
	server   http.Server
}

func runTest(server string) (report, error) {
	r := report{}

	client, err := speedtest.NewDefaultClient()
	if err != nil {
		return r, err
	}

	s, err := client.GetServer(server)
	if err != nil {
		return r, err
	}
	r.server = s

	dl, err := client.Download(s)
	if err != nil {
		return r, err
	}
	r.download = dl

	ul, err := client.Upload(s)
	if err != nil {
		return r, err
	}
	r.upload = ul

	return r, nil
}

func main() {
	var opts struct {
		Interval        int    `long:"interval" description:"how often in seconds to run a speedtest and track the results. If a interval is not provided, this will track once then exit"`
		InfluxDBURL     string `long:"influxdb-url" description:"URL to influxdb server" required:"true"`
		InfluxDBToken   string `long:"influxdb-token" description:"API token for influxdb server"`
		InfluxDBOrg     string `long:"influxdb-org" description:"influxdb org to write to" required:"true"`
		InfluxDBBucket  string `long:"influxdb-bucket" description:"influxdb bucket to write to" required:"true"`
		SpeedTestServer string `long:"speedtest-server" description:"speedtest server id, if none is provided the nearest server will be used"`
	}

	_, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		os.Exit(2)
	}

	log.Println("Starting speedtest2influx")
	db := influxdb2.NewClient(opts.InfluxDBURL, opts.InfluxDBToken)
	writeAPI := db.WriteAPI(opts.InfluxDBOrg, opts.InfluxDBBucket)
	log.Println("Connected to database")

	for {
		log.Printf("Running speedtest")
		r, _ := runTest(opts.SpeedTestServer)
		log.Printf("Results (%s): Download %3.2f Mbps | Upload: %3.2f | Ping %3.2f ms \n", r.server.Name, r.download, r.upload, r.server.Latency)
		p := influxdb2.NewPoint(
			"speedtest",
			map[string]string{
				"name":    r.server.Name,
				"id":      r.server.ID,
				"country": r.server.Country,
			},
			map[string]interface{}{
				"ping":     r.server.Latency,
				"download": r.download,
				"upload":   r.upload,
				"distance": r.server.Distance,
			},
			time.Now())
		writeAPI.WritePoint(p)
		writeAPI.Flush()

		if opts.Interval == 0 {
			break
		} else {
			time.Sleep(time.Duration(opts.Interval) * time.Second)
		}
	}

	db.Close()
}
