package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jxsl13/stromgedacht/api"
	"github.com/jxsl13/stromgedacht/client"
	"log"
	"net/http"
	"net/netip"
	"sort"
	"time"
)

const (
	MetricsPrefix           = "stromgedacht"
	DescNow                 = "Current state (supergreen, green, yellow, orange or red)"
	DescLoad                = "Current in kWh"
	DescRenewableEnergy     = "Current supply of renewables in kWh"
	DescResidualLoad        = "Current residual load in kWh"
	DescSuperGreenThreshold = "Current threshold on kWh for supergreen state"
)

type MetricsHandler struct {
	apiClient *client.Client
}

func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var metrics Metrics

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// query params
	zip := r.URL.Query().Get("zip")

	if zip == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("zip parameter is required"))
		return
	}

	// metric: now state
	metricNowState, err := h.fetchNowState(zip)
	if err != nil {
		fmt.Errorf("failed to fetch now state for zip %s, %v\n", zip, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	metrics = append(metrics, metricNowState)

	// metrics: forecast stats
	metricsForecastData, err := h.fetchForecastData(zip)
	if err != nil {
		fmt.Errorf("failed to fetch now state for zip %s, %v\n", zip, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	metrics = append(metrics, metricsForecastData...)

	sort.Sort(metrics)

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Write([]byte(metrics.Print()))
}

func (h *MetricsHandler) fetchNowState(zip string) (Metric, error) {
	nowState, err := h.apiClient.GetNow(zip)
	if err != nil {
		return nil, err
	}

	return &GaugeMetric{
		Name:   MetricsPrefix + "_state_now",
		Value:  int64(*(*nowState).State),
		Desc:   DescNow,
		Labels: nil,
	}, nil
}

func (h *MetricsHandler) fetchForecastData(zip string) (Metrics, error) {
	now := time.Now()

	forecast, err := h.apiClient.GetForecast(zip, &now, &now)
	if err != nil {
		return nil, err
	}

	getClosest := func(points *[]api.ForecastPointInTimeViewModel) *api.ForecastPointInTimeViewModel {
		p := *points
		if len(p) == 0 {
			return nil
		}
		sort.Slice(p, func(i, j int) bool {
			d1 := (*p[i].DateTime).Sub(now).Abs()
			d2 := (*p[j].DateTime).Sub(now).Abs()
			return d1 < d2
		})
		return &(p[0])
	}

	pointLoad := getClosest(forecast.Load)
	pointRenewableEnergy := getClosest(forecast.RenewableEnergy)
	pointResidualLoad := getClosest(forecast.ResidualLoad)
	pointSuperGreenThreshold := getClosest(forecast.SuperGreenThreshold)

	if pointLoad == nil || pointRenewableEnergy == nil || pointResidualLoad == nil || pointSuperGreenThreshold == nil {
		return nil, errors.New("failed to fetch one or more forecast metrics")
	}

	return []Metric{
		&GaugeMetric{
			Name:   MetricsPrefix + "_load",
			Value:  int64(*pointLoad.Value),
			Desc:   DescLoad,
			Labels: nil,
		},
		&GaugeMetric{
			Name:   MetricsPrefix + "_renewable_energy",
			Value:  int64(*pointRenewableEnergy.Value),
			Desc:   DescRenewableEnergy,
			Labels: nil,
		},
		&GaugeMetric{
			Name:   MetricsPrefix + "_residual_load",
			Value:  int64(*pointResidualLoad.Value),
			Desc:   DescResidualLoad,
			Labels: nil,
		},
		&GaugeMetric{
			Name:   MetricsPrefix + "_supergreen_threshold",
			Value:  int64(*pointSuperGreenThreshold.Value),
			Desc:   DescSuperGreenThreshold,
			Labels: nil,
		},
	}, nil
}

func main() {
	var (
		bind string
	)

	flag.StringVar(&bind, "web.listen-address", "127.0.0.1:9321", "Address to listen on")
	flag.Parse()

	listenAddr, err := netip.ParseAddrPort(bind)
	if err != nil {
		log.Fatal(err)
	}

	apiClient, err := client.New()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", &MetricsHandler{apiClient: apiClient})

	log.Printf("listening at http://%s/metrics", listenAddr.String())
	if err := http.ListenAndServe(listenAddr.String(), nil); err != nil {
		log.Fatal(err)
	}
}
