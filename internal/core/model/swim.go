package model

import "fmt"

type CourseTarget string

const (
	CourseTargetAll   CourseTarget = "all"
	CourseTargetAdult CourseTarget = "adult"
)

type CourseData struct {
	CourseName          string
	CourseTarget        string
	WeekName            string `example:"월화목금"`
	RegisterCount       int
	InCartCart          int
	MaxParticipantCount int
}

func (d *CourseData) VacancyCount() int {
	return d.MaxParticipantCount - d.RegisterCount - d.InCartCart
}

func (d *CourseData) VacantAlarmMessage(registerURL string) string {
	tmpl := "- Course name: %s\n- Vacancy count: %d\n- Register URL: %s"
	return fmt.Sprintf(tmpl, d.CourseName, d.VacancyCount(), registerURL)
}
