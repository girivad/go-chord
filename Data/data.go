package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	pb "github.com/girivad/go-chord/Proto"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/types/known/anypb"
)

type Value struct {
	Val any
}

type DataServer struct {
	KVMap          map[string]Value
	Lock           sync.RWMutex
	RegisterKey    func(string)
	RegisterDelete func(string)
}

func NewDataServer(registerKey func(string), registerDelete func(string)) *DataServer {
	return &DataServer{KVMap: make(map[string]Value), RegisterKey: registerKey, RegisterDelete: registerDelete}
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
	dataServer.Lock.RLock()
	value, found := dataServer.KVMap[key]
	dataServer.Lock.RUnlock()

	// Handle not found
	if !found {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Else, put in JSON and return

	log.Printf("[INFO] GET: K: %s, V: %v", key, value)

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
	defer r.Body.Close()
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

	log.Printf("[INFO] PUT key: %s, value: %v", key, value)

	// Edit the key-value pair
	dataServer.Lock.Lock()
	_, found := dataServer.KVMap[key]
	dataServer.KVMap[key] = value
	dataServer.Lock.Unlock()

	var status int
	if found {
		status = http.StatusAccepted
	} else {
		dataServer.RegisterKey(key)
		status = http.StatusCreated
	}

	w.WriteHeader(status)
	w.Write(([]byte)(http.StatusText(status)))
}

func (dataServer *DataServer) DeleteKV(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	log.Printf("[INFO] DEL key: %s, value:%v", key, dataServer.KVMap[key])

	if key == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": No key provided.", http.StatusBadRequest)
		return
	}

	_, found := dataServer.KVMap[key]

	if !found {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	dataServer.Lock.Lock()
	delete(dataServer.KVMap, key)
	dataServer.Lock.Unlock()

	dataServer.RegisterDelete(key)

	w.WriteHeader(http.StatusOK)
	w.Write(([]byte)(http.StatusText(http.StatusOK)))
}

func (dataServer *DataServer) Serve(port int) {
	router := mux.NewRouter()
	router.HandleFunc("/data/{key}", dataServer.GetValue).Methods("GET")
	router.HandleFunc("/data/{key}", dataServer.PutValue).Methods("PUT")
	router.HandleFunc("/data/{key}", dataServer.DeleteKV).Methods("DELETE")

	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func (dataServer *DataServer) GetValuesForTransfer(keys []string) (*pb.KVMap, error) {
	kvMap := make(map[string]*pb.Value)
	dataServer.Lock.RLock()
	for _, key := range keys {
		value, ok := dataServer.KVMap[key]

		if !ok {
			continue
		}

		byteValue, err := json.Marshal(value)

		if err != nil {
			fmt.Printf("byteValue marshalling failed: %v", err)
			return nil, err
		}

		kvMap[key] = &pb.Value{Val: &anypb.Any{Value: byteValue}}
	}
	dataServer.Lock.RUnlock()

	return &pb.KVMap{Kvmap: kvMap}, nil
}

func (dataServer *DataServer) PutValuesForTransfer(data *pb.KVMap) error {
	dataServer.Lock.Lock()

	for key, value := range data.Kvmap {
		parsedValue := &Value{}
		err := json.Unmarshal(value.Val.Value, parsedValue)

		if err != nil {
			fmt.Printf("Unmarshalling value failed: %v", err)
			return err
		}

		dataServer.KVMap[key] = *parsedValue
	}

	dataServer.Lock.Unlock()

	return nil
}
