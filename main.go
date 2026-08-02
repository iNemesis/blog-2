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
	outputDir   = "public"
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

type PageData struct {
	Title       string
	Posts       []Post
	Post        *Post
	AllTags     []string
	CurrentTag  string
	GeneratedAt time.Time
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

	if err := cleanOutput(); err != nil {
		return err
	}
	if err := copyStatic(); err != nil {
		return err
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"join":  strings.Join,
		"slug":  slugify,
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
		Title:       "Blog",
		Posts:       posts,
		AllTags:     tags,
		GeneratedAt: now,
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
			Post:        &post,
			Posts:       posts,
			AllTags:     tags,
			GeneratedAt: now,
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
			Posts:       filtered,
			AllTags:     tags,
			CurrentTag:  tag,
			GeneratedAt: now,
		}); err != nil {
			return err
		}
	}

	return nil
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
	htmlBody := mdToHTML(bodyStr)
	excerpt := meta.Summary
	if excerpt == "" {
		excerpt = firstParagraph(bodyStr)
	}

	words := len(strings.Fields(bodyStr))
	mins := words / 200
	if mins < 1 {
		mins = 1
	}

	return Post{
		FrontMatter: meta,
		Slug:        slug,
		Content:     template.HTML(htmlBody),
		Excerpt:     excerpt,
		ParsedDate:  parsedDate,
		URL:         "/posts/" + slug + "/",
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
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formats acceptés: YYYY-MM-DD")
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
			if strings.HasPrefix(trim, "#") ||
				strings.HasPrefix(trim, "```") ||
				strings.HasPrefix(trim, "![") ||
				strings.HasPrefix(trim, "- ") ||
				strings.HasPrefix(trim, "* ") ||
				strings.HasPrefix(trim, "> ") {
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
	if len(ex) > 280 {
		ex = ex[:277] + "…"
	}
	return ex
}

// Minimal Markdown → HTML (enough for a personal blog).
func mdToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var out strings.Builder
	inCode := false
	inList := false
	inBlockquote := false
	var para strings.Builder

	flushPara := func() {
		if para.Len() == 0 {
			return
		}
		out.WriteString("<p>")
		out.WriteString(inlineMD(para.String()))
		out.WriteString("</p>\n")
		para.Reset()
	}
	closeList := func() {
		if inList {
			out.WriteString("</ul>\n")
			inList = false
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
			closeList()
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
			closeList()
			closeBQ()
			continue
		}

		// headings
		if strings.HasPrefix(trim, "### ") {
			flushPara()
			closeList()
			closeBQ()
			out.WriteString("<h3>" + inlineMD(strings.TrimPrefix(trim, "### ")) + "</h3>\n")
			continue
		}
		if strings.HasPrefix(trim, "## ") {
			flushPara()
			closeList()
			closeBQ()
			out.WriteString("<h2>" + inlineMD(strings.TrimPrefix(trim, "## ")) + "</h2>\n")
			continue
		}
		if strings.HasPrefix(trim, "# ") {
			flushPara()
			closeList()
			closeBQ()
			out.WriteString("<h1>" + inlineMD(strings.TrimPrefix(trim, "# ")) + "</h1>\n")
			continue
		}

		// hr
		if trim == "---" || trim == "***" || trim == "___" {
			flushPara()
			closeList()
			closeBQ()
			out.WriteString("<hr>\n")
			continue
		}

		// image alone on line
		if strings.HasPrefix(trim, "![") {
			if alt, src, ok := parseImage(trim); ok {
				flushPara()
				closeList()
				closeBQ()
				out.WriteString(`<figure><img src="` + htmlEscape(src) + `" alt="` + htmlEscape(alt) + `" loading="lazy"></figure>` + "\n")
				continue
			}
		}

		// unordered list
		if strings.HasPrefix(trim, "- ") || strings.HasPrefix(trim, "* ") {
			flushPara()
			closeBQ()
			item := trim[2:]
			if !inList {
				out.WriteString("<ul>\n")
				inList = true
			}
			out.WriteString("<li>" + inlineMD(item) + "</li>\n")
			continue
		}

		// blockquote
		if strings.HasPrefix(trim, "> ") {
			flushPara()
			closeList()
			if !inBlockquote {
				out.WriteString("<blockquote>\n")
				inBlockquote = true
			}
			out.WriteString("<p>" + inlineMD(strings.TrimPrefix(trim, "> ")) + "</p>\n")
			continue
		}

		closeList()
		closeBQ()
		if para.Len() > 0 {
			para.WriteByte(' ')
		}
		para.WriteString(trim)
	}

	flushPara()
	closeList()
	closeBQ()
	if inCode {
		out.WriteString("</code></pre>\n")
	}
	return out.String()
}

func parseImage(s string) (alt, src string, ok bool) {
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
	src = rest[:endSrc]
	return alt, src, true
}

func inlineMD(s string) string {
	// order matters: code, links/images, bold, italic
	s = htmlEscape(s)

	// `code`
	s = replaceDelim(s, "`", func(inner string) string {
		return "<code>" + inner + "</code>"
	})

	// [text](url)
	s = replaceLinks(s)

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

func replaceLinks(s string) string {
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
		url := rest[:endURL]
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
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
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
