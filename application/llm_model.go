package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/chat"
	"smart/tools/enums"
	"strconv"
)

type LlmModel struct {
}

// LlmModelList 获取大模型列表
func (l *LlmModel) LlmModelList(ctx context.Context, req *typespec.LlmModelListReq, res *typespec.LlmModelListRes) error {
	var llmModelSrv services.LlmModel

	llmModels, total, err := llmModelSrv.GetLlmModelList(ctx, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}

	res.List = make([]typespec.LlmModelItem, 0)
	for _, model := range llmModels {
		res.List = append(res.List, typespec.LlmModelItem{
			ID:            model.ID,
			ModelName:     model.ModelName,
			Platform:      model.Platform,
			PlatformDesc:  enums.LlmEnums.GetLlmModelPlatformName(model.Platform),
			ApiUrl:        model.ApiUrl,
			ApiKey:        model.ApiKey,
			ModelID:       model.ModelID,
			ModelType:     model.ModelType,
			ModelTypeDesc: enums.LlmEnums.GetLlmModelTypeName(model.ModelType),
			IsDefault:     model.IsDefault,
			Status:        model.Status,
			CreateTime:    model.CreateTime,
			UpdateTime:    model.UpdateTime,
		})
	}
	res.Total = total

	return nil
}

// LlmModelDetail 获取大模型详情
func (l *LlmModel) LlmModelDetail(ctx context.Context, req *typespec.LlmModelDetailReq, res *typespec.LlmModelDetailRes) error {
	var llmModelSrv services.LlmModel

	model, err := llmModelSrv.GetLlmModelDetail(ctx, req.ID)
	if err != nil {
		return err
	}

	res.ID = model.ID
	res.ModelName = model.ModelName
	res.Platform = model.Platform
	res.PlatformDesc = enums.LlmEnums.GetLlmModelPlatformName(model.Platform)
	res.ApiUrl = model.ApiUrl
	res.ApiKey = model.ApiKey
	res.ModelID = model.ModelID
	res.ModelType = model.ModelType
	res.ModelTypeDesc = enums.LlmEnums.GetLlmModelTypeName(model.ModelType)
	res.IsDefault = model.IsDefault
	res.Status = model.Status
	res.CreateTime = model.CreateTime
	res.UpdateTime = model.UpdateTime

	return nil
}

// LlmModelSave 保存大模型（添加/编辑）
func (l *LlmModel) LlmModelSave(ctx context.Context, req *typespec.LlmModelSaveReq, res *typespec.LlmModelSaveRes) error {
	var llmModelSrv services.LlmModel

	if req.IsDefault != enums.LlmModelIsDefault {
		req.IsDefault = enums.LlmModelNotDefault
	}

	// 如果ID为0或未提供，则为添加操作
	if req.ID == 0 {
		err := llmModelSrv.AddLlmModel(ctx, req.ModelName, req.Platform, req.ApiUrl, req.ApiKey, req.ModelID, req.ModelType, req.IsDefault)
		if err != nil {
			return err
		}
	} else {
		// 如果ID不为0，则为编辑操作
		err := llmModelSrv.UpdateLlmModel(ctx, req.ID, req.ModelName, req.Platform, req.ApiUrl, req.ApiKey, req.ModelID, req.ModelType, req.IsDefault, req.Status)
		if err != nil {
			return err
		}
	}

	res.Success = true
	return nil
}

// LlmModelDelete 删除大模型
func (l *LlmModel) LlmModelDelete(ctx context.Context, req *typespec.LlmModelDeleteReq, res *typespec.LlmModelDeleteRes) error {
	var llmModelSrv services.LlmModel

	err := llmModelSrv.DeleteLlmModels(ctx, req.IDs)
	if err != nil {
		return err
	}

	res.Success = true
	return nil
}

