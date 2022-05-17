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
	NextRunDateTime *prometheus.Desc
	LastRunFailure  *prometheus.Desc
	IsScheduled     *prometheus.Desc
	IsLocked        *prometheus.Desc
	IsDisabled      *prometheus.Desc
	IsRunning       *prometheus.Desc
}

func NewSyncoveryCollector(client *SyncoveryClient) *SyncoveryCollector {

	return &SyncoveryCollector{
		client: client,

		LastRunDateTime: prometheus.NewDesc(
			namespace+"last_run_timestamp_seconds",
			"Timestamp of last run of profile",
			[]string{"number", "name"},
			nil,
		),
		NextRunDateTime: prometheus.NewDesc(
			namespace+"next_run_timestamp_seconds",
			"Timestamp of next run of profile",
			[]string{"number", "name"},
			nil,
		),
		LastRunFailure: prometheus.NewDesc(
			namespace+"last_run_failure",
			"Displays whether or not the last profile run was a failure",
			[]string{"number", "name"},
			nil,
		),
		IsScheduled: prometheus.NewDesc(
			namespace+"scheduled",
			"Displays whether or not the profile is scheduled",
			[]string{"number", "name"},
			nil,
		),
		IsRunning: prometheus.NewDesc(
			namespace+"running",
			"Displays whether or not the profile is running",
			[]string{"number", "name"},
			nil,
		),
		IsDisabled: prometheus.NewDesc(
			namespace+"disabled",
			"Displays whether or not the profile is disabled",
			[]string{"number", "name"},
			nil,
		),
		IsLocked: prometheus.NewDesc(
			namespace+"locked",
			"Displays whether or not the profile is locked",
			[]string{"number", "name"},
			nil,
		),
	}
}

func (sc *SyncoveryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sc.LastRunDateTime
	ch <- sc.LastRunFailure
	ch <- sc.NextRunDateTime
	ch <- sc.IsRunning
	ch <- sc.IsDisabled
	ch <- sc.IsScheduled
	ch <- sc.IsLocked

}

func (sc *SyncoveryCollector) Collect(ch chan<- prometheus.Metric) {
	profiles, err := sc.client.GetProfiles()
	if err != nil {
		log.Printf("failed to get profiles: %v", err)
		return
	}

	for _, v := range *profiles {
		ch <- prometheus.MustNewConstMetric(
			sc.LastRunDateTime,
			prometheus.GaugeValue,
			float64(v.LastRunDateTime.Unix()),
			strconv.Itoa(v.Number), v.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.LastRunFailure,
			prometheus.GaugeValue,
			float64Bool(v.LastRunHadError),
			strconv.Itoa(v.Number), v.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.IsLocked,
			prometheus.GaugeValue,
			float64Bool(v.IsLocked),
			strconv.Itoa(v.Number), v.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.IsDisabled,
			prometheus.GaugeValue,
			float64Bool(v.IsDisabled),
			strconv.Itoa(v.Number), v.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.IsRunning,
			prometheus.GaugeValue,
			float64Bool(v.IsRunning),
			strconv.Itoa(v.Number), v.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.IsScheduled,
			prometheus.GaugeValue,
			float64Bool(v.IsScheduled),
			strconv.Itoa(v.Number), v.Name,
		)
		if v.NextRunDateTime != nil {
			ch <- prometheus.MustNewConstMetric(
				sc.NextRunDateTime,
				prometheus.GaugeValue,
				float64(v.NextRunDateTime.Unix()),
				strconv.Itoa(v.Number), v.Name,
			)
		}
	}

}

func float64Bool(b bool) float64 {
	if b {
		return 1
	} else {
		return 0
	}
}
