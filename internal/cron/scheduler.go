package cron

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/isutare412/swim-vacancy-alarm/internal/core/port"
)

type Scheduler struct {
	courseService port.CourseService

	client   gocron.Scheduler
	lifeErrs chan error
}

func NewScheduler(cfg Config, courseService port.CourseService) (*Scheduler, error) {
	client, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	scheduler := &Scheduler{
		courseService: courseService,

		client:   client,
		lifeErrs: make(chan error, 1),
	}

	client.NewJob(
		gocron.DurationJob(cfg.ListCoursesInterval),
		gocron.NewTask(scheduler.findSwimVacancies),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)

	return scheduler, nil
}

func (s *Scheduler) Run() <-chan error {
	s.client.Start()
	return s.lifeErrs
}

func (s *Scheduler) Shutdown(ctx context.Context) error {
	shutdownDone := make(chan struct{})
	go func() {
		s.client.Shutdown()
		close(shutdownDone)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-shutdownDone:
	}
	return nil
}

func (s *Scheduler) findSwimVacancies() {
	defer func() {
		if v := recover(); v != nil {
			slog.Error("cron goroutine panicked", "recover", v, "stack", string(debug.Stack()))
			s.lifeErrs <- fmt.Errorf("panic during swim vacancies search: %v", v)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.courseService.FindSwimVacancies(ctx); err != nil {
		slog.Error("failed to find swim vacancies", "error", err)
		return
	}
}
