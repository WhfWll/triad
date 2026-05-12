package data

//int数组去重
func ArrayIntUnique(list []int) []int {
	set := make(map[int]struct{}, len(list))
	j := 0
	for _, v := range list {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		list[j] = v
		j++
	}
	return list[:j]
}

//删除指定元素int[]
func RemoveInt(nums []int, v int) []int {
	for i := 0; i < len(nums); {
		if nums[i] == v {
			nums = append(nums[:i], nums[i+1:]...)
		} else {
			i++
		}

	}
	return nums
}

//白名单-过滤srcArray，只返回存在whiteArray中的元素
func StringArrayWhiteList(srcArray []string, whiteArray []string) []string {
	var result = make([]string, 0)
	set := make(map[string]struct{}, len(whiteArray))
	for _, v := range whiteArray {
		set[v] = struct{}{}
	}
	for i := 0; i < len(srcArray); i++ {
		if _, ok := set[srcArray[i]]; ok {
			result = append(result, srcArray[i])
		}
	}
	return result
}

//黑名单-过滤srcArray,只返回不存在blackArray中的元素
func StringArrayBlackList(srcArray []string, blackArray []string) []string {
	var result = make([]string, 0)
	set := make(map[string]struct{}, len(blackArray))
	for _, v := range blackArray {
		set[v] = struct{}{}
	}
	for i := 0; i < len(srcArray); i++ {
		if _, ok := set[srcArray[i]]; !ok {
			result = append(result, srcArray[i])
		}
	}
	return result
}
