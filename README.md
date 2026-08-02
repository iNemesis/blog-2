# Blog statique (Go)

Génère un site statique à partir de fichiers Markdown.

## Structure

```
content/          → tes articles (.md)
templates/        → HTML (base, index, post)
static/           → CSS, images, favicon…
public/           → site généré (à uploader)
main.go           → générateur
```

## Front matter

```yaml
---
title: Titre de l'article
date: 2026-07-31
tags:
  - go
  - blog
image: /images/ma-photo.jpg
summary: Optionnel — sinon le 1er paragraphe est utilisé
draft: false
---

Contenu Markdown ici.
```

## Générer

```bash
go run .
```

Le site est écrit dans `public/`. Pour prévisualiser :

```bash
# PowerShell
cd public; python -m http.server 8080
```

Puis ouvre http://localhost:8080

## Publier

Envoie le contenu de `public/` sur ton hébergement (OVH, Netlify, etc.).
