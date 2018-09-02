DEP=dep
GOLANG=go
GOBINDATA=go-bindata

all: dependencies go_bindata go_app

dependencies:
	$(DEP) ensure

go_bindata:
	$(GOBINDATA) -pkg pagestore -prefix default/pages -o src/lib/pagestore/default.go default/pages/ default/pages/docs/

go_app:
	$(GOLANG) build -o go-wiki main.go

test:
	$(GOLANG) test ./...

clean:
	rm -rf vendor