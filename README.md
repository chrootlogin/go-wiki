# Go-Wiki

**This is work in progress!**

A wiki software written in Go.

## Usage

The easiest way to run go-wiki is with docker:
```
$ docker run -p 80:8000 -e SESSION_KEY=AVerySecureString rootlogin/go-wiki
```

### Environment variables

* **SESSION_KEY**: Sets the session key for the auth cookie encryption. ***(required)***
* **REPOSITORY_PATH**: Sets the data repository path (default: $PWD/data).

### Default login

The following user is available on a new go-wiki installation:

* **Username:** admin
* **Password:** admin1234
 
**Please change this credentials on your first login or delete the user!**

## Documentation

The documentation is available on every new instance by default. [But you can also find it in the repository](default/pages/docs/index.md).

## Development

To work on go-wiki, you need to have Golang and NodeJS installed.

### Dependencies

* [NodeJS](https://nodejs.org) 8.x and NPM
* [Golang](https://golang.org/) 1.10.x
* [Dep](https://golang.github.io/dep/) 0.43.x

### Running

```bash
# Run backend
dep ensure
go run wiki.go

# Run frontend
cd /frontend
npm install
npm run dev
```