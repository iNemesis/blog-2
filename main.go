package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	contentDir  = "content"
	templateDir = "templates"
	staticDir   = "static"
	outputDir   = "docs"

	// Domaine public du site, ex. "https://monblog.fr" (sans / final).
	// Vide, il désactive : flux Atom, sitemap.xml, robots.txt, canonical
	// et URLs Open Graph absolues.
	siteURL = "https://www.atmaaa.fr"
	// Nom et description du site (flux Atom, meta description de l'accueil).
	siteName        = "ATMA's blog"
	siteDescription = "Parlons tech, parlons lead, parlons cuisine."
)

type FrontMatter struct {
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Tags    []string `yaml:"tags"`
	Image   string   `yaml:"image"`
	Draft   bool     `yaml:"draft"`
	Summary string   `yaml:"summary"`
}

type Post struct {
	FrontMatter
	Slug        string
	Content     template.HTML
	Excerpt     string
	ParsedDate  time.Time
	URL         string
	ReadingMins int
}

type Pensee struct {
	Date time.Time
	HTML template.HTML
}

type PageData struct {
	Title       string
	Description string
	Posts       []Post
	Post        *Post
	Pensees     []Pensee
	AllTags     []string
	CurrentTag  string
	GeneratedAt time.Time
	Path        string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Site généré dans ./" + outputDir + "/")
}

func run() error {
	posts, err := loadPosts(contentDir)
	if err != nil {
		return err
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ParsedDate.After(posts[j].ParsedDate)
	})

	pensees, err := loadPensees(filepath.Join(contentDir, "pensees.md"))
	if err != nil {
		return err
	}

	if err := cleanOutput(); err != nil {
		return err
	}
	if err := copyStatic(); err != nil {
		return err
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"join": strings.Join,
		"slug": slugify,
		// "url" est redéfini par page dans render() : relURL dépend du chemin
		// de la page en cours de rendu.
		"url": func(p string) string { return p },
		"abs": func(p string) string {
			if siteURL == "" {
				return ""
			}
			return siteURL + p
		},
		"fmtDate": func(t time.Time) string {
			months := []string{
				"janvier", "février", "mars", "avril", "mai", "juin",
				"juillet", "août", "septembre", "octobre", "novembre", "décembre",
			}
			return fmt.Sprintf("%d %s %d", t.Day(), months[t.Month()-1], t.Year())
		},
		"fmtDateShort": func(t time.Time) string {
			return t.Format("02/01/2006")
		},
	}).ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return fmt.Errorf("templates: %w", err)
	}

	tags := collectTags(posts)
	now := time.Now()

	if err := render(tmpl, "index.html", filepath.Join(outputDir, "index.html"), PageData{
		Title:       "Accueil",
		Description: siteDescription,
		Posts:       posts,
		AllTags:     tags,
		GeneratedAt: now,
		Path:        "/",
	}); err != nil {
		return err
	}

	if err := render(tmpl, "pensees.html", filepath.Join(outputDir, "pensees", "index.html"), PageData{
		Title:       "Pensées",
		Description: "Micro-posts, style fil d'actu.",
		Pensees:     pensees,
		AllTags:     tags,
		GeneratedAt: now,
		Path:        "/pensees/",
	}); err != nil {
		return err
	}

	for _, p := range posts {
		post := p
		dir := filepath.Join(outputDir, "posts", post.Slug)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		if err := render(tmpl, "post.html", filepath.Join(dir, "index.html"), PageData{
			Title:       post.Title,
			Description: post.Excerpt,
			Post:        &post,
			Posts:       posts,
			AllTags:     tags,
			GeneratedAt: now,
			Path:        "/posts/" + post.Slug + "/",
		}); err != nil {
			return err
		}
	}

	for _, tag := range tags {
		var filtered []Post
		for _, p := range posts {
			for _, t := range p.Tags {
				if t == tag {
					filtered = append(filtered, p)
					break
				}
			}
		}
		dir := filepath.Join(outputDir, "tags", slugify(tag))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		if err := render(tmpl, "index.html", filepath.Join(dir, "index.html"), PageData{
			Title:       "Tag: " + tag,
			Description: "Articles tagués « " + tag + " ».",
			Posts:       filtered,
			AllTags:     tags,
			CurrentTag:  tag,
			GeneratedAt: now,
			Path:        "/tags/" + slugify(tag) + "/",
		}); err != nil {
			return err
		}
	}

	if siteURL == "" {
		fmt.Println("  siteURL vide dans main.go : Atom, sitemap, robots et canonical désactivés")
		return nil
	}
	if err := writeAtom(posts, now); err != nil {
		return err
	}
	if err := writeSitemap(posts, tags); err != nil {
		return err
	}
	return writeRobots()
}

