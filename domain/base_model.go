package domain

type (
	HttpResponse struct {
		Code       int         `json:"code"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data"`
		Pagination *Pagination `json:"pagination,omitempty"`
	}
)
