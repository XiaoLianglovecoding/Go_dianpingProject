@echo off
set GOPATH=%cd%\.cache\gopath
set GOMODCACHE=%cd%\.cache\gomod
set GOCACHE=%cd%\.cache\gobuild
echo Environment configured! Starting server...
go run ./cmd/server