func loadPosts(dir string) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("dossier %q introuvable — crée-le et ajoute des .md", dir)
		}
		return nil, err
	}

	var posts []Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
			continue
		}
		if e.Name() == "pensees.md" {
			continue // fil de pensées, rendu à part
		}
		path := filepath.Join(dir, e.Name())
		post, err := parsePost(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if post.Draft {
			fmt.Printf("  skip draft: %s\n", e.Name())
			continue
		}
		posts = append(posts, post)
		fmt.Printf("  + %s\n", post.Title)
	}
	return posts, nil
}

// loadPensees lit content/pensees.md : une suite de mini-articles
// "--- / date: … / --- / markdown", écrits du plus ancien (haut) au plus
// récent (bas). Renvoie les pensées de la plus récente à la plus ancienne.
func loadPensees(path string) ([]Pensee, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	s := strings.ReplaceAll(string(raw), "\r\n", "\n")
	// Chaque pensée ouvre par "---" suivi de "date:" : ce motif de découpe
	// tolère un "---" isolé (hr markdown) dans le corps d'une pensée.
	const sep = "\n---\ndate"
	parts := strings.Split(s, sep)
	var out []Pensee
	for i, part := range parts {
		if i > 0 {
			part = "---\ndate" + part
		}
		if strings.TrimSpace(part) == "" {
			continue
		}
		p, err := parsePensee(part, i+1)
		if err != nil {
			return nil, err
		}
		if p != nil {
			out = append(out, *p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Date.After(out[j].Date) })
	return out, nil
}

func parsePensee(doc string, n int) (*Pensee, error) {
	fm, body, err := splitFrontMatter([]byte(doc))
	if err != nil {
		return nil, fmt.Errorf("pensée %d: %w", n, err)
	}
	var meta FrontMatter
	if err := yaml.Unmarshal(fm, &meta); err != nil {
		return nil, fmt.Errorf("pensée %d: front matter: %w", n, err)
	}
	if meta.Draft {
		return nil, nil
	}
	date, err := parseDate(meta.Date)
	if err != nil {
		return nil, fmt.Errorf("pensée %d: date %q: %w", n, meta.Date, err)
	}
	return &Pensee{
		Date: date,
		HTML: template.HTML(mdToHTML(string(body), "/pensees/")),
	}, nil
}

func parsePost(path string) (Post, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Post{}, err
	}

	fm, body, err := splitFrontMatter(raw)
	if err != nil {
		return Post{}, err
	}

	var meta FrontMatter
	if err := yaml.Unmarshal(fm, &meta); err != nil {
		return Post{}, fmt.Errorf("front matter: %w", err)
	}
	if meta.Title == "" {
		return Post{}, fmt.Errorf("title manquant dans le front matter")
	}
	parsedDate, err := parseDate(meta.Date)
	if err != nil {
		return Post{}, fmt.Errorf("date %q: %w", meta.Date, err)
	}

	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	slug := slugify(base)

	bodyStr := string(body)
	pagePath := "/posts/" + slug + "/"
	htmlBody := mdToHTML(bodyStr, pagePath)
	excerpt := meta.Summary
	if excerpt == "" {
		excerpt = firstParagraph(bodyStr)
	}

	words := len(strings.Fields(bodyStr))
	mins := max(words/200, 1)

	return Post{
		FrontMatter: meta,
		Slug:        slug,
		Content:     template.HTML(htmlBody),
		Excerpt:     excerpt,
		ParsedDate:  parsedDate,
		URL:         pagePath,
		ReadingMins: mins,
	}, nil
}