// LlmModelSetDefault 设置默认大模型
func (l *LlmModel) LlmModelSetDefault(ctx context.Context, req *typespec.LlmModelSetDefaultReq, res *typespec.LlmModelSetDefaultRes) error {
	var llmModelSrv services.LlmModel

	err := llmModelSrv.SetDefaultLlmModel(ctx, req.ID)
	if err != nil {
		return err
	}

	res.Success = true
	return nil
}

type openAiErrorMessage struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
	RequestId string `json:"request_id"`
}

// LlmModelEnabledTest 测试大模型是否可用
func (l *LlmModel) LlmModelEnabledTest(ctx context.Context, req *typespec.LlmModelEnabledTestReq, resp *typespec.LlmModelEnabledTestRes) error {
	var srv services.LlmModel
	model, err := srv.GetLlmModelDetail(ctx, req.Id)
	if err != nil {
		return err
	}
	if model.ID == 0 {
		return errors.New("模型不存在")
	}
	resp.Enabled = enums.LlmModelEnabledYes
	client := chat.NewCommonClient(model.ApiUrl, model.ApiKey, model.ModelID, false, false)
	err = client.Run(ctx, chat.NoneCallBack, chat.AiMessage{
		Role:    "user",
		Content: "你是谁",
	})
	if err != nil {
		resp.Enabled = enums.LlmModelEnabledNo
		var errMsg openAiErrorMessage
		_ = json.Unmarshal([]byte(err.Error()), &errMsg)
		if errMsg.Error.Message != "" {
			resp.ErrorMsg = err.Error()
		} else {
			resp.ErrorMsg = fmt.Sprintf("%s: %s", errMsg.Error.Code, errMsg.Error.Message)
		}
	}
	resp.EnabledDesc = enums.LlmEnums.GetLlmModelEnableName(resp.Enabled)
	return nil
}

func (l *LlmModel) LlmModelEnums(resp *typespec.LlmModelEnumsResp) {
	resp.ModelType = enums.LlmEnums.LlmModelTypeEnums()
	resp.Platform = enums.LlmEnums.LlmModelPlatformEnums()
	resp.Enhance = enums.LlmEnums.LlmModelEnhancementEnums()
	resp.DefaultLlm = enums.LlmEnums.LlmModelIsDefaultEnums()
}

func (l *LlmModel) LlmModelEnhancementDetail(ctx context.Context, resp *typespec.LlmModelEnhancementDetailResp) error {
	var mapSetSrv services.MapSet
	objValue, err := mapSetSrv.GetMapValue(ctx, enums.LlmModelEnhancementObjKey)
	if err != nil {
		return err
	}
	status, _ := strconv.Atoi(objValue)
	if status == 0 {
		err = mapSetSrv.Create(ctx,
			enums.LlmModelEnhancementObjKey,
			strconv.Itoa(enums.LlmModelEnhancementClose),
			enums.LlmModelEnhancementContent)
		if err != nil {
			return err
		}
		status = enums.LlmModelEnhancementClose
	}
	resp.Status = status
	resp.StatusDesc = enums.LlmEnums.GetLlmModelEnhancementName(status)
	return nil
}

func (l *LlmModel) LlmModelEnhancementEdit(ctx context.Context, req *typespec.LlmModelEnhancementEditReq) error {
	var (
		mapSetSrv   services.MapSet
		llmModelSrv services.LlmModel
	)
	if req.Status == enums.LlmModelEnhancementOpen {
		model, _ := llmModelSrv.GetDefaultLlmModel(ctx, enums.LlmModelTypeText)
		if model.ID == 0 {
			return errors.New("至少需要配置一个文本模型才可以开启模型增强功能")
		}
	}
	objValue, _ := mapSetSrv.GetMapValue(ctx, enums.LlmModelEnhancementObjKey)
	if objValue == "" {
		return mapSetSrv.Create(ctx,
			enums.LlmModelEnhancementObjKey,
			strconv.Itoa(req.Status),
			enums.LlmModelEnhancementContent)
	}
	return mapSetSrv.UpdateMapValue(ctx, enums.LlmModelEnhancementObjKey, strconv.Itoa(req.Status))
}
