package services

import (
	"api/models"
	"strings"
)

type Authservice struct{}

// 检测菜单权限
func (s *Authservice) Check(id int, groupId int) bool {
	// 超级管理员
	if groupId == 1 {
		return true
	}
	model := models.CreateAuthGroupModel()
	// 获取用户所有的权限
	var authGroupInfo models.AuthGroup
	authGroupInfoErr := model.QueryTable(new(models.AuthGroup)).Filter("id", groupId).One(&authGroupInfo)
	if authGroupInfoErr != nil {
		return false
	}
	ruleIds := strings.Split(authGroupInfo.Rules, ",")
	var ruleInfos []models.AuthRule
	_, ruleInfoErr := model.QueryTable(new(models.AuthRule)).Filter("id__in", ruleIds).All(&ruleInfos)
	if ruleInfoErr != nil {
		return false
	}

	for _, ruleInfo := range ruleInfos {
		if ruleInfo.Id == id {
			return true
		}
	}

	return false
}
