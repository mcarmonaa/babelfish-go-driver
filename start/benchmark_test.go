package start

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/src-d/babelfish-go-driver/msg"

	"github.com/ugorji/go/codec"
)

var (
	reqBench = loadFile("../testfiles/test4.source")
	resBench = &msg.Response{
		Status:          msg.Ok,
		Driver:          "betatesting",
		Language:        "Go",
		LanguageVersion: "gotesting",
		AST:             getTree(reqBench.Content),
	}
)

// loadFile generates a msg.Request with the content from a file.
func loadFile(name string) *msg.Request {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	source, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return &msg.Request{
		Action:  msg.ParseAst,
		Content: string(source),
	}
}

// getTree get the ast from a source.
func getTree(source string) *ast.File {
	fset := token.NewFileSet()
	tree, _ := parser.ParseFile(fset, "source.go", source, parser.ParseComments)
	ast.Inspect(tree, setObjNil)

	return tree
}

func BenchmarkSerializeMsgpckResponse(b *testing.B) {
	var handle codec.MsgpackHandle
	enc := codec.NewEncoder(ioutil.Discard, &handle)
	b.Run("Msgpack serialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			enc.MustEncode(resBench)
		}
	})

}

func BenchmarkSerializeJSONResponse(b *testing.B) {
	var handle codec.JsonHandle
	enc := codec.NewEncoder(ioutil.Discard, &handle)
	b.Run("JSON serialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			enc.MustEncode(resBench)
		}
	})
}

func BenchmarkSerializeStdJSONResponse(b *testing.B) {
	enc := json.NewEncoder(ioutil.Discard)
	b.Run("StdJSON serialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			enc.Encode(resBench)
		}
	})
}

func BenchmarkDeserializedMsgpckRequest(b *testing.B) {
	buf := &bytes.Buffer{}
	var eHandle codec.MsgpackHandle
	enc := codec.NewEncoder(buf, &eHandle)
	enc.MustEncode(reqBench)
	in := buf.Bytes()
	var dHandle codec.MsgpackHandle
	b.Run("Msgpack deserialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec := codec.NewDecoderBytes(in, &dHandle)
			var gotReq interface{}
			if err := dec.Decode(&gotReq); err != nil {
				if err != io.EOF {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkDeserializeJSONRequest(b *testing.B) {
	buf := &bytes.Buffer{}
	var eHandle codec.JsonHandle
	enc := codec.NewEncoder(buf, &eHandle)
	enc.MustEncode(reqBench)
	in := buf.Bytes()
	var dHandle codec.JsonHandle
	b.Run("JSON deserialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec := codec.NewDecoderBytes(in, &dHandle)
			var gotReq interface{}
			if err := dec.Decode(&gotReq); err != nil {
				if err != io.EOF {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkDeserializeStdJSONRequest(b *testing.B) {
	buf := &bytes.Buffer{}
	var eHandle codec.JsonHandle
	enc := codec.NewEncoder(buf, &eHandle)
	enc.MustEncode(reqBench)
	in := buf.Bytes()
	b.Run("StdJSON deserialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec := json.NewDecoder(bytes.NewBuffer(in))
			var gotReq interface{}
			if err := dec.Decode(&gotReq); err != nil {
				if err != io.EOF {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkCompleteMsgPack(b *testing.B) {
	buf := &bytes.Buffer{}
	var eHandle codec.MsgpackHandle
	enc := codec.NewEncoder(buf, &eHandle)
	enc.MustEncode(reqBench)
	in := buf.Bytes()
	b.Run("Full cycle Msgpck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			StartMsgpck(bytes.NewBuffer(in), ioutil.Discard)
		}
	})
}

func BenchmarkCompleteJSON(b *testing.B) {
	buf := &bytes.Buffer{}
	var eHandle codec.JsonHandle
	enc := codec.NewEncoder(buf, &eHandle)
	enc.MustEncode(reqBench)
	in := buf.Bytes()
	b.Run("Full cycle JSON", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			StartJSON(bytes.NewBuffer(in), ioutil.Discard)
		}
	})
}

func BenchmarkCompleteStdJSON(b *testing.B) {
	buf := &bytes.Buffer{}
	var eHandle codec.JsonHandle
	enc := codec.NewEncoder(buf, &eHandle)
	enc.MustEncode(reqBench)
	in := buf.Bytes()
	b.Run("Full cycle Standard JSON", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			StartStdJSON(bytes.NewBuffer(in), ioutil.Discard)
		}
	})
}
