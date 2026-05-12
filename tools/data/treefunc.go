package data

var Func funcStruct

type funcStruct struct {
}

type TreeList struct {
	Id       int         `json:"id"`
	PId      int         `json:"pid"` //上级Id
	Value    string      `json:"value"`
	Label    string      `json:"label"`
	Children []*TreeList `json:"children"` //这里必须是指针
}

func (funcStruct funcStruct) Tree(resources []TreeList) []*TreeList {
	//定义根树，既id=0的根节点，我用的时候并不存在于数据库
	var rootResouce = TreeList{}
	//创建一个map，把父级相同的地址归纳起来
	DataMap := make(map[int][]*TreeList, len(resources))
	// 记录所有存在的ID，用于判断孤儿节点
	ExistMap := make(map[int]bool, len(resources))

	//寻找对应的父级，添加子节点集合
	for key, _ := range resources {
		pid := resources[key].PId
		DataMap[pid] = append(DataMap[pid], &resources[key])
		ExistMap[resources[key].Id] = true
	}
	for key, _ := range resources {
		// 如果是根节点，或者父节点不存在（孤儿节点），都作为根节点展示
		if resources[key].PId == 0 || !ExistMap[resources[key].PId] {
			rootResouce.Children = append(rootResouce.Children, &resources[key])
		}
		resources[key].Children = DataMap[resources[key].Id]
	}
	// 添加完成，既建立树形关系完成
	return rootResouce.Children
}
