package main

import (
	"github.com/Hrubon/pivot/sources"
	"github.com/Hrubon/pivot/webui"
	"log"
)

func main() {
	s := sources.NewBIRDSource("/tmp/bird.1.ctl")
	rlist, err := s.GetRoutes()
	if err != nil {
		log.Fatal(err)
	}
	//gv := drawing.NewGraphviz("/tmp/map.pdf")
	//if err := gv.Draw(rlist); err != nil {
	//	log.Fatal(err)
	//}
	ui := webui.NewServer("localhost", 9000)
	ui.Push(rlist)
	ui.Start()
}
