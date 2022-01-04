package mackerelnullbridge

import (
	"context"
	"time"
)

type App struct {
	serviceMetrics []*ServiceMetric
}

func New(cfg *Config, apikey string, deploy bool) *App {
	client := NewMackerelClient(apikey, deploy)
	app := &App{
		serviceMetrics: make([]*ServiceMetric, 0, len(cfg.Targets)),
	}
	for _, target := range cfg.Targets {
		app.serviceMetrics = append(app.serviceMetrics, &ServiceMetric{
			ServiceName:  target.Service,
			MetricName:   target.MetricName,
			DefaultValue: target.Value,
			DelaySeconds: target.DelaySeconds,
			client:       client,
		})
	}
	return app
}

func (app *App) Run(ctx context.Context) error {
	now := time.Now()
	to := now.Unix()
	from := now.Add(-15 * time.Minute).Unix()

	for _, metric := range app.serviceMetrics {
		if err := metric.FillConst(ctx, from, to); err != nil {
			return err
		}
	}
	return nil
}
