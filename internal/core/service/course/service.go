package course

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/isutare412/swim-vacancy-alarm/internal/core/model"
	"github.com/isutare412/swim-vacancy-alarm/internal/core/port"
)

type Service struct {
	seongnamSDCClient port.SeongnamSDCClient
	telegramClient    port.TelegramClient

	swimCourseNames        []string
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

		swimCourseNames:        coursesToSearch(cfg.SwimCourseNames),
		seongnamSDCRegisterURL: cfg.SeongnamSDCRegisterURL,
	}
}

func (s *Service) FindSwimVacancies(ctx context.Context) error {
	coursesChan := make(chan []*model.CourseData, len(s.swimCourseNames))

	eg, _ := errgroup.WithContext(ctx)
	for _, name := range s.swimCourseNames {
		eg.Go(func() error {
			courses, err := s.seongnamSDCClient.FetchSwimCourseData(
				ctx, model.CourseTargetAdult, name)
			if err != nil {
				return fmt.Errorf("fetching swim course data: %w", err)
			}

			s.increaseFetchCount()
			coursesChan <- courses

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("waiting for errgroup: %w", err)
	}

	var courseDataList []*model.CourseData
	close(coursesChan)
	for courses := range coursesChan {
		courseDataList = append(courseDataList, courses...)
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

	slog.Info("send vacant course alaram message", "vacancyCount", len(vacantCourses))

	message := buildVacancyAlarmMessage(vacantCourses, s.seongnamSDCRegisterURL)
	if err := s.telegramClient.SendMessage(ctx, message); err != nil {
		return fmt.Errorf("sending telegram message: %w", err)
	}

	return nil
}

func (s *Service) increaseFetchCount() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.fetchCount++
}

func (s *Service) tryLogCourseList(courseList []*model.CourseData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.lastLogTime.IsZero() && time.Now().Sub(s.lastLogTime) < time.Hour {
		return
	}

	slog.Info("fetched swim course data", "fetchCount", s.fetchCount, "courseCount", len(courseList))
	s.lastLogTime = time.Now()
}

func coursesToSearch(names []string) []string {
	if len(names) == 0 {
		return []string{""} // Empty course name means to search ALL.
	}
	return names
}

func buildVacancyAlarmMessage(courses []*model.CourseData, registerURL string) string {
	messages := make([]string, 0, len(courses))
	for _, c := range courses {
		msg := c.VacantAlarmMessage(registerURL)
		messages = append(messages, msg)
	}
	return strings.Join(messages, "\n\n---\n\n")
}
