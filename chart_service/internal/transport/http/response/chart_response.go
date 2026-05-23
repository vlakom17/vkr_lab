package response

type ChartResponse struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Title         string `json:"title"`
	Genre         string `json:"genre"`
	Description   string `json:"description"`
	PositionCount int    `json:"position_count"`
}
