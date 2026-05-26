package command

import (
	"api/services"
)

type InitDataCommand struct{}

func init() {
	// logs.Info("脚本运行")
	// cmd := &InitDataCommand{}
	// cmd.initConfig()
	// cmd.initActivity()
	// cmd.UpdateGameTypeName()
}

// 初始化配置列表
func (c *InitDataCommand) initConfig() {
	configService := new(services.ConfigService)
	err := configService.InitPackageConfig(0)
	if err != nil {
		panic(err)
	}
}

// 初始化活动，规则，奖励
func (c *InitDataCommand) initActivity() {
	err := initActivity()
	if err != nil {
		panic(err)
	}
}

// 更新游戏类型名称
// func (c *InitDataCommand) UpdateGameTypeName() {
// 	logs.Info("修复游戏类型名称")
// 	var gameList []models.Game
// 	model := models.CreateGameModel()
// 	model.QueryTable(new(models.Game)).All(&gameList)

// 	for _, game := range gameList {
// 		var gameTypeInfo models.GameType
// 		model.QueryTable(new(models.GameType)).Filter("id", game.GameTypeId).One(&gameTypeInfo)
// 		_, err := model.QueryTable(new(models.Game)).Filter("id", game.Id).Update(orm.Params{
// 			"game_type_name": gameTypeInfo.Name,
// 		})
// 		if err != nil {
// 			logs.Error("更新游戏类型名称失败:%v", err.Error())
// 		}
// 	}
// }
