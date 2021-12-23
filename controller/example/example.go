package example

import (
	"go-zentao-task/core"
	"go-zentao-task/service/example"
)

var serviceExample = example.InitializeServiceExample()

// var serviceExample = example.InitializeServiceExample2()

func Wire(c *core.Context) {
	// invoke service and model
	model1Result := serviceExample.W1.A()
	model2Result := serviceExample.W2.A()
	c.Success(model1Result + " and " + model2Result)
}
