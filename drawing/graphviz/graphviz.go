package graphviz

import (
	"github.com/Hrubon/pivot/model"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"fmt"
)

const (
	gvizBin = "/usr/bin/neato"
)

func writeGraph(rlist *model.RouteList, stdin io.WriteCloser) {
	fmt.Fprintf(stdin, "graph G {\n")
	fmt.Fprintf(stdin, "\toverlap = false;\n")
	fmt.Fprintf(stdin, "\tsep = \"+20\";\n")
	nets := make(map[string]string, 0)
	routers := make(map[string]string, 0)
	for i, r := range rlist.Routes {
		if _, ok := nets[r.Network.String()]; !ok {
			n := fmt.Sprintf("net%d", i)
			nets[r.Network.String()] = n
			fmt.Fprintf(stdin, "\t%s [label=\"%v\", shape=box]\n", n, r.Network)
		}
		if _, ok := routers[r.RouterID]; !ok {
			n := fmt.Sprintf("router%d", i)
			routers[r.RouterID] = n
			fmt.Fprintf(stdin, "\t%s [label=\"%s\"]\n", n, r.RouterID)
		}
		fmt.Fprintf(stdin, "\t\"%s\" -- \"%v\"\n [label=%d]",
			routers[r.RouterID], nets[r.Network.String()], r.Metric)
	}
	fmt.Fprintf(stdin, "}\n")
}

func Draw(rlist *model.RouteList, term string, w io.Writer) error {
	cmd := exec.Command(gvizBin, "-T", term)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.Wrapf(err, "cannot pipe stdin of %s", gvizBin)
	}
	cmd.Stdout = w
	cmd.Stderr = os.Stderr // TODO
	writeGraph(rlist, stdin)
	stdin.Close()
	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "cannot start %s", gvizBin)
	}
	if err := cmd.Wait(); err != nil {
		return errors.Wrapf(err, "error running %s", gvizBin)
	}
	return nil
}
//
//type Graphviz struct {
//	Drawer
//	filename string
//}
//
//func NewGraphviz(filename string) *Graphviz {
//	return &Graphviz{
//		filename: filename,
//	}
//}
//
//func (g *Graphviz) Draw(rlist *model.RouteList) error {
//}
