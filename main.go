package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/aca/neuron-language-server/neuron"
	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
)

type DebugLogger struct {
	conn net.Conn
}

func main() {
	var connOpt []jsonrpc2.ConnOpt

	logConn, err := net.Dial("tcp", "localhost:3000")
	var logger *log.Logger
	if err == nil {
		logger = log.New(logConn, "", log.LstdFlags|log.Lshortfile)
	} else {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	connOpt = append(connOpt, jsonrpc2.LogMessages(logger))

	s := &server{
		logger:     logger,
		documents:  make(map[lsp.DocumentURI]string),
		neuronMeta: make(map[string]neuron.Result),
	}

	s.logger.Print("start")

	handler := jsonrpc2.HandlerWithError(s.handle)

	<-jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}),
		handler, connOpt...).DisconnectNotify()

	s.logger.Print("shutdown")
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (c stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (c stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
