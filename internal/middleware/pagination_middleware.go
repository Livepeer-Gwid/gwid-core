package middleware

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type QueryParams struct {
	Search  string            `json:"search"`
	Page    int               `json:"page"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	Sort    string            `json:"sort"`
	Order   string            `json:"order"`
	Filters map[string]string `json:"filters"`
	Fields  []string          `json:"fields"`
}

type QueryConfig struct {
	DefaultPage   int
	DefaultLimit  int
	MaxLimit      int
	DefaultOrder  string
	AllowedOrders []string
	FilterPrefix  string
}

func DefaultQueryConfig() QueryConfig {
	return QueryConfig{
		DefaultPage:   1,
		DefaultLimit:  10,
		MaxLimit:      100,
		DefaultOrder:  "created_at desc",
		AllowedOrders: []string{"asc", "desc"},
		FilterPrefix:  "filter_",
	}
}

func QueryMiddleware(config ...QueryConfig) gin.HandlerFunc {
	cfg := DefaultQueryConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *gin.Context) {
		params := parseQueryParams(c, cfg)

		c.Set("queryParams", params)

		c.Set("search", params.Search)
		c.Set("page", params.Page)
		c.Set("limit", params.Limit)
		c.Set("offset", params.Offset)
		c.Set("order", params.Order)
		c.Set("filters", params.Filters)
		c.Set("fields", params.Fields)

		c.Next()
	}
}

func parseQueryParams(c *gin.Context, cfg QueryConfig) QueryParams {
	params := QueryParams{
		Filters: make(map[string]string),
	}

	params.Search = strings.TrimSpace(c.Query("search"))

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		} else {
			params.Page = cfg.DefaultPage
		}
	} else {
		params.Page = cfg.DefaultPage
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = min(limit, cfg.MaxLimit)
		} else {
			params.Limit = cfg.DefaultLimit
		}
	} else {
		params.Limit = cfg.DefaultLimit
	}

	params.Offset = (params.Page - 1) * params.Limit

	if order := strings.ToLower(c.Query("order")); order != "" {
		if isValidOrder(order, cfg.AllowedOrders) {
			params.Order = order
		} else {
			params.Order = cfg.DefaultOrder
		}
	} else {
		params.Order = cfg.DefaultOrder
	}

	if fields := c.Query("fields"); fields != "" {
		params.Fields = strings.Split(fields, ",")
		for i, field := range params.Fields {
			params.Fields[i] = strings.TrimSpace(field)
		}
	}

	for key, values := range c.Request.URL.Query() {
		if after, ok := strings.CutPrefix(key, cfg.FilterPrefix); ok {
			filterKey := after
			if len(values) > 0 {
				params.Filters[filterKey] = values[0]
			}
		}
	}

	return params
}

func isValidSort(sort string, allowedSorts []string) bool {
	return slices.Contains(allowedSorts, sort)
}

func isValidOrder(order string, allowedOrders []string) bool {
	return slices.Contains(allowedOrders, order)
}

func GetQueryParams(c *gin.Context) (*QueryParams, bool) {
	if params, exists := c.Get("queryParams"); exists {
		if queryParams, ok := params.(QueryParams); ok {
			return &queryParams, true
		}
	}
	return nil, false
}

func GetSearch(c *gin.Context) string {
	if search, exists := c.Get("search"); exists {
		return search.(string)
	}
	return ""
}

func GetPage(c *gin.Context) int {
	if page, exists := c.Get("page"); exists {
		return page.(int)
	}
	return 1
}

func GetLimit(c *gin.Context) int {
	if limit, exists := c.Get("limit"); exists {
		return limit.(int)
	}
	return 10
}

func GetOffset(c *gin.Context) int {
	if offset, exists := c.Get("offset"); exists {
		return offset.(int)
	}
	return 0
}

func GetSort(c *gin.Context) string {
	if sort, exists := c.Get("sort"); exists {
		return sort.(string)
	}
	return "created_at"
}

func GetOrder(c *gin.Context) string {
	if order, exists := c.Get("order"); exists {
		return order.(string)
	}
	return "desc"
}

func GetFilters(c *gin.Context) map[string]string {
	if filters, exists := c.Get("filters"); exists {
		return filters.(map[string]string)
	}
	return make(map[string]string)
}

func GetFields(c *gin.Context) []string {
	if fields, exists := c.Get("fields"); exists {
		return fields.([]string)
	}
	return []string{}
}
