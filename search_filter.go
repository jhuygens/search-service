package main

import (
	"regexp"
	"strings"

	searcher "github.com/jhuygens/searcher-engine"
)

var (
	valuesLayout = "(\\w+|\"(\\w+| )+\")+"
)

func parseQueryToSearchFilter(q string) searcher.Filter {
	var filter searcher.Filter
	var name = false
	fields := strings.Split(q, ",")
	for _, field := range fields {
		value := strings.Split(field, ":")
		if len(value) == 2 {
			if name {
				continue
			}
			switch value[0] {
			case "name":
				filter.Name = getFieldValues(value[1])
			case "artist":
				filter.Artist = getFieldValues(value[1])
			case "album":
				filter.Album = getFieldValues(value[1])
			case "genre":
				filter.Genre = getFieldValues(value[1])
			case "year":
				filter.Year = getFieldValues(value[1])
			case "country":
				filter.Country = getFieldValues(value[1])
			default:
				continue
			}
		}
		if len(value) == 1 {
			name = true
			filter.Name = getFieldValues(value[0])
		}
	}
	return filter
}

func getFieldValues(s string) []searcher.FieldValue {
	var fieldValues []searcher.FieldValue
	var excludeNext = false
	for _, value := range getFieldValuesStr(s) {
		if value == "NOT" {
			excludeNext = true
			continue
		}
		fieldValues = append(
			fieldValues,
			searcher.FieldValue{
				Value:   value,
				Exclude: excludeNext,
			},
		)
		excludeNext = false
	}
	return fieldValues
}

func getFieldValuesStr(s string) []string {
	re := regexp.MustCompile(valuesLayout)
	values := re.FindAllString(s, -1)
	return values
}
