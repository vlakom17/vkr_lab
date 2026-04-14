package response

type TrackResponse struct {
	ID     string `json:"id"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
}
