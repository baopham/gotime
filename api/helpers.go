package api

import (
	"github.com/baopham/gotime/gotime"
	"log"
	"net/http"
)

func handleError(err error, w http.ResponseWriter, service gotime.GoTimer) {
	log.Println("error: ", err)

	message := "Something went wrong"

	if service.IsRateLimitError(err) {
		message = "API rate limit. Supply a token to increase the limit"
	}

	http.Error(w, message, http.StatusBadRequest)
}
