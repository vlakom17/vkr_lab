package utilits

import (
	"charts-chart-service/internal/domain/chart"
	"charts-chart-service/internal/transport/http/response"
)

func MapChartToResponse(c chart.Chart) response.ChartResponse {
	return response.ChartResponse{
		ID:            c.ID.String(),
		Title:         c.Title,
		Genre:         c.Genre,
		PositionCount: c.PositionCount,
		Description:   c.Description,
		UserID:        c.UserID.String(),
	}
}
