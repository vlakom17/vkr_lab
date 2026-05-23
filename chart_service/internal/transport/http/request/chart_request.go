package request

type CreateChartRequest struct {
	Title         string `json:"title" binding:"required"`
	Genre         string `json:"genre" binding:"required"`
	Description   string `json:"description"`
	PositionCount int    `json:"position_count" binding:"required"`
}

type PatchChartRequest struct {
	Title         *string `json:"title"`
	Genre         *string `json:"genre"`
	Description   *string `json:"description"`
	PositionCount *int    `json:"position_count"`
}
