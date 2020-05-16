### build outside of GOPATH
go mod init github.com/zrthstr/gozones
go build
go mod tidy


Note: By default for your native platform cgo will likely be enabled for the compiler and thus link against the system's domain resovler instead of using the pure Go one. If you would like to stop that use:

export CGO_ENABLED=0
