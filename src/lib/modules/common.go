package modules

// Module is the interface implemented by types that
// register themselves as modular plug-ins.
type Module interface {
	Init() error
	Version() string
}