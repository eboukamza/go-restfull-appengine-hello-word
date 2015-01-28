package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"strconv"
	"log"

	"github.com/gorilla/mux"

	"appengine"
	"appengine/datastore"
)

type Hello struct {
	Id          int64 `json:"id" datastore:"-"`
	Name        string `json:"name"`
	TimeStamp   time.Time `json:"timeStamp"`
}

type Hellos []Hello

func init() {

	router := mux.NewRouter().StrictSlash(false)

	hellos := router.Path("/api/hellos").Subrouter()
	hellos.Methods("GET").HandlerFunc(GetHellosHandler)
	hellos.Methods("POST").HandlerFunc(CreateHelloHandler)

	hello := router.Path("/api/hellos/{id}").Subrouter()
	hello.Methods("GET").HandlerFunc(GetHelloHandler)
	hello.Methods("DELETE").HandlerFunc(DeleteHelloHandler)
	hello.Methods("PUT").HandlerFunc(UpdateHelloHander)

	http.Handle("/", router)
}

func GetHellosHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("hello")

	var hellos [] Hello

	keys, err := q.GetAll(c, &hellos)
	for i, key := range keys {
		hellos[i].Id = key.IntID()
	}

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(hellos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func GetHelloHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)

	c := appengine.NewContext(r)
	helloKey := datastore.NewKey(c, "hello", "", id, nil)

	var hello Hello
	if err := datastore.Get(c, helloKey, &hello); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(hello); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreateHelloHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var hello Hello
	if err := json.NewDecoder(r.Body).Decode(&hello); err != nil {
		panic(err)
	}

	hello.TimeStamp = time.Now();

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "hello", nil), &hello)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hello.Id = key.IntID();
	log.Println(fmt.Sprintf("Entity with key='%d' created", hello.Id));
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Location", fmt.Sprintf("api/hellos/%d", +key.IntID()))
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(hello); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func UpdateHelloHander(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)

	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "hello", "", id, nil)

	var hello Hello
	if err := json.NewDecoder(r.Body).Decode(&hello); err != nil {
		panic(err)
	}

	key, err := datastore.Put(c, key, &hello)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteHelloHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)

	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "hello", "", id, nil)
	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
