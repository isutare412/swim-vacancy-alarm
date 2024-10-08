package sportscenter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/isutare412/swim-vacancy-alarm/internal/core/model"
)

const (
	sdcSportsAPIListCourses = "https://spo.isdc.co.kr/courseListAjax.ajax"
)

type SeongnamSDCClient struct{}

func (c *SeongnamSDCClient) FetchSwimCourseData(
	ctx context.Context,
	target model.CourseTarget,
	className string,
) ([]*model.CourseData, error) {
	bodyValues := buildListCoursesFormBody(
		sportsCenterIDPangyo, categoryIDSwim, smallCategoryIDAll, courseTargetIDAdult, className)
	req, err := http.NewRequestWithContext(
		ctx, "POST", sdcSportsAPIListCourses, strings.NewReader(bodyValues.Encode()))
	if err != nil {
		return nil, fmt.Errorf("building http request: %w", err)
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing http request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		msg := simplifyResponseBody(string(bodyBytes))
		return nil, fmt.Errorf("unexpected status code '%s': %s", resp.Status, msg)
	}

	var coursesData sdcListSportsCourseResponse
	if err := json.Unmarshal(bodyBytes, &coursesData); err != nil {
		return nil, fmt.Errorf("unmarshaling courses data: %w", err)
	}

	return coursesData.toCourseDataList(), nil
}

func buildListCoursesFormBody(
	centerID sportsCenterID,
	categoryID categoryID,
	smallCategoryID smallCategoryID,
	courseTargetID courseTargetID,
	courseName string,
) url.Values {
	values := url.Values{}
	values.Add("up_id", string(centerID))
	values.Add("bas_cd", string(categoryID))
	values.Add("item_cd", string(smallCategoryID))
	values.Add("pgm_nm", courseName)
	values.Add("week_nm", "")
	values.Add("target_cd", string(courseTargetID))
	values.Add("page", "0")
	values.Add("perPageNum", "100")
	return values
}

func simplifyResponseBody(msg string) string {
	switch {
	case strings.Contains(msg, "서비스 일시 사용불가") && strings.Contains(msg, "계속될 경우, 관리자에게 문의 바랍니다."):
		return "service temporarily unavailable, ask administrator if continues"
	}
	return msg
}
