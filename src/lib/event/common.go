package event

import "sync"

type GoWikiEvents struct {
	events map[string][]func(...interface{})
}

var instance *GoWikiEvents
var once sync.Once

func Events() *GoWikiEvents {
	once.Do(func() {
		instance = &GoWikiEvents{
			events: make(map[string][]func(...interface{})),
		}
	})
	return instance
}

func (e *GoWikiEvents) On(event string, do func(...interface{})) {
	e.events[event] = append(e.events[event], do)
}

func (e *GoWikiEvents) Emit(event string, args ...interface{}) {
	for _, f := range e.events[event] {
		f(args...)
	}
}