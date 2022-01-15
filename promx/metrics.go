package promx

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/waytohome/lightning/logx"
)

var (
	registry = prometheus.NewRegistry()

	serverOnce sync.Once
	mysqlOnce  sync.Once

	serverHV *prometheus.HistogramVec
	mysqlHV  *prometheus.HistogramVec
)

func StartPusher(address, port, job string, interval int) {
	registry.MustRegister(collectors.NewGoCollector())
	pusher := push.New(fmt.Sprintf("%s:%s", address, port), job).Gatherer(registry)
	go func() {
		time.Sleep(time.Duration(interval) * time.Second)
		if err := pusher.Add(); err != nil {
			logx.Error("push metrics to pushgateway failed", logx.String("err", err.Error()))
		}
	}()
}

func CollectServerRequest(app, meth, url, status, ua string, begin time.Time) {
	if hv := getServerHV(); hv != nil {
		hv.WithLabelValues(app, meth, url, status, ua).
			Observe(getDuration(begin))
	}
}

func CollectMySQLRequest(db, table, key string, begin time.Time) {
	if hv := getMysqlHV(); hv != nil {
		hv.WithLabelValues(db, table, key).Observe(getDuration(begin))
	}
}

func getDuration(begin time.Time) float64 {
	return float64(time.Since(begin).Nanoseconds()/1e4) / 100
}

func getServerHV() *prometheus.HistogramVec {
	serverOnce.Do(func() {
		if serverHV != nil {
			return
		}
		serverHV = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "server",
			Name:      "request_duration",
			Help:      "服务端请求时长",
			Buckets:   []float64{100, 300, 500, 1000, 3000},
		}, []string{"app", "method", "url", "status", "ua"})
		registry.MustRegister(serverHV)
	})
	return serverHV
}

func getMysqlHV() *prometheus.HistogramVec {
	mysqlOnce.Do(func() {
		if mysqlHV != nil {
			return
		}
		mysqlHV = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "mysql",
			Name:      "sql_duration",
			Help:      "MySQL请求时长",
			Buckets:   []float64{100, 500, 1000, 3000},
		}, []string{"db", "table", "key"})
		registry.MustRegister(mysqlHV)
	})
	return mysqlHV
}
