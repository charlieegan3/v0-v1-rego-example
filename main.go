package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func main() {
	examples := []struct {
		RegoVersion ast.RegoVersion
		Module      string
	}{
		{
			RegoVersion: ast.RegoV1,
			Module: `package example

messages contains "foo"
`,
		},
		{
			RegoVersion: ast.RegoV0,
			Module: `package example

messages[msg] {
	msg := "foo"
}
`,
		},
	}

	for _, ex := range examples {
		fmt.Println("Rego Version:", ex.RegoVersion)
		r := rego.New(
			rego.Query("data.example.messages"),
			rego.SetRegoVersion(ex.RegoVersion),
			rego.Module("example.rego", ex.Module),
		)

		rs, err := r.Eval(context.TODO())
		if err != nil {
			panic(err)
		}

		bs, err := json.Marshal(rs)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(bs))
	}
}
