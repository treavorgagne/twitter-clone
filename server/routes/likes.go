package routes

import (
	"net/http"
)

func likeTweet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented) // 501 Not Implemented
}