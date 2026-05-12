package invoke

//
////buildBruteParam 进行弱口令脚本调用
//func (c *CallInfo) buildBruteParam(ctx context.Context, scriptParamList []string, fName string, params []string) ([]string, []string) {
//	scriptParamList = append(scriptParamList, fName)
//	tempFileList := make([]string, 0)
//	var dictModel mysqls.Dictionary
//	serviceEnum := enums.DictionaryServiceEnum()
//	var service int
//	//第一步 根据参数获取到服务名字
//	paramMap := make(map[string]string, 0)
//	for _, param := range c.ToolParamList {
//		if strings.HasSuffix(param.ParamName, "_port") {
//			for key, value := range serviceEnum {
//				if value == strings.ReplaceAll(param.ParamName, "_port", "") {
//					service = key
//					break
//				}
//			}
//		}
//		paramMap[param.ParamName] = param.ParamValue
//	}
//	//第二步 根据类型获取字典内容
//	var userContent, passContent string
//	if paramMap["dictType"] == strconv.Itoa(enums.TaskConfigurationWeakPassDictTypeCommon) && paramMap["commonUserDict"] != "" && paramMap["commonPassDict"] != "" {
//		dictUser, _ := dictModel.GetDictById(ctx, paramMap["commonUserDict"])
//		userContent = dictUser.Content
//		dictPass, _ := dictModel.GetDictById(ctx, paramMap["commonPassDict"])
//		passContent = dictPass.Content
//	} else if paramMap["dictType"] == strconv.Itoa(enums.TaskConfigurationWeakPassDictTypeAdd) {
//		if paramMap["onlyUseAdd"] == "1" {
//			userContent = paramMap["addAccount"]
//			passContent = paramMap["addPass"]
//		}
//	}
//	//第三步 使用默认用户字典和密码字典
//	if userContent == "" && passContent == "" {
//		dictList := dictModel.GetDictByServiceAndIsDefault(ctx, service, 1)
//		for _, dict := range dictList {
//			if dict.Types == enums.DictionaryTypeUser {
//				if paramMap["onlyUseAdd"] == "0" {
//					userContent = dict.Content + "\n" + paramMap["addAccount"]
//				} else {
//					userContent = dict.Content
//				}
//			} else if dict.Types == enums.DictionaryTypePassword {
//				if paramMap["onlyUseAdd"] == "0" {
//					passContent = dict.Content + "\n" + paramMap["addPass"]
//				} else {
//					passContent = dict.Content
//				}
//			}
//		}
//	}
//
//	//第四步 构建参数
//	fBruteUserFile, err := ioutil.TempFile("", "tempfile-*.txt")
//	if err != nil {
//		return scriptParamList, tempFileList
//	}
//	fBruteUserFile.WriteString(userContent)
//	fBruteUserFile.Close()
//	scriptParamList = append(scriptParamList, "--user-list-file")
//	scriptParamList = append(scriptParamList, fBruteUserFile.Name())
//	tempFileList = append(tempFileList, fBruteUserFile.Name())
//	fBrutePassFile, err := ioutil.TempFile("", "tempfile-*.txt")
//	if err != nil {
//		return scriptParamList, tempFileList
//	}
//	fBrutePassFile.WriteString(passContent)
//	fBrutePassFile.Close()
//	scriptParamList = append(scriptParamList, "--pass-list-file")
//	scriptParamList = append(scriptParamList, fBrutePassFile.Name())
//	tempFileList = append(tempFileList, fBrutePassFile.Name())
//
//	scriptParamList = append(scriptParamList, params...)
//	return scriptParamList, tempFileList
//}
