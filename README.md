Comandos para las instalaciones de proto gen

- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

- protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative studentpb/student.proto

- protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative testspb/tests.proto

- protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=./proto/generated --go-grpc_opt=paths=source_relative proto/student.proto

- go get google.golang.org/protobuf
- go get github.com/lib/pq
- go get google.golang.org/grpc
- docker run -d -p 5432:5432 gogrpc
- docker build . -t gogrpc

- go build -o server-student.exe
- go build -o server-student.exe server-student/main.go

## PROXY

- https://github.com/grpc/grpc-web
- https://www.envoyproxy.io/

## MAKEFILE

- Windows installatio of mingw
- https://sourceforge.net/projects/mingw/
- mingw-get install mingw32-make
- make clean generate
