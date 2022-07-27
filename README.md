# HelloWorld

## 环境依赖
mysql5.7

redis@latest

Go1.18.4

Kratos v1.0.52

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

7. Protobuf查看

   访问：```http://localhost:8000/helloworld/say_hello?name=Soul```

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



## V1.0 版本内容更新2022.7.26

1. ```Internal/server/http/server.go```文件下增加了一个handler方法showParam

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

2. ```api/api/proto```文件下增加了自定义接口Login、AddUser、UpdateUser、GetUser、GetUserList。相应的```inernal/service/service.go```也需要添加相应的方法

   - 更新之后要生成新的pb.go文件需要执行以下命令

     ``` bash
     cd helloworld/api
     kratos tool protoc #生成新的pb.go文件
     cd helloworld/internal/di
     go generate #生成新的go接口，其实也不用自己来，启动项目的时候它会自动执行go generate
     ```

   - 打开浏览器访问：```http://localhost:8000/user/login?username=Soul&passwd=111111```

   - 输出内容：

     ```bash
     {
         "code": 0,
         "message": "0",
         "ttl": 1,
         "data": {
             "content": "login:Soul, passwd: 111111"
         }
     }
     ```

3. 修改了```configs/mysql/toml```中的内容、创建了测试数据库，数据库内容如下：

   ```sql
   create database kratos_demo;
   use kratos_demo;
   
   CREATE TABLE `users` (
     `uid` int(10) unsigned NOT NULL AUTO_INCREMENT,
     `nickname` varchar(100) NOT NULL DEFAULT '' COMMENT '昵称',
     `age` smallint(5) unsigned NOT NULL COMMENT '年龄',
     `uptime` int(10) unsigned NOT NULL DEFAULT '0',
     `addtime` int(10) unsigned NOT NULL DEFAULT '0',
     PRIMARY KEY (`uid`)
   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
   ```

   在```model/model.go```中添加了结构体```User```

   在```dao/dao.go```中新增了四个接口

   创建了文件```dao/dao.user.go```，实现了四个接口

   在```api/api.proto```中增加了http接口

   在```internal/service/service.go```，增加了接口实现

   - 完成上述内容后，重新生成pd文件

     ```bash
     cd helloworld
     kratos tool protoc
     ```

   - 打开浏览器：

     - 添加用户：```http://localhost:8000/adduser?nickname=soul&age=22```
     - 删除用户：```http://localhost:8000/deleteuser?uid=3```
     - 更新用户：```http://localhost:8000/updateuser?uid=3&nickname=soul&age=22```
     - 查询单个用户：```http://localhost:8000/getuser?uid=3```
     - 查询用户列表：```http://localhost:8000/getuserlist```

## V1.1 版本内容更新2022.7.27

1. 在增删改查上做了redis中间件的缓存处理，用来优化操作体验
2. 考虑到更新和删除动作在多线程情况下可能会读取脏数据，本人采用延时双删的策略
3. 对于查询用户列表的情况，因为这种情况操作不是很多，且数据量巨大，因此本人没有做redis的处理，直接请求数据库即可
