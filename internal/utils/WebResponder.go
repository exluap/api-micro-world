/**
 * Project api-microworld created by exluap
 * Date: 02.08.2020 18:51
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package utils

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
