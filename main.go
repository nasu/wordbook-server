package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
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
	if os.Getenv("REDISTOGO_URL") == "" {
		os.Setenv("REDISTOGO_URL", "redis://127.0.0.1:6379/")
	}
	os.Setenv("REDIS_URL", os.Getenv("REDISTOGO_URL"))
	redisConn, err := redisurl.Connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer redisConn.Close()
	word := strings.Trim(strings.ToLower(c.URLParams["word"]), " .,!?")
	key := "nasu:" + word
	val, err := redis.Int(redisConn.Do("GET", key))
	if err != nil {
		val = 0
	}
	redisConn.Do("SET", key, val+1)

	// response
	res, err := json.Marshal(map[string]string{
		"word": word,
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
