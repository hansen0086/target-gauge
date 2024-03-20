package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lrita/cmap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	//使用GaugeVec类型可以为监控指标设置标签，这里为监控指标增加一个标签"target"
	gaugeMap    cmap.Map[string, *prometheus.GaugeVec]
	targetGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "target_gauge",
		Help: "set value to /update to change value",
	}, []string{"target", "target_label1"})
)

func main() {
	prometheus.MustRegister(targetGauge)
	gaugeMap.Store("target_gauge", targetGauge)
	targetGauge.WithLabelValues("any", "target_value").Set(0)
	http.Handle("/metrics", promhttp.Handler())
	go monitoring()
	log.Info("begin to server on port 8080")
	// listen on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))

}
func monitoring() {
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/create", create)
}
func update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateMetricValue")
	value := r.URL.Query().Get("value")
	metricName := r.URL.Query().Get("metric_name")
	if value == "" {
		fmt.Println("value not given")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("value is required\n"))
		return
	}
	updatedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Print("error is ", err)
		return
	}
	if metricName == "" {
		fmt.Println("metric_name not given")
		targetGauge.WithLabelValues("any", "target_value").Set(updatedValue)
		w.Write([]byte("setting default target_gauge\n"))
		return
	}
	fmt.Println("value is ", value)
	fmt.Println("metric_name is ", metricName)

	// get metric
	gauge, ok := gaugeMap.Load(metricName)
	if !ok {
		fmt.Printf("metric name %s not found\n", metricName)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%s not found\n", metricName)))
		return
	}

	gauge.WithLabelValues("any", "target_value").Set(updatedValue)
	w.Write([]byte("updated"))
}
func delete(w http.ResponseWriter, r *http.Request) {
	metricName := r.URL.Query().Get("metric_name")
	fmt.Println("metric_name is ", metricName)
	if metricName == "" {
		w.Write([]byte("metric_name not given\n"))
		return
	}
	if metricName == "target_gauge" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("default metric cant be deleted\n"))
		return
	}
	// get metric
	gauge, ok := gaugeMap.Load(metricName)
	if !ok {
		fmt.Printf("metric name %s not found\n", metricName)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%s not found\n", metricName)))
		return
	}
	prometheus.Unregister(gauge)
	gaugeMap.Delete(metricName)
	w.Write([]byte("deleted"))
}

func create(w http.ResponseWriter, r *http.Request) {
	metricName := r.URL.Query().Get("metric_name")
	fmt.Println("metric_name is ", metricName)
	if metricName == "" {
		//prometheus.DefaultRegisterer.Register(targetGauge)
		//targetGauge.WithLabelValues("any", "target_value").Set(0)
		w.Write([]byte("metric_name not given\n"))
		return
	}
	newGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricName,
		Help: "set value to /update to change value",
	}, []string{"target", "target_label1"})
	err := prometheus.Register(newGauge)
	newGauge.WithLabelValues("any", "target_value").Set(0)
	//err := prometheus.DefaultRegisterer.Register(newGauge)
	if err != nil && err.Error() != "" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	gaugeMap.Store(metricName, newGauge)
	w.Write([]byte("created\n"))
}
