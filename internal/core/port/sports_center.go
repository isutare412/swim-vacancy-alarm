package port

import (
	"context"

	"github.com/isutare412/swim-vacancy-alarm/internal/core/model"
)

type SeongnamSDCClient interface {
	FetchSwimCourseData(ctx context.Context, target model.CourseTarget, className string) ([]*model.CourseData, error)
}
