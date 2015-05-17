package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func main() {
	goji.Get("/word/:word", save)
	goji.Post("/word/:word", save)
	goji.Delete("/word/:word", remove)
	goji.Serve()
}

func get(c web.C, w http.ResponseWriter, r *http.Request) {
}
func save(c web.C, w http.ResponseWriter, r *http.Request) {
	// save
	redis_conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer redis_conn.Close()
	key := "nasu:" + c.URLParams["word"]
	val, err := redis.Int(redis_conn.Do("GET", key))
	if err != nil {
		val = 0
	}
	redis_conn.Do("SET", key, val+1)

	// response
	res, err := json.Marshal(map[string]string{
		"word": c.URLParams["word"],
		"cnt":  strconv.Itoa(val + 1),
	})
	if err != nil {
		http.Error(w, "System Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.Write(res)
}
func remove(c web.C, w http.ResponseWriter, r *http.Request) {
}
