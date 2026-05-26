// @APIVersion 1.0.0
// @Summary beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	controllers "api/controllers/adminApi"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	if beego.BConfig.RunMode == "dev" {
		beego.SetStaticPath("/swagger", "swagger")
	}

	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/login",
			beego.NSInclude(
				&controllers.LoginController{},
			),
		),
		beego.NSNamespace("/system",
			beego.NSInclude(
				&controllers.SystemController{},
			),
		),
		beego.NSNamespace("/agent",
			beego.NSInclude(
				&controllers.AgentController{},
			),
		),
		beego.NSNamespace("/team",
			beego.NSInclude(
				&controllers.TeamController{},
			),
		),
		beego.NSNamespace("/platform",
			beego.NSInclude(
				&controllers.PlatformController{},
			),
		),
		beego.NSNamespace("/game",
			beego.NSInclude(
				&controllers.GameController{},
			),
		),
		beego.NSNamespace("/gameType",
			beego.NSInclude(
				&controllers.GameTypeController{},
			),
		),
		beego.NSNamespace("/image",
			beego.NSInclude(
				&controllers.ImageController{},
			),
		),
		beego.NSNamespace("/broadcast",
			beego.NSInclude(
				&controllers.BroadcastController{},
			),
		),
		beego.NSNamespace("/active",
			beego.NSInclude(
				&controllers.ActiveController{},
			),
		),
		beego.NSNamespace("/activeSign",
			beego.NSInclude(
				&controllers.ActiveSignController{},
			),
		),
		beego.NSNamespace("/activeInvite",
			beego.NSInclude(
				&controllers.ActiveInviteController{},
			),
		),
		beego.NSNamespace("/activeFirstRecharge",
			beego.NSInclude(
				&controllers.ActiveFirstRechargeRuleController{},
			),
		),
		beego.NSNamespace("/activeRelief",
			beego.NSInclude(
				&controllers.ActiveReliefController{},
			),
		),
		beego.NSNamespace("/activeVip",
			beego.NSInclude(
				&controllers.ActiveVipController{},
			),
		),
		beego.NSNamespace("/config",
			beego.NSInclude(
				&controllers.ConfigController{},
			),
		),
		beego.NSNamespace("/gameTag",
			beego.NSInclude(
				&controllers.GameTagController{},
			),
		),
		beego.NSNamespace("/betRate",
			beego.NSInclude(
				&controllers.BetRateController{},
			),
		),
		beego.NSNamespace("/alert",
			beego.NSInclude(
				&controllers.AlertController{},
			),
		),
		beego.NSNamespace("/banner",
			beego.NSInclude(
				&controllers.BannerController{},
			),
		),
		beego.NSNamespace("/activeLuckyWheel",
			beego.NSInclude(
				&controllers.ActiveLuckyWheelController{},
			),
		),
		beego.NSNamespace("/activeInterest",
			beego.NSInclude(
				&controllers.ActiveInterestController{},
			),
		),
		beego.NSNamespace("/log",
			beego.NSInclude(
				&controllers.LogController{},
			),
		),
		beego.NSNamespace("/payment",
			beego.NSInclude(
				&controllers.PaymentController{},
			),
		),
		beego.NSNamespace("/withdraw",
			beego.NSInclude(
				&controllers.WithdrawController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/customer",
			beego.NSInclude(
				&controllers.CustomerController{},
			),
		),
		beego.NSNamespace("/serves",
			beego.NSInclude(
				&controllers.ServesController{},
			),
		),
		beego.NSNamespace("/rechargeOrder",
			beego.NSInclude(
				&controllers.RechargeOrderController{},
			),
		),
		beego.NSNamespace("/withdrawOrder",
			beego.NSInclude(
				&controllers.WithdrawOrderController{},
			),
		),
		beego.NSNamespace("/floatLogo",
			beego.NSInclude(
				&controllers.FloatLogoController{},
			),
		),
		beego.NSNamespace("/dailyStatistics",
			beego.NSInclude(
				&controllers.DailyStatisticsController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
