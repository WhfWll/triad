package application

import (
	"context"
	"fmt"
	"strings"

	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
)

type DatasecDBTargetApp struct {
	crypto *services.DatasecDBTargetCrypto
}

func NewDatasecDBTargetApp() *DatasecDBTargetApp {
	return &DatasecDBTargetApp{crypto: services.NewDatasecDBTargetCrypto()}
}

func (a *DatasecDBTargetApp) List(ctx context.Context, uid int, req *typespec.DatasecDBTargetListReq) (*typespec.DatasecDBTargetListResp, error) {
	page, size := normalizeDataSecPage(req.Page, req.Size)
	var model mysqls.DatasecDBTarget
	rows, total, err := model.ListByUser(ctx, uid, req.DBType, strings.TrimSpace(req.Search), page, size)
	if err != nil {
		return nil, err
	}
	list := make([]typespec.DatasecDBTargetListItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, typespec.DatasecDBTargetListItem{
			ID:          r.ID,
			Name:        r.Name,
			GroupName:   r.GroupName,
			DBType:      r.DBType,
			DBHost:      r.DBHost,
			DBPort:      r.DBPort,
			DBName:      r.DBName,
			DBUser:      r.DBUser,
			HasPassword: r.DBPassword != "",
			Remark:      r.Remark,
			CreateTime:  formatAppSecTime(r.CreateTime),
			UpdateTime:  formatAppSecTime(r.UpdateTime),
		})
	}
	return &typespec.DatasecDBTargetListResp{List: list, Total: int(total)}, nil
}

func (a *DatasecDBTargetApp) Save(ctx context.Context, uid int, req *typespec.DatasecDBTargetSaveReq) error {
	host := strings.TrimSpace(req.DBHost)
	user := strings.TrimSpace(req.DBUser)
	if host == "" {
		return fmt.Errorf("请填写数据库地址")
	}
	dbType := req.DBType.Int()
	if user == "" && dbType != enums.DBSupportTypeRedis {
		return fmt.Errorf("请填写用户名")
	}
	if dbType < 1 {
		return fmt.Errorf("请选择数据库类型")
	}
	port := req.DBPort.Int()
	if port <= 0 {
		port = defaultDataSecDBPort(dbType)
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = formatDataSecTargetURL(host, port, strings.TrimSpace(req.DBName))
	}
	row := &mysqls.DatasecDBTarget{
		UserID:    uid,
		Name:      name,
		GroupName: strings.TrimSpace(req.GroupName),
		DBType:    dbType,
		DBHost:    host,
		DBPort:    port,
		DBName:    strings.TrimSpace(req.DBName),
		DBUser:    user,
		Remark:    strings.TrimSpace(req.Remark),
	}
	if pwd := req.DBPassword; pwd != "" {
		row.DBPassword = a.crypto.EncryptPassword(pwd)
	}
	var model mysqls.DatasecDBTarget
	if req.ID > 0 {
		old, err := model.GetByID(ctx, req.ID, uid)
		if err != nil || old.ID == 0 {
			return fmt.Errorf("目标不存在")
		}
		row.ID = old.ID
		if pwd := req.DBPassword; pwd == "" {
			row.DBPassword = old.DBPassword
		}
		return model.Update(ctx, row)
	}
	if req.DBPassword == "" && dbType != enums.DBSupportTypeRedis {
		return fmt.Errorf("请填写密码")
	}
	return model.Create(ctx, row)
}

func (a *DatasecDBTargetApp) Delete(ctx context.Context, uid, id int) error {
	var model mysqls.DatasecDBTarget
	return model.Delete(ctx, id, uid)
}

func (a *DatasecDBTargetApp) Import(ctx context.Context, uid int, req *typespec.DatasecDBTargetImportReq) (int, error) {
	if len(req.Items) == 0 {
		return 0, fmt.Errorf("导入列表为空")
	}
	ok := 0
	for i, item := range req.Items {
		saveReq := &typespec.DatasecDBTargetSaveReq{
			Name:       item.Name,
			GroupName:  item.GroupName,
			DBType:     item.DBType,
			DBHost:     item.DBHost,
			DBPort:     item.DBPort,
			DBName:     item.DBName,
			DBUser:     item.DBUser,
			DBPassword: item.DBPassword,
			Remark:     item.Remark,
		}
		if err := a.Save(ctx, uid, saveReq); err != nil {
			return ok, fmt.Errorf("第 %d 条：%w", i+1, err)
		}
		ok++
	}
	return ok, nil
}

