package main

import (
	"encoding/json"
	"flag"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var pool = &redis.Pool{MaxIdle: 3, Dial: func() (redis.Conn, error) { return redis.Dial("tcp", ":6379") }}
var scripts = make(map[string]*redis.Script)

type Company struct {
	Name        string `redis:"name" json:"name"`
	Website     string `redis:"website" json:"website"`
	Description string `redis:"description" json:"description"`
}

func main() {
	flag.Parse()
	scriptNames, err := ioutil.ReadDir("lua")

	if err != nil {
		log.Fatal("could not load lua scripts: ", err)
	}

	for _, s := range scriptNames {
		scriptName := s.Name()
		scriptData, err := ioutil.ReadFile("lua/" + scriptName)

		if err != nil {
			log.Fatal("could not load lua script: ", err)
		}

		log.Println("loaded script: ", scriptName)
		scripts[scriptName] = redis.NewScript(0, string(scriptData))
	}

	http.HandleFunc("/", html)
	http.HandleFunc("/api", api)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func html(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func api(w http.ResponseWriter, r *http.Request) {
	conn := pool.Get()

	data, err := redis.Values(scripts["get_companies.lua"].Do(conn))
	if err != nil {
		log.Print("failed to get companies: ", err)
		w.WriteHeader(500)
		return
	}

	companies := make([]*Company, 0)

	for len(data) > 0 {
		company := &Company{}
		data, _ = redis.Scan(data, &company.Name, &company.Website, &company.Description)
		companies = append(companies, company)
	}

	companiesJSON, err := json.Marshal(companies)

	if err != nil {
		log.Print("failed to serialize companies: ", err)
		w.WriteHeader(500)
		return
	}

	w.Write(companiesJSON)
}
