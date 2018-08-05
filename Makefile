DEP=dep
GOLANG=go
NPM=npm
GOBINDATA=go-bindata

all: dependencies web_app go_app

dependencies:
	$(DEP) ensure

go_app:
	$(GOBINDATA) -pkg frontend -prefix frontend/dist/ -o src/frontend/frontend.go frontend/dist/
	$(GOBINDATA) -pkg repo -prefix default/ -o src/repo/default.go default/pages/ default/prefs/
	$(GOLANG) build -o go-wiki wiki.go

web_app:
	cd frontend && $(NPM) install && \
	  $(NPM) run build

clean:
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf vendor