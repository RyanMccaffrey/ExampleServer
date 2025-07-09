package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type ExampleServerService struct {
	addr string
}

func NewExampleServerService(addr string) *ExampleServerService {
	return &ExampleServerService{
		addr: addr,
	}
}
func (s *ExampleServerService) Serve() error {

	if s.addr == "" {
		s.addr = ":8080"
	}
	for {
		log.Printf("Starting server on %s", s.addr)
		fmt.Println("Server is running at", s.addr)
		err := http.ListenAndServe(s.addr, nil)
		if err != nil {
			fmt.Println("Server failed:", err)
			log.Printf("Server failed: %v. Restarting in 5 seconds...", err)
			time.Sleep(5 * time.Second)
		}
	}
}
