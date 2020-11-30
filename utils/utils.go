package utils

import (
	"encoding/json"
	"net/http"
)

//Message - message handler
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"Status": status, "Message": message}
}

//Respond - json respond handler
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
