package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Khabi/chromatic/internal/chromatic"
	"github.com/gorilla/mux"
)

type service struct {
	command chan chromatic.State
	status  chan chromatic.ServerStatus
}

func Run(bind string, command chan chromatic.State, status chan chromatic.ServerStatus) {
	s := service{
		command: command,
		status:  status,
	}

	r := mux.NewRouter()
	r.HandleFunc("/action/{key}", s.Action)
	r.HandleFunc("/status", s.Status)
	srv := &http.Server{
		Handler:      r,
		Addr:         bind,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s service) Action(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch vars["key"] {
	case "start":
		s.command <- chromatic.Running
	case "pause":
		s.command <- chromatic.Paused
	case "stop":
		s.command <- chromatic.Stop
	}
}

func (s service) Status(w http.ResponseWriter, r *http.Request) {
	s.command <- chromatic.Status
	status := <-s.status

	resp := make(map[string]interface{})
	resp["state"] = status.State
	resp["fps"] = status.FPS
	json.NewEncoder(w).Encode(status)
}
