// Code generated by protoc-gen-bm v0.1, DO NOT EDIT.
// source: helloworld/api/api.proto

/*
Package api is a generated blademaster stub package.
This code was generated with kratos/tool/protobuf/protoc-gen-bm v0.1.

package 命名使用 {appid}.{version} 的方式, version 形如 v1, v2 ..

It is generated from these files:
	helloworld/api/api.proto
*/
package api

import (
	"context"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

// to suppressed 'imported but not used warning'
var _ *bm.Context
var _ context.Context
var _ binding.StructValidator

var PathDemoPing = "/demo.service.v1.Demo/Ping"
var PathDemoSayHello = "/demo.service.v1.Demo/SayHello"
var PathDemoSayHelloURL = "/kratos-demo/say_hello"
var PathDemoLogin = "/user/login"
var PathDemoAddUser = "/adduser"
var PathDemoUpdateUser = "/updateuser"
var PathDemoGetUser = "/getuser"
var PathDemoGetUserList = "/getuserlist"

// DemoBMServer is the server API for Demo service.
// 定义服务接口
type DemoBMServer interface {
	// google.protobuf.Empty返回空参
	Ping(ctx context.Context, req *google_protobuf1.Empty) (resp *google_protobuf1.Empty, err error)

	SayHello(ctx context.Context, req *HelloReq) (resp *google_protobuf1.Empty, err error)

	SayHelloURL(ctx context.Context, req *HelloReq) (resp *HelloResp, err error)

	// 新增登录服务接口
	Login(ctx context.Context, req *LoginReq) (resp *LoginResp, err error)

	AddUser(ctx context.Context, req *AddReq) (resp *Response, err error)

	UpdateUser(ctx context.Context, req *UpdateReq) (resp *Response, err error)

	GetUser(ctx context.Context, req *GetReq) (resp *Response, err error)

	GetUserList(ctx context.Context, req *google_protobuf1.Empty) (resp *Response, err error)
}

var DemoSvc DemoBMServer

func demoPing(c *bm.Context) {
	p := new(google_protobuf1.Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.Ping(c, p)
	c.JSON(resp, err)
}

func demoSayHello(c *bm.Context) {
	p := new(HelloReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.SayHello(c, p)
	c.JSON(resp, err)
}

func demoSayHelloURL(c *bm.Context) {
	p := new(HelloReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.SayHelloURL(c, p)
	c.JSON(resp, err)
}

func demoLogin(c *bm.Context) {
	p := new(LoginReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.Login(c, p)
	c.JSON(resp, err)
}

func demoAddUser(c *bm.Context) {
	p := new(AddReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.AddUser(c, p)
	c.JSON(resp, err)
}

func demoUpdateUser(c *bm.Context) {
	p := new(UpdateReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.UpdateUser(c, p)
	c.JSON(resp, err)
}

func demoGetUser(c *bm.Context) {
	p := new(GetReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.GetUser(c, p)
	c.JSON(resp, err)
}

func demoGetUserList(c *bm.Context) {
	p := new(google_protobuf1.Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.GetUserList(c, p)
	c.JSON(resp, err)
}

// RegisterDemoBMServer Register the blademaster route
func RegisterDemoBMServer(e *bm.Engine, server DemoBMServer) {
	DemoSvc = server
	e.GET("/demo.service.v1.Demo/Ping", demoPing)
	e.GET("/demo.service.v1.Demo/SayHello", demoSayHello)
	e.GET("/kratos-demo/say_hello", demoSayHelloURL)
	e.GET("/user/login", demoLogin)
	e.GET("/adduser", demoAddUser)
	e.GET("/updateuser", demoUpdateUser)
	e.GET("/getuser", demoGetUser)
	e.GET("/getuserlist", demoGetUserList)
}
