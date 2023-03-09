package prometheus

import (
	"aweb/web"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

type MiddlewareBuilder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        m.Name,
		Subsystem:   m.Subsystem,
		ConstLabels: m.ConstLabels,
		Help:        m.Help,
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(summaryVec)

	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) {
			log.Println("prometheus middleware log")
			startTime := time.Now()
			next(c)
			endTime := time.Now()
			go report(endTime.Sub(startTime), c, summaryVec)
		}
	}
}

func report(dur time.Duration, ctx *web.Context, vec prometheus.ObserverVec) {
	vec.WithLabelValues(ctx.R.URL.Path, ctx.R.Method, dur.String())
}
