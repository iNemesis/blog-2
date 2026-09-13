.PHONY: all run serve pensee fmt fmt-html fmt-go test optimize optimize-lossy deploy help

## 🚀 Génère le site (par défaut)
all: run help

## 🚀 Génère le site
run:
	go run .

## 👀 Génère puis prévisualise sur http://localhost:8080
serve: run
	cd docs && python -m http.server 8080

## ✏️ Ajoute une pensée via http://localhost:8081
pensee:
	go run ./tools/pensee

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
	go run ./tools/optimize

## 🖼️ Compresse les images avec perte
optimize-lossy:
	go run ./tools/optimize -lossy

## 📦 Déploie docs/ sur le serveur SFTP
deploy:
	go run ./tools/deploy

## 🎢 Lance toutes les commandes stylées, du fmt au déploiement
full: fmt optimize-lossy run deploy

## ❓ Affiche cette aide
help:
	@echo "  🚀 all              Generer le site (defaut)"
	@echo "  👀 serve            Previsualiser sur http://localhost:8080"
	@echo "  ✏️ pensee           Ajouter une pensee sur localhost:8081"
	@echo "  🧹 fmt              Formater le code Go et HTML"
	@echo "  🧹 fmt-html         Formater les templates HTML"
	@echo "  🧹 fmt-go           Formater les fichiers Go"
	@echo "  🧪 test             Executer les tests"
	@echo "  🖼️ optimize         Optimiser les images sans perte"
	@echo "  🖼️ optimize-lossy   Compresser les images avec perte"
	@echo "  📦 deploy           Deployer docs/ sur le serveur SFTP"
	@echo "  🎢 full             Lance toutes les commandes stylées, du fmt au déploiement"
	@echo "  ❓ help              Affiche cette aide"
