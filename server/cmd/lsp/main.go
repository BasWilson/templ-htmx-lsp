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

		contentLength := regex.FindSubmatch(buffer)[1]
		startOfBody := regex.FindSubmatchIndex(buffer)[1] + 4
		contentLengthInt, err := strconv.Atoi(string(contentLength))
		if err != nil {
			methods.WriteLog(err.Error())
			continue
		}

		if len(buffer) < startOfBody+contentLengthInt {
			continue
		}

		raw := buffer[startOfBody : startOfBody+contentLengthInt]
		methods.WriteLog(string(raw))

		var message methods.RequestMessage
		if err := json.Unmarshal(raw, &message); err != nil {
			methods.WriteLog(err.Error())
			continue
		}

		var result interface{}
		var methodError error

		switch message.Method {
		case "initialize":
			result, methodError = methods.Initialize(message)
		}

		if methodError != nil {
			methods.WriteLog("ERROR: " + methodError.Error())
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
