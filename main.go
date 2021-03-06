package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/src-d/babelfish-go-driver/msg"
)

const (
	lang = "Go"
)

var (
	langVersion   = runtime.Version()
	driverVersion string
)

func main() {
	in := os.Stdin
	out := os.Stdout

	if err := start(in, out); err != nil {
		log.Fatal(err)
	}
}

// start launchs a loop to read requests and write responses.
func start(in io.Reader, out io.Writer) error {
	dec := json.NewDecoder(in)
	enc := json.NewEncoder(out)
	req := &msg.Request{}
	var res *msg.Response

	for {
		if err := dec.Decode(req); err != nil {
			if err == io.EOF {
				break
			}

			res = &msg.Response{
				Status:          msg.Fatal,
				Errors:          []string{err.Error()},
				Language:        lang,
				LanguageVersion: langVersion,
				Driver:          driverVersion,
			}

			if encErr := enc.Encode(res); encErr != nil {
				return fmt.Errorf("%v: %v", err, encErr)
			}

			return err
		}

		res = getResponse(req)
		if err := enc.Encode(res); err != nil {
			return err
		}
	}

	return nil
}

// getResponse always generates a msg.Response. The response will have the properly status (Ok, Error, Fatal).
func getResponse(m *msg.Request) *msg.Response {
	res := &msg.Response{
		Language:        lang,
		LanguageVersion: langVersion,
		Driver:          driverVersion,
	}

	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, "source.go", m.Content, parser.ParseComments|parser.AllErrors)
	if err != nil {
		if tree == nil {
			res.Status = msg.Fatal
			res.Errors = []string{err.Error()}
			return res
		}

		res.Status = msg.Error
		errList := err.(scanner.ErrorList)
		res.Errors = getErrors(errList)
	} else {
		res.Status = msg.Ok
	}

	ast.Inspect(tree, setObjNil)
	res.AST = tree

	return res
}

// getErrors build a []string with the err.Error() from a scanner.ErrorList.
func getErrors(errList scanner.ErrorList) []string {
	list := make([]string, 0, len(errList))
	for _, err := range errList {
		list = append(list, err.Error())
	}

	return list
}
