package utils

import (
	"strconv"
	"strings"
)

type QueryParam struct {
	Filters  map[string]string
	Limit    uint64
	Page   uint64
	Ordering []string
	Search   string
}

func ParseQueryParam(queryParams map[string][]string) (*QueryParam, []string) {
	params := QueryParam{
		Filters:  make(map[string]string),
		Limit:    10,
		Page:   1,
		Ordering: []string{},
		Search:   "",
	}
	var errStr []string
	var err error

	for key, value := range queryParams {
		if key == "page" {
			params.Page, err = strconv.ParseUint(value[0], 10, 64)
			if err != nil {
				errStr = append(errStr, "invalid `page` param")
			}
			continue
		}

		if key == "limit" {
			params.Limit, err = strconv.ParseUint(value[0], 10, 64)
			if err != nil {
				errStr = append(errStr, "invalid `limit` param")
			}
			continue
		}

		if key == "search" {
			params.Search = value[0]
			continue
		}
		if key == "ordering" {
			params.Ordering = strings.Split(value[0], ",")
			continue
		}
		params.Filters[key] = value[0]
	}
	return &params, errStr
}
