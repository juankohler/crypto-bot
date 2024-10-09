package bots

import (
	"github.com/juankohler/crypto-bot/common"
)

// var errorsToCode = map[error]int{
// 	domain.ErrInternal: http.StatusInternalServerError,
// 	domain.ErrNotFound: http.StatusNotFound,
// 	domain.ErrInvalid:  http.StatusBadRequest,
// }

type Handlers struct {
}

func NewHandlers(cfg *common.Config, deps *Dependencies) *Handlers {
	return &Handlers{}
}
