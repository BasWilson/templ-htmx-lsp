package main

import (
	"regexp"
)

// this file i just test code snippets in isolation. part of the process of learning go and lsp.

func main() {

	testString := `Content-Length: 1234\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"initialize\",\"params\":{\"processId\":null,\"rootPath\":\"/home/baswilson/Projects/templ-lsp\",\"rootUri\":\"file:///home/baswilson/Projects/templ-lsp\",\"capabilities\":{\"workspace\":{\"applyEdit\":true,\"workspaceEdit\":{\"documentChanges\":true,\"resourceOperations\":[\"create\",\"rename\",\"delete\"]},\"didChangeConfiguration\":{\"dynamicRegistration\":true},\"didChangeWatchedFiles\":{\"dynamicRegistration\":true},\"symbol\":{\"dynamicRegistration\":true,\"symbolKind\":{\"valueSet\":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]}}},\"textDocument\":{\"publishDiagnostics\":{\"relatedInformation\":true},\"synchronization\":{\"dynamicRegistration\":true,\"willSave\":true,\"willSaveWaitUntil\":true,\"didSave\":true},\"completion\":{\"dynamicRegistration\":true,\"completionItem\":{\"snippetSupport\":true,\"commitCharactersSupport\":true,\"documentationFormat\":[\"markdown\",\"plaintext\"],\"deprecatedSupport\":true,\"preselectSupport\":true},\"completionItemKind\":{\"valueSet\":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]}}}}}}`
	buffer := []byte(testString)

	regex := regexp.MustCompile(`Content-Length: (\d+)`)

	matches := regex.FindSubmatchIndex(buffer)
	contentLength := buffer[matches[2]:matches[3]] // [startOfSubmatch endOfSubmatch] [16 20]
	startOfBody := matches[1] + 4                  // +4 to skip \r\n\r\n (they are seen as 1 byte each.)

	println(string(contentLength))
	println(startOfBody)
}
