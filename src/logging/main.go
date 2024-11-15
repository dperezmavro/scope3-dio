package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	LevelError = "error"
	LevelInfo  = "info"
)

// Data represents logging data in structured logging.
type Data map[string]interface{}

// StructuredLog is used for JSON-formatted log output
type StructuredLog struct {
	Service     string `json:"service"`
	Environment string `json:"environment"`
	Message     string `json:"message"`
	Level       string `json:"level"`
	Error       string `json:"error"`
}

func Fatal(ctx context.Context, err error, data Data, m string) {
	sLog := StructuredLog{
		Error:   err.Error(),
		Message: m,
	}
	msg, err := json.Marshal(sLog)
	if err != nil {
		log.Fatalf("unable to generate log entry for %s %+v %+v", m, err, data)
	}
	fmt.Println(string(msg))
	os.Exit(1)
}
