package sportscenter

import (
	"cmp"
	"fmt"
	"strconv"

	"github.com/isutare412/swim-vacancy-alarm/internal/core/model"
)

type sportsCenterID string

const (
	sportsCenterIDPangyo sportsCenterID = "04"
)

type categoryID string

const (
	categoryIDSwim categoryID = "01"
)

type smallCategoryID string

const (
	smallCategoryIDAll smallCategoryID = "00"
)

type courseTargetID string

const (
	courseTargetIDAll   courseTargetID = "99"
	courseTargetIDAdult courseTargetID = "01"
)

type sdcSportsCourse struct {
	CourseName          string `json:"pgm_nm"`
	CourseTarget        string `json:"target_nm"`
	WeekName            string `json:"week_nm"`
	RegisterCount       string `json:"group_regi_inwon"`
	InCartCart          string `json:"cart_inwon"`
	MaxParticipantCount string `json:"group_rgl_qty"`
}

func (r *sdcSportsCourse) toCourseData() *model.CourseData {
	return &model.CourseData{
		CourseName:          r.CourseName,
		CourseTarget:        r.CourseTarget,
		WeekName:            r.WeekName,
		RegisterCount:       mustAtoi(r.RegisterCount),
		InCartCart:          mustAtoi(r.InCartCart),
		MaxParticipantCount: mustAtoi(r.MaxParticipantCount),
	}
}

func mustAtoi(s string) int {
	num, err := strconv.Atoi(cmp.Or(s, "0"))
	if err != nil {
		panic(fmt.Errorf("failed to convert %s: %w", s, err))
	}
	return num
}

type sdcListSportsCourseResponse struct {
	Data []sdcSportsCourse `json:"data"`
}

func (r *sdcListSportsCourseResponse) toCourseDataList() []*model.CourseData {
	data := make([]*model.CourseData, 0, len(r.Data))
	for _, d := range r.Data {
		data = append(data, d.toCourseData())
	}
	return data
}
