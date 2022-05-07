package main

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace string = "syncovery_profile_"
)

type SyncoveryCollector struct {
	client          *SyncoveryClient
	LastRunDateTime *prometheus.Desc
}

func NewSyncoveryCollector(client *SyncoveryClient) *SyncoveryCollector {

	return &SyncoveryCollector{
		client: client,

		LastRunDateTime: prometheus.NewDesc(
			namespace+"last_run",
			"Timestamp of last run of profile",
			[]string{"number", "name"},
			nil,
		),
	}
}

func (sc *SyncoveryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sc.LastRunDateTime
}

func (sc *SyncoveryCollector) Collect(ch chan<- prometheus.Metric) {
	profiles, err := sc.client.GetProfiles()
	if err != nil {
		log.Printf("failed to get torrents: %v", err)
		return
	}

	for _, v := range *profiles {
		ch <- prometheus.MustNewConstMetric(
			sc.LastRunDateTime,
			prometheus.GaugeValue,
			float64(v.LastRunDateTime.Unix()),
			strconv.Itoa(v.Number), v.Name,
		)
	}

}
