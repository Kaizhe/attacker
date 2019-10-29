package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/attacker/types"
	util "github.com/attacker/utils"
	"github.com/gorilla/mux"
)

func launchAttack(w http.ResponseWriter, r *http.Request) {
	var ac types.AttackConfig
	var err error
	var name string

	params, ok := r.URL.Query()["tool"]

	if !ok || len(params) != 1 {
		name = "metasploit"
	} else {
		name = params[0]
	}

	if r.Method != "POST" {
		msg := fmt.Sprintf("Invalid request method: %s", r.Method)
		util.LogPrint(w, msg, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err = decoder.Decode(&ac)
	fmt.Println(ac)
	if err != nil {
		util.LogPrint(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = types.LaunchNewAttack(name, ac)
	if err != nil {
		util.LogPrint(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf("Attack '%s' success.", ac.Name)
	util.LogPrint(w, msg, http.StatusOK)
}

func getAttackTool(w http.ResponseWriter, r *http.Request) {
	tools := []string{}
	for t := range types.Attackers {
		tools = append(tools, t)
	}

	msg := fmt.Sprintf("Attack tools: %s", tools)
	util.LogPrint(w, msg, http.StatusOK)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/attack", launchAttack).Methods("POST")
	router.HandleFunc("/attack", getAttackTool).Methods("GET")
	fmt.Println("Listening on 8080.")
	http.ListenAndServe(":8080", router)
}
