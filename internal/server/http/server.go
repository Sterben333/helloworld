package http

import (
	"go-common/library/net/http/blademaster/render"
	"net/http"

	"go-common/library/conf/paladin.v2"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	pb "helloworld/api"
	"helloworld/internal/model"
	"helloworld/internal/service"
)

var svc *service.Service

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine, err error) {
	var (
		cfg bm.ServerConfig
		ct  paladin.TOML
	)
	//日志处理，监听处理
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	svc = s
	// DefaultServer returns an Engine instance with the Recovery, Logger and CSRF middleware already attached.
	engine = bm.DefaultServer(&cfg) //创建引擎
	pb.RegisterDemoBMServer(engine, s)
	initRouter(engine)

	err = engine.Start()
	return
}

//注册路由
//Router用于根据请求的路径分发请求
func initRouter(e *bm.Engine) {
	e.Ping(ping)                // engine自带的"/ping"借口，用于负载均衡检测服务将康状态
	g := e.Group("/helloworld") // e.Grop创建一组"/helloworld"起始的路由组
	{
		g.GET("/start", howToStart) // g.GET 创建一个 "helloworld/start" 的路由，使用GET方式请求，默认处理Handler为howToStart方法

		// NOTE: 可以拿到一个key为name的参数。注意只能匹配到/param1/soul，无法匹配/param1/soul/hao(该路径会404)
		g.GET("/param1/:name", showParam)
		// NOTE: 可以拿到多个key参数。注意只能匹配到/param2/soul/male/hello，无法匹配/param2/soul或/param2/soul/hello
		g.GET("/param2/:name/:gender/:say", showParam)
		// NOTE: 可以拿到一个key为name的参数 和 一个key为action的路径。
		// NOTE: 如/params3/soul/hello，action的值为"/hello"
		// NOTE: 如/params3/soul/hello/hi，action的值为"/hello/hi"
		// NOTE: 如/params3/soul/hello/hi/，action的值为"/hello/hi/"
		g.GET("/param3/:name/*action", showParam)
	}
}

func showParam(c *bm.Context) {
	name, _ := c.Params.Get("name")
	gender, _ := c.Params.Get("gender")
	say, _ := c.Params.Get("say")
	action, _ := c.Params.Get("action")
	path := c.RoutePath // NOTE: 获取注册的路由原始地址，如: /httpdemo/param1/:name
	c.Render(http.StatusOK, render.JSON{
		Code:    2000,
		Message: "请求成功",
		Data: map[string]interface{}{
			"name":   name,
			"gender": gender,
			"say":    say,
			"action": action,
			"path":   path,
		},
	})
}

//engine自带Ping方法，用于设置"/ping"路由的handler，该路由统一提供于负载均衡服务做健康检测。服务是否健康，
//可自定义 ping handler 进行逻辑判断，如检测DB是否正常等。
func ping(ctx *bm.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// example for http request handler.主要改这个部分，启动起来这个部分就是显示在前端的内容
// bm的handler方法.
func howToStart(c *bm.Context) {
	k := &model.Kratos{
		Id:    1,
		Title: "我是梁志超他奶奶",
	}
	c.Render(http.StatusOK, render.JSON{
		Code:    2000,
		Message: "请求成功",
		Data:    k,
	})
}
