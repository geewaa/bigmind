package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/geewaa/bigmind/env"
	"github.com/geewaa/bigmind/mind"
)

func main() {
	myenv := env.NewEnv("test")
	myenv.Set(env.EnvVar{
		"times": 1000,
		"max":   1000000000000000000,
		"min":   1,
		"step":  1,
		"name":  "zhang",
	})

	m := &mind.BaseMinder{}
	m.SetAction(calint)

	myenv.AddMinder(m)

	m1 := &Minder1{}
	myenv.AddMinder(m1)
	m1.SetAction(calint1)

	myenv.Start()
	myenv.Snapshot()
}

type MyData struct {
	Times int64  `json:"times"`
	Max   int64  `json:"max"`
	Min   int64  `json:"min"`
	Step  int64  `json:"step"`
	Name  string `json:"name" required:"true" default:"bob"`
	Text  string `json:"text" required:"true" default:"hello world"`
}

func calint(m mind.Minder) error {
	vm := m.(*mind.BaseMinder)
	d := vm.GetEnv()
	v, err := env.LoadEnv[MyData](d)
	if err != nil {
		// log.Println(err)
		return fmt.Errorf("加载数据出错:%w", err)
	}

	log.Printf("env:%#v", v)
	i := v.Times
	for ; i > 0; i -= v.Step {
		n := rand.Int63()
		if n < int64(v.Min) {
			v.Min = n
		} else if n > int64(v.Max) {
			v.Max = n
		}
	}
	log.Println(v.Min, v.Max)
	vm.UpdateEnv(
		map[string]any{
			"min":   v.Min,
			"max":   v.Max,
			"times": i,
			"name":  "li",
		})
	log.Println("li done")

	return nil
}

// /////////////////////////////////////////////////////

type MyData1 struct {
	Times int64  `json:"times"`
	Max   int64  `json:"max"`
	Min   int64  `json:"min"`
	Step  int64  `json:"step"`
	Name  string `json:"name" required:"true" default:"niel"`
	Text  string `json:"text" required:"true" default:"hello world"`
}

func calint1(m mind.Minder) error {
	vm := m.(*Minder1)
	d := vm.GetEnv()
	v, err := env.LoadEnv[MyData1](d)
	if err != nil {
		// log.Println(err)
		return fmt.Errorf("加载数据出错:%w", err)
	}

	v.Step = 3
	log.Printf("env:%#v", v)
	i := v.Times

	for ; i > 0; i -= v.Step {
		n := rand.Int63()
		if n < int64(v.Min) {
			v.Min = n
		} else if n > int64(v.Max) {
			v.Max = n
		}
	}
	vm.Test()
	log.Println(v.Name, v.Min, v.Max)
	vm.UpdateEnv(
		map[string]any{
			"min":   v.Min,
			"max":   v.Max,
			"times": i,
			"name":  "xu",
		})
	log.Println("xu done")

	return nil
}
