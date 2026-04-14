package response

import "github.com/google/uuid"

type ChartIDsResponse struct {
	ChartIDs []uuid.UUID `json:"chart_ids"`
}
