DEP=dep
GOLANG=go
GOBINDATA=go-bindata

all: dependencies go_bindata go_app

dependencies:
	$(DEP) ensure

go_bindata:
	$(GOBINDATA) -pkg repo -prefix default/ -o src/lib/repo/default.go default/pages/ default/pages/docs/ default/prefs/

go_app:
	$(GOLANG) build -o go-wiki main.go

test:
	$(GOLANG) test -json -coverprofile=coverage.out ./... > test-report.out

clean:
	rm -rf vendor