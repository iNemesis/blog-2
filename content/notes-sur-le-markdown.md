---
title: Notes sur le Markdown
date: 2026-07-28
tags:
  - markdown
  - ecriture
image: /images/kuzco.jpg
summary: Un rappel concis de la syntaxe supportée par le générateur, pour écrire sans friction.
---

Ce résumé vient du champ `summary` du front matter, pas du premier paragraphe. Pratique quand l'intro est trop longue ou trop technique.

## Syntaxe supportée

Tu disposes de :

- Titres `#` `##` `###`
- **Gras**, *italique*, `code inline`
- Listes à puces
- Citations
- Images `![alt](/images/kuzco.jpg)`
- Liens `[texte](url)`
- Blocs de code avec langage


![alt](/images/kuzco.jpg)

### Astuce image de couverture

Le champ `image` du front matter sert de couverture sur la home **et** en tête d'article. Mets tes fichiers dans `static/images/` : ils seront copiés tels quels dans `public/`.
