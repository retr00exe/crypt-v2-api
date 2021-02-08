run:
	go run main.go
build:
	go build -o bin/main main.go
tidy:
	go mod tidy
compile:
	GOOS=linux GOARCH=386 go build -o bin/main-linux-i386 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-i386 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 main.go
	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-i386 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-i386 main.go
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 main.go