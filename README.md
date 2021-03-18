# rock-fasthttp-go

## 说明
- 主要是以fasthttp为基础实现的web服务器框架 目前主要是拦截中心的功能

### 服务器配置

- 主要配置参数
```lua
    local http = rock.fasthttp.server{
        protocol = "tcp",
        listen = "0.0.0.0:9090",
        router = "resource/conf.d/fasthttp",
        handler    = "resource/conf.d/fasthttp/handlers",
        default    = "default",
        access_log = "time,server_addr,server_port,remote_addr,host,path,query,ua,referer,status,content-lenght,region_info,http_risk",
        --access_log = "off",
        access_format = "json",
        access_region = "x-real-ip",
        access = rock.kafka{},
        region = rock.ip2region{ db = "resource/share/ip2.db" },
        daemon = "on",
    }
```

- 主要功能函数
```lua
    -- 如果daemon功能开启需要利用  rock.base.notify()
    http.start()
```

- 参数说明
- listen   监听端口
- protocol 监听协议 , tcp , udp , unix
  
- router 通过主机名获取对应应用脚本 脚本路径 , 如sec.edunx.com.lua
```lua
    --框架主要是根据流量中host查找对应的app处理的应用逻辑
    local r = http.router --http 是内置全局变量
    r.not_found("*" , "not_found")
```

- handler 通过处理函数名称加载处理逻辑, 如上面的not_found.lua
```lua
    local kv = http.base.keyval
    return http.handler{
        rule = "xss", --rule名称 命中规则 , 多个用，分开
        code = 200, --返回状态码
        header = {kv("server" , "IIS/9.0") , kv("helo", "edunx")}, --返回heanders内容
        eof = "on", --是否关闭直接退出，不继续往下匹配handler, 有on , off , 默认：on          
        hook = function() end
    }

```
- default 默认没有发现的主机名处理 脚本 , 功能如sec.edunx.com.lua , 主要控制路由

- access_log 日志字段 默认是如上 字段分割 目前没有response body

- access_format 日志格式 目前只实现了JSON

- access 日志保存方式 满足transport.tunnel接口 都可以

- region ip地址查询库 功能见rock-ip2region-go

- daemon 是否后台运行 on , off 默认：off

### 主要函数功能
- http.router.GET 
- http.router.HEAD
- http.router.POST
- http.router.PUT
- http.router.PATCH
- http.router.DELETE
- http.router.CONNECT
- http.router.OPTIONS
- http.router.TRACE
- http.router.POST
- http.router.ANY 忽略发方法名
- 语法:  r.GET(path string , http.handler ... )
- 参数 path： 代表路径的 完全兼容 fasthttp.router的路径语法 如:/api/{name}/{val:*}
- 参数 handler: 就是用http.handler构造的对象 
- `注意 handler为字符时候默认搜索路径`

```lua
    --上面的函数接口形式都是一样的 用法如下
    local r = http.router
    local kv = http.base.keyval

    local say = http.response.say
    local exit = http.response.exit
    local append = http.respnose.append
    local var = http.var

    local function test()
        if var.arg_name == "hello" then
            say("hello good " , var.arg_val)
            exit(200)
        end
        exit(404)
    end
    
    local function insert()
        say("hello ")
        append(" world")
        exit(200)
    end
    
    r.GET("/" , 
        http.handler{
            code = 200 , 
            header = {kv("server" , "IIS/9.0") },
            body = "helo readme.md",
            hook = test --这里是注入方法的地方
        },
        "xss",
        "sqli",
        "not_found"
    )

```
- http.router.file
- 作用:  添加文件返回
- 语法： http.router.file(path string , root string , hook function)
- 参数： path 必须格式如/{filepath:*} , /api/{filepath:*}  必须是{filepath:*}结尾
- 参数： root 文件查找路径
- 参数： hook 是否根据不同的API路径改写函数 , 用户跟hook类似
``` lua
    local r = http.router
    local var = http.var
    local kv = http.base.keyval
    local header = http.response.header
    local function rewrite()
        if var.param_name == "guba" then    
            header({ kv("server" , "guba/1.0") } )
            var.uri = "/1" ..  var.param_filepath
        end
    end
    
    r.file("/api/{name}/{filepath:*} , "root/html" , rewirte)
```

- http.router.region
- 修改获取IP地理位置方式 
- 语法: http.router.region(addr string)
- 参数: addr 获取ip的方式 , 定义根据用户的请求获取字段获取地址信息 
```lua
    local r = http.router
    r.region("x-real-ip")
```

- http.router.access_push_off
- 作用： 取消access保存
- 语法： http.router.access_push_off()
```lua
    local r = http.router
    r.access_push_off()
```

- http.logger.json
- 作用： 定义各种不同域名的记录日志模式
- 语法： http.logger.json(format string)
- 参数： format 跟server 中的access_format一样 只是 这个作用只在 当前域名下生效 优先级大于全局
```lua
    http.logger.json("time,server_addr,server_port,remote_addr,host,path,headers,body,status,content-length,region_info,http_risk")
```

