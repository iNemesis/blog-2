package main

import (
	"strings"
	"testing"
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
