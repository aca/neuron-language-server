package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aca/neuron-language-server/neuron"
	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
)

type server struct {
	conn       *jsonrpc2.Conn
	rootURI    string // from initliaze param
	rootDir    string
	logger     *log.Logger
	neuronMeta map[string]neuron.Result
	documents  map[lsp.DocumentURI]string
}

func (s *server) handleTextDocumentDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.DidOpenTextDocumentParams
	err = json.Unmarshal(*req.Params, &params)
	if err != nil {
		return nil, err
	}

	s.documents[params.TextDocument.URI] = params.TextDocument.Text
	return nil, nil
}

func (s *server) handleTextDocumentDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.DidChangeTextDocumentParams
	err = json.Unmarshal(*req.Params, &params)
	if err != nil {
		return nil, err
	}

	if len(params.ContentChanges) != 1 {
		return nil, fmt.Errorf("len(params.ContentChanges) = %v", len(params.ContentChanges))
	}

	s.documents[params.TextDocument.URI] = params.ContentChanges[0].Text
	return nil, nil
}

func (s *server) handleTextDocumentCompletion(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.CompletionParams
	err = json.Unmarshal(*req.Params, &params)
	if err != nil {
		return nil, err
	}

	items := make([]lsp.CompletionItem, 0)

	w := s.WordAt(params.TextDocument.URI, params.Position)

	w = strings.ReplaceAll(w, "<", "")
	w = strings.ReplaceAll(w, ">", "")

	for id, m := range s.neuronMeta {
		if w == "" {
			item := lsp.CompletionItem{
        Label:      fmt.Sprintf("%v:%v", id, m.ZettelTitle),
				InsertText: id,
				Detail:     m.ZettelDay,
			}
			items = append(items, item)
			continue
		}

		if strings.Contains(m.ZettelID, w) {
			item := lsp.CompletionItem{
        Label:      fmt.Sprintf("%v:%v", id, m.ZettelTitle),
				InsertText: id,
				Detail:     m.ZettelDay,
			}
			items = append(items, item)
			continue
		}

		if strings.Contains(m.ZettelTitle, w) {
			item := lsp.CompletionItem{
        Label:      fmt.Sprintf("%v:%v", id, m.ZettelTitle),
				InsertText: id,
				Detail:     m.ZettelDay,
			}
			items = append(items, item)
			continue
		}
	}

	return items, nil
}

func (s *server) handleTextDocumentDefinition(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.TextDocumentPositionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	w := s.WordAt(params.TextDocument.URI, params.Position)
	w = strings.ReplaceAll(w, "<", "")
	w = strings.ReplaceAll(w, ">", "")

	neuronResult, ok := s.neuronMeta[w]
	if !ok {
		s.logger.Printf("%v doesn't exist", w)
		return nil, nil
	}

	p := filepath.Join(s.rootDir, neuronResult.ZettelPath)

	return lsp.Location{
		URI: "file://" + lsp.DocumentURI(p),
	}, nil
}

func (s *server) handleTextDocumentHover(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.TextDocumentPositionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	w := s.WordAt(params.TextDocument.URI, params.Position)
	w = strings.ReplaceAll(w, "<", "")
	w = strings.ReplaceAll(w, ">", "")

	neuronResult, ok := s.neuronMeta[w]
	if !ok {
		s.logger.Printf("%v doesn't exist", w)
		return nil, nil
	}

	msgl := []string{
		fmt.Sprintf("[%s](%v)\n", neuronResult.ZettelTitle, neuronResult.ZettelPath),
	}

	if len(neuronResult.ZettelTags) != 0 {
		msgl = append(msgl, fmt.Sprintf("tags: %v", strings.Join(neuronResult.ZettelTags, ",")))
	}

	if neuronResult.ZettelDay != "" {
		msgl = append(msgl, fmt.Sprintf("date: %v", neuronResult.ZettelDay))
	}

	msg := strings.Join(msgl, "\n")

	return lsp.Hover{
		Contents: []lsp.MarkedString{
			{
				Language: `markdown`,
				Value:    msg,
			},
		},
	}, nil
}

func (s *server) handleInitialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	s.logger.Print("handleInitialize")
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	s.conn = conn

	var params lsp.InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	u, err := url.ParseRequestURI(string(params.RootURI))
	if err != nil {
		return nil, err
	}

	s.rootDir = u.EscapedPath()
	s.logger.Printf("neuron: query -d %v", s.rootDir)
	queryResult, err := neuron.Query("-d", s.rootDir)
	if err != nil {
		s.logger.Println(queryResult)
		return nil, err
	}

	s.logger.Printf("neuron: %v found", len(queryResult.Result))

	for _, result := range queryResult.Result {
		s.logger.Printf("neuron: added %s", result.ZettelID)
		s.neuronMeta[result.ZettelID] = result
	}

	initializeResult := lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Options: &lsp.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    lsp.TDSKFull,
				},
			},
			DefinitionProvider: true,
			HoverProvider:      true,
			CompletionProvider: &lsp.CompletionOptions{
				ResolveProvider:   true,
				TriggerCharacters: []string{"<"},
			},
		},
	}

	return initializeResult, nil
}

func (s *server) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(ctx, conn, req)
	case "initialized":
		return
	case "shutdown":
		os.Exit(0)
		return
	case "textDocument/didOpen":
		return s.handleTextDocumentDidOpen(ctx, conn, req)
	case "textDocument/didChange":
		return s.handleTextDocumentDidChange(ctx, conn, req)
	case "textDocument/completion":
		return s.handleTextDocumentCompletion(ctx, conn, req)
	case "textDocument/definition":
		return s.handleTextDocumentDefinition(ctx, conn, req)
	case "textDocument/hover":
		return s.handleTextDocumentHover(ctx, conn, req)
		// case "textDocument/didSave":
		// 	return s.handleTextDocumentDidSave(ctx, conn, req)
		// case "textDocument/didClose":
		// 	return s.handleTextDocumentDidClose(ctx, conn, req)
		// case "textDocument/formatting":
		// 	return s.handleTextDocumentFormatting(ctx, conn, req)
		// case "textDocument/documentSymbol":
		// 	return s.handleTextDocumentSymbol(ctx, conn, req)
		// case "textDocument/codeAction":
		// 	return s.handleTextDocumentCodeAction(ctx, conn, req)
		// case "workspace/executeCommand":
		// 	return s.handleWorkspaceExecuteCommand(ctx, conn, req)
		// case "workspace/didCs.ngeConfiguration":
		// 	return s.handleWorkspaceDidChangeConfiguration(ctx, conn, req)
		// case "workspace/workspaceFolders":
		// 	return s.handleWorkspaceWorkspaceFolders(ctx, conn, req)
	}

	return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
}
