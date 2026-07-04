package request

type UpdateService struct {
	Name *string `json:"name"`
}

type CreateEndpoint struct {
	Method string `json:"method" validate:"required"`
	URL    string `json:"url"    validate:"required"`
	Secure bool   `json:"secure"`
}

type CreateService struct {
	Name      string           `json:"name"      validate:"required"`
	Endpoints []CreateEndpoint `json:"endpoints" validate:"required,dive"`
}
