package methods

import (
	"os"
	"time"
)

func WriteLog(msg string) {
	f, err := os.OpenFile("/tmp/lsp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err = f.WriteString(time.Now().Format(time.RFC3339) + " " + msg + "\n")
	if err != nil {
		panic(err)
	}

	f.Sync()
}
