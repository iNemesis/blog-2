.PHONY: all run serve fmt fmt-html fmt-go test optimize optimize-lossy deploy help

## 🚀 Génère le site (par défaut)
all: run

## 🚀 Génère le site
run:
	go run .

## 👀 Génère puis prévisualise sur http://localhost:8080
serve: run
	cd docs && python -m http.server 8080

## 🧹 Formate les templates HTML
fmt-html:
	npx prettier --write "templates/**/*.html"

## 🧹 Formate les fichiers Go
fmt-go:
	go fmt ./...

## 🧹 Formate le code Go et HTML
fmt: fmt-html fmt-go

## 🧪 Exécute les tests
test:
	go test ./...

## 🖼️ Optimise les images sans perte
optimize:
	go run tools/optimize.go

## 🖼️ Compresse les images avec perte
optimize-lossy:
	go run tools/optimize.go -lossy

## 📦 Déploie docs/ sur le serveur SFTP
deploy:
	go run tools/deploy.go

## ❓ Affiche cette aide
help:
	@awk '/^## / { desc = substr($$0, 4); next } desc && /^[a-zA-Z_-]+:/ { sub(/:.*/, "", $$1); printf "  \033[36m%-16s\033[0m %s\n", $$1, desc; desc = "" }' $(MAKEFILE_LIST)
