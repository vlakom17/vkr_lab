package request

type CreateChartRequest struct {
	Title         string `json:"title" binding:"required"`
	Genre         string `json:"genre" binding:"required"`
	PositionCount int    `json:"position_count" binding:"required"`
	Description   string `json:"description"`
}

type PatchChartRequest struct {
	Title         *string `json:"title"`
	Genre         *string `json:"genre"`
	PositionCount *int    `json:"position_count"`
	Description   *string `json:"description"`
}
