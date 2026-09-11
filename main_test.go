package main

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestRelURL(t *testing.T) {
	cases := []struct {
		from, to, want string
	}{
		{"/", "/", "./"},
		{"/", "/css/style.css", "css/style.css"},
		{"/", "/posts/foo/", "posts/foo/"},
		{"/", "/images/atma.jpg", "images/atma.jpg"},
		{"/posts/foo/", "/", "../../"},
		{"/posts/foo/", "/css/style.css", "../../css/style.css"},
		{"/posts/foo/", "/posts/foo/", "./"},
		{"/posts/foo/", "/posts/bar/", "../bar/"},
		{"/posts/foo/", "/tags/go/", "../../tags/go/"},
		{"/tags/go/", "/posts/foo/", "../../posts/foo/"},
		{"/tags/go/", "/", "../../"},
		{"/tags/go/", "/css/style.css", "../../css/style.css"},
		{"/posts/foo/", "https://go.dev/doc/", "https://go.dev/doc/"},
		{"/posts/foo/", "http://example.com/x", "http://example.com/x"},
		{"/posts/foo/", "//cdn.example.com/x", "//cdn.example.com/x"},
		{"/posts/foo/", "mailto:hi@example.com", "mailto:hi@example.com"},
		{"/posts/foo/", "#section", "#section"},
		{"/posts/foo/", "relative.jpg", "relative.jpg"},
	}
	for _, c := range cases {
		got := relURL(c.from, c.to)
		if got != c.want {
			t.Errorf("relURL(%q, %q) = %q, want %q", c.from, c.to, got, c.want)
		}
	}
}

func TestMdToHTMLRewritesRootRelativeURLs(t *testing.T) {
	html := mdToHTML("[texte](/posts/foo/)\n\n![alt](/images/x.jpg)", "/posts/slug/")
	if !strings.Contains(html, `href="../foo/"`) {
		t.Errorf("markdown link not relative:\n%s", html)
	}
	if !strings.Contains(html, `src="../../images/x.jpg"`) {
		t.Errorf("markdown image not relative:\n%s", html)
	}
	if strings.Contains(html, `href="/posts/foo/"`) {
		t.Errorf("root-relative link still present:\n%s", html)
	}
}

func TestMdToHTMLLeavesExternalURLs(t *testing.T) {
	html := mdToHTML("[doc](https://go.dev/doc/)", "/posts/slug/")
	if !strings.Contains(html, `href="https://go.dev/doc/"`) {
		t.Errorf("external URL was rewritten:\n%s", html)
	}
}

func TestMdToHTMLLists(t *testing.T) {
	html := mdToHTML("- un\n- deux\n\n1. premier\n2. deuxième", "/")
	for _, want := range []string{
		"<ul>", "<li>un</li>", "<li>deux</li>", "</ul>",
		"<ol>", "<li>premier</li>", "<li>deuxième</li>", "</ol>",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("manque %q dans :\n%s", want, html)
		}
	}
}

func TestMdToHTMLNestedList(t *testing.T) {
	html := mdToHTML("- a\n  - b\n- c", "/")
	if !strings.Contains(html, "<li>a<ul>") {
		t.Errorf("liste imbriquée hors du <li> parent :\n%s", html)
	}
	if !strings.Contains(html, "</ul>\n</li>") {
		t.Errorf("fermeture de liste imbriquée incorrecte :\n%s", html)
	}
}

func TestFirstParagraphTruncatesOnRunes(t *testing.T) {
	ex := firstParagraph(strings.Repeat("é", 300))
	if !utf8.ValidString(ex) {
		t.Errorf("extrait avec UTF-8 invalide : %q", ex)
	}
	if !strings.HasSuffix(ex, "…") {
		t.Errorf("extrait non tronqué : %q", ex)
	}
}
