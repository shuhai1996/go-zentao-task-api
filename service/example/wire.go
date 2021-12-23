//go:build wireinject
// +build wireinject

package example

import (
	"github.com/google/wire"
	"go-zentao-task/model/example"
)

type ServiceExample struct {
	W1 *example.WireExample1
	W2 *example.WireExample2
}

func NewServiceExample(w1 *example.WireExample1, w2 *example.WireExample2) *ServiceExample {
	return &ServiceExample{W1: w1, W2: w2}
}

func InitializeServiceExample() *ServiceExample {
	wire.Build(NewServiceExample, example.NewWireExample1, example.NewWireExample2)
	return &ServiceExample{}
}

// same as above
var Set = wire.NewSet(
	example.NewWireExample1,
	example.NewWireExample2,
	wire.Struct(new(ServiceExample), "*"),
)

func InitializeServiceExample2() *ServiceExample {
	wire.Build(Set)
	return &ServiceExample{}
}
