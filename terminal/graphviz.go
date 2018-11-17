package terminal

import (
	"github.com/Hrubon/pivot/model"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"fmt"
)

const (
	gvizBin = "/usr/bin/neato"
)

type Graphviz struct {
	Terminal
	filename string
}

func NewGraphviz(filename string) *Graphviz {
	return &Graphviz{
		filename: filename,
	}
}

func (g *Graphviz) Draw(rlist *model.RouteList) error {
	cmd := exec.Command(gvizBin, "-Tpdf", "-o", g.filename)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.Wrapf(err, "cannot pipe stdin to %s", gvizBin)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "cannot start %s", gvizBin)
	}
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
	stdin.Close()
	if err := cmd.Wait(); err != nil {
		return errors.Wrapf(err, "error running %s", gvizBin)
	}
	return nil
}
