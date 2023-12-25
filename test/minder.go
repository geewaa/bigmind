package main

import (
	"log"

	"github.com/geewaa/bigmind/mind"
)

func (m *Minder1) Test() {
	log.Println("ret1111111111111")
}

type Minder1 struct {
	mind.BaseMinder
	// action mind.ACTION
	// envget mind.ENVGET
	// envset mind.ENVSET
	// notice chan string
	// ev     *sync.Map
}
