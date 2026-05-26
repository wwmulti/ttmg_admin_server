package middleware

import (
	"api/utils"
	"log"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// JSONResponse JSON响应结构
type JSONResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func CheckToken(ctx *context.Context) {
	// 1. 获取Token
	authHeader := ctx.Input.Header("Authorization")

	if authHeader == "" {
		if err := ctx.Output.JSON(JSONResponse{
			Code: 401,
			Msg:  "Missing token",
		}, false, false); err != nil {
			// 处理错误，例如记录日志
			log.Printf("写入响应体失败: %v", err)
			ctx.Output.SetStatus(500)
		}
		return
	}

	// 2. 解析Bearer Token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		if err := ctx.Output.JSON(JSONResponse{
			Code: 401,
			Msg:  "Invalid token format",
		}, false, false); err != nil {
			// 处理错误，例如记录日志
			log.Printf("写入响应体失败: %v", err)
			ctx.Output.SetStatus(500)
		}
		return
	}

	// 3. 验证Token
	tokenString := parts[1]
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		if err := ctx.Output.JSON(JSONResponse{
			Code: 401,
			Msg:  "Please login again",
		}, false, false); err != nil {
			// 处理错误，例如记录日志
			log.Printf("写入响应体失败: %v", err)
			ctx.Output.SetStatus(500)
		}
		return
	}
	// 4. 将用户信息存储到Context中
	ctx.Input.SetData("userId", claims.UserId)
	ctx.Input.SetData("userName", claims.Username)
}
