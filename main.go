package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/bundle"
	"github.com/open-policy-agent/opa/rego"
)

func main() {
	// first, the simple case where we just need to evaludate some rego in
	// different versions init'ing Rego each time
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

	fmt.Println("-----------------")

	// second, the more involved case where rego can be initialized once for
	// bundles of different versions.

	mod0, err := ast.ParseModuleWithOpts(
		"policy.rego",
		`package v0

messages[msg] {
	msg := "foo"
}
`,
		ast.ParserOptions{
			RegoVersion: ast.RegoV0,
		},
	)
	if err != nil {
		panic(err)
	}

	mod1, err := ast.ParseModuleWithOpts(
		"policy.rego",
		`package v1

messages contains "bar"
`,
		ast.ParserOptions{
			RegoVersion: ast.RegoV1,
		},
	)
	if err != nil {
		panic(err)
	}

	bv0 := bundle.Bundle{
		Manifest: bundle.Manifest{
			Roots: &[]string{"examplev0"},
		},
		Modules: []bundle.ModuleFile{
			{
				Parsed: mod0,
			},
		},
	}

	bv0.SetRegoVersion(ast.RegoV0)

	bv1 := bundle.Bundle{
		Manifest: bundle.Manifest{
			Roots: &[]string{"examplev1"},
		},
		Modules: []bundle.ModuleFile{
			{
				Parsed: mod1,
			},
		},
	}

	bv1.SetRegoVersion(ast.RegoV1)

	r := rego.New(
		rego.ParsedBundle("v0", &bv0),
		rego.ParsedBundle("v1", &bv1),
		// query for both v0 and v1 policies at the same time
		rego.Query("data.v1.messages | data.v0.messages"),
	)

	rs, err := r.Eval(context.Background())
	if err != nil {
		panic(err)
	}

	bs, err := json.Marshal(rs)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))
}
