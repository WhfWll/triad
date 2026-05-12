package services

import (
	"context"
	"encoding/json"
	"errors"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

type Report struct {
}

// GetReportEnum 获取报告类型/报告状态/报告格式枚举
func (r *Report) GetReportEnum() (interface{}, interface{}, interface{}) {
	var reportEnums enums.ReportEnums
	return reportEnums.GetReportTypeEnumArray(), reportEnums.GetReportStatusEnumArray(), reportEnums.GetReportFormatEnumArray()
}

// GetReportContentEnum 获取报告内容枚举
func (r *Report) GetReportContentEnum(ctx context.Context) ([]ReportContentItem, []ReportContentItem, error) {
	var mapSet mysqls.MapSet
	mapSetRes, err := mapSet.GetsByObjKey(ctx, enums.ReportContentMapSetKey)
	if err != nil {
		return nil, nil, err
	}
	if len(mapSetRes.ObjValue) == 0 {
		return nil, nil, errors.New("找不到报告内容枚举")
	}
	var tmp ReportContentMapSet
	err = json.Unmarshal([]byte(mapSetRes.ObjValue), &tmp)
	if err != nil {
		return nil, nil, err
	}
	return tmp.TaskContent, tmp.TargetContent, nil
}

// ReportList 报告清单列表及筛选
func (r *Report) ReportList(ctx context.Context, search string, page, size int, userIdList []int) ([]mysqls.Reportrecord, int64, error) {
	var reportMysqls mysqls.Reportrecord
	return reportMysqls.GetReportrecordList(ctx, "id,name,type,status,format,update_time", search, page, size, userIdList)
}

// ReportDownload 报告下载
func (r *Report) ReportDownload(ctx context.Context, reportId int) (mysqls.Reportrecord, error) {
	var reportMysqls mysqls.Reportrecord
	return reportMysqls.GetReportrecord(ctx, reportId)
}

// ReportDel 报告删除
func (r *Report) ReportDel(ctx context.Context, ids any) error {
	var reportMysqls mysqls.Reportrecord
	return reportMysqls.DeleteReportrecord(ctx, ids)
}

// ReportSave 报告保存
func (r *Report) ReportSave(ctx context.Context, name string, retype int, configjson string, format int, userId int) (int, error) {
	var reportMysqls = mysqls.Reportrecord{
		Name:       name,
		Type:       retype,
		Status:     enums.ReportStatusWait,
		ConfigJSON: configjson,
		Format:     format,
		Content:    "",
		UserID:     userId,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	return reportMysqls.AddReportrecord(ctx)
}

// ReportUpdateContent 更新报告内容
func (r *Report) ReportUpdateContent(ctx context.Context, id int, content string) error {
	var reportModel mysqls.Reportrecord
	return reportModel.UpdateContent(ctx, id, content)
}

// UpdateStatus 更新报告状态
func (r *Report) UpdateStatus(ctx context.Context, id int, status int) error {
	var reportModel mysqls.Reportrecord
	return reportModel.UpdateStatus(ctx, id, status)
}

// GetsByStatus 依据状态获取所有报告
func (r *Report) GetsByStatus(ctx context.Context, status int) []mysqls.Reportrecord {
	var reportModel mysqls.Reportrecord
	return reportModel.GetsByStatus(ctx, status)
}
