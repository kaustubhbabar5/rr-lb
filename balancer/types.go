package balancer

type RegisterRequest struct {
	//TODO url validations
	Endpoint string `json:"endpoint" validate:"required"`
}
