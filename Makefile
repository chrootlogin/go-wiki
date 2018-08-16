DEP=dep
GOLANG=go
GOBINDATA=go-bindata

all: dependencies go_app

dependencies:
	$(DEP) ensure

go_app:
	$(GOBINDATA) -pkg repo -prefix default/ -o src/lib/repo/default.go default/pages/ default/pages/docs/ default/prefs/
	$(GOLANG) build -o go-wiki wiki.go

clean:
	rm -rf vendor