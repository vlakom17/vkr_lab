package request

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password"`
	About    *string `json:"about"`
}
