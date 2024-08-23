package port

import "context"

type CourseService interface {
	FindSwimVacancies(context.Context) error
}
