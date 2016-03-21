//go:generate protoc -I. lh/lh.proto --go_out=plugins=grpc:pb
//go:generate protoc -I. driver.proto runner.proto --go_out=plugins=grpc:pb
package driver
