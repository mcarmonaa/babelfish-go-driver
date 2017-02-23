package start

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"

	"github.com/src-d/babelfish-go-driver/msg"

	"github.com/ugorji/go/codec"
)

const (
	lang          = "Go"
	langVersion   = "go-testing-version"
	driverVersion = "testing-version"
)

// startMsgpck launchs a loop to read requests and write responses. Msgpack serialize.
func StartMsgpck(in io.Reader, out io.Writer) error {
	var handle codec.MsgpackHandle
	dec := codec.NewDecoder(in, &handle)
	enc := codec.NewEncoder(out, &handle)
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

			enc.MustEncode(res)
			return err
		}

		res = getResponse(req)
		enc.MustEncode(res)
	}

	return nil
}

// startJSON launchs a loop to read requests and write responses. JSON serialize.
func StartJSON(in io.Reader, out io.Writer) error {
	var handle codec.JsonHandle
	dec := codec.NewDecoder(in, &handle)
	enc := codec.NewEncoder(out, &handle)
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

			enc.MustEncode(res)
			return err
		}

		res = getResponse(req)
		enc.MustEncode(res)
	}

	return nil
}

// startStdJSON launchs a loop to read requests and write responses. Standard JSON serialize.
func StartStdJSON(in io.Reader, out io.Writer) error {
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

			enc.Encode(res)
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

// setObjNil looks for the elements that can't be serialized and set it to nil.
// It has the properly signature to be a parameter of ast.Inspect function.
func setObjNil(node ast.Node) bool {
	switch node.(type) {
	case *ast.BranchStmt:
		n := node.(*ast.BranchStmt)
		safeIdent(n.Label)
	case *ast.File:
		n := node.(*ast.File)
		safeImportSpecList(n.Imports)
		safeIdent(n.Name)
		safeScope(n.Scope)
		safeIdentList(n.Unresolved)
	case *ast.FuncDecl:
		n := node.(*ast.FuncDecl)
		safeIdent(n.Name)
		safeFieldList(n.Recv)
		safeFuncTye(n.Type)
	case *ast.FuncLit:
		n := node.(*ast.FuncLit)
		safeFuncTye(n.Type)
	case *ast.FuncType:
		n := node.(*ast.FuncType)
		safeFuncTye(n)
	case *ast.Ident:
		n := node.(*ast.Ident)
		safeIdent(n)
	case *ast.ImportSpec:
		n := node.(*ast.ImportSpec)
		safeImportSpec(n)
	case *ast.InterfaceType:
		n := node.(*ast.InterfaceType)
		safeFieldList(n.Methods)
	case *ast.LabeledStmt:
		n := node.(*ast.LabeledStmt)
		safeIdent(n.Label)
	case *ast.Package:
		n := node.(*ast.Package)
		n.Files = nil
		n.Imports = nil
		safeScope(n.Scope)
	case *ast.SelectorExpr:
		n := node.(*ast.SelectorExpr)
		safeIdent(n.Sel)
	case *ast.StructType:
		n := node.(*ast.StructType)
		safeFieldList(n.Fields)
	case *ast.TypeSpec:
		n := node.(*ast.TypeSpec)
		safeIdent(n.Name)
	}

	return true
}

// safeIdent set to nil conflictives fields of a ast.Ident.
func safeIdent(node *ast.Ident) {
	if node == nil {
		return
	}

	node.Obj = nil
}

// safeIdentList iterates over a slice of ast.Ident and calls safeIdent.
func safeIdentList(list []*ast.Ident) {
	for i := range list {
		safeIdent(list[i])
	}
}

// safeImportSpec set to nil conflictives fields of a ast.ImportSpect.
func safeImportSpec(is *ast.ImportSpec) {
	safeIdent(is.Name)
}

// safeImportSpecList iterates over a slice of ast.ImportSpec and calls safeImportSpec.
func safeImportSpecList(list []*ast.ImportSpec) {
	for i := range list {
		safeImportSpec(list[i])
	}
}

// safeField set to nil conflictives fields of a ast.Field.
func safeField(field *ast.Field) {
	safeIdentList(field.Names)
}

// safeFieldList iterates over a slice of ast.Field and calls safeField.
func safeFieldList(flist *ast.FieldList) {
	if flist == nil {
		return
	}

	for i := range flist.List {
		safeField(flist.List[i])
	}
}

// safeFuncTye set to nil conflictives fields of a ast.FuncType.
func safeFuncTye(ftype *ast.FuncType) {
	safeFieldList(ftype.Params)
	safeFieldList(ftype.Results)
}

// safeScope set to nil conflictives fields of a ast.Scope.
func safeScope(scope *ast.Scope) {
	scope.Objects = nil
	scope.Outer = nil
}
