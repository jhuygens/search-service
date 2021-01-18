package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jgolang/api"
	"github.com/jgolang/config"
	"github.com/jgolang/log"
	"github.com/jhuygens/cache"
	searcher "github.com/jhuygens/searcher-engine"
)

func search(filter searcher.Filter, offset, limit int, url *url.URL) api.Response {
	searchKey, err := searcher.GenerateSearchKey(filter)
	if err != nil {
		log.Error(err)
		return api.Error{
			Message: "Por favor intenta mas tarde",
		}
	}
	items, total, err := getCacheItems(searchKey, offset, limit)
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
		items, total, err = getCacheItems(searchKey, offset, limit)
		if err != nil {
			log.Error(err)
			return api.Error{
				Message:    "Por favor intenta mas tarde",
				ErrorCode:  "6",
				StatusCode: http.StatusInternalServerError,
			}
		}
		if items == nil {
			log.Warn("Not get items")
		}
	}
	return api.Success{
		Message: fmt.Sprintf("Search key: %v", searchKey),
		Data: Paging{
			Href:     getCurrentURL(url.RequestURI()),
			Items:    items,
			Limit:    limit,
			Next:     getNextURL(url.Query(), offset, limit, total),
			Offset:   offset,
			Previous: getPreviousURL(url.Query(), offset, limit, total),
			Total:    total,
		},
	}
}

func getCacheItems(searchKey string, offset, limit int) ([]searcher.Item, int, error) {
	result, err := cache.Get(searchKey)
	if err != nil {
		return nil, 0, err
	}
	if result == "" {
		return nil, 0, nil
	}
	var items []searcher.Item
	err = json.Unmarshal([]byte(result), &items)
	if err != nil {
		return nil, 0, err
	}
	end := offset + limit
	total := len(items)
	if end >= total {
		end = total - 1
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		offset = total - 1
	}
	itemsResponse := items[offset:end]
	return itemsResponse, total, nil
}

func getCurrentURL(uri string) string {
	return fmt.Sprintf("%v%v", config.GetString("general.host"), uri)
}

func getNextURL(queryValues url.Values, offset, limit, total int) string {
	nextOffset := getNextOffset(offset, limit, total)
	if nextOffset == total {
		return ""
	}
	return generateSearchURL(queryValues, nextOffset, limit)
}

func getNextOffset(currentOffset, limit, total int) int {
	offset := currentOffset + limit
	if offset > total {
		offset = total
	}
	return offset
}

func generateSearchURL(queryValues url.Values, offset, limit int) string {
	q := queryValues.Get("q")
	typeResource := queryValues.Get("type")
	library := queryValues.Get("library")
	queryString := url.PathEscape(fmt.Sprintf("q=%s&type=%v&library=%v&offset=%v&limit=%v", q, typeResource, library, offset, limit))
	return fmt.Sprintf("%v?%v", config.GetString("services.search.paging.search.url"), queryString)
}

func getPreviousURL(queryValues url.Values, offset, limit, total int) string {
	if offset <= 0 {
		return ""
	}
	if total == 0 {
		return ""
	}
	return generateSearchURL(queryValues, getPreviousOffset(offset, limit), limit)
}

func getPreviousOffset(currentOffset, limit int) int {
	offset := currentOffset - limit
	if offset < 0 {
		offset = 0
	}
	return offset
}
