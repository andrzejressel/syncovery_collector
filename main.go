package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	WebAddr string `arg:"env:WEB_ADDR"`
	WebPath string `arg:"env:WEB_PATH"`
}

func main() {
	var syncoveryUrl string
	flag.StringVar(&syncoveryUrl, "url", "", "Syncovery instance url")
	flag.Parse()

	if len(syncoveryUrl) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	client, err := CreateSyncoveryClient(syncoveryUrl)
	if err != nil {
		panic(err)
	}

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(NewSyncoveryCollector(client))

	c := Config{
		WebPath: "/metrics",
		WebAddr: ":8080",
	}

	http.Handle(c.WebPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Syncovery Exporter</title></head>
			<body>
			<h1>Syncovery Exporter</h1>
			<p><a href="` + c.WebPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Println("Collector launched on port 8080")
	log.Fatal(http.ListenAndServe(c.WebAddr, nil))

}
