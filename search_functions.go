package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jgolang/api"
	"github.com/jgolang/log"
	"github.com/jhuygens/cache"
	searcher "github.com/jhuygens/searcher-engine"
	"go.mnc.gt/config"
)

func search(filter searcher.Filter, offset, limit int, url *url.URL) api.Response {
	searchKey, err := searcher.GenerateSearchKey(filter)
	if err != nil {
		log.Error(err)
		return api.Error{
			Message: "Por favor intenta mas tarde",
		}
	}
	items, err := getCacheItems(searchKey, offset, limit)
	if err != nil {
		log.Error(err)
		return api.Error{
			Message:    "Por favor intenta mas tarde",
			ErrorCode:  "6",
			StatusCode: http.StatusInternalServerError,
		}
	}
	if items == nil {
		searchKey, err = searcher.Search(filter)
		if err != nil {
			log.Error(err)
			return api.Error{
				Message:    "Por favor intenta mas tarde",
				ErrorCode:  "7",
				StatusCode: http.StatusInternalServerError,
			}
		}
		items, err = getCacheItems(searchKey, offset, limit)
		if err != nil {
			log.Error(err)
			return api.Error{
				Message:    "Por favor intenta mas tarde",
				ErrorCode:  "6",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}
	total := len(items)
	return api.Success{
		Message: fmt.Sprintf("Search key: %v", searchKey),
		Data: Paging{
			Href:     getCurrentURL(url.RequestURI()),
			Items:    items,
			Limit:    limit,
			Next:     getNextURL(url.Query(), offset, limit, total),
			Offset:   offset,
			Previous: getPreviousURL(url.Query(), offset, limit),
			Total:    total,
		},
	}
}

func getCurrentURL(uri string) string {
	return fmt.Sprintf("%v%v", config.GetString("general.host"), uri)
}

func getPreviousURL(queryValues url.Values, offset, limit int) string {
	if offset <= 0 {
		return ""
	}
	return generateSearchURL(queryValues, getPreviousOffset(offset, limit), limit)
}

func getNextURL(queryValues url.Values, offset, limit, total int) string {
	nextOffset := getNextOffset(offset, limit, total)
	if nextOffset == total {
		return ""
	}
	return generateSearchURL(queryValues, nextOffset, limit)
}

func generateSearchURL(queryValues url.Values, offset, limit int) string {
	q := queryValues.Get("q")
	typeResource := queryValues.Get("type")
	library := queryValues.Get("library")
	queryString := url.PathEscape(fmt.Sprintf("q=%s&type=%v&library=%v&offset=%v&limit=%v", q, typeResource, library, offset, limit))
	return fmt.Sprintf("%v/v1/search?%v", config.GetString("general.host"), queryString)
}

func getNextOffset(currentOffset, limit, total int) int {
	offset := currentOffset + limit
	if offset > total {
		offset = total
	}
	return offset
}
func getPreviousOffset(currentOffset, limit int) int {
	offset := currentOffset - limit
	if offset < 0 {
		offset = 0
	}
	return offset
}

func getCacheItems(searchKey string, offset, limit int) ([]searcher.Item, error) {
	result, err := cache.Get(searchKey)
	if err != nil {
		return nil, err
	}
	var items []searcher.Item
	err = json.Unmarshal([]byte(result), &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
