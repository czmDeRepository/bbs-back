# bbs-back

[论坛后端](https://github.com/czmDeRepository/bbs-back)

[论坛前端](https://github.com/czmDeRepository/bbs-front)

> 注：启动前请确保本地已安装go，mysql和redis

## 初始化数据库表

1. 将`bbs.sql`导入mysql数据库中

> 配置文件默认数据库名bbs

## 管理项目依赖

```bash
go mod tidy
```

## 配置文件

`conf/app.conf`

```conf
### 数据库配置
dbhost = localhost
dbport = 3306
dbuser = root
dbpwd = 123456
dbname = bbs
# 时区
loc = Asia%2FShanghai
# 最大空闲连接
maxIdle = 5
# 最大连接
maxConn = 30

# token加密密钥
secretKey = czmDeBBS

# 文件最大上传限制 1024 * 1024 * 10   10M
maxFileLimit = 10485760

# 版本
version = v1

# Redis配置
RedisConn = 127.0.0.1:6379
redisMaxIdle = 1
redisMaxActive = 2

# 系统邮箱
# SMTP服务器
email.host = smtp.126.com
email.port = 25
# 邮箱账号
email.username = 
# 邮箱授权码，如qq邮箱可查看https://service.mail.qq.com/cgi-bin/help?subtype=1&&no=1001256&&id=28
email.password = 
```

## 启动

```bash
# 项目根目录
go build
bbs-back.exe
# 或直接由beego管理
bee run 
```

## 自动化生成API文档（swagger）

> 非必需

```bash
# 修改conf/app.conf
EnableDocs = true
# 项目根目录执行
# -gendoc=true 表示每次自动化的 build 文档
# -downdoc=true 就会自动的下载 swagger 文档查看器
bee run -gendoc=true -downdoc=true
# 访问http://localhost:8081/swagger/
```

## 默认账号

超级管理员 

> 账号：admin
>
> 密码：123456

## 效果展示

### 登录页面

![image-20220702174122302](https://img2022.cnblogs.com/blog/1961904/202207/1961904-20220702180859913-2069702432.png)

### 编辑页面

![image-20220702180947058](https://img2022.cnblogs.com/blog/1961904/202207/1961904-20220702180949764-986765398.png)

### 浏览页面

![image-20220702181008714](https://img2022.cnblogs.com/blog/1961904/202207/1961904-20220702181011820-310728876.png)

### 群聊页面

![image-20220702200620211](https://img2022.cnblogs.com/blog/1961904/202207/1961904-20220702200623166-1312581414.png)

### 个人主页

![image-20220702181032429](https://img2022.cnblogs.com/blog/1961904/202207/1961904-20220702181035420-1646718355.png)

### 管理页面

![image-20220702181054365](https://img2022.cnblogs.com/blog/1961904/202207/1961904-20220702181057093-1706869674.png)

### 监控大盘

![image-20220702181132709](https://img2022.cnblogs.com/blog/1961904/202207/1961904-20220702181135672-728974879.png)

## END