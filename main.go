package main

import (
	"github.com/Hrubon/pivot/sources"
	"github.com/Hrubon/pivot/terminal"
	"log"
)

func main() {
	s := sources.NewBIRDSource("/tmp/bird.2.ctl")
	rlist, err := s.GetRoutes()
	if err != nil {
		log.Fatal(err)
	}
	gv := terminal.NewGraphviz("/tmp/map.pdf")
	if err := gv.Draw(rlist); err != nil {
		log.Fatal(err)
	}
}
