# bbs-back

论坛后端

## 框架

[Beego](https://beego.vip)

## 管理依赖

```bash
go mod tidy
```

## 配置文件

`conf/app.conf`

## 启动

```bash
# 项目根目录
go build
bbs-back.exe
# 或直接由beego管理
bee run 
```

## 自动化生成API文档

```bash
# 修改conf/app.conf
EnableDocs = true
# 项目根目录执行
# -gendoc=true 表示每次自动化的 build 文档
# -downdoc=true 就会自动的下载 swagger 文档查看器
bee run -gendoc=true -downdoc=true
# 访问http://localhost:8081/swagger/
```

