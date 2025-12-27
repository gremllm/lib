package converter

import (
	"bytes"

	"golang.org/x/net/html"
)

type StripConfig struct {
	ElementsToStrip []string
}

// StripElements removes specified HTML elements from the DOM
func StripElements(n *html.Node, tags ...string) {
	tagSet := make(map[string]bool)
	for _, tag := range tags {
		tagSet[tag] = true
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		// Collect nodes to remove (can't remove while iterating)
		var toRemove []*html.Node

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && tagSet[c.Data] {
				shouldKeep := false

				for _, attr := range c.Attr {
					// Check for data-llm attribute, which overrides stripping
					if attr.Key == "data-llm" && attr.Val == "keep" {
						shouldKeep = true
					}
				}

				// If we decided not to keep, mark for removal
				if !shouldKeep {
					toRemove = append(toRemove, c)
				}
			} else {
				shouldKeep := true

				for _, attr := range c.Attr {
					// Check for data-llm attribute, which overrides stripping
					if attr.Key == "data-llm" && attr.Val == "drop" {
						shouldKeep = false
					}
				}

				if !shouldKeep {
					toRemove = append(toRemove, c)
				} else {
					f(c)
				}
			}
		}

		// Remove collected nodes
		for _, node := range toRemove {
			n.RemoveChild(node)
		}
	}
	f(n)
}

// ProcessHTML strips specified tags from HTML based on options
func ProcessHTML(htmlContent []byte, stripConfig StripConfig) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	// Add default elements to strip, if folks want to keep these they can use
	// data-llm="keep" on them
	elementsToStrip := stripConfig.ElementsToStrip
	elementsToStrip = append(elementsToStrip, "header", "footer")

	// Strip specified tags
	StripElements(doc, elementsToStrip...)

	// Serialize back to HTML
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