func splitFrontMatter(raw []byte) (fm []byte, body []byte, err error) {
	const delim = "---"
	s := string(raw)
	s = strings.TrimPrefix(s, "\uFEFF")
	if !strings.HasPrefix(s, delim) {
		return nil, nil, fmt.Errorf("front matter attendu (commence par ---)")
	}
	rest := s[len(delim):]
	// skip optional newline after opening ---
	if strings.HasPrefix(rest, "\r\n") {
		rest = rest[2:]
	} else if strings.HasPrefix(rest, "\n") {
		rest = rest[1:]
	}
	idx := strings.Index(rest, "\n"+delim)
	if idx < 0 {
		return nil, nil, fmt.Errorf("front matter non fermé (--- manquant)")
	}
	fm = []byte(rest[:idx])
	after := rest[idx+1+len(delim):]
	after = strings.TrimPrefix(after, "\r\n")
	after = strings.TrimPrefix(after, "\n")
	return fm, []byte(after), nil
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formats acceptés: YYYY-MM-DD ou YYYY-MM-DD HH:MM")
}

func firstParagraph(md string) string {
	lines := strings.Split(md, "\n")
	var b strings.Builder
	started := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			if started {
				break
			}
			continue
		}
		// skip headings, lists, code fences, images at start
		if !started {
			_, ordered := orderedItem(trim)
			if strings.HasPrefix(trim, "#") ||
				strings.HasPrefix(trim, "```") ||
				strings.HasPrefix(trim, "![") ||
				strings.HasPrefix(trim, "- ") ||
				strings.HasPrefix(trim, "* ") ||
				strings.HasPrefix(trim, "> ") ||
				ordered {
				continue
			}
		}
		started = true
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(trim)
	}
	ex := b.String()
	// strip simple markdown emphasis for excerpt
	ex = strings.ReplaceAll(ex, "**", "")
	ex = strings.ReplaceAll(ex, "__", "")
	ex = strings.ReplaceAll(ex, "*", "")
	ex = strings.ReplaceAll(ex, "_", "")
	ex = strings.ReplaceAll(ex, "`", "")
	if runes := []rune(ex); len(runes) > 280 {
		ex = string(runes[:277]) + "…"
	}
	return ex
}

