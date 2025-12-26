package http

import (
	"github.com/hammer-code/lms-be/domain"
)

type Handler struct {
	usecase domain.BlogPostUsecase
}

var (
	handlr *Handler
)

func NewHandler(usecase domain.BlogPostUsecase) domain.BlogPostHandler {
	if handlr == nil {
		handlr = &Handler{
			usecase: usecase,
		}

	}
	return *handlr
}
