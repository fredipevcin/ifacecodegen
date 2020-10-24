package examples

import "context"

// Entity is a sample struct
type Entity struct{}

// Foo demo interface
type Foo interface {
	Load(context.Context) (Entity, error)
	Save([]Entity) error
	IsValid() bool
	ValidateMulti(...Entity)
	Multi(p1, p2 string) (r1, r2 string)
}
