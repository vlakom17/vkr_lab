package response

type ChartResponse struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Genre         string `json:"genre"`
	PositionCount int    `json:"position_count"`
	Description   string `json:"description"`
	UserID        string `json:"user_id"`
}
