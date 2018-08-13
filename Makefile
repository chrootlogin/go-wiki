DEP=dep
GOLANG=go
GOBINDATA=go-bindata

all: dependencies go_app

dependencies:
	$(DEP) ensure

go_app:
	$(GOBINDATA) -pkg frontend -prefix frontend/dist/ -o src/frontend/frontend.go frontend/dist/ frontend/dist/assets/
	$(GOLANG) build -o go-wiki wiki.go

clean:
	rm -rf vendor