// Minimal Markdown → HTML (enough for a personal blog).
func mdToHTML(md, fromPage string) string {
	lines := strings.Split(md, "\n")
	var out strings.Builder
	inCode := false
	inBlockquote := false
	var lists []openList // pile des listes ouvertes (ul/ol imbriquées)
	var para strings.Builder

	flushPara := func() {
		if para.Len() == 0 {
			return
		}
		out.WriteString("<p>")
		out.WriteString(inlineMD(para.String(), fromPage))
		out.WriteString("</p>\n")
		para.Reset()
	}
	closeItem := func() {
		if n := len(lists); n > 0 && lists[n-1].liOpen {
			out.WriteString("</li>\n")
			lists[n-1].liOpen = false
		}
	}
	closeLists := func() {
		for len(lists) > 0 {
			closeItem()
			out.WriteString("</" + lists[len(lists)-1].tag + ">\n")
			lists = lists[:len(lists)-1]
		}
	}
	closeBQ := func() {
		if inBlockquote {
			out.WriteString("</blockquote>\n")
			inBlockquote = false
		}
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trim := strings.TrimSpace(line)

		if strings.HasPrefix(trim, "```") {
			flushPara()
			closeLists()
			closeBQ()
			if !inCode {
				lang := strings.TrimPrefix(trim, "```")
				lang = strings.TrimSpace(lang)
				if lang != "" {
					out.WriteString(`<pre><code class="language-` + htmlEscape(lang) + `">`)
				} else {
					out.WriteString("<pre><code>")
				}
				inCode = true
			} else {
				out.WriteString("</code></pre>\n")
				inCode = false
			}
			continue
		}
		if inCode {
			out.WriteString(htmlEscape(line))
			out.WriteByte('\n')
			continue
		}

		if trim == "" {
			flushPara()
			closeLists()
			closeBQ()
			continue
		}

		// headings
		if strings.HasPrefix(trim, "### ") {
			flushPara()
			closeLists()
			closeBQ()
			out.WriteString("<h3>" + inlineMD(strings.TrimPrefix(trim, "### "), fromPage) + "</h3>\n")
			continue
		}
		if strings.HasPrefix(trim, "## ") {
			flushPara()
			closeLists()
			closeBQ()
			out.WriteString("<h2>" + inlineMD(strings.TrimPrefix(trim, "## "), fromPage) + "</h2>\n")
			continue
		}
		if strings.HasPrefix(trim, "# ") {
			flushPara()
			closeLists()
			closeBQ()
			out.WriteString("<h1>" + inlineMD(strings.TrimPrefix(trim, "# "), fromPage) + "</h1>\n")
			continue
		}

		// hr
		if trim == "---" || trim == "***" || trim == "___" {
			flushPara()
			closeLists()
			closeBQ()
			out.WriteString("<hr>\n")
			continue
		}

		// image alone on line
		if strings.HasPrefix(trim, "![") {
			if alt, src, ok := parseImage(trim, fromPage); ok {
				flushPara()
				closeLists()
				closeBQ()
				out.WriteString(`<figure><img src="` + htmlEscape(src) + `" alt="` + htmlEscape(alt) + `" loading="lazy"></figure>` + "\n")
				continue
			}
		}

		// list items : "- ", "* ", "1. " — imbrication par indentation (2 espaces)
		if typ, text, depth, ok := parseListItem(line, trim); ok {
			flushPara()
			closeBQ()
			for len(lists) > depth+1 {
				closeItem()
				out.WriteString("</" + lists[len(lists)-1].tag + ">\n")
				lists = lists[:len(lists)-1]
			}
			if len(lists) == depth+1 {
				closeItem()
				if lists[len(lists)-1].tag != typ {
					out.WriteString("</" + lists[len(lists)-1].tag + ">\n")
					lists[len(lists)-1].tag = typ
					out.WriteString("<" + typ + ">\n")
				}
			}
			for len(lists) < depth+1 {
				out.WriteString("<" + typ + ">\n")
				lists = append(lists, openList{tag: typ})
			}
			out.WriteString("<li>" + inlineMD(text, fromPage))
			lists[len(lists)-1].liOpen = true
			continue
		}

		// blockquote
		if strings.HasPrefix(trim, "> ") {
			flushPara()
			closeLists()
			if !inBlockquote {
				out.WriteString("<blockquote>\n")
				inBlockquote = true
			}
			out.WriteString("<p>" + inlineMD(strings.TrimPrefix(trim, "> "), fromPage) + "</p>\n")
			continue
		}

		closeLists()
		closeBQ()
		if para.Len() > 0 {
			para.WriteByte(' ')
		}
		para.WriteString(trim)
	}

	flushPara()
	closeLists()
	closeBQ()
	if inCode {
		out.WriteString("</code></pre>\n")
	}
	return out.String()
}

func parseImage(s, fromPage string) (alt, src string, ok bool) {
	// ![alt](src)
	if !strings.HasPrefix(s, "![") {
		return "", "", false
	}
	endAlt := strings.Index(s, "](")
	if endAlt < 0 {
		return "", "", false
	}
	alt = s[2:endAlt]
	rest := s[endAlt+2:]
	endSrc := strings.Index(rest, ")")
	if endSrc < 0 {
		return "", "", false
	}
	src = relURL(fromPage, rest[:endSrc])
	return alt, src, true
}

func inlineMD(s, fromPage string) string {
	// order matters: code, links/images, bold, italic
	s = htmlEscape(s)

	// `code`
	s = replaceDelim(s, "`", func(inner string) string {
		return "<code>" + inner + "</code>"
	})

	// [text](url)
	s = replaceLinks(s, fromPage)

	// **bold**
	s = replaceWrapped(s, "**", "strong")
	// *italic*
	s = replaceWrapped(s, "*", "em")

	return s
}

