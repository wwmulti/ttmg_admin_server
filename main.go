package main

import (
	_ "api/command" // 修复脚本
	"api/middleware"
	_ "api/routers"
	"api/task"
	"fmt"
	"strings"

	"github.com/beego/i18n"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

var RedisCache *redis.Client

func init() {
	setSuportLanguageFile()
	// 获取redis client 对象
	// RedisCache = utils.GetRedisClient()

	if beego.BConfig.RunMode == "dev" {
		beego.InsertFilter("*", beego.BeforeRouter, middleware.CORS())
		logs.SetLogger("console", "")
		orm.Debug = true
	}

	beego.SetStaticPath("/public", "public")
}

func setSuportLanguageFile() {
	languagesStr, err := beego.AppConfig.String("languages")
	if err != nil {
		panic(err)
	}

	languages := strings.Split(languagesStr, ",")
	for _, language := range languages {
		err := i18n.SetMessage(language, fmt.Sprintf("conf/locale/%v.conf", language))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// userModel := models.CreateUserModel();
	// id, err := userModel.Insert(&models.User{ Username: "321" });
	// if nil != err {
	// 	panic("===========>>>" + err.Error());
	// }
	// fmt.Println("===========>>>", id);
	// var tbUser []*models.User;
	// if count, err := userModel.QueryTable(&models.User{}).Filter("id__gt", 1).All(&tbUser); nil == err {
	// if count, err := userModel.QueryTable(&models.User{}).SetCond(orm.NewCondition().And("id__gt", 1)).All(&tbUser); nil == err {
	// 	fmt.Println("===========>>>", count, tbUser);
	// }
	// user := &models.User{ Uid: 1 };
	// if err := userDb.Read(user); nil == err {
	// 	fmt.Println("===========>>>", user);
	// }
	// err := RedisCache.Set(context.Background(), "my_key", "Hello, Redis!", time.Second*60) // 缓存 60 秒
	// if err != nil {
	// 	fmt.Println("Error setting cache:", err)
	// }

	// 获取缓存值
	// value := RedisCache.Get(context.Background(), "my_key")
	// if value == nil {
	// 	fmt.Println("Error getting cache")
	// } else {
	// 	fmt.Println("Cache Value:", value.Val()) // 将 []byte 转为 string
	// }

	// 创建任务管理器
	taskMgr := task.TaskJobManager()
	// 启动所有任务
	taskMgr.Start()
	beego.Run()
}
