package env

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/geewaa/bigmind/mind"
)

type EnvVar map[string]any

type Env struct {
	name    string
	ev      *sync.Map     // ENV VALUE STORE
	ms      []mind.Minder // MINDER
	notices []chan string // DATA CHANGE NOTICE
}

func NewEnv(name string) *Env {
	return &Env{
		name:    name,
		ev:      &sync.Map{},
		notices: []chan string{},
	}
}

func (e *Env) AddMinder(m mind.Minder) {
	notice := make(chan string, 100)
	m.Init(notice, e.Get, e.Set)
	// m.SetENVGET(e.Get)
	// m.SetENVSET(e.Set)
	e.ms = append(e.ms, m)
	e.notices = append(e.notices, notice)
}

func (e *Env) Get() map[string]any {
	ev := make(EnvVar)
	e.ev.Range(
		func(key, value any) bool {
			ev[key.(string)] = value
			return true
		})
	return ev
}

func (e *Env) Set(v map[string]any) {
	for k, v := range v {
		e.Update(k, v)
	}
	for _, ntc := range e.notices {
		ntc <- "upgrade"
	}
}

func (e *Env) Update(key string, value any) {
	e.ev.Store(key, value)
}

func (e *Env) Start() {
	wg := &sync.WaitGroup{}
	wg.Add(len(e.ms))

	for _, m := range e.ms {
		go m.Work(wg, m)
	}
	wg.Wait()
}

func (e *Env) Snapshot() {
	e.ev.Range(
		func(key, value any) bool {
			fmt.Println(key, ":", value)
			return true
		})
}

func LoadEnv[T any](ev *sync.Map) (*T, error) {
	v := new(T)
	s := reflect.TypeOf(v).Elem()
	if s.Kind() != reflect.Struct {
		panic("LoadEnv: not a struct")
	}
	num := s.NumField()
	for i := range num {
		f := s.Field(i)
		key := f.Tag.Get("json")
		if key == "" || key == "-" {
			continue
		}
		vl, ok := ev.Load(key)
		if ok {
			reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(vl).Convert(f.Type))
		} else {
			if f.Tag.Get("required") == "true" && f.Tag.Get("default") == "" {
				return nil, fmt.Errorf("LoadEnv: %s is required", key)
			}
			if f.Tag.Get("default") != "" {
				reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(f.Tag.Get("default")).Convert(f.Type))
			} else if f.Tag.Get("default") == "" && f.Tag.Get("required") != "true" {
				reflect.ValueOf(v).Elem().Field(i).Set(reflect.Zero(f.Type))
			}
		}

	}
	return v, nil
}
