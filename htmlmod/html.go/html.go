package htmlmod

import "golang.org/x/net/html"

type PWASettings struct {
	BaseURL string
}

func EnhanceHTML(n *html.Node, pwa *PWASettings) bool {
	if n.Type == html.ElementNode && n.Data == "head" {
		n.InsertBefore(&html.Node{
			Type: html.ElementNode,
			Data: "base",
			Attr: []html.Attribute{
				html.Attribute{
					Key: "href",
					Val: pwa.BaseURL,
				},
			},
		}, n.FirstChild)

		return true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if ok := EnhanceHTML(c, pwa); ok {
			return ok
		}
	}

	return false
}