- http.var.arg_*
- 作用: 读取用户请求参数的值
- 语法： http.var.arg_name , http.var.arg_a
```lua
    -- http://sec.edunx.com/a?name=edx&a=123
    local var = http.var 
    local name = var.arg_name or "" --name=edx
    local a = var.arg_a or "" -- a=123
```

- http.var.post_*
- 作用： 读取用户POST参数
- 语法： http.var.post_value
```lua
    -- POST /api HTTP/1.1
    -- Host: sec.edunx.com
    -- Rule: sqli-deny
    --
    -- name=edunx&value=123
    
    local var = http.var
    local say = http.response.say

    say(var.post_name , " " , var.post_value)
```

- http.var.param_*
- 作用： 读取路由中的param
- 语法： http.var.param_name , http.var.param_val
```lua
 -- http://sec.edunx.com/api/admin/123456
 -- 路由 r.GET("/api/{name}/{val:*}

    local var = http.var
    local say = http.response.say

    local name = var.param_name 
    local val = var.param_val

    say(name , " " , val)
```

- http.var.http_*
- 作用：获取请求header头里面的参数
- 语法： http.var.http_rule , http.var.http_user_agent
```lua
    -- GET / HTTP/1.1
    -- Host: sec.edunx.com
    -- Rule: sqli-deny

    local var = http.var
    local rule = var.http_rule
```
- http.var.cookie_*
- 作用： 获取cookie里的某个具体字段
- 语法： http.var.cookie_session

```lua
    -- GET / HTTP/1.1
    -- Host: sec.edunx.com
    -- cookie:session=123x

    local say = http.response.say
    local var = http.var
    local v = var.cookie_session
    say(v)
```

- http.var.remote_addr 
- 作用: 获取connection的四层IP地址
- 语法: http.var.remote_addr
```lua
    local addr = http.var.remote_addr
``` 
- http.var.remote_port
- 作用: 获取connection的四层端口
- 语法: http.var.remote_port
```lua
    local port = http.var.remote_port
``` 
- http.var.server_addr
- 作用: 获取本地服务器的IP地址
- 语法: http.var.server_addr
```lua
    local addr = http.var.server_addr
``` 
- http.var.server_port
- 作用: 获取本地服务器的端口
- 语法: http.var.server_addr
```lua
    local port = http.var.server_port
``` 
- http.var.host        
- 作用: 获取用户请求的主机名
- 语法: http.var.host
```lua
    local host = http.var.host
``` 
- http.var.uri
- 作用: 获取获用户请求的URI
- 语法: http.var.uri
```lua
    -- http://sec.edunx.com/api/info
    local uri = http.var.uri -- /api/info
``` 
- http.var.args
- 作用: 获取用户请求的args字符串
- 语法: http.var.args
```lua
    -- http://sec.edunx.com/api/info?name=admin&val=123
    local args = http.var.args -- name=admin&val=123
``` 
- http.var.request_uri
- 作用: 获取完整的请求URI
- 语法: http.var.request_uri
```lua
    -- http://sec.edunx.com/api/info?name=admin&val=123
    local request = http.var.request_uri -- /api/info?name=admin&val=123
``` 

- http.var.http_time
- 作用: 获取请求时间
- 语法: http.var.http_time
```lua
    local ht = http.var.http_time --2020-01-01 01:02:03.00
``` 
- http.var.cookie_raw
- 作用: 获取cooie的子完整自字符串
- 语法: http.var.cookie_raw
```lua
    local cookie = http.var.cookie_raw
``` 

- http.var.header_raw
- 作用: 获取完整的header字符串
- 语法: http.var.header_raw
```lua
    local raw = http.var.header_raw
``` 
- http.var.content_length
- 作用: 获取获取用户请求的包大小
- 语法: http.var.content_length
```lua
    local len = http.var.content_length
``` 
- http.var.content_type
- 作用: 获取获取用户请求的content_type
- 语法: http.var.content_type
```lua
    local ct = http.var.content_type
``` 
- http.var.body
- 作用: 获取用户的请求的body请求体
- 语法: http.var.body
```lua
    local body = http.var.body
``` 
- http.var.region_cityid
- 作用: 获取用户所在城市的ID
- 语法: http.var.region_cityid
```lua
    local id = http.var.region_cityid
``` 
- http.var.region_info
- 作用: 获取获IP地址位置信息
- 语法: http.var.region_info
```lua
    local info = http.var.region_info
``` 
- http.var.ua
- 作用: 获取user_agent
- 语法: http.var.ua
```lua
    local ua = http.var.ua
``` 
- http.var.referer
- 作用: 获取referer
- 语法: http.var.referer
```lua
    local ref = http.var.referer
``` 
- http.var.status
- 作用: 获取返回状态码
- 语法: http.var.status
```lua
    local status = http.var.status
``` 
- http.var.sent 
- 作用: 获取发送数据包的大小
- 语法: http.var.sent
```lua
    local sent = http.var.sent
``` 
