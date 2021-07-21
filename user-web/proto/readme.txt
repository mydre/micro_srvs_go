//核心命令
protoc -I . user.proto --go_out=plugins=grpc:.


acat@acat-xx:proto$ protoc -I . user.proto --go_out=plugins=grpc:.
protoc-gen-go: program not found or is not executable
Please specify a program using absolute path or make sure the program is available in your PATH system variable
--go_out: protoc-gen-go: Plugin failed with status code 1.

//发现安装失败，这个时候可以安装一些插件
acat@acat-xx:~$ go get google.golang.org/grpc
go: downloading github.com/golang/protobuf v1.4.3
go: downloading google.golang.org/protobuf v1.25.0
go: downloading golang.org/x/text v0.3.0
acat@acat-xx:~$ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go: downloading google.golang.org/protobuf v1.26.0
acat@acat-xx:~$