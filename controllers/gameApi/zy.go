package gameApi

import (
	"api/config"
	"api/models"
	"api/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ZY_GAME_DOMAIN         string
	ZY_GAME_OPERATOR_TOKEN string
	ZY_GAME_SECRET         string
	ZY_GAME_CURRENCY_CODE  string
)

type ZyGame struct {
	CreateUserUrl string
	GameListUrl   string
	GameUrl       string
	SetRtpUrl     string // 设置玩家rtp
	UserInfoUrl   string // 获取玩家信息
}

func init() {
	model := models.CreateMerchantModel()
	var info models.Merchant
	err := model.QueryTable(new(models.Merchant)).Filter("status", 1).Filter("type", config.PgGameType).Filter("is_deleted", 0).One(&info)
	if err != nil {
		panic("自研商户未配置")
	}

	ZY_GAME_DOMAIN = info.Domain
	ZY_GAME_OPERATOR_TOKEN = info.Token
	ZY_GAME_SECRET = info.Secret
	ZY_GAME_CURRENCY_CODE = info.Currency
}

func (g *ZyGame) GetGameConfig() ZyGame {
	return ZyGame{
		CreateUserUrl: ZY_GAME_DOMAIN + "/api/web/user_session/",
		GameListUrl:   ZY_GAME_DOMAIN + "/api/web/game_list/",
		GameUrl:       ZY_GAME_DOMAIN + "/api/web/game_url/",
		SetRtpUrl:     ZY_GAME_DOMAIN + "/api/web/set_user_rtp/",
		UserInfoUrl:   ZY_GAME_DOMAIN + "/api/web/user_info/",
	}
}

// 创建用户
type ZyCreateUserParams struct {
	UserId   int         `json:"user_id"`
	UserName string      `json:"user_name"`
	IsMock   int         // 是否模拟
	TeamInfo models.Team // 团队配置
}

type ZyCreateUserData struct {
	PlayerId int    `json:"player_id"`
	Balance  int    `json:"balance"`
	IsNew    bool   `json:"is_new"`
	Token    string `json:"token"`
	UserId   string `json:"user_id"`
	IsMock   int    `json:"is_mock"` // 是否模拟
}

var ZyCreateUserResponse struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data ZyCreateUserData `json:"data"`
}

func (g *ZyGame) CreateUser(params ZyCreateUserParams) (string, string) {
	gameConfig := g.GetGameConfig()

	paramMap := make(map[string]interface{})
	paramMap["user_id"] = strconv.Itoa(params.UserId)
	paramMap["user_name"] = params.UserName
	paramMap["ts"] = time.Now().Unix()

	paramMap["operator_token"] = ZY_GAME_OPERATOR_TOKEN
	secret := ZY_GAME_SECRET
	if params.IsMock == 1 {
		paramMap["operator_token"] = params.TeamInfo.Token
		secret = params.TeamInfo.Secret
	}
	paramMap["sign"] = g.generateSign(paramMap, secret)
	resp := utils.PostJSON(gameConfig.CreateUserUrl, paramMap)
	// 获取并解析JSON响应

	if resp.OK() {
		resp.JSON(&ZyCreateUserResponse)
		if ZyCreateUserResponse.Code == 0 {
			return ZyCreateUserResponse.Data.UserId, ZyCreateUserResponse.Data.Token
		}
		return "", ZyCreateUserResponse.Msg
	}
	return "", ""
}

// 获取游戏列表 gameType支持pg,jdb,kess,zy,wg

type ZyGameInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var ZyGameListResponse struct {
	Code int          `json:"code"`
	Data []ZyGameInfo `json:"data"`
}

func (g *ZyGame) GetGameList(gameType string) []ZyGameInfo {
	gameConfig := g.GetGameConfig()
	paramMap := make(map[string]interface{})
	paramMap["operator_token"] = ZY_GAME_OPERATOR_TOKEN
	paramMap["type"] = gameType
	paramMap["ts"] = time.Now().Unix()
	paramMap["sign"] = g.generateSign(paramMap, ZY_GAME_SECRET)

	resp := utils.PostJSON(gameConfig.GameListUrl, paramMap)

	if resp.OK() {
		resp.JSON(&ZyGameListResponse)
		if ZyGameListResponse.Code == 0 {
			return ZyGameListResponse.Data
		}
	}
	return []ZyGameInfo{}
}

// 获取游戏链接
type ZyGameLinkParams struct {
	GameCode string      `json:"game_code"`
	Language string      `json:"language"`
	Type     string      `json:"type"`
	UserId   int         `json:"user_id"`
	UserName string      `json:"user_name"`
	Amount   int         `json:"amount"`
	IsMock   int         // 是否模拟
	TeamInfo models.Team // 团队配置
}

type ZyGameLinkResponseData struct {
	Url string `json:"url"`
}

var ZyGameLinkResponse struct {
	Code int                    `json:"code"`
	Data ZyGameLinkResponseData `json:"data"`
}