func replaceDelim(s, delim string, wrap func(string) string) string {
	var b strings.Builder
	for {
		i := strings.Index(s, delim)
		if i < 0 {
			b.WriteString(s)
			break
		}
		b.WriteString(s[:i])
		s = s[i+len(delim):]
		j := strings.Index(s, delim)
		if j < 0 {
			b.WriteString(delim)
			b.WriteString(s)
			break
		}
		b.WriteString(wrap(s[:j]))
		s = s[j+len(delim):]
	}
	return b.String()
}

func replaceWrapped(s, delim, tag string) string {
	return replaceDelim(s, delim, func(inner string) string {
		return "<" + tag + ">" + inner + "</" + tag + ">"
	})
}

func replaceLinks(s, fromPage string) string {
	var b strings.Builder
	for {
		i := strings.Index(s, "[")
		if i < 0 {
			b.WriteString(s)
			break
		}
		b.WriteString(s[:i])
		s = s[i:]
		endText := strings.Index(s, "](")
		if endText < 0 {
			b.WriteByte('[')
			s = s[1:]
			continue
		}
		text := s[1:endText]
		rest := s[endText+2:]
		endURL := strings.Index(rest, ")")
		if endURL < 0 {
			b.WriteByte('[')
			s = s[1:]
			continue
		}
		url := relURL(fromPage, rest[:endURL])
		b.WriteString(`<a href="` + url + `">` + text + `</a>`)
		s = rest[endURL+1:]
	}
	return b.String()
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
	)
	return replacer.Replace(s)
}

func relURL(fromPage, target string) string {
	if target == "" {
		target = "/"
	}
	if strings.Contains(target, "://") ||
		strings.HasPrefix(target, "//") ||
		strings.HasPrefix(target, "mailto:") ||
		strings.HasPrefix(target, "#") ||
		!strings.HasPrefix(target, "/") {
		return target
	}

	fromParts := splitURLPath(strings.Trim(fromPage, "/"))
	keepSlash := strings.HasSuffix(target, "/")
	toParts := splitURLPath(strings.Trim(target, "/"))

	i := 0
	for i < len(fromParts) && i < len(toParts) && fromParts[i] == toParts[i] {
		i++
	}

	var parts []string
	for j := i; j < len(fromParts); j++ {
		parts = append(parts, "..")
	}
	parts = append(parts, toParts[i:]...)
	if len(parts) == 0 {
		if keepSlash {
			return "./"
		}
		return "."
	}
	out := strings.Join(parts, "/")
	if keepSlash {
		out += "/"
	}
	return out
}

func splitURLPath(p string) []string {
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		case r == ' ' || r == '_' || r == '-' || r == '.':
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		case r == 'é' || r == 'è' || r == 'ê' || r == 'ë':
			b.WriteByte('e')
			prevDash = false
		case r == 'à' || r == 'â' || r == 'ä':
			b.WriteByte('a')
			prevDash = false
		case r == 'ù' || r == 'û' || r == 'ü':
			b.WriteByte('u')
			prevDash = false
		case r == 'ô' || r == 'ö':
			b.WriteByte('o')
			prevDash = false
		case r == 'î' || r == 'ï':
			b.WriteByte('i')
			prevDash = false
		case r == 'ç':
			b.WriteByte('c')
			prevDash = false
		}
	}
	return strings.Trim(b.String(), "-")
}

func collectTags(posts []Post) []string {
	seen := map[string]bool{}
	var tags []string
	for _, p := range posts {
		for _, t := range p.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	sort.Strings(tags)
	return tags
}

func render(tmpl *template.Template, name, outPath string, data PageData) error {
	// Clone par page : la fonction "url" dépend du chemin de la page rendue.
	pageTmpl, err := tmpl.Clone()
	if err != nil {
		return err
	}
	pagePath := data.Path
	pageTmpl.Funcs(template.FuncMap{
		"url": func(p string) string { return relURL(pagePath, p) },
	})
	var buf bytes.Buffer
	if err := pageTmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return fmt.Errorf("render %s: %w", name, err)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, buf.Bytes(), 0o644)
}

func cleanOutput() error {
	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}
	return os.MkdirAll(outputDir, 0o755)
}

