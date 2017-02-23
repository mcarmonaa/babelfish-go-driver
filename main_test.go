package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"

	"github.com/src-d/babelfish-go-driver/msg"

	"github.com/stretchr/testify/require"
)

const (
	driverTestVersion = "beta-testing-driver"
)

var tests = []*myTest{
	0: newMyTest("statusError", &msg.Request{Action: msg.ParseAst},
		msg.Error, []string{"source.go:1:1: expected ';', found 'EOF'", "source.go:1:1: expected 'IDENT', found 'EOF'", "source.go:1:1: expected 'package', found 'EOF'"}),
	1: newMyTest("test1.source", loadFile("testfiles/test1.source"), msg.Ok, nil),
	2: newMyTest("test2.source", loadFile("testfiles/test2.source"), msg.Ok, nil),
	3: newMyTest("test3.source", loadFile("testfiles/test3.source"), msg.Ok, nil),
	4: newMyTest("test4.source", loadFile("testfiles/test4.source"), msg.Ok, nil),
	5: newMyTest("test5.source", loadFile("testfiles/test5.source"), msg.Ok, nil),
	6: newMyTest("test6.source", loadFile("testfiles/test6.source"), msg.Ok, nil),
}

func TestGetResponse(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getResponse(test.req)
			require.Equal(t, test.res, got, fmt.Sprintf("getResponse() = %v, want %v", got, test.res))
		})
	}
}

func TestStart(t *testing.T) {
	input := &bytes.Buffer{}
	output := &bytes.Buffer{}
	want := &bytes.Buffer{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				input.Reset()
				output.Reset()
				want.Reset()
			}()

			// encode request
			enc := json.NewEncoder(input)
			err := enc.Encode(test.req)
			require.NoError(t, err)

			// execute start()
			err = start(input, output)
			require.NoError(t, err, fmt.Sprintf("start(): error = %v, want nil", err))

			// encode desired response
			encWant := json.NewEncoder(want)
			err = encWant.Encode(test.res)
			require.NoError(t, err)

			// Comapare output(encoded generated response) against want(encoded desired response)
			require.Equal(t, want.String(), output.String(), "start(): output != want")

		})
	}
}

func TestCmd(t *testing.T) {
	input := &bytes.Buffer{}
	output := &bytes.Buffer{}
	test := tests[4]
	test.res.Driver = driverTestVersion
	t.Run(test.name, func(t *testing.T) {
		// encode request
		enc := json.NewEncoder(input)
		err := enc.Encode(test.req)
		require.NoError(t, err)

		// run command
		dv := fmt.Sprintf("-X main.driverVersion=%v", driverTestVersion)
		cmd := exec.Command("go", "run", "-ldflags", dv, "main.go", "conf_nodes.go")
		cmd.Stdin = input
		cmd.Stdout = output
		err = cmd.Run()
		require.NoError(t, err, fmt.Sprintf("exit command with errors: %v", err))

		// encode desired response
		want := &bytes.Buffer{}
		encWant := json.NewEncoder(want)
		err = encWant.Encode(test.res)
		require.NoError(t, err)

		// Comapare output(encoded generated response) against want(encoded desired response)
		require.Equal(t, want.String(), output.String(), "start(): output != want")
	})
}
