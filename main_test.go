package main

import (
	"strings"
	"testing"
)

func TestNormalizeBasePath(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"   ", ""},
		{"/", ""},
		{"blog-2", "/blog-2"},
		{"/blog-2", "/blog-2"},
		{"/blog-2/", "/blog-2"},
		{"blog-2/", "/blog-2"},
		{"/foo/bar/", "/foo/bar"},
	}
	for _, c := range cases {
		got := normalizeBasePath(c.in)
		if got != c.want {
			t.Errorf("normalizeBasePath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestJoinBasePath(t *testing.T) {
	cases := []struct {
		prefix, path, want string
	}{
		{"", "/", "/"},
		{"", "/css/style.css", "/css/style.css"},
		{"", "/posts/foo/", "/posts/foo/"},
		{"/blog-2", "/", "/blog-2/"},
		{"/blog-2", "/css/style.css", "/blog-2/css/style.css"},
		{"/blog-2", "/posts/foo/", "/blog-2/posts/foo/"},
		{"/blog-2", "/tags/go/", "/blog-2/tags/go/"},
		{"/blog-2", "/images/welcome.jpg", "/blog-2/images/welcome.jpg"},
		{"/blog-2", "https://go.dev/doc/", "https://go.dev/doc/"},
		{"/blog-2", "http://example.com/x", "http://example.com/x"},
		{"/blog-2", "//cdn.example.com/x", "//cdn.example.com/x"},
		{"/blog-2", "mailto:hi@example.com", "mailto:hi@example.com"},
		{"/blog-2", "#section", "#section"},
		{"/blog-2", "relative.jpg", "relative.jpg"},
		{"/blog-2", "", "/blog-2/"},
	}
	for _, c := range cases {
		got := joinBasePath(c.prefix, c.path)
		if got != c.want {
			t.Errorf("joinBasePath(%q, %q) = %q, want %q", c.prefix, c.path, got, c.want)
		}
	}
}

func TestAbsURLReadsBASE_PATH(t *testing.T) {
	t.Setenv("BASE_PATH", "/blog-2")
	got := absURL("/css/style.css")
	want := "/blog-2/css/style.css"
	if got != want {
		t.Errorf("absURL(/css/style.css) with BASE_PATH=/blog-2 = %q, want %q", got, want)
	}

	t.Setenv("BASE_PATH", "")
	got = absURL("/css/style.css")
	if got != "/css/style.css" {
		t.Errorf("absURL without BASE_PATH = %q, want %q", got, "/css/style.css")
	}
}

func TestMdToHTMLPrefixesRootRelativeURLs(t *testing.T) {
	t.Setenv("BASE_PATH", "/blog-2")

	html := mdToHTML("[texte](/posts/foo/)\n\n![alt](/images/x.jpg)")
	if !strings.Contains(html, `href="/blog-2/posts/foo/"`) {
		t.Errorf("markdown link not prefixed:\n%s", html)
	}
	if !strings.Contains(html, `src="/blog-2/images/x.jpg"`) {
		t.Errorf("markdown image not prefixed:\n%s", html)
	}
	if strings.Contains(html, `href="/posts/foo/"`) {
		t.Errorf("unprefixed link still present:\n%s", html)
	}
}

func TestMdToHTMLLeavesExternalURLs(t *testing.T) {
	t.Setenv("BASE_PATH", "/blog-2")
	html := mdToHTML("[doc](https://go.dev/doc/)")
	if !strings.Contains(html, `href="https://go.dev/doc/"`) {
		t.Errorf("external URL was rewritten:\n%s", html)
	}
}
