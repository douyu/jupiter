
# Usage
## Install swag

```
$ go get -u github.com/swaggo/swag/cmd/swag
```

## Generate docs
Run the Swag in your Go project root folder which contains main.go file, Swag will parse comments and generate required files(docs folder and docs/doc.go).
```
$ swag init
```
## Run
```
go run main.go --config=config.toml
```
## Open swagger website
http://localhost:8080/swagger/index.html#/