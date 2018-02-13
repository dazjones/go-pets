# go-pets
Simple API for testing. Written in Go.

## Usage
``go run main.go```

## Get pets
```curl http://127.0.0.1:8080/pets```

## Create pet
```curl -X POST '{"type":"dog", "name":"spot"} http://127.0.0.1:8080/pets```
