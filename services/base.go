package services

type BaseService struct{}

// 只查看自己名下的分包
func (s *BaseService) LimitPackageId(condition map[string]interface{}, packageIds []int) map[string]interface{} {
	// 不指定分包每个人只能查看自己分包下面的数据
	if _, ok := condition["package_id"]; !ok {
		condition["package_id__in"] = packageIds
	}

	return condition
}
