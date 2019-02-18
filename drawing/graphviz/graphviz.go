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

type Pos struct {
	X float64
	Y float64
}

type PinList map[string]Pos

func writeGraph(rlist *model.RouteList, pins PinList, stdin io.WriteCloser) {
	fmt.Fprintf(stdin, "graph G {\n")
	fmt.Fprintf(stdin, "\toverlap = false;\n")
	fmt.Fprintf(stdin, "\tsep = \"+20\";\n")
	//fmt.Fprintf(stdin, "\tnotranslate = true;\n")
	nets := make(map[string]string, 0)
	routers := make(map[string]string, 0)
	for i, r := range rlist.Routes {
		if _, ok := nets[r.Network.String()]; !ok {
			n := fmt.Sprintf("net%d", i)
			nets[r.Network.String()] = n
			fmt.Fprintf(stdin, "\t%s [", n)
			fmt.Fprintf(stdin, "label=%q,", r.Network)
			fmt.Fprintf(stdin, "shape=box,")
			//if pos, ok := pins[n]; ok {
			//	fmt.Fprintf(stdin, "pos=\"%f,%f!\",", pos.X, pos.Y)
			//}
			fmt.Fprintf(stdin, "]\n")
		}
		if _, ok := routers[r.RouterID]; !ok {
			n := fmt.Sprintf("router%d", i)
			routers[r.RouterID] = n
			fmt.Fprintf(stdin, "\t%s [", n)
			fmt.Fprintf(stdin, "label=%q,", r.RouterID)
			//if pos, ok := pins[n]; ok {
			//	fmt.Fprintf(stdin, "pos=\"%f,%f!\",", pos.X, pos.Y)
			//}
			fmt.Fprintf(stdin, "]")
		}
		fmt.Fprintf(stdin, "\t\"%s\" -- \"%v\"\n [label=%d]",
			routers[r.RouterID], nets[r.Network.String()], r.Metric)
	}
	fmt.Fprintf(stdin, "}\n")
}

func Draw(rlist *model.RouteList, pins PinList, term string, w io.Writer) error {
	cmd := exec.Command(gvizBin, "-T", term)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.Wrapf(err, "cannot pipe stdin of %s", gvizBin)
	}
	cmd.Stdout = w
	cmd.Stderr = os.Stderr // TODO
	writeGraph(rlist, pins, stdin)
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
