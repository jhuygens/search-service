package main

import (
	"net/http"
	"strings"

	"github.com/jgolang/api"
	"github.com/jgolang/log"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q, response := api.GetQueryParamValueString("q", r)
	if response != nil {
		response.Write(w)
		return
	}
	typeResource, response := api.GetQueryParamValueString("type", r)
	if response != nil {
		response.Write(w)
		return
	}
	var limit = 20
	var offset = 0
	if r.URL.Query().Get("limit") != "" {
		limit, response = api.GetQueryParamValueInt("limit", r)
		if response != nil {
			response.Write(w)
			return
		}
	}
	if r.URL.Query().Get("offset") != "" {
		offset, response = api.GetQueryParamValueInt("offset", r)
		if response != nil {
			response.Write(w)
			return
		}
	}
	library := r.URL.Query().Get("library")
	if library == "" {
		library = "all"
	}
	types := strings.Split(typeResource, ",")
	filter := parseQueryToSearchFilter(q)
	filter.Types = types
	filter.Library = library
	log.Info(q)
	log.Info(filter)
	log.Info(offset)
	log.Info(limit)
	response = search(filter, offset, limit, r.URL)
	response.Write(w)
	return
}
