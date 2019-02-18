package webui

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Hrubon/pivot/drawing/graphviz"
	"github.com/Hrubon/pivot/model"
	"github.com/pkg/errors"
	"net/http"
)

func init() {
}

type Server struct {
	addr  string
	port  int
	rlist *model.RouteList
	pins graphviz.PinList
}

func NewServer(addr string, port int) *Server {
	return &Server{
		addr: addr,
		port: port,
		pins: graphviz.PinList{},
	}
}

type MoveRequest struct {
	Name string
	X    float64
	Y    float64
}

type GraphResponse struct {
	GvizData json.RawMessage
}

func decodePOST(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Method != "POST" {
		err := errors.New("end-point must be called using POST")
		http.Error(w, err.Error(), 405)
		return err
	}
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		err2 := errors.Wrap(err, "cannot decode request body as JSON")
		http.Error(w, err2.Error(), 400)
		return err2
	}
	return nil
}

func sendJSON(w http.ResponseWriter, src interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(src); err != nil {
		err2 := errors.Wrap(err, "cannot encode response body as JSON")
		panic(err2)
	}
	return nil
}

func (s *Server) graphJSON() []byte {
	var gvizData bytes.Buffer
	graphviz.Draw(s.rlist, s.pins, "json", bufio.NewWriter(&gvizData))
	return gvizData.Bytes()
}

func (s *Server) getGraph(w http.ResponseWriter, r *http.Request) {
	if s.rlist == nil {
		fmt.Fprintf(w, "{}\n")
		return
	}
	sendJSON(w, GraphResponse{
		GvizData: s.graphJSON(),
	})
}

func (s *Server) postMove(w http.ResponseWriter, r *http.Request) {
	var move MoveRequest
	if decodePOST(w, r, &move) != nil {
		return
	}
	s.pins[move.Name] = graphviz.Pos{move.X, move.Y}
	sendJSON(w, GraphResponse{
		GvizData: s.graphJSON(),
	})
}

func (s *Server) Start() {
	// TODO relpath
	fs := http.FileServer(http.Dir("/home/d/go/src/github.com/Hrubon/pivot/webui/static"))
	http.HandleFunc("/graph.json", s.getGraph)
	http.HandleFunc("/graph.png", s.getGraph)
	http.HandleFunc("/move", s.postMove)
	http.Handle("/", fs)
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.addr, s.port), nil)
}

func (s *Server) Push(rlist *model.RouteList) error {
	s.rlist = rlist
	return nil
}

func (s *Server) Stop() error {
	return nil
}
