package main

import (
	"encoding/json"
	"net/http"

	"github.com/AmreeshTyagi/goldflake"
)

var gf *goldflake.Goldflake

func getMachineID() (uint16, error) {
	return 1234, nil
}

func init() {
	var st goldflake.Settings
	st.MachineID = getMachineID
	gf = goldflake.NewGoldflake(st)
	if gf == nil {
		panic("goldflake not created")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	id, err := gf.NextID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(goldflake.Decompose(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	w.Write(body)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
