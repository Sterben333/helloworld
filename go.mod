module helloworld

go 1.12

require (
	github.com/gogo/protobuf v1.3.0
	github.com/golang/protobuf v1.5.2
	github.com/google/wire v0.3.0
	go-common v1.37.8
	go.uber.org/automaxprocs v1.4.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1
	google.golang.org/grpc v1.43.0
)

replace go-common => git.bilibili.co/platform/go-common v1.37.8
