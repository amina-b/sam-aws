#Command for invoking lambda: sam local invoke -e events/event.json
#Command for zip lambda: GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
