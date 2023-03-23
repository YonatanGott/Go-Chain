package utils

import "encoding/json"

func JsonStatus(message string) []byte {
	marshal, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return marshal
}
