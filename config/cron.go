package config

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/fx"
)

func provideCronScheduler() (gocron.Scheduler, error) {
	return gocron.NewScheduler()
}

func handleSchedulerLifecycle(
	lifecycle fx.Lifecycle,
	scheduler gocron.Scheduler,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			scheduler.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return scheduler.Shutdown()
		},
	})
}

func provideCron() fx.Option {
	return fx.Module(
		"cron",
		fx.Provide(provideCronScheduler),
		fx.Invoke(handleSchedulerLifecycle),
	)
}
