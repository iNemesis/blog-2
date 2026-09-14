---
title: Bienvenue sur mon blog
date: 2026-07-20
tags:
  - meta
  - go
image: /images/welcome.jpg
draft: true
---

Voici le premier paragraphe : c'est lui qui apparaît en extrait sur la page d'accueil, juste sous le titre. Tu peux aussi forcer un résumé avec le champ `summary` du front matter.

## Pourquoi un générateur statique ?

Parce que c'est simple, rapide, et que tu contrôles tout. Tu écris du Markdown, tu lances une commande Go, tu uploades le dossier `public/` sur OVH (ou ailleurs).

- Pas de base de données
- Pas de runtime côté serveur
- Du HTML tout bête

### Un peu de code

```go
package main

func main() {
    println("hello blog")
}
```

Et un lien utile : [Documentation Go](https://go.dev/doc/).

> Écris, génère, publie. Rinse and repeat.
