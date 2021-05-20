package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func API(buildNo string, shutdown chan os.Signal, log *log.Logger) http.Handler {
	api := http.NewServeMux()
	api.HandleFunc("/status/ping", pingHandler)
	return api
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	pingResponse := fmt.Sprintf("pong: %v", time.Now().UTC().Format(time.RFC3339))
	w.Write([]byte(pingResponse))
}
