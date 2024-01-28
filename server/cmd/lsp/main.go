package main

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"strconv"

	"github.com/baswilson/templ-lsp/cmd/lsp/methods"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	buffer := []byte{}

	// Split on each byte, otherwise scanner will split on newline and we'll
	// get certain messages too late.
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) == 0 {
			return 0, nil, nil
		}
		return 1, data[:1], nil
	})

	for scanner.Scan() {
		buffer = append(buffer, scanner.Text()...)

		regex := regexp.MustCompile(`Content-Length: (\d+)`)
		if !regex.Match(buffer) {
			continue
		}

		// this function returns a slice of 4 ints:
		// [startOfMatch endOfMatch startOfSubmatch endOfSubmatch]
		// [0 20 16 20] for example (where 20 is the length of the match)
		// 0-20: Content-Length: 1234 (the match) (20 chars). 16-20: 1234 (the submatch) (4 chars). Total length of match is 20.
		matches := regex.FindSubmatchIndex(buffer)
		contentLength := buffer[matches[2]:matches[3]] // [startOfSubmatch endOfSubmatch] [16 20]
		startOfBody := matches[1] + 4                  // +4 to skip \r\n\r\n (they are seen as 1 byte each.)
		contentLengthInt, err := strconv.Atoi(string(contentLength))
		if err != nil {
			methods.WriteLog(err.Error())
			continue
		}

		if len(buffer) < startOfBody+contentLengthInt {
			continue
		}

		rawMessage := buffer[startOfBody : startOfBody+contentLengthInt]
		methods.WriteLog(string(rawMessage))

		var message methods.RequestMessage
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			methods.WriteLog(err.Error())
			continue
		}

		var result interface{}
		var methodError error

		switch message.Method {
		case "initialize":
			result, methodError = methods.Initialize(message)
		case "textDocument/completion":
			result, methodError = methods.TextDocumentCompletion(rawMessage)
		case "textDocument/didChange":
			methodError = methods.TextDocumentDidChange(rawMessage)
		}

		if methodError != nil {
			methods.WriteLog("ERROR: " + methodError.Error())
			continue
		}

		if result == nil {
			buffer = buffer[startOfBody+contentLengthInt:]
			continue
		}

		err = methods.SendResponse(writer, message.Id, result)

		if err != nil {
			methods.WriteLog("ERROR: " + err.Error())
			continue
		}

		buffer = buffer[startOfBody+contentLengthInt:]
	}

	if err := scanner.Err(); err != nil {
		methods.WriteLog(err.Error())
	}
}
