package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"testing"

	"github.com/ugorji/go/codec"
)

var (
	resBench = tests[4].res
	reqBench = tests[4].req
)

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
	b.Run("Full cycle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			start(bytes.NewBuffer(in), ioutil.Discard)
		}
	})
}
