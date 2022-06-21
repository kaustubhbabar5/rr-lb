package balancer

type RegisterRequest struct {
	//TODO url validations
	Url string `json:"url" validate:"required"`
}
