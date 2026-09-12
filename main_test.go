package main

import (
	"os"
	"path/filepath"
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

func TestLoadPensees(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pensees.md")
	src := "---\ndate: 2026-09-12 14:30\n---\n\nPremière.\n\n---\n\nFin après le hr.\n\n" +
		"---\ndate: 2026-09-11 09:15\ndraft: true\n---\n\nBrouillon.\n\n" +
		"---\ndate: 2026-09-10 08:00\n---\n\nDeuxième [lien](https://go.dev/).\n"
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	ps, err := loadPensees(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 2 {
		t.Fatalf("got %d pensées, want 2 (draft ignorée)", len(ps))
	}
	if !strings.Contains(string(ps[0].HTML), "Première") ||
		!strings.Contains(string(ps[0].HTML), "<hr>") {
		t.Errorf("pensée la plus récente mal rendue : %q", ps[0].HTML)
	}
	if ps[0].Date.Format("2006-01-02 15:04") != "2026-09-12 14:30" {
		t.Errorf("date et heure mal parsées : %v", ps[0].Date)
	}
	if !strings.Contains(string(ps[1].HTML), `href="https://go.dev/"`) {
		t.Errorf("markdown non rendu : %q", ps[1].HTML)
	}
}

func TestLoadPenseesMissingFile(t *testing.T) {
	ps, err := loadPensees(filepath.Join(t.TempDir(), "absent.md"))
	if err != nil || ps != nil {
		t.Fatalf("fichier absent : want nil, nil ; got %v, %v", ps, err)
	}
}
