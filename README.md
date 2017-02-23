# babelfish-go-driver

[babelfish doc](https://github.com/src-d/babelfish)

To get a request and reply a response, go-driver uses standard input and output.

To generate the binary to add to the container and build the docker image:

* $ make clean && make

To run inside a container:

* $ docker run --rm -i babelfish-go-driver

Directory driverclient/ contains a program to generate a single request and feed babelfish-go-driver for testing. You can set language 
and language version by flags. See go run driverclient/main.go --help

Directory start/ contains benchmarks to compare serialization with JSON against Msgpack

Demo(in the project directory):

* $ go run driverclient/main.go -f testfiles/test1.source | docker run --rm -i babelfish-go-driver