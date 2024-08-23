package course

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/isutare412/swim-vacancy-alarm/internal/core/model"
	"github.com/isutare412/swim-vacancy-alarm/internal/core/port"
)

type Service struct {
	seongnamSDCClient port.SeongnamSDCClient
	telegramClient    port.TelegramClient

	swimCourseName         string
	seongnamSDCRegisterURL string

	mu          sync.Mutex
	fetchCount  int
	lastLogTime time.Time
}

func NewService(
	cfg Config,
	seongnamSDCClient port.SeongnamSDCClient,
	telegramClient port.TelegramClient,
) *Service {
	return &Service{
		seongnamSDCClient: seongnamSDCClient,
		telegramClient:    telegramClient,

		swimCourseName:         cfg.SwimCourseName,
		seongnamSDCRegisterURL: cfg.SeongnamSDCRegisterURL,
	}
}

func (s *Service) FindSwimVacancies(ctx context.Context) error {
	courseDataList, err := s.seongnamSDCClient.FetchSwimCourseData(ctx, model.CourseTargetAdult, s.swimCourseName)
	if err != nil {
		return fmt.Errorf("fetching swim course data: %w", err)
	}
	s.tryLogCourseList(courseDataList)

	var vacantCourses []*model.CourseData
	for _, course := range courseDataList {
		if course.VacancyCount() > 0 {
			vacantCourses = append(vacantCourses, course)
		}
	}
	if len(vacantCourses) == 0 {
		return nil
	}

	slog.Debug("send vacant course alaram message", "vacancyCount", len(vacantCourses))

	message := buildVacancyAlarmMessage(vacantCourses, s.seongnamSDCRegisterURL)
	if err := s.telegramClient.SendMessage(ctx, message); err != nil {
		return fmt.Errorf("sending telegram message: %w", err)
	}

	return nil
}

func (s *Service) tryLogCourseList(courseList []*model.CourseData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.fetchCount++

	if !s.lastLogTime.IsZero() && time.Now().Sub(s.lastLogTime) < time.Hour {
		return
	}

	slog.Info("fetched swim course data", "fetchCount", s.fetchCount, "courseCount", len(courseList))
	s.lastLogTime = time.Now()
}

func buildVacancyAlarmMessage(courses []*model.CourseData, registerURL string) string {
	messages := make([]string, 0, len(courses))
	for _, c := range courses {
		msg := c.VacantAlarmMessage(registerURL)
		messages = append(messages, msg)
	}
	return strings.Join(messages, "\n\n---\n\n")
}
