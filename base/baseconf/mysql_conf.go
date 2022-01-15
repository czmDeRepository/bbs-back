package baseconf

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

	dbUser, _ := beego.AppConfig.String("dbuser")
	dbPwd, _ := beego.AppConfig.String("dbpwd")
	dbHost, _ := beego.AppConfig.String("dbhost")
	dbPort, _ := beego.AppConfig.String("dbport")
	dbName, _ := beego.AppConfig.String("dbname")
	loc, _ := beego.AppConfig.String("loc")
	maxIdle, _ := beego.AppConfig.Int("maxIdle")
	maxConn, _ := beego.AppConfig.Int("maxConn")
	// 参数1        数据库的别名，用来在 ORM 中切换数据库使用
	// 参数2        driverName
	// 参数3        对应的链接字符串
	// 参数4(可选)  设置最大空闲连接
	// 参数5(可选)  设置最大数据库连接 (go >= 1.2)
	dbUrl := dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?loc=" + loc
	orm.RegisterDataBase("default", "mysql", dbUrl,
		orm.MaxIdleConnections(maxIdle), orm.MaxOpenConnections(maxConn))
}