func (a *DatasecDBTargetApp) Export(ctx context.Context, uid int, ids []int, includePassword bool) (*typespec.DatasecDBTargetExportResp, error) {
	var model mysqls.DatasecDBTarget
	var rows []mysqls.DatasecDBTarget
	var err error
	if len(ids) > 0 {
		rows, err = model.ListByIDs(ctx, uid, ids)
	} else {
		rows, _, err = model.ListByUser(ctx, uid, 0, "", 1, 10000)
	}
	if err != nil {
		return nil, err
	}
	out := make([]typespec.DatasecDBTargetExportItem, 0, len(rows))
	for _, r := range rows {
		item := typespec.DatasecDBTargetExportItem{
			Name:      r.Name,
			GroupName: r.GroupName,
			DBType:    r.DBType,
			DBHost:    r.DBHost,
			DBPort:    r.DBPort,
			DBName:    r.DBName,
			DBUser:    r.DBUser,
			Remark:    r.Remark,
		}
		if includePassword && r.DBPassword != "" {
			item.DBPassword = a.crypto.DecryptPassword(r.DBPassword)
		}
		out = append(out, item)
	}
	return &typespec.DatasecDBTargetExportResp{Version: 1, Items: out}, nil
}

func (a *DatasecDBTargetApp) ResolveTargets(ctx context.Context, uid int, dbType int, ids []int) ([]datasecTaskConfig, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var model mysqls.DatasecDBTarget
	rows, err := model.ListByIDs(ctx, uid, ids)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("未找到可用的目标库记录")
	}
	out := make([]datasecTaskConfig, 0, len(rows))
	for _, r := range rows {
		if dbType > 0 && r.DBType != dbType {
			continue
		}
		port := r.DBPort
		if port <= 0 {
			port = defaultDataSecDBPort(r.DBType)
		}
		out = append(out, datasecTaskConfig{
			DBType:     r.DBType,
			DBHost:     r.DBHost,
			DBPort:     port,
			DBName:     r.DBName,
			DBUser:     r.DBUser,
			DBPassword: a.crypto.DecryptPassword(r.DBPassword),
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("所选目标与当前数据库类型不匹配")
	}
	return out, nil
}

func (a *DatasecDBTargetApp) SaveTargetsFromTask(ctx context.Context, uid int, dbType int, targets []typespec.DataSecDBTargetInput, groupName string) (int, error) {
	if len(targets) == 0 {
		return 0, fmt.Errorf("没有可保存的目标")
	}
	ok := 0
	for _, t := range targets {
		req := &typespec.DatasecDBTargetSaveReq{
			Name:       formatDataSecTargetURL(strings.TrimSpace(t.DBHost), t.DBPort.Int(), strings.TrimSpace(t.DBName)),
			GroupName:  groupName,
			DBType:     typespec.FlexInt(dbType),
			DBHost:     t.DBHost,
			DBPort:     t.DBPort,
			DBName:     t.DBName,
			DBUser:     t.DBUser,
			DBPassword: t.DBPassword,
		}
		if err := a.Save(ctx, uid, req); err != nil {
			continue
		}
		ok++
	}
	return ok, nil
}

func (a *DatasecDBTargetApp) TestConnection(ctx context.Context, uid int, req *typespec.DatasecDBTargetTestReq) (*typespec.DataSecDBTestConnResp, error) {
	cfg, err := a.resolveTestConfig(ctx, uid, req)
	if err != nil {
		return nil, err
	}
	testReq := &typespec.DataSecDBTestConnReq{
		DBType:     typespec.FlexInt(cfg.DBType),
		DBHost:     cfg.DBHost,
		DBPort:     typespec.FlexInt(cfg.DBPort),
		DBName:     cfg.DBName,
		DBUser:     cfg.DBUser,
		DBPassword: cfg.DBPassword,
	}
	var scan DataSecScan
	return scan.TestDBConnection(ctx, testReq)
}

func (a *DatasecDBTargetApp) BatchTestConnection(ctx context.Context, uid int, req *typespec.DatasecDBTargetBatchTestReq) (*typespec.DatasecDBTargetBatchTestResp, error) {
	if len(req.IDs) == 0 {
		return nil, fmt.Errorf("请至少选择一个目标")
	}
	var model mysqls.DatasecDBTarget
	rows, err := model.ListByIDs(ctx, uid, req.IDs)
	if err != nil {
		return nil, err
	}
	rowByID := make(map[int]mysqls.DatasecDBTarget, len(rows))
	for _, r := range rows {
		rowByID[r.ID] = r
	}
	resp := &typespec.DatasecDBTargetBatchTestResp{
		Total:   len(req.IDs),
		Results: make([]typespec.DatasecDBTargetBatchTestItem, 0, len(req.IDs)),
	}
	var scan DataSecScan
	for _, id := range req.IDs {
		r, ok := rowByID[id]
		item := typespec.DatasecDBTargetBatchTestItem{ID: id}
		if !ok {
			item.OK = false
			item.Message = "目标不存在或无权访问"
			resp.Fail++
			resp.Results = append(resp.Results, item)
			continue
		}
		item.Name = r.Name
		item.DBHost = r.DBHost
		cfg, cfgErr := a.resolveTestConfig(ctx, uid, &typespec.DatasecDBTargetTestReq{ID: r.ID})
		if cfgErr != nil {
			item.OK = false
			item.Message = cfgErr.Error()
			resp.Fail++
			resp.Results = append(resp.Results, item)
			continue
		}
		connReq := &typespec.DataSecDBTestConnReq{
			DBType: typespec.FlexInt(cfg.DBType), DBHost: cfg.DBHost, DBPort: typespec.FlexInt(cfg.DBPort),
			DBName: cfg.DBName, DBUser: cfg.DBUser, DBPassword: cfg.DBPassword,
		}
		result, testErr := scan.TestDBConnection(ctx, connReq)
		if testErr != nil {
			item.OK = false
			item.Message = testErr.Error()
			resp.Fail++
		} else if result != nil && result.OK {
			item.OK = true
			item.Message = result.Message
			resp.OK++
		} else {
			item.OK = false
			if result != nil {
				item.Message = result.Message
			}
			resp.Fail++
		}
		resp.Results = append(resp.Results, item)
	}
	return resp, nil
}

func (a *DatasecDBTargetApp) resolveTestConfig(ctx context.Context, uid int, req *typespec.DatasecDBTargetTestReq) (datasecTaskConfig, error) {
	cfg := datasecTaskConfig{}
	if req.ID > 0 {
		var model mysqls.DatasecDBTarget
		row, err := model.GetByID(ctx, req.ID, uid)
		if err != nil || row.ID == 0 {
			return cfg, fmt.Errorf("目标不存在")
		}
		port := row.DBPort
		if port <= 0 {
			port = defaultDataSecDBPort(row.DBType)
		}
		cfg = datasecTaskConfig{
			DBType: row.DBType, DBHost: row.DBHost, DBPort: port,
			DBName: row.DBName, DBUser: row.DBUser,
			DBPassword: a.crypto.DecryptPassword(row.DBPassword),
		}
	}
	if strings.TrimSpace(req.DBHost) != "" {
		cfg.DBHost = strings.TrimSpace(req.DBHost)
	}
	if req.DBPort.Int() > 0 {
		cfg.DBPort = req.DBPort.Int()
	}
	if req.ID == 0 {
		cfg.DBName = strings.TrimSpace(req.DBName)
	} else if strings.TrimSpace(req.DBName) != "" {
		cfg.DBName = strings.TrimSpace(req.DBName)
	}
	if strings.TrimSpace(req.DBUser) != "" {
		cfg.DBUser = strings.TrimSpace(req.DBUser)
	}
	if req.DBType.Int() > 0 {
		cfg.DBType = req.DBType.Int()
	}
	if req.DBPassword != "" {
		cfg.DBPassword = req.DBPassword
	}
	if cfg.DBHost == "" {
		return cfg, fmt.Errorf("请填写数据库地址")
	}
	if cfg.DBUser == "" && cfg.DBType != enums.DBSupportTypeRedis {
		return cfg, fmt.Errorf("请填写用户名")
	}
	if cfg.DBType < 1 {
		return cfg, fmt.Errorf("请选择数据库类型")
	}
	if cfg.DBPort <= 0 {
		cfg.DBPort = defaultDataSecDBPort(cfg.DBType)
	}
	if cfg.DBPassword == "" && cfg.DBType != enums.DBSupportTypeRedis {
		return cfg, fmt.Errorf("请填写密码，或使用已保存的目标进行测试")
	}
	return cfg, nil
}
