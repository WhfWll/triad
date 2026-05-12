package enums

const (
	//验证报告类型
	ReportVerifyProducerUnKnown      = iota //未知报告类型
	ReportVerifyProducerTianJing            //启明天镜报告
	ReportVerifyProducerNsfocus             //绿盟nsfocus报告
	ReportVerifyProducerTianJing2025        //启明天镜2025报告
)

const (
	// 报告验证文件类型
	ReportVerifyFileTypeUnKnown = iota //未知文件类型
	ReportVerifyFileTypeHtml           //html文件
	ReportVerifyFileTypeZip            //zip文件
)

const (
	// 报告验证漏洞验证状态
	ReportVerifyStatusUnknown  = iota //未知文件类型
	ReportVerifyStatusUnVerify        //未能验证
	ReportVerifyStatusVerify          //验证成功
	ReportVerifyStatusFailed          //验证失败
	ReportVerifyStatusExp             //利用成功
)

var ReportVerifyEnum ReportVerify

type ReportVerify struct {
}

func (r ReportVerify) AllReportVerifyTypeEnum() map[int]string {
	enum := map[int]string{
		ReportVerifyProducerTianJing: "启明天境",
		ReportVerifyProducerNsfocus:  "绿盟NsFocus",
	}
	return enum
}

func (r ReportVerify) ReportVerifyProducerEnum(verifyType int) string {
	enum := r.AllReportVerifyTypeEnum()
	if res, ok := enum[verifyType]; ok {
		return res
	}
	return ""
}

func (r ReportVerify) AllReportVerifyStatusEnum() map[int]string {
	enum := map[int]string{
		//ReportVerifyStatusUnknown:  "未知状态",
		ReportVerifyStatusUnVerify: "未能验证",
		ReportVerifyStatusVerify:   "验证成功",
		ReportVerifyStatusFailed:   "验证失败",
		ReportVerifyStatusExp:      "利用成功",
	}
	return enum
}

func (r ReportVerify) ReportVerifyStatusEnum(status int) string {
	enum := r.AllReportVerifyStatusEnum()
	if res, ok := enum[status]; ok {
		return res
	}
	return ""
}
