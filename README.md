# Go-Wiki (WIP)

**This is work in progress!**

A wiki software written in Go.

## Usage

The easiest way to run go-wiki is with docker:
```
$ docker run -p 80:8000 -e SESSION_KEY=AVerySecureString rootlogin/go-wiki
```

### Environment variables

* **SESSION_KEY**: Sets the session key for the auth cookie encryption. ***(required)***

## Development

To work on agw, you need to have Golang and NodeJS installed.

### Dependencies
 * [NodeJS](https://nodejs.org) 8.x and NPM
 * [Golang](https://golang.org/) 1.10.x
 * [Dep](https://golang.github.io/dep/) 0.43.x

### Running

```bash
# Run backend
dep ensure
go run main.go

# Run frontend
cd /frontend
npm install
npm run dev
```