package webui

import (
	"github.com/Hrubon/pivot/model"
	"github.com/Hrubon/pivot/drawing/graphviz"
	"net/http"
	"fmt"
)

func init() {
}

type Server struct {
	addr   string
	port   int
	rlist  *model.RouteList
}

func NewServer(addr string, port int) *Server {
	return &Server{
		addr: addr,
		port: port,
	}
}

func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
	if s.rlist == nil {
		fmt.Fprintf(w, "{}\n")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	graphviz.Draw(s.rlist, "json", w)
}

func (s *Server) Start() {
	fs := http.FileServer(http.Dir("/home/d/go/src/github.com/Hrubon/pivot/webui/static"))
	http.Handle("/", fs)
	http.HandleFunc("/graph.json", s.handleGraph);
	http.HandleFunc("/graph.png", s.handleGraph);
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.addr, s.port), nil)
}

func (s *Server) Push(rlist *model.RouteList) error {
	s.rlist = rlist
	return nil
}

func (s *Server) Stop() error {
	return nil
}
