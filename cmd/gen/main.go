package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"text/template"
)

type Handler struct {
	Name        string // Name of the original function to execute
	RequestName string // Type name of the request object passed into the original function
}

type Params struct {
	Handlers []Handler
}

func parseFile(path string) []Handler {
	handlers := []Handler{}

	fset := token.NewFileSet()
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	pf, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range pf.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Naive way to filter out functions which don't match the gRPC function-like criteria
		if len(fd.Type.Params.List) != 2 || fd.Type.Results.List == nil || len(fd.Type.Results.List) != 2 {
			continue
		}

		name := fd.Name.Name
		request := fd.Type.Params.List[1].Type
		// We can also use types.ExprString(request) to grab the type name
		// Do not subtract 1 from Pos() to avoid grabbing the asterisk (*)
		requestName := string(file)[request.Pos() : request.End()-1]

		handlers = append(handlers, Handler{Name: name, RequestName: requestName})
	}

	return handlers
}

func main() {
	params := Params{Handlers: []Handler{}}
	tmpl, err := template.New("template.go.tmpl").ParseFiles("cmd/gen/template.go.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	params.Handlers = append(params.Handlers, parseFile("internal/server/songs.go")...)
	params.Handlers = append(params.Handlers, parseFile("internal/server/svelte.go")...)

	outFile, err := os.Create("internal/server/handlers.gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, params)
	if err != nil {
		log.Fatal(err)
	}
}
