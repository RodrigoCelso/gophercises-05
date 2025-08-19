package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RodrigoCelso/gophercises-05/internal/domain"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func SearchForLinks(node *html.Node, currentPage string, links *[]domain.Link) error {
	// Se não tem mais nodes então chegou no final da página
	if node.FirstChild == nil {
		return nil
	}

	// Procura por uma tag <a> na página
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.DataAtom == atom.A {
			// Extrai as informações da tag
			link, err := FindLink(child, currentPage, *links)
			if err != nil {
				return fmt.Errorf("link não encontrado: %w", err)
			}
			isLinkEmpty := link == domain.Link{}
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
		SearchForLinks(child, currentPage, links)
	}
	return nil
}

func FindLink(node *html.Node, currentPage string, links []domain.Link) (domain.Link, error) {
	link := domain.Link{}
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
	SearchForLinks(node, currentPage, &links)

	return link, nil
}

func AccessPage(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w", errors.New("url inválida: "+url))
	}
	return resp, nil
}
