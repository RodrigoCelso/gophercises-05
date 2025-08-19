package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/RodrigoCelso/gophercises-05/internal/domain"
	"github.com/RodrigoCelso/gophercises-05/internal/service"
	"github.com/RodrigoCelso/gophercises-05/scripts"
)

// func generateXML() {}

func main() {
	// Lê o argumento passado
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Erro:", fmt.Errorf("%v", errors.New("este sistema requer um link como argumento para construir o mapa")))
		return
	}
	urlArg := args[1]

	var visitedLinks []string
	siteMapRaw, err := service.CreateSiteMap(urlArg, &visitedLinks)
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
