package main

import (
	"github.com/Hrubon/pivot/sources"
	"github.com/Hrubon/pivot/terminal"
	"log"
)

func main() {
	s := sources.NewBIRDSource("/tmp/bird.1.ctl")
	rlist, err := s.GetRoutes()
	if err != nil {
		log.Fatal(err)
	}
	gv := terminal.NewGraphviz("/tmp/map.pdf")
	if err := gv.Draw(rlist); err != nil {
		log.Fatal(err)
	}
	ui := terminal.NewWebUI("localhost", 9000)
	ui.Start()
}
