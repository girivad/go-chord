package data

import (
	"encoding/json"
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
	key := mux.Vars(r)["key"]

	if key == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": Key not provided.", http.StatusBadRequest)
		return
	}

	val, ok := dataServer.KVMap[key]

	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	// Optimize, given that the marshalling is only for writing.
	valueBytes, err := json.Marshal(val)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": Error encoding value for transmission.", http.StatusInternalServerError)
		return
	}

	// Automatically sends Status OK Header.
	w.Write(valueBytes)
}

func (dataServer *DataServer) PutValue(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	if key == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check if key is in the map => 201 if no, 200 if yes.
	// val, ok := dataServer.KVMap[key]
}
