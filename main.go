package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jgolang/api"
	"github.com/jgolang/config"
	"github.com/jgolang/log"
	"github.com/jgolang/redis"
	"github.com/jhuygens/cache"
	_ "github.com/jhuygens/itunes-searcher"
	"github.com/jhuygens/security"
)

var (
	defaultCacheExpire = config.GetInt("general.default_cache_expire")
	searchersRegister  = make(map[int]string)
)

func init() {
	// Load searchers-register names (library names)
	searchersRegister[0] = config.GetString("searchers.itunes")
	searchersRegister[1] = config.GetString("searchers.tvmaze")
	searchersRegister[2] = config.GetString("searchers.crcind")

	// Api package custom config
	api.CustomTokenValidatorFunc = security.ValidateAccessTokenFunc
	api.Print = log.Printf
	api.PrintError = log.Error
	api.Fatal = log.Fatal

	// Register redis in cache package
	_, err := redis.DefaultClient(config.GetString("cache.host"))
	if err != nil {
		log.Fatal(err)
	}
	cache.Register(redis.RConnect{})
}

func main() {
	router := mux.NewRouter()
	port := config.GetInt("services.search.port")
	// tokenAuthMiddlewares := api.MiddlewaresChain(api.ContentExtractor, api.CustomToken)
	tokenAuthMiddlewares := api.MiddlewaresChain(OmitRequestGET, api.ContentExtractor)

	router.HandleFunc("/v1/search", tokenAuthMiddlewares(searchHandler)).Methods(http.MethodGet)

	log.Info("Starting server, listen on port: ", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), router); err != nil {
		log.Panic(err)
	}
}

// OmitRequestGET doc ...
func OmitRequestGET(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody api.JSONRequest
		if r.Method == http.MethodGet {
			rawRequestBody, _ := json.Marshal(requestBody)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(rawRequestBody))
		}
		next.ServeHTTP(w, r)
	}
}
