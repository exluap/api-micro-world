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

type Message struct {
	Result  bool
	Message interface{}
} //@name Default Response

func (m *Message) Respond(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}
