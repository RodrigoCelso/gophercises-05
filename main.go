package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"sort"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	XMLName xml.Name `xml:"url"`
	Href    string   `xml:"loc"`
}

func parsePath(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w", errors.New("url inválida: "+url))
	}
	return resp, nil
}

func findLink(node *html.Node, currentPage string, links []Link) (Link, error) {
	link := Link{}
	// Percorre por todos os atribudos do node procurando pelo href
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			value := attr.Val

			// Converte a url de entrada e a url do href para *url.URL
			hrefUrl, err := url.Parse(value)
			if err != nil {
				return link, fmt.Errorf("erro ao converter o href para url: %w", err)
			}
			hostUrl, err := url.Parse(currentPage)
			if err != nil {
				return link, fmt.Errorf("erro ao converter o host para url: %w", err)
			}

			// Verifica se a url do href faz não é uma url para outro domínio
			if hrefUrl.Host == "" || hrefUrl.Host == hostUrl.Host {
				link.Href = hostUrl.Scheme + "://" + hostUrl.Host + hrefUrl.Path
			} else {
				return link, nil
			}
		}
	}
	// Verifica se não tem outro link dentro deste link
	searchForLinks(node, currentPage, &links)

	return link, nil
}

func searchForLinks(node *html.Node, currentPage string, links *[]Link) error {
	// Se não tem mais nodes então chegou no final da página
	if node.FirstChild == nil {
		return nil
	}

	// Procura por uma tag <a> na página
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.DataAtom == atom.A {
			// Extrai as informações da tag
			link, err := findLink(child, currentPage, *links)
			if err != nil {
				return fmt.Errorf("link não encontrado: %w", err)
			}
			isLinkEmpty := link == Link{}
			if isLinkEmpty {
				continue
			}

			// Para todos os resultados já encontrados
			linkExists := false
			for _, l := range *links {
				// O link encontrado já existe?
				if link.Href == l.Href {
					linkExists = true
				}
			}
			// Se o link existir, ignore e continue para o próximo child
			if linkExists {
				continue
			}

			*links = append(*links, link)
		}

		// Se não encontrou a tag então continua procurando
		searchForLinks(child, currentPage, links)
	}
	return nil
}

func createPageMap(pageUrl string, result *[]Link) error {
	var err error

	// Acessa a página html
	fmt.Println("Página acessada: " + pageUrl)
	htmlPage, err := parsePath(pageUrl)
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
	err = searchForLinks(nodes, pageUrl, result)
	if err != nil {
		return fmt.Errorf("erro ao buscar os links: %w", err)
	}
	return nil
}

func createSiteMap(path string, visitedLinks *[]string) ([]Link, error) {
	siteMap := []Link{}

	if slices.Contains(*visitedLinks, path) {
		return nil, nil
	}
	*visitedLinks = append(*visitedLinks, path)

	err := createPageMap(path, &siteMap)
	if err != nil {
		err = fmt.Errorf("erro ao criar mapa da página: %w", err)
		return nil, err
	}

	for _, site := range siteMap {
		links := *visitedLinks
		if !slices.Contains(links, site.Href) {
			newMap, err := createSiteMap(site.Href, visitedLinks)
			if err != nil {
				err = fmt.Errorf("erro ao criar o mapa do site: %w", err)
				return nil, err
			}
			siteMap = append(siteMap, newMap...)
		}
	}
	return siteMap, nil
}

func main() {
	// Lê o argumento passado
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Erro:", fmt.Errorf("%v", errors.New("este sistema requer um link como argumento para construir o mapa")))
	}
	urlArg := args[1]

	var visitedLinks []string
	siteMap, err := createSiteMap(urlArg, &visitedLinks)
	if err != nil {
		fmt.Println("Erro:", err)
	}
	sort.Slice(siteMap, func(i, j int) bool {
		return siteMap[i].Href < siteMap[j].Href
	})

	// Gera o XML
	urlset := "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n"
	encoder := xml.NewEncoder(os.Stdout)
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write([]byte(urlset))
	encoder.Indent("", "  ")
	err = encoder.Encode(siteMap)
	if err != nil {
		fmt.Println("Erro:", errors.New("erro ao gerar o XML"))
		return
	}
}
