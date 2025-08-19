package service

import (
	"flag"
	"fmt"
	"slices"

	"github.com/RodrigoCelso/gophercises-05/internal/domain"
	"golang.org/x/net/html"
)

var depth int

func ParseFlags() {
	flag.IntVar(&depth, "d", -1, "Maximum depth to search")
}

func CreateSiteMap(path string, visitedLinks *[]string, currentDepth int) ([]domain.Link, error) {
	siteMap := []domain.Link{}
	if currentDepth < depth || depth == -1 {
		if slices.Contains(*visitedLinks, path) {
			return nil, nil
		}
		*visitedLinks = append(*visitedLinks, path)

		err := CreatePageMap(path, &siteMap)
		if err != nil {
			err = fmt.Errorf("erro ao criar mapa da página: %w", err)
			return nil, err
		}

		for _, site := range siteMap {
			links := *visitedLinks
			if !slices.Contains(links, site.Href) {
				newMap, err := CreateSiteMap(site.Href, visitedLinks, currentDepth+1)
				if err != nil {
					err = fmt.Errorf("erro ao criar o mapa do site: %w", err)
					return nil, err
				}
				siteMap = append(siteMap, newMap...)
			}
		}
	}
	return siteMap, nil
}

func CreatePageMap(pageUrl string, result *[]domain.Link) error {
	var err error

	// Acessa a página html
	fmt.Println("Página acessada: " + pageUrl)
	htmlPage, err := AccessPage(pageUrl)
	if err != nil {
		err = fmt.Errorf("erro ao acessar a página: %w", err)
		return err
	}

	// Transforma a página em nós de tags html
	nodes, err := html.Parse(htmlPage.Body)
	if err != nil {
		err = fmt.Errorf("página inválida: %w", err)
		return err
	}

	// Constrói um slice com os links encontrados para o mesmo domínio na página
	err = SearchForLinks(nodes, pageUrl, result)
	if err != nil {
		return fmt.Errorf("erro ao buscar os links: %w", err)
	}
	return nil
}
