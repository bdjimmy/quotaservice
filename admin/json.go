// Licensed under the Apache License, Version 2.0
// Details: https://raw.githubusercontent.com/maniksurtani/quotaservice/master/LICENSE

package admin

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/maniksurtani/quotaservice/logging"
)

func writeJSONError(w http.ResponseWriter, err *HttpError) {
	response := make(map[string]string)
	response["error"] = http.StatusText(err.status)
	response["description"] = err.message

	w.WriteHeader(err.status)
	writeJSON(w, response)
}

func writeJSON(w http.ResponseWriter, object interface{}) {
	b, e := json.Marshal(object)

	if e != nil {
		writeJSONError(w, &HttpError{e.Error(), http.StatusBadRequest})
		return
	}

	_, e = w.Write(b)

	if e != nil {
		logging.Printf("Error writing JSON! %+v", e)
	}
}

func unmarshalJSON(r io.Reader, object interface{}) error {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, object)
}
