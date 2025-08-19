package scripts

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/RodrigoCelso/gophercises-05/internal/domain"
)

func generateXML(file *os.File, siteMap []domain.Link) error {
	writer := bufio.NewWriter(file)

	urlset := "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n"
	encoder := xml.NewEncoder(writer)

	if _, err := writer.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("erro ao escrever o header no XML: %w", err)
	}

	if _, err := writer.Write([]byte(urlset)); err != nil {
		return fmt.Errorf("erro ao escrever o urlset no XML: %w", err)
	}

	encoder.Indent("", "  ")

	if err := encoder.Encode(siteMap); err != nil {
		return fmt.Errorf("erro ao escrever os dados no XML: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("erro ao escrever o XML no arquivo: %w", err)
	}

	return nil
}

func MakeFile(siteMap []domain.Link) error {
	file, err := os.Create("sitemap.xml")
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo: %w", err)
	}
	defer file.Close()

	if err = generateXML(file, siteMap); err != nil {
		return fmt.Errorf("erro ao gerar o XML: %w", err)
	}

	return nil
}
