package middleware

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func CORS() web.FilterFunc {
	return func(ctx *context.Context) {
		// 允许的域名，生产环境可以设置为具体域名
		ctx.Output.Header("Access-Control-Allow-Origin", "*")

		// 允许的请求方法
		ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 允许的请求头
		ctx.Output.Header("Access-Control-Allow-Headers", "*")

		// 预检请求的缓存时间（秒）
		ctx.Output.Header("Access-Control-Max-Age", "86400")

		// 处理 OPTIONS 预检请求
		if ctx.Input.Method() == "OPTIONS" {
			ctx.Output.Header("Access-Control-Allow-Origin", "*")

			// 允许的请求方法
			ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			// 允许的请求头
			ctx.Output.Header("Access-Control-Allow-Headers", "*")

			// 预检请求的缓存时间（秒）
			ctx.Output.Header("Access-Control-Max-Age", "86400")
			ctx.Output.SetStatus(200)
			ctx.Output.Body([]byte(""))
			return
		}
	}
}
