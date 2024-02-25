package data

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Value struct {
	Val any
}

type DataServer struct {
	KVMap    map[string]Value
	KeyIndex any // Should be BST/RBT once implemented, used to identify keys that need to be transferred.
}

func NewDataServer() *DataServer {
	return &DataServer{KVMap: make(map[string]Value), KeyIndex: nil}
}

func (dataServer *DataServer) GetValue(w http.ResponseWriter, r *http.Request) {
	// Collect the key parameter.
	key := mux.Vars(r)["key"]

	// Handle no key
	if key == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Retrieve value from map
	value, found := dataServer.KVMap[key]

	// Handle not found
	if !found {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Else, put in JSON and return

	valBytes, err := json.Marshal(value)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(valBytes)
}

func (dataServer *DataServer) PutValue(w http.ResponseWriter, r *http.Request) {
	// Collect the key parameter.
	key := mux.Vars(r)["key"]

	if key == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": No key provided.", http.StatusBadRequest)
	}

	// Retrieve Value from Body
	requestBody, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	value := Value{}
	err = json.Unmarshal(requestBody, &value)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": invalid or no value provided.", http.StatusBadRequest)
		return
	}

	// Edit the key-value pair
	_, found := dataServer.KVMap[key]
	dataServer.KVMap[key] = value

	var status int
	if found {
		status = http.StatusAccepted
	} else {
		status = http.StatusCreated
	}

	w.WriteHeader(status)
	w.Write(([]byte)(http.StatusText(status)))
}

func (dataServer *DataServer) DeleteKV(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	if key == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": No key provided.", http.StatusBadRequest)
		return
	}

	_, found := dataServer.KVMap[key]

	if !found {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	delete(dataServer.KVMap, key)

	w.WriteHeader(http.StatusOK)
	w.Write(([]byte)(http.StatusText(http.StatusOK)))
}

func (dataServer *DataServer) Serve(port int) {

}
