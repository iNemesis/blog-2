package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	addr        = ":8081"
	penseesPath = "content/pensees.md"
	dateLayout  = "2006-01-02 15:04"
)

const page = `<!doctype html>
<html lang="fr">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Nouvelle pensée</title>
<style>
	body{font-family:system-ui,sans-serif;max-width:40rem;margin:3rem auto;padding:0 1rem}
	h1{font-size:1.4rem}
	textarea{width:100%;height:14rem;box-sizing:border-box;resize:vertical;font:0.95rem ui-monospace,monospace;padding:.6rem}
	button{margin-top:.6rem;padding:.5rem 1.2rem;font-size:1rem;cursor:pointer}
	#msg{min-height:1.4rem}
	.ok{color:#080}.err{color:#c00}
</style>
</head>
<body>
<h1>💭 Nouvelle pensée</h1>
<textarea id="md" placeholder="Markdown… (Ctrl+Entrée pour ajouter)"></textarea>
<button onclick="ajouter()">Ajouter</button>
<p id="msg"></p>
<script>
	const md = document.getElementById('md'), msg = document.getElementById('msg');
	md.focus();
	md.addEventListener('keydown', e => { if (e.ctrlKey && e.key == 'Enter') ajouter(); });
	async function ajouter() {
		const r = await fetch('/ajouter', {method: 'POST', body: md.value});
		if (r.ok) {
			md.value = '';
			msg.textContent = '✓ Ajouté — relance make run pour régénérer le site';
			msg.className = 'ok';
		} else {
			msg.textContent = '✗ ' + await r.text();
			msg.className = 'err';
		}
	}
</script>
</body>
</html>
`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, page)
	})
	http.HandleFunc("/ajouter", ajouter)
	fmt.Printf("✏️  Éditeur de pensées : http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ajouter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode interdite", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	content := strings.TrimSpace(strings.ReplaceAll(string(body), "\r\n", "\n"))
	if content == "" {
		http.Error(w, "contenu vide", http.StatusBadRequest)
		return
	}

	raw, err := os.ReadFile(penseesPath)
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var b strings.Builder
	if prev := strings.TrimRight(string(raw), "\n"); prev != "" {
		b.WriteString(prev + "\n\n")
	}
	fmt.Fprintf(&b, "---\ndate: %s\n---\n%s\n", time.Now().Format(dateLayout), content)

	if err := os.WriteFile(penseesPath, []byte(b.String()), 0o644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
