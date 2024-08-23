package config

import (
	"time"

	"github.com/isutare412/swim-vacancy-alarm/internal/core/service/course"
	"github.com/isutare412/swim-vacancy-alarm/internal/cron"
	"github.com/isutare412/swim-vacancy-alarm/internal/telegram"
)

type Config struct {
	Search   SearchConfig   `koanf:"search"`
	Register RegisterConfig `koanf:"register"`
	Telegram TelegramConfig `koanf:"telegram"`
}

type SearchConfig struct {
	SwimCourse struct {
		Every       time.Duration `koanf:"every"`
		CourseNames []string      `koanf:"course-names"`
	} `koanf:"swim-course"`
}

type RegisterConfig struct {
	SeongnameSDCURL string `koanf:"seongnam-sdc-url"`
}

type TelegramConfig struct {
	BotToken string `koanf:"bot-token"`
	ChatID   string `koanf:"chat-id"`
}

func (c *Config) ToCronConfig() cron.Config {
	return cron.Config{
		ListCoursesInterval: c.Search.SwimCourse.Every,
	}
}

func (c *Config) ToCourseConfig() course.Config {
	return course.Config{
		SwimCourseNames:        c.Search.SwimCourse.CourseNames,
		SeongnamSDCRegisterURL: c.Register.SeongnameSDCURL,
	}
}

func (c *Config) TelegramConfig() telegram.Config {
	return telegram.Config{
		BotToken: c.Telegram.BotToken,
		ChatID:   c.Telegram.ChatID,
	}
}
