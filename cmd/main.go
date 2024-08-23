package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/isutare412/swim-vacancy-alarm/internal/config"
	"github.com/isutare412/swim-vacancy-alarm/internal/core/service/course"
	"github.com/isutare412/swim-vacancy-alarm/internal/cron"
	"github.com/isutare412/swim-vacancy-alarm/internal/http"
	_ "github.com/isutare412/swim-vacancy-alarm/internal/log"
	"github.com/isutare412/swim-vacancy-alarm/internal/telegram"
)

var configPath = flag.String("configs", ".", "path to config directory")

func init() {
	flag.Parse()
}

func main() {
	slog.Info(`Swim Vacancy Alarm
      ████████                                          
    ██████████████████████                              
    ████████████████████████                            
      ██████████████████████                            
            ██████████████████          ████████        
                    ██████████        ████████████      
                      ██████████    ████████████████    
                      ██████████    ████████████████    
                        ██████████  ████████████████    
                    ██████████████  ████████████████    
              ██████████████████████  ████████████      
                        ████████████    ████████        
          ██████████████    ██████████                  
      ██████████████████████    ██████            ██████
    ████████████████████████████      ██        ████████
  ██████████          ████████████████      ████████████
████████                  ██████████▓▓████████████████  
██▓▓██                        ████████████████████      
░░                            ░░    ██████████          `)

	cfg, err := config.LoadValidated(*configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}

	slog.Debug("loaded config", "config", cfg)

	telegramClient := telegram.NewClient(cfg.TelegramConfig())
	seongnamSDCClient := http.SeongnamSDCClient{}
	courseService := course.NewService(cfg.ToCourseConfig(), &seongnamSDCClient, telegramClient)

	cronScheduler, err := cron.NewScheduler(cfg.ToCronConfig(), courseService)
	if err != nil {
		slog.Error("failed to create cron scheduler", "error", err)
		return
	}

	cronSchedulerErrs := cronScheduler.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case s := <-signals:
		slog.Info("received signal", "signal", s.String())
	case err := <-cronSchedulerErrs:
		slog.Error("fatal error from cron scheduler", "error", err)
	}

	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cronScheduler.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown cron scheduler", "error", err)
	}
}
