package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	redis "github.com/redis/go-redis/v9"
	"github.com/urfave/negroni"
)

var (
	masterPool  *redis.Client
	replicaPool *redis.Client
)

func ListRangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	key := mux.Vars(r)["key"]
	values := handleError(
		replicaPool.LRange(ctx, key, 0, -1).Result(),
	).([]string)

	valueMarshalled := handleError(
		json.Marshal(values),
	).([]byte)

	w.Header().Add("content-type", "application/json")
	w.Write(valueMarshalled)
}

func ListPushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	muxVars := mux.Vars(r)
	key := muxVars["key"]
	value := muxVars["value"]

	handleError(
		masterPool.LPush(ctx, key, value).Result(),
	)

	ListRangeHandler(w, r)
}

func EnvHandler(w http.ResponseWriter, r *http.Request) {
	environment := make(map[string]string)

	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		value := strings.Join(splits[1:], "=")
		environment[key] = value
	}

	json := handleError(
		json.Marshal(environment),
	).([]byte)

	w.Header().Add("content-type", "application/json")
	w.Write(json)
}

func handleError(result interface{}, e error) interface{} {
	if e != nil {
		panic(e)
	}
	return result
}

func main() {
	masterPool = redis.NewClient(&redis.Options{
		Addr:     "redis-master.default:6379",
		Password: "",
		DB:       0,
	})
	defer masterPool.Close()
	replicaPool = redis.NewClient(&redis.Options{
		Addr:     "redis-replica.default:6379",
		Password: "",
		DB:       0,
	})
	defer replicaPool.Close()

	router := mux.NewRouter()
	router.Path("/lrange/{key}").Methods("GET").HandlerFunc(ListRangeHandler)
	router.
		Path("/rpush/{key}/{value}").
		Methods("GET").
		HandlerFunc(ListPushHandler)
	router.Path("/env").Methods("GET").HandlerFunc(EnvHandler)

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
