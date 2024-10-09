package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juankohler/crypto-bot/libs/go/errors"
)

func RenderErrorResponse(w http.ResponseWriter, r *http.Request, err error, codeToStatus map[error]int) {
	if err == nil {
		// TODO: add wrn log ?
		return
	}

	statusCode := http.StatusInternalServerError
	for e, s := range codeToStatus {
		if errors.Is(err, e) {
			statusCode = s
			break
		}
	}

	body := map[string]string{
		"code":    "internal",
		"message": err.Error(),
	}
	if err, ok := err.(*errors.Error); ok {
		body = map[string]string{
			"code":    err.Code().String(),
			"message": err.Message(),
		}
	}

	RenderReponse(w, r, body, statusCode)

	logRequestError(err)
}

func RenderReponse(w http.ResponseWriter, r *http.Request, body interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	if statusCode == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.Write([]byte(err.Error()))
		logRequestError(err)
	}
}

func logRequestError(err error) {
	if loggingEnabled {
		// TODO: use logging lib
		fmt.Printf("error response: %v\n", err)
	}
}
