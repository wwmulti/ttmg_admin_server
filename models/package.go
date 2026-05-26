package models

type PackageModel struct {
	*Base
}

type Package struct {
	Id           int    `json:"id" orm:"auto;column(id)"`                             // 主键
	Title        string `json:"title" orm:"column(title);size(100)"`                  // 平台名称
	Domain       string `json:"domain" orm:"column(domain);size(200)"`                // 域名地址
	ApiDomain    string `json:"api_domain" orm:"column(api_domain);size(200)"`        // api域名地址
	Icon         string `json:"icon" orm:"column(icon);type(text);null"`              // icon图片地址
	Logo         string `json:"logo" orm:"column(logo);type(text);null"`              // logo图片地址
	LoadPic      string `json:"load_pic" orm:"column(load_pic);type(text);null"`      // 平台加载图片
	Pg           string `json:"pg" orm:"column(pg);size(100)"`                        // pg商户 结构1,2
	Pp           string `json:"pp" orm:"column(pp);size(100)"`                        // pp商户 结构1,2
	OpenRegister int    `json:"open_register" orm:"column(open_register);default(0)"` // 打开注册 0关闭 1打开
	CreateTime   int64  `json:"create_time" orm:"column(create_time);bigint"`         // 创建时间
	UpdateTime   int64  `json:"update_time" orm:"column(update_time);bigint"`         // 更新时间
	IsDeleted    int    `json:"-" orm:"column(is_deleted);int"`                       // 是否删除 0未删 1删除
	GroupId      int    `json:"group_id" orm:"column(group_id);int"`                  // 分组id 0未分配
}

func CreatePackageModel() *PackageModel {
	return &PackageModel{CreateBase()}
}
