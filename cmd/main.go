package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"parser/internal/services"
	"parser/pkg/sites"
)

func main() {
	cookie, exist := os.LookupEnv("LENTA_COOKIE")
	if !exist {
		panic("LENTA_COOKIE переменная не задана")
	}
	cookie = strings.Trim(cookie, " ")
	if len(cookie) == 0 {
		panic("LENTA_COOKIE переменная не задана")
	}
	categoryName, exist := os.LookupEnv("CATEGORY_NAME")
	if !exist {
		panic("CATEGORY_NAME must be present")
	}
	categoryName = strings.Trim(categoryName, " ")
	defaultLimit := 100
	if limit, exist := os.LookupEnv("LIMIT"); exist {
		if lim, err := strconv.Atoi(limit); err == nil {
			defaultLimit = lim
		}
	}

	lenta := sites.NewLenta()
	service := services.NewService(lenta)
	result, err := service.Parse(categoryName, &defaultLimit)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, res := range result {
		fmt.Printf("Категория: %s\nНазвание: %s\nЦена: %d\nСсылка на страницу: %s\n\n", res.Category, res.Name, int(res.Price), res.Link)
	}
}
