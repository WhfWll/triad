package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type VulScripts struct {
	ID           int       `gorm:"column:id;primary_key" json:"id"`
	DataType     int       `gorm:"column:data_type" json:"dataType"`
	UserId       int       `gorm:"column:user_id" json:"userId"`
	ScriptName   string    `gorm:"column:script_name" json:"scriptName"`
	LibName      string    `gorm:"column:lib_name" json:"libName"`
	Type         string    `gorm:"column:type" json:"type"`
	Content      string    `gorm:"column:content" json:"content"`
	VulID        string    `gorm:"column:vul_id" json:"vulID"`
	VerifyType   string    `gorm:"column:verify_type" json:"verifyType"`
	EvidenceType int       `gorm:"column:evidence_type" json:"evidenceType"`
	Params       string    `gorm:"column:params" json:"params"`
	TargetList   string    `gorm:"column:target_list" json:"targetList"`
	TargetStatus int       `gorm:"column:target_status" json:"targetStatus"`
	Status       int       `gorm:"column:status" json:"status"`
	StatusMsg    string    `gorm:"column:status_msg" json:"statusMsg"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
}

// TableName sets insert table name for this struct type
func (v *VulScripts) TableName() string {
	return "vul_scripts"
}

// GetVulScriptsList 获取脚本列表
func (v *VulScripts) GetVulScriptsList(ctx context.Context, page, limit int, scriptName, libraryName, scriptType, scriptVerifyType, ScriptTargetStatus string) ([]VulScripts, int64, error) {
	var (
		vulScriptsList []VulScripts
		count          int64
		db             = mysql.FromContext(ctx).Model(&VulScripts{})
	)

	if scriptName != "" {
		db = db.Where("script_name like ?", "%"+scriptName+"%")
	}

	if libraryName != "" {
		db = db.Where("lib_name like ?", "%"+libraryName+"%")
	}

	if scriptType != "" {
		db = db.Where("type = ?", scriptType)
	}

	if scriptVerifyType != "" {
		db = db.Where("verify_type = ?", scriptVerifyType)
	}

	if ScriptTargetStatus != "" {
		db = db.Where("target_status = ?", ScriptTargetStatus)
	}
	db.Where("type != 'awvs'")
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&vulScriptsList)
	//db.Limit(limit).Offset(limit * (page - 1)).Order("id").Find(&vulScriptsList)

	return vulScriptsList, count, nil
}

// Get retrieves a single record of vulScripts from database
func (v *VulScripts) GetVulScripts(ctx context.Context) (VulScripts, error) {
	var (
		vulScripts VulScripts
		err        error
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)

	curErr := db.Where("id = ?", v.ID).First(&vulScripts).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return vulScripts, err
}

// 通过IDs获取多个脚本
func (v *VulScripts) GetVulScriptsForIds(ctx context.Context, ids []int) ([]VulScripts, error) {
	var (
		vulScripts []VulScripts
		err        error
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)

	curErr := db.Where("id in ?", ids).Find(&vulScripts).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return vulScripts, err
}

// Add persists vulScripts to database
func (v *VulScripts) AddVulScripts(ctx context.Context) (int, error) {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})

	if err := db.Create(v).Error; err != nil {
		return 0, err
	}

	return v.ID, nil
}

// Update changes vulScripts by id
func (v *VulScripts) UpdateVulScripts(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})

	if err := db.Select("*").Where("id = ?", v.ID).Updates(v).Error; err != nil {
		return err
	}

	return nil
}

// 更新pocname
func (v *VulScripts) UpdateVulScriptsPocnameForId(ctx context.Context, id uint, scriptName string) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})

	if err := db.Where("id = ?", id).Update("script_name", scriptName).Error; err != nil {
		return err
	}

	return nil
}

// 更新取证类型
func (v *VulScripts) UpdateEvidenceTypeForId(ctx context.Context, id int, evidenceType int) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})

	if err := db.Where("id = ?", id).Update("evidence_type", evidenceType).Error; err != nil {
		return err
	}

	return nil
}

// 更新靶站地址
func (v *VulScripts) UpdateTargetListForId(ctx context.Context, id int, targetList string, targetStatus int) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})

	data := make(map[string]interface{})
	data["target_list"] = targetList
	data["target_status"] = targetStatus

	if err := db.Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// Delete vulScripts by id
func (v *VulScripts) DeleteVulScripts(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})
	if err := db.Where("id = ?", v.ID).Delete(v).Error; err != nil {
		return err
	}
	return nil
}

// Delete vulScripts by ids
func (v *VulScripts) DeleteVulScriptsForIds(ctx context.Context, ids []int) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})
	if err := db.Where("id in ?", ids).Delete(v).Error; err != nil {
		return err
	}
	return nil
}

// GetVulScriptForScriptName 通过pocname获取脚本数据
func (v *VulScripts) GetVulScriptForScriptName(ctx context.Context, scriptName string) (VulScripts, error) {
	var (
		vulScripts VulScripts
		err        error
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)
	curErr := db.Where("script_name = ?", scriptName).First(&vulScripts).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return vulScripts, err
}

// GetVulScriptsByscriptType 通过type获取脚本数据
func (v *VulScripts) GetVulScriptsByscriptType(ctx context.Context, scriptType string, scriptName []string) ([]VulScripts, error) {
	var (
		vulScripts []VulScripts
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)
	db.Where("type = ? and script_name IN ?", scriptType, scriptName).Find(&vulScripts)
	return vulScripts, nil
}

// 获取所有数据
func (v *VulScripts) All(ctx context.Context, dataType int, fields string) ([]VulScripts, error) {
	var (
		vulScripts []VulScripts
		err        error
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)

	if fields == "" {
		fields = "*"
	}

	curErr := db.Select(fields).Where("data_type = ?", dataType).Find(&vulScripts).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return vulScripts, err
}

// GetVulScriptByDataTypeAndScriptName 根据数据类型和脚本名获取数据
func (v *VulScripts) GetVulScriptByDataTypeAndScriptName(ctx context.Context, dataType int, scriptName []string) []VulScripts {
	var (
		vulScripts []VulScripts
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)
	db.Where("data_type = ? and script_name IN ?", dataType, scriptName).Find(&vulScripts)
	return vulScripts
}

// 通过vul_ids获取漏洞
func (v *VulScripts) AllVulScriptsForVulIds(ctx context.Context, vulId []string) ([]VulScripts, error) {
	var (
		vulScripts []VulScripts
		err        error
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)

	db.Where("vul_id in ?", vulId).Find(&vulScripts)

	return vulScripts, err
}

// 通过vul_ids获取漏洞
func (v *VulScripts) GetScriptNameByVulIds(ctx context.Context, vulId []string) ([]VulScripts, error) {
	var (
		vulScripts []VulScripts
		err        error
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)

	db.Where("vul_id in ?", vulId).Find(&vulScripts)

	return vulScripts, err
}

// 通过vul_id获取漏洞
func (v *VulScripts) GetVulScriptsByVulId(ctx context.Context, vulId string) ([]VulScripts, error) {
	var (
		vulScripts []VulScripts
		err        error
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)

	curErr := db.Where("vul_id = ?", vulId).Find(&vulScripts).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return vulScripts, err
}

func (v *VulScripts) GetScriptNamesByVulId(ctx context.Context, vulIDs []string) []VulScripts {
	var (
		scriptList []VulScripts
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)
	db.Select("script_name, type").Where("vul_id IN ?", vulIDs).Find(&scriptList)
	return scriptList
}

// 通过verify_type获取漏洞 vulId
func (v *VulScripts) AllVulIdByVerifyType(ctx context.Context, verifyType []string) (res []string) {
	var (
		vulScripts []VulScripts
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)

	db.Select("vul_id").Where("verify_type in ?", verifyType).Find(&vulScripts)

	for _, item := range vulScripts {
		res = append(res, item.VulID)
	}

	return
}

func (v *VulScripts) UpdateStatusById(ctx context.Context, id int, status int, statusMsg string) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})

	if err := db.Where("id = ?", id).Updates(map[string]interface{}{"status": status, "status_msg": statusMsg}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteAllVulScripts 清空vulScripts
func (v *VulScripts) DeleteAllVulScripts(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})
	if err := db.Where("id >= 0").Delete(v).Error; err != nil {
		return err
	}
	return nil
}

// GetVulScriptsByType 通过漏洞类型获取脚本
func (v *VulScripts) GetVulScriptsByType(ctx context.Context, scriptType string, scripName any) []VulScripts {
	var (
		vulScripts []VulScripts
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)
	if scripName == nil {
		db.Where("type = ? ", scriptType).Find(&vulScripts)
	} else {
		db.Where("type = ? and script_name IN ?", scriptType, scripName).Find(&vulScripts)
	}
	return vulScripts
}

// GetScriptPocnameByScriptType 通过脚本类型获取script表中的vul_id
func (v *VulScripts) GetScriptPocnameByScriptType(ctx context.Context, scriptType string) []VulScripts {
	var (
		vulScripts []VulScripts
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)
	db.Select("script_name").Where("type = ? ", scriptType).Find(&vulScripts)
	return vulScripts
}

// 更新pocname
func (v *VulScripts) UpdateVulScriptsById(ctx context.Context, id int, scriptType, verifyType, scriptContent string) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})
	if err := db.Where("id = ?", id).Updates(VulScripts{
		DataType:   1,
		Type:       scriptType,
		VerifyType: verifyType,
		Content:    scriptContent,
		UpdateTime: time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}

// 获取所有数据
func (v *VulScripts) AllScript(ctx context.Context, filter string) ([]VulScripts, error) {
	var (
		vulScripts []VulScripts
		err        error
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
		curErr     error
	)
	if filter != "" {
		curErr = db.Where(filter).Find(&vulScripts).Error
	} else {
		curErr = db.Find(&vulScripts).Error
	}
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return vulScripts, err
}

// DeleteVulScriptsByPocName 通过漏洞名称删除漏洞
func (v *VulScripts) DeleteVulScriptsByPocName(ctx context.Context, name string) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})
	if err := db.Where("script_name = ?", name).Delete(VulScripts{}).Error; err != nil {
		return err
	}
	return nil
}

// GetVulScriptsByScriptTypeForFlow 专门给被动流量提供获取脚本的函数
func (v *VulScripts) GetVulScriptsByScriptTypeForFlow(ctx context.Context, scriptType string, scriptNameList []string) ([]VulScripts, error) {
	var (
		vulScripts []VulScripts
		db         = mysql.FromContext(ctx).Model(&VulScripts{})
	)
	scriptTypeList := make([]string, 0)
	if scriptType == "" {
		scriptTypeList = append(scriptTypeList, "mitm")
		scriptTypeList = append(scriptTypeList, "folderMitm")
	} else {
		scriptTypeList = append(scriptTypeList, scriptType)
	}
	if len(scriptNameList) > 0 {
		db.Where("type in ? and script_name IN ?", scriptTypeList, scriptNameList).Find(&vulScripts)
	} else {
		db.Where("type in ?", scriptTypeList).Find(&vulScripts)
	}
	return vulScripts, nil
}

// 查询多少个脚本
func (v *VulScripts) Count(ctx context.Context) (int64, error) {
	var (
		number int64
		err    error
		db     = mysql.FromContext(ctx).Model(&VulScripts{})
	)
	curErr := db.Count(&number).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return number, err
}

// 删除不在名称列表中的脚本
func (v *VulScripts) DeleteVulScriptsNotInNameList(ctx context.Context, nameList []string) error {
	var db = mysql.FromContext(ctx).Model(&VulScripts{})
	if err := db.Where("script_name not in ?", nameList).Delete(v).Error; err != nil {
		return err
	}
	return nil
}
