# HelloWorld

## 环境依赖
mysql

## 部署步骤

1. Environment
    
    ```LOG_NO_AGENT=true;DEPLOY_ENV=uat```
    
1. Program arguments
    
    ```-conf=./configs```

3. 安装kratos

   ```sudo /bin/bash -c "$(curl -x 172.22.21.10:80 -fsSL http://kratos.bilibili.co/x/kratos/install.sh)"```

4. 新建一个kratos项目

   ```kratos new helloworld```

5. 启动项目

   ```cd helloworld/cmd```

   ```go build```

   ```./cmd -conf ../configs```

6. 查看已注册路由信息

   访问：```http://localhost:8000/metadata```

## 目录结构描述

```
.
├── CHANGELOG.md           # CHANGELOG
├── CONTRIBUTORS.md        # CONTRIBUTORS
├── README.md              # README
├── api                    # api目录为对外保留的proto文件及生成的pb.go文件，注：需要"--proto"参数
│   ├── api.proto
│   ├── api.pb.go          # 通过go generate生成的pb.go文件
│   └── generate.go
├── cmd                    # cmd目录为main所在
│   └── main.go            # main.go
├── configs                # configs为配置文件目录
│   ├── application.toml   # 应用的自定义配置文件，可能是一些业务开关如：useABtest = true
│   ├── grpc.toml          # grpc相关配置 
│   ├── http.toml          # http相关配置
│   ├── log.toml           # log相关配置
│   ├── memcache.toml      # memcache相关配置
│   ├── mysql.toml         # mysql相关配置
│   └── redis.toml         # redis相关配置
├── go.mod                 # go.mod
└── internal               # internal为项目内部包，包括以下目录：
    ├── dao                # dao层，用于数据库、cache、MQ、依赖某业务grpc|http等资源访问
    │   └── dao.go
    ├── model              # model层，用于声明业务结构体
    │   └── model.go
    ├── server             # server层，用于初始化grpc和http server
    │   └── http           # http层，用于初始化http server和声明handler
    │       └── http.go
    │   └── grpc           # grpc层，用于初始化grpc server和定义method
    │       └── grpc.go
    └── service            # service层，用于业务逻辑处理，且为方便http和grpc共用方法，建议入参和出参保持grpc风格，且使用pb文件生成代码
        └── service.go
```



## V1.0 版本内容更新

1. ```Internal/server/http/server.go```文件夹下增加了一个handler方法showParam

   打开浏览器访问：```http://localhost:8000/httpdemo/param2/Soul/male/hello```

   输出内容：

   ```xml
   {
       "action": "",
       "code": 0,
       "gender": "male",
       "message": "0",
       "name": "Soul",
       "path": "/httpdemo/param2/:name/:gender/:say",
       "say": "hello"
   }
   ```

   
