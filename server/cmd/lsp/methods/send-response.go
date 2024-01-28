package methods

import (
	"bufio"
	"encoding/json"
	"strconv"
)

type ResponseMessage struct {
	Id     int         `json:"id"`
	Result interface{} `json:"result"`
}

func SendResponse(writer *bufio.Writer, id int, result interface{}) error {
	response, err := json.Marshal(ResponseMessage{
		Id:     id,
		Result: result,
	})
	if err != nil {
		WriteLog(err.Error())
		return err
	}

	responseLength := len(response)
	responseLengthString := strconv.Itoa(responseLength)
	header := "Content-Length: " + responseLengthString + "\r\n\r\n"
	responseBytes := append([]byte(header), response...)

	_, err = writer.Write(responseBytes)
	if err != nil {
		WriteLog("Error writing to os.Stdout: " + err.Error())
		return err
	}

	// Flush standard output
	if err := writer.Flush(); err != nil {
		WriteLog("Error flushing os.Stdout: " + err.Error())
		return err
	}

	WriteLog("Sent response for id " + strconv.Itoa(id))

	return nil
}
