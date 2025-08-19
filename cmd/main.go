package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/RodrigoCelso/gophercises-05/internal/domain"
	"github.com/RodrigoCelso/gophercises-05/internal/service"
	"github.com/RodrigoCelso/gophercises-05/scripts"
)

var url string

func parseFlags() {
	service.ParseFlags()
	flag.StringVar(&url, "url", "", "URL para fazer o Sitemap")
	flag.Parse()
}

func main() {
	parseFlags()

	var visitedLinks []string
	siteMapRaw, err := service.CreateSiteMap(url, &visitedLinks, 0)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	sort.Slice(siteMapRaw, func(i, j int) bool {
		return siteMapRaw[i].Href < siteMapRaw[j].Href
	})

	siteMap := []domain.Link{siteMapRaw[0]}
	for i := 1; i < len(siteMapRaw); i++ {
		if siteMapRaw[i] != siteMapRaw[i-1] {
			siteMap = append(siteMap, siteMapRaw[i])
		}
	}

	if err = scripts.MakeFile(siteMap); err != nil {
		fmt.Println("Erro:", err)
		return
	}
}