func copyStatic() error {
	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		rel, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(outputDir, rel)
		if info.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		return copyFile(path, dst)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// openList : liste ouverte pendant le rendu Markdown (tag + <li> en cours).
type openList struct {
	tag    string
	liOpen bool
}

// parseListItem reconnaît un item de liste ("- ", "* ", "1. ") et retourne son
// type, son texte et sa profondeur (2 espaces ou 1 tabulation par niveau).
func parseListItem(line, trim string) (typ, text string, depth int, ok bool) {
	lead := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
	depth = (strings.Count(lead, " ") + 2*strings.Count(lead, "\t")) / 2
	switch {
	case strings.HasPrefix(trim, "- "), strings.HasPrefix(trim, "* "):
		return "ul", trim[2:], depth, true
	}
	if rest, ok := orderedItem(trim); ok {
		return "ol", rest, depth, true
	}
	return "", "", 0, false
}

// orderedItem reconnaît un item de liste ordonnée ("12. " ou "12) ").
func orderedItem(s string) (string, bool) {
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == 0 || i > 3 || i+1 >= len(s) ||
		(s[i] != '.' && s[i] != ')') ||
		(s[i+1] != ' ' && s[i+1] != '\t') {
		return "", false
	}
	return strings.TrimLeft(s[i+2:], " \t"), true
}

// writeAtom génère le flux Atom des 20 derniers articles (URLs absolues).
func writeAtom(posts []Post, now time.Time) error {
	updated := now
	if len(posts) > 0 {
		updated = posts[0].ParsedDate
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="utf-8"?>` + "\n")
	b.WriteString(`<feed xmlns="http://www.w3.org/2005/Atom">` + "\n")
	fmt.Fprintf(&b, "  <title>%s</title>\n", htmlEscape(siteName))
	fmt.Fprintf(&b, "  <id>%s/</id>\n", siteURL)
	fmt.Fprintf(&b, "  <link href=\"%s/\"/>\n", siteURL)
	fmt.Fprintf(&b, "  <link rel=\"self\" href=\"%s/atom.xml\"/>\n", siteURL)
	fmt.Fprintf(&b, "  <updated>%s</updated>\n", updated.Format(time.RFC3339))
	fmt.Fprintf(&b, "  <author><name>%s</name></author>\n", htmlEscape(siteName))
	for i, p := range posts {
		if i == 20 {
			break
		}
		fmt.Fprintf(&b, "  <entry>\n")
		fmt.Fprintf(&b, "    <title>%s</title>\n", htmlEscape(p.Title))
		fmt.Fprintf(&b, "    <id>%s%s</id>\n", siteURL, p.URL)
		fmt.Fprintf(&b, "    <link href=\"%s%s\"/>\n", siteURL, p.URL)
		fmt.Fprintf(&b, "    <updated>%s</updated>\n", p.ParsedDate.Format(time.RFC3339))
		if p.Excerpt != "" {
			fmt.Fprintf(&b, "    <summary>%s</summary>\n", htmlEscape(p.Excerpt))
		}
		b.WriteString("  </entry>\n")
	}
	b.WriteString("</feed>\n")
	return os.WriteFile(filepath.Join(outputDir, "atom.xml"), []byte(b.String()), 0o644)
}

func writeSitemap(posts []Post, tags []string) error {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="utf-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	entry := func(p string, lastmod time.Time) {
		b.WriteString("<url><loc>" + siteURL + p + "</loc>")
		if !lastmod.IsZero() {
			b.WriteString("<lastmod>" + lastmod.Format("2006-01-02") + "</lastmod>")
		}
		b.WriteString("</url>\n")
	}
	entry("/", time.Time{})
	for _, p := range posts {
		entry(p.URL, p.ParsedDate)
	}
	for _, t := range tags {
		entry("/tags/"+slugify(t)+"/", time.Time{})
	}
	b.WriteString("</urlset>\n")
	return os.WriteFile(filepath.Join(outputDir, "sitemap.xml"), []byte(b.String()), 0o644)
}

func writeRobots() error {
	robots := "User-agent: *\nAllow: /\nSitemap: " + siteURL + "/sitemap.xml\n"
	return os.WriteFile(filepath.Join(outputDir, "robots.txt"), []byte(robots), 0o644)
}
