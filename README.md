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

## Commandes

| Commande | Description |
|---|---|
| `make` ou `go run .` | Génère le site dans `docs/` |
| `make test` | Lance les tests (`go test ./...`) |
| `make fmt` | Formate le code Go et les templates HTML |
| `make optimize` | Optimise les images sans perte |
| `make optimize-lossy` | Compresse les images avec perte |
| `make deploy` | Déploie `docs/` sur le serveur SFTP |
| `make help` | Liste toutes les cibles du Makefile |

## Prévisualisation locale

```bash
cd docs && python -m http.server 8080
```

Accéder ensuite à : http://localhost:8080
