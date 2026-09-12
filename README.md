# Blog statique (Go)

Générateur de site statique en Go, simple et sans dépendance superflue. Compile des articles Markdown en HTML statique.

## Structure

```text
content/          → Articles Markdown (.md)
templates/        → Templates HTML (base, index, post)
static/           → Assets statiques copiés tels quels (CSS, images...)
docs/             → Site généré (⚠️ ne pas modifier manuellement)
tools/            → Outils de déploiement SFTP et optimisation d'images
main.go           → Moteur de rendu
Makefile          → Raccourcis de développement
```

## Front matter (`content/*.md`)

```yaml
---
title: Titre de l'article
date: 2026-09-11
tags:
  - go
  - blog
image: /images/ma-photo.jpg
summary: Optionnel — le 1er paragraphe est utilisé par défaut
draft: false
---

Contenu Markdown ici.
```

## Pensées (`content/pensees.md`)

Micro-posts style fil d'actu, publiés sur `/pensees/` (plus récent en
premier). Le fichier est une suite de mini-articles, écrits du plus ancien
(haut) au plus récent (bas) — pour ajouter une pensée, on ajoute un bloc en
bas du fichier :

```yaml
---
date: 2026-09-12 14:30
---

Texte markdown libre : liens, **gras**, images, code, listes…
```

- `date` accepte `YYYY-MM-DD HH:MM` (ou la date seule).
- `draft: true` masque une pensée.
- Le corps est du Markdown complet, rendu comme un article.
- Ne pas insérer de ligne vide entre `---` et `date:` (c'est le séparateur
  de pensées).

## Configuration (`main.go`)

- `siteURL` : domaine public du site, ex. `"https://monblog.fr"` (sans `/` final).
  Tant qu'il est vide, le flux Atom (`atom.xml`), `sitemap.xml`, `robots.txt`,
  la balise `canonical` et les URLs Open Graph absolues ne sont pas générés.

## Déploiement

En plus de `SFTP_HOST`, `SFTP_USER` et (`SFTP_PASSWORD` ou `SFTP_KEY_PATH`), `.env`
doit épingler la clé d'hôte du serveur :

```
SFTP_HOST_KEY=<sortie de : ssh-keyscan -p <port> <host>>
```

La connexion échoue si la clé du serveur change (protection MITM).

## Commandes

| Commande | Description |
|---|---|
| `make` ou `go run .` | Génère le site dans `docs/` |
| `make serve` | Génère puis sert `docs/` sur http://localhost:8080 |
| `make test` | Lance les tests (`go test ./...`) |
| `make fmt` | Formate le code Go et les templates HTML |
| `make optimize` | Optimise les images sans perte |
| `make optimize-lossy` | Compresse les images avec perte |
| `make deploy` | Déploie `docs/` sur le serveur SFTP |
| `make help` | Liste toutes les cibles du Makefile |

## Prévisualisation locale

```bash
make serve
```

Puis ouvrir : http://localhost:8080
