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
	Name        string   `redis:"name" json:"name"`
	Website     string   `redis:"website" json:"website"`
	Description string   `redis:"description" json:"description"`
	Languages   []string `redis:"languages" json:"languages"`
}

type Tool struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func main() {
	flag.Parse()

	loadScripts()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api", companyListHandler)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func loadScripts() {
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
		scripts[scriptName] = redis.NewScript(-1, string(scriptData))
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.Redirect(w, r, "/", 301)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func companyListHandler(w http.ResponseWriter, r *http.Request) {
	conn := pool.Get()
	defer conn.Close()

	data, err := redis.Values(scripts["get_companies.lua"].Do(conn, 0))
	if err != nil {
		log.Print("failed to get companies: ", err)
		w.WriteHeader(500)
		return
	}

	companies := make([]*Company, 0, len(data))

	for _, companyData := range data {
		company := &Company{}

		redis.ScanStruct(companyData.([]interface{})[0].([]interface{}), company)
		company.Languages, _ = redis.Strings(companyData.([]interface{})[1], nil)

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
