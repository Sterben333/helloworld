package service

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/library/conf/paladin.v2"
	pb "helloworld/api"
	"helloworld/internal/dao"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.DemoServer), new(*Service)))

// Service .
type Service struct {
	ac  *paladin.Map //paladin.Map 通过atomic.Value支持自动热加载
	dao dao.Dao
}

// New 新建 a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

// SayHello grpc demo func.
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)
	fmt.Printf("我是hello %s", req.Name) // 控制台打印
	return
}

// SayHelloURL bm demo func.这里是后台控制台打印的内容
func (s *Service) SayHelloURL(ctx context.Context, req *pb.HelloReq) (reply *pb.HelloResp, err error) {
	reply = &pb.HelloResp{
		Content: "你好啊小可爱 " + req.Name,
	}
	fmt.Printf("我是hello-url %s", req.Name)
	return
}

// Login 新增登录接口
func (s *Service) Login(ctx context.Context, req *pb.LoginReq) (reply *pb.LoginResp, err error) {
	reply = &pb.LoginResp{
		Content: "login:" + req.Username + ", passwd: " + req.Passwd,
	}
	fmt.Printf("login url %s", req.Username)
	return
}

// AddUser 添加用户
func (s *Service) AddUser(ctx context.Context, req *pb.AddReq) (reply *pb.Response, err error) {
	fmt.Printf("AddUser: %s, %d\n", req.Nickname, req.Age)
	user, err := s.dao.AddUser(ctx, req.Nickname, req.Age)
	if err != nil {
		fmt.Printf("AddUser %s, %d Error", req.Nickname, req.Age)
		return
	}
	res, _ := json.Marshal(user)
	reply = &pb.Response{
		Content: string(res),
	}
	return
}

// DeleteUser 更新用户信息
func (s *Service) DeleteUser(ctx context.Context, req *pb.DeleteReq) (reply *pb.Response, err error) {
	fmt.Printf("DeleteUser:  %d\n", req.Uid)
	rows, err := s.dao.DeleteUser(ctx, req.Uid)
	if err != nil {
		fmt.Printf("DeleteUser Uid = %d Error", req.Uid)
		return
	}
	reply = &pb.Response{
		Content: fmt.Sprintf("删除行数: %d", rows),
	}
	return
}

// UpdateUser 更新用户信息
func (s *Service) UpdateUser(ctx context.Context, req *pb.UpdateReq) (reply *pb.Response, err error) {
	fmt.Printf("UpdateUser: %s, %d\n", req.Nickname, req.Age)
	rows, err := s.dao.UpdateUser(ctx, req.Uid, req.Nickname, req.Age)
	if err != nil {
		fmt.Printf("UpdateUser %s, %d Error", req.Nickname, req.Age)
		return
	}
	reply = &pb.Response{
		Content: fmt.Sprintf("更新行数: %d", rows),
	}
	return
}

//GetUser 获取用户信息
func (s *Service) GetUser(ctx context.Context, req *pb.GetReq) (reply *pb.Response, err error) {
	fmt.Printf("GetUser: %d\n", req.Uid)
	user, err := s.dao.GetUser(ctx, req.Uid)
	if err != nil {
		fmt.Printf("GetUser %s Error", req.Uid)
		return
	}
	res, _ := json.Marshal(user)
	reply = &pb.Response{
		Content: string(res),
	}
	return
}

//GetUserList 获取用户列表
func (s *Service) GetUserList(ctx context.Context, req *empty.Empty) (reply *pb.Response, err error) {
	fmt.Printf("GetUserList")
	userlist, err := s.dao.GetUserList(ctx)
	if err != nil {
		fmt.Printf("GetUserList Error")
		return
	}
	res, _ := json.Marshal(userlist)
	reply = &pb.Response{
		Content: string(res),
	}
	return
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}
