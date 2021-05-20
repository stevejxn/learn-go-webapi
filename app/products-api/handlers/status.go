package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func statusPingHandler(w http.ResponseWriter, r *http.Request) {
	pingResponse := fmt.Sprintf("pong: %v", time.Now().UTC().Format(time.RFC3339))
	w.Write([]byte(pingResponse))
}
