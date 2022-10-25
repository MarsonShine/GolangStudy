go env -w GOOS=linux
go build -buildmode=plugin -o plugin1.so .\plugin1.go
go env -w GOOS=windows