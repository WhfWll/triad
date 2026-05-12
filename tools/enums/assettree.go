package enums

const (
	AssetLogTypeOne                  int    = 1     //资产日志类型-创建资产渗透
	AssetTaskResultSubObjTypeService string = "1_1" //资产结果表subobjtype-1_1
	AssetTreeNodeTypeOne             int    = 1     //资产树节点类型-资产组
	AssetTreeNodeTypeTwo             int    = 2     //资产树节点类型-资产

	AssetStatusWait = 1 // 资产同步状态 待同步
	AssetStatusSync = 2 // 资产同步状态 同步中
	AssetStatusDone = 3 // 资产同步状态 同步完成

	AssetIsReplaceN = 1 // 是否替换漏洞，1-否
	AssetIsReplaceY = 2 // 是否替换漏洞，2-是

	DefaultAssetGroup = 1 // 默认资产组

	NormalAssetExportType  = 1 // 正常资产导出
	SpecialAssetExportType = 2 // 特殊资产导出

)

var AssetTreeEnum asset

type assetTree struct {
}

func (a *assetTree) AllAssetFilingLevelEnum() map[int]string {
	res := map[int]string{
		AssetFilingLevelZero:  "无",
		AssetFilingLevelOne:   "一级",
		AssetFilingLevelTwo:   "二级",
		AssetFilingLevelThree: "三级",
		AssetFilingLevelFour:  "四级",
		AssetFilingLevelFive:  "五级",
	}
	return res
}

func (a *assetTree) GetAssetFilingLevel(filingLevel int) string {
	return a.AllAssetFilingLevelEnum()[filingLevel]
}

func (a *assetTree) GetAssetFilingLevelEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: AssetFilingLevelZero,
		Label: a.GetAssetFilingLevel(AssetFilingLevelZero),
	}, {
		Value: AssetFilingLevelOne,
		Label: a.GetAssetFilingLevel(AssetFilingLevelOne),
	}, {
		Value: AssetFilingLevelTwo,
		Label: a.GetAssetFilingLevel(AssetFilingLevelTwo),
	}, {
		Value: AssetFilingLevelThree,
		Label: a.GetAssetFilingLevel(AssetFilingLevelThree),
	}, {
		Value: AssetFilingLevelFour,
		Label: a.GetAssetFilingLevel(AssetFilingLevelFour),
	}, {
		Value: AssetFilingLevelFive,
		Label: a.GetAssetFilingLevel(AssetFilingLevelFive),
	}}
	return result
}
