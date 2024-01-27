package server

import (
	"log"
	"net/http"

	"github.com/hellofresh/health-go/v5"
)

type Server struct{}

func NewServer() *Server {

	return &Server{}
}

func (s *Server) Run(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		// status ok
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from server " + addr))

		log.Println("Request received on server ", addr)
	})

	// add some checks on instance creation
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    "node" + addr,
		Version: "v1.0",
	}))

	http.Handle("/status", h.Handler())

	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal(err)
	}
}
