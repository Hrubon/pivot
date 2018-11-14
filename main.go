package main

import (
	"fmt"
	"github.com/Hrubon/pivot/sources"
	"log"
)

func main() {
	s := sources.NewBirdSource("/tmp/bird.1.ctl")
	rlist, err := s.GetRoutes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rlist)
}
