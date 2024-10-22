package htmlmod

import "golang.org/x/net/html"

// Settings settings for html modifications
type Settings struct {
	// BaseURL used for adding the <base href='BaseURL'> to <head>
	BaseURL string
}

// EnhanceHTML enhanced specified bits of the HTML file, based on the settings
func EnhanceHTML(n *html.Node, pwa *Settings) bool {
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