func (g *ZyGame) GetGameLink(params ZyGameLinkParams) (string, error) {
	if params.IsMock == 1 {
		rtpParams := ZySetRtpParams{
			UserId:   params.UserId,
			Rtp:      params.TeamInfo.Rtp,
			IsMock:   params.IsMock,
			TeamInfo: params.TeamInfo,
		}
		status, _ := g.SetRtp(rtpParams)
		if !status {
			return "", fmt.Errorf("SheZhiGaiLvShiBai")
		}
	}
	userParams := ZyCreateUserParams{
		UserId:   params.UserId,
		UserName: params.UserName,
		IsMock:   params.IsMock,
		TeamInfo: params.TeamInfo,
	}
	uid, token := g.CreateUser(userParams)
	if token == "" {
		return "", fmt.Errorf("ChuangJianYouXiYongHuShiBai")
	}

	gameConfig := g.GetGameConfig()
	paramMap := make(map[string]interface{})

	paramMap["operator_token"] = ZY_GAME_OPERATOR_TOKEN
	secret := ZY_GAME_SECRET
	if params.IsMock == 1 {
		paramMap["operator_token"] = params.TeamInfo.Token
		secret = params.TeamInfo.Secret
	}
	paramMap["user_id"] = uid
	paramMap["user_token"] = token
	paramMap["game_code"] = params.GameCode
	paramMap["language"] = params.Language
	paramMap["type"] = params.Type
	paramMap["ts"] = time.Now().Unix()
	paramMap["sign"] = g.generateSign(paramMap, secret)

	resp := utils.PostJSON(gameConfig.GameUrl, paramMap)
	if resp.OK() {
		resp.JSON(&ZyGameLinkResponse)
		if ZyGameLinkResponse.Code == 0 {
			return ZyGameLinkResponse.Data.Url, nil
		}
	}
	return "", fmt.Errorf("HuoQuYouXiLianJieShiBai")
}

// 设置玩家rtp
type ZyRtpResponseData struct {
	Rtp int `json:"rtp"`
}

var ZyRtpResponse struct {
	Code int               `json:"code"`
	Data ZyRtpResponseData `json:"data"`
	Msg  string            `json:"msg"`
}

type ZySetRtpParams struct {
	UserId   int
	Rtp      int
	IsMock   int
	TeamInfo models.Team
}

// 0 20 40 50 60 65 70 75 80 85 87 89 90 91 93 95 96 97 98 99 100 105 120 150
func (g *ZyGame) SetRtp(params ZySetRtpParams) (bool, int) {
	gameConfig := g.GetGameConfig()
	paramMap := make(map[string]interface{})
	paramMap["operator_token"] = ZY_GAME_OPERATOR_TOKEN
	secret := ZY_GAME_SECRET
	if params.IsMock == 1 {
		paramMap["operator_token"] = params.TeamInfo.Token
		secret = params.TeamInfo.Secret
	}

	paramMap["user_id"] = strconv.Itoa(params.UserId)
	paramMap["rtp"] = params.Rtp
	paramMap["ts"] = time.Now().Unix()
	paramMap["sign"] = g.generateSign(paramMap, secret)

	resp := utils.PostJSON(gameConfig.SetRtpUrl, paramMap)
	if resp.OK() {
		resp.JSON(&ZyRtpResponse)
		if ZyRtpResponse.Code == 0 {
			return true, ZyRtpResponse.Data.Rtp
		}
	}
	return false, 0
}

// 获取用户信息
type ZyUserInfoData struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	RemoteUid string `json:"remote_uid"`
	Rtp       int    `json:"rtp"`
}

var ZyUserInfo struct {
	Code    int            `json:"code"`
	Data    ZyUserInfoData `json:"data"`
	Message string         `json:"message"`
}

type ZyGetUserInfoParams struct {
	UserId   int
	Username string
	IsMock   int
	TeamInfo models.Team
}

func (g *ZyGame) GetUserInfo(params ZyGetUserInfoParams) ZyUserInfoData {
	userParams := ZyCreateUserParams{
		UserId:   params.UserId,
		UserName: params.Username,
		IsMock:   params.IsMock,
		TeamInfo: params.TeamInfo,
	}
	uid, token := g.CreateUser(userParams)
	if token == "" {
		return ZyUserInfoData{}
	}

	gameConfig := g.GetGameConfig()
	paramMap := make(map[string]interface{})

	paramMap["operator_token"] = ZY_GAME_OPERATOR_TOKEN
	secret := ZY_GAME_SECRET
	if params.IsMock == 1 {
		paramMap["operator_token"] = params.TeamInfo.Token
		secret = params.TeamInfo.Secret
	}

	paramMap["user_id"] = uid
	paramMap["ts"] = time.Now().Unix()
	paramMap["sign"] = g.generateSign(paramMap, secret)
	resp := utils.PostJSON(gameConfig.UserInfoUrl, paramMap)
	if resp.OK() {
		resp.JSON(&ZyUserInfo)
		if ZyUserInfo.Code == 0 {
			return ZyUserInfo.Data
		}
	}
	return ZyUserInfoData{}
}

// 生成签名
func (g *ZyGame) generateSign(paramMap map[string]interface{}, secret string) string {
	// 1. 复制一份map，避免修改原map
	tempMap := make(map[string]string)
	for k, v := range paramMap {
		// 排除 sign 字段
		if k != "sign" {
			tempMap[k] = fmt.Sprintf("%v", v)
		}
	}

	// 2. 获取所有key并排序
	keys := make([]string, 0, len(tempMap))
	for k := range tempMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 3. 构建 key=value 字符串，用 & 连接
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, tempMap[k]))
	}

	// 4. 用 & 连接所有参数
	signStr := strings.Join(parts, "&")

	// 5. 最后添加 &key=secret
	signStr = signStr + "&key=" + secret

	// 6. 计算小写 MD5
	hash := md5.Sum([]byte(signStr))

	// // 调试输出（可选）
	// fmt.Printf("参与签名的key: %v\n", keys)
	// fmt.Printf("待签名字符串: %s\n", signStr)
	// fmt.Printf("生成的签名: %s\n", hex.EncodeToString(hash[:]))

	return hex.EncodeToString(hash[:])
}
