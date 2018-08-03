DEP=dep
GOLANG=go
NPM=npm
GOBINDATA=go-bindata

all: dependencies go_app

dependencies:
	$(DEP) ensure

go_app:
	$(GOBINDATA) -pkg repo -prefix default/ -o src/repo/default.go default/pages/
	$(GOLANG) build -o wiki wiki.go

web_app:
	cd frontend && $(NPM) install && \
	  $(NPM) run build

clean:
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf vendor