package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Item struct {
	ID        string      `json:"id"`
	Created   time.Time   `json:"created"`
	URL       string      `json:"url"`
	Method    string      `json:"method"`
	Status    int         `json:"status"`
	UserIP    string      `json:"userIp"`
	RemoteIP  string      `json:"remoteIp"`
	Referer   string      `json:"referer"`
	UserAgent string      `json:"userAgent"`
	Meta      interface{} `json:"meta"`
}

type RequestCount struct {
	Total int    `json:"total"`
	Date  string `json:"date"`
}

var (
	items             []Item
	mutex             sync.Mutex
	requestCountMap   = make(map[string]int)
	requestCountMutex sync.Mutex
)

func addItem(item Item) {
	mutex.Lock()
	items = append(items, item)
	mutex.Unlock()
}

func updateRequestCount() {
	requestCountMutex.Lock()
	defer requestCountMutex.Unlock()

	now := time.Now()
	minute := now.Format("2006-01-02 15:04")
	requestCountMap[minute]++
}

func startRequestCountUpdater() {
	interval := 5 * time.Minute

	for {
		time.Sleep(interval)
		updateRequestCount()
	}
}

func countRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		updateRequestCount()
		return next(c)
	}
}

func main() {

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(countRequests)
	go startRequestCountUpdater()

	api := e.Group("/api")
	api.POST("/log", logar)
	api.GET("/logs/requests", getItems)
	api.GET("/logs/requests/stats", getRequestCountsHandler)

	e.Static("/", "ui/dist")

	e.Logger.Fatal(e.Start(":8080"))

}

func getRequestCountsHandler(c echo.Context) error {
	counts := getRequestCounts()
	return c.JSON(http.StatusOK, counts)
}

func getItems(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("perPage"))

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	startIndex := (page - 1) * perPage
	endIndex := startIndex + perPage

	mutex.Lock()
	defer mutex.Unlock()

	var filteredItems []Item
	if endIndex > len(items) {
		endIndex = len(items)
	}
	if startIndex >= len(items) {
		startIndex = len(items) - 1
	}

	filteredItems = items[startIndex:endIndex]

	totalPages := (len(items) + perPage - 1) / perPage

	apiResponse := struct {
		Page       int    `json:"page"`
		PerPage    int    `json:"perPage"`
		TotalItems int    `json:"totalItems"`
		TotalPages int    `json:"totalPages"`
		Items      []Item `json:"items"`
	}{
		Page:       page,
		PerPage:    perPage,
		TotalItems: len(items),
		TotalPages: totalPages,
		Items:      filteredItems,
	}

	return c.JSON(http.StatusOK, apiResponse)
}

func logar(c echo.Context) error {
	var meta interface{}

	err := c.Bind(&meta)
	if err != nil {
		return c.String(http.StatusBadRequest, "Erro ao decodificar o JSON do corpo da requisição")
	}

	newItem := Item{
		ID:        strconv.Itoa(len(items) + 1),
		Created:   time.Now(),
		URL:       c.Path(),
		Method:    c.Request().Method,
		Status:    http.StatusOK,
		UserIP:    c.RealIP(),
		RemoteIP:  "[demo_redact]",
		Referer:   c.Request().Referer(),
		UserAgent: c.Request().UserAgent(),
		Meta:      meta,
	}

	addItem(newItem)

	return c.String(http.StatusOK, "Item adicionado com sucesso")
}

func getRequestCounts() []RequestCount {
	requestCountMutex.Lock()
	defer requestCountMutex.Unlock()

	var counts []RequestCount
	for date, total := range requestCountMap {
		counts = append(counts, RequestCount{
			Total: total,
			Date:  date,
		})
	}
	return counts
}
