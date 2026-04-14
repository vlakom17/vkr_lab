package request

type SetReactionRequest struct {
	Type string `json:"type" binding:"required"`
}
