package drawing

import (
	"github.com/Hrubon/pivot/model"
	"net/http"
	"fmt"
)

func init() {
	http.HandleFunc("/graph", handleGraph);
}

type WebUI struct {
	Drawer
	addr   string
	port   int
	err    error
}

func NewWebUI(addr string, port int) *WebUI {
	return &WebUI{
		addr: addr,
		port: port,
	}
}

func handleGraph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{}\n");
}

func (w *WebUI) Start() {
	http.ListenAndServe(fmt.Sprintf("%s:%d", w.addr, w.port), nil)
}

func (w *WebUI) Stop() error {
	return nil
}

func (w *WebUI) Draw(rlist *model.RouteList) error {
	return nil
}
