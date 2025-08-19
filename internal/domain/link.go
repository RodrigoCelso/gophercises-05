package domain

import "encoding/xml"

type Link struct {
	XMLName xml.Name `xml:"url"`
	Href    string   `xml:"loc"`
}
