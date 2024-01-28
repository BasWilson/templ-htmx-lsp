package methods

import (
	"encoding/json"

	"github.com/baswilson/templ-lsp/cmd/lsp/globals"
)

type DidChangeParams struct {
	TextDocument   TextDocumentIdentifier           `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	Range       Range  `json:"range"`
	RangeLength int    `json:"rangeLength"`
	Text        string `json:"text"`
}

type TextDocumentIdentifier struct {
	Uri string `json:"uri"`
}

type VersionedTextDocumentIdentifier struct {
	Version int `json:"version"`
	TextDocumentIdentifier
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type DidChangeMessage struct {
	Message
	Params DidChangeParams `json:"params"`
}

func TextDocumentDidChange(rawMessage []byte) error {

	message := DidChangeMessage{}
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		WriteLog(err.Error())
		return err
	}
	params := message.Params
	globals.Documents[params.TextDocument.Uri] = params.ContentChanges[0].Text
	return nil
}
