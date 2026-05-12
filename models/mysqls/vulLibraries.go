package mysqls

import (
	"context"
	"fmt"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"strings"
	"time"
)

type VulLibraries struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`
	DataType        int       `gorm:"column:data_type" json:"dataType"`
	VulID           string    `gorm:"column:vul_id" json:"VulID"`
	Name            string    `gorm:"column:name" json:"name"`
	Cve             string    `gorm:"column:cve" json:"cve"`
	Risk            int       `gorm:"column:risk" json:"risk"`
	Type            int       `gorm:"column:type" json:"type"`
	Class           int       `gorm:"column:class" json:"class"`
	PublishedTime   string    `gorm:"column:published_time" json:"publishedTime"`
	Description     string    `gorm:"column:description" json:"description"`
	AffectRange     string    `gorm:"column:affect_range" json:"affectRange"`
	ExploitImpact   int       `gorm:"column:exploit_impact" json:"exploitImpact"`
	FixSuggest      string    `gorm:"column:fix_suggest" json:"fixSuggest"`
	Cnvd            string    `gorm:"column:cnvd" json:"cnvd"`
	Cnnvd           string    `gorm:"column:cnnvd" json:"cnnvd"`
	Cncve           string    `gorm:"column:cncve" json:"cncve"`
	Bugtraq         string    `gorm:"column:bugtraq" json:"bugtraq"`
	Component       string    `gorm:"column:component" json:"component"`
	Status          int       `gorm:"column:status" json:"status"`
	StatusMsg       string    `gorm:"column:status_msg" json:"statusMsg"`
	Priority        string    `gorm:"column:priority" json:"priority"`
	OperatingSystem int       `gorm:"column:operating_system" json:"operatingSystem"`
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`
	Pocname         string    `gorm:"column:pocname" json:"pocname"`
	VerifyType      int       `gorm:"column:verify_type" json:"verifyType"`
	PatchUrl        string    `gorm:"column:patch_url" json:"patch_url"`      // 补丁地址
	CvssVersion     string    `gorm:"column:cvss_version" json:"cvssVersion"` // 补丁地址
	CvssScore       string    `gorm:"column:cvss_score" json:"cvssScore"`     // 补丁地址
	PocOrExp        string    `gorm:"column:poc_or_exp" json:"poc_or_exp"`    // poc或者exp
	ScriptType      string    `gorm:"column:script_type" json:"script_type"`  // 脚本类型
	VulSource       int       `gorm:"column:vul_source" json:"vulSource"`     // 漏洞来源（是否国产化）
}

// TableName sets insert table name for this struct type
func (v *VulLibraries) TableName() string {
	return "vul_libraries"
}

// Get retrieves a list of vulLibraries from database
func (v *VulLibraries) GetVulLibrariesList(ctx context.Context, page, limit int) ([]VulLibraries, int64, error) {
	var (
		vulLibrariesList []VulLibraries
		count            int64
		db               = mysql.FromContext(ctx).Model(&VulLibraries{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&vulLibrariesList)
	db.Count(&count)

	return vulLibrariesList, count, nil
}

// Get retrieves a single record of vulLibraries from database
func (v *VulLibraries) GetVulLibraries(ctx context.Context) (VulLibraries, error) {
	var (
		vulLibraries VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
	)

	curErr := db.Where("id = ?", v.ID).First(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return vulLibraries, err
}

// Add persists vulLibraries to database
func (v *VulLibraries) AddVulLibraries(ctx context.Context) (int, error) {
	var db = mysql.FromContext(ctx).Model(&VulLibraries{})

	if err := db.Create(v).Error; err != nil {
		return 0, err
	}

	return v.ID, nil
}

// Count 获取漏洞库总数
func (v *VulLibraries) Count(ctx context.Context) (int64, error) {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&VulLibraries{})
	)
	err := db.Count(&count).Error
	return count, err
}

// Update changes vulLibraries by id
func (v *VulLibraries) UpdateVulLibraries(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&VulLibraries{})

	if err := db.Select("*").Where("id = ?", v.ID).Updates(v).Error; err != nil {
		return err
	}

	return nil
}

// Delete vulLibraries by id
func (v *VulLibraries) DeleteVulLibraries(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&VulLibraries{})

	//v.Estate = "deleted"
	v.UpdateTime = time.Now()
	if err := db.Where("id = ?", v.ID).Updates(v).Error; err != nil {
		return err
	}

	return nil
}

// 删除多个 依据IDs
func (v *VulLibraries) DeleteVulLibrariesForVulIds(ctx context.Context, vulIds []string) error {
	var db = mysql.FromContext(ctx).Model(&VulLibraries{})

	v.UpdateTime = time.Now()
	if err := db.Where("vul_id in ?", vulIds).Delete(v).Error; err != nil {
		return err
	}

	return nil
}

// GetVulLibrariesForName 通过name获取漏洞
func (v *VulLibraries) GetVulLibrariesForName(ctx context.Context, name string) (VulLibraries, error) {
	var (
		vulLibraries VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
	)
	curErr := db.Where("name = ?", name).First(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return vulLibraries, err
}

// 通过vul_id获取漏洞
func (v *VulLibraries) GetVulLibrariesForVulId(ctx context.Context, vulId string) (VulLibraries, error) {
	var (
		vulLibraries VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
	)

	curErr := db.Where("vul_id = ?", vulId).First(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return vulLibraries, err
}

// 通过vul_id获取漏洞
func (v *VulLibraries) AllVulLibrariesForVulIds(ctx context.Context, vulId []string) ([]VulLibraries, error) {
	var (
		vulLibraries []VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
	)

	curErr := db.Where("vul_id in ?", vulId).Find(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return vulLibraries, err
}

// 通过ids获取漏洞
func (v *VulLibraries) AllVulLibrariesForIds(ctx context.Context, ids []int) ([]VulLibraries, error) {
	var (
		vulLibraries []VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
	)

	curErr := db.Where("id in ?", ids).Find(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return vulLibraries, err
}

// 通过ids获取漏洞
// 通过ids获取漏洞的vul_id集合
func (v *VulLibraries) GetVulIdsByIds(ctx context.Context, ids []int) ([]string, error) {
	var (
		vulIds []string
		err    error
		db     = mysql.FromContext(ctx).Model(&VulLibraries{})
	)

	curErr := db.Where("id in ?", ids).Pluck("vul_id", &vulIds).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return vulIds, err
}

// Update changes vulLibraries by id
func (v *VulLibraries) UpdateVulIdForId(ctx context.Context, id int, vulId string) error {
	var db = mysql.FromContext(ctx).Model(&VulLibraries{})

	if err := db.Where("id = ?", id).Update("vul_id", vulId).Error; err != nil {
		return err
	}

	return nil
}

func (v *VulLibraries) UpdateStatusById(ctx context.Context, id, status int) error {
	var db = mysql.FromContext(ctx).Model(&VulLibraries{})

	if err := db.Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func (v *VulLibraries) GetVulLibCount(ctx context.Context) (int64, error) {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&VulLibraries{})
	)
	err := db.Count(&count).Error
	return count, err
}

// GetVulLibList 漏洞 - 列表页
func (v *VulLibraries) GetVulLibList(ctx context.Context, page, limit, verifyType int, libIds []int, vulIds []string, search string, libTypes []int, libClasses []int, libRisks []int,
	exploitImpact []int, operateSystem []int, ptOrder string, status []int, code string) ([]VulLibraries, int64, error) {
	var (
		vulLibraries []VulLibraries
		count        int64
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
		query        string
		args         []interface{}
	)

	if len(libIds) > 0 {
		if query != "" {
			query += " and id in ?"
		} else {
			query += "id in ?"
		}
		args = append(args, libIds)
	}

	if len(vulIds) > 0 {
		if query != "" {
			query += " and vul_id in ?"
		} else {
			query += "vul_id in ?"
		}
		args = append(args, vulIds)
	}

	if search != "" {
		if query != "" {
			query += " and (name LIKE ? or cve like ?) "
		} else {
			query += "(name LIKE ? or cve like ?) "
		}
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	if code != "" {
		if query != "" {
			query += " and (cve like ? or cnvd like ? or cnnvd like ? or cncve like ? or bugtraq like ?)"
		} else {
			query += "(cve like ? or cnvd like ? or cnnvd like ? or cncve like ? or bugtraq like ?)"
		}
		args = append(args, "%"+code+"%", "%"+code+"%", "%"+code+"%", "%"+code+"%", "%"+code+"%")
	}

	if len(libTypes) > 0 {
		if query != "" {
			query += " and type in ?"
		} else {
			query += "type in ?"
		}
		args = append(args, libTypes)
	}

	if len(libClasses) > 0 {
		if query != "" {
			query += " and class in ?"
		} else {
			query += "class in ?"
		}
		args = append(args, libClasses)
	}

	if len(libRisks) > 0 {
		if query != "" {
			query += " and risk in ?"
		} else {
			query += "risk in ?"
		}
		args = append(args, libRisks)
	}

	if len(status) > 0 {
		if query != "" {
			query += " and status in ?"
		} else {
			query += "status in ?"
		}
		args = append(args, status)
	}

	if len(exploitImpact) != 0 {
		if query != "" {
			query += " and exploit_impact in ?"
		} else {
			query += "exploit_impact in ?"
		}
		args = append(args, exploitImpact)
	}

	if len(operateSystem) > 0 {
		if query != "" {
			query += " and operating_system in ?"
		} else {
			query += "operating_system in ?"
		}
		args = append(args, operateSystem)
	}
	if verifyType != 0 {
		if query != "" {
			query += " and verify_type = ?"
		} else {
			query += "verify_type = ?"
		}
		args = append(args, verifyType)
	}

	db = db.Where(query, args...)

	switch ptOrder {
	case "-published_time":
		db = db.Order("published_time DESC")
	case "published_time":
		db = db.Order("published_time ASC")
	}

	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&vulLibraries)

	return vulLibraries, count, nil
}

// GetVulIdList 获取漏洞id列表
func (v *VulLibraries) GetVulIdList(ctx context.Context, libIds []int, vulIds []string, libName string, libTypes []int, libClasses []int, libRisks []int,
	exploitImpact []int, operateSystem []int, ptOrder string, status []int, verifyType int) ([]int, error) {
	var (
		db    = mysql.FromContext(ctx).Model(&VulLibraries{})
		query string
		args  []interface{}
	)

	if len(libIds) > 0 {
		if query != "" {
			query += " and id in ?"
		} else {
			query += "id in ?"
		}
		args = append(args, libIds)
	}

	if len(vulIds) > 0 {
		if query != "" {
			query += " and vul_id in ?"
		} else {
			query += "vul_id in ?"
		}
		args = append(args, vulIds)
	}

	if libName != "" {
		if query != "" {
			query += " and name LIKE ?"
		} else {
			query += "name LIKE ?"
		}
		args = append(args, "%"+libName+"%")
	}

	if len(libTypes) > 0 {
		if query != "" {
			query += " and type in ?"
		} else {
			query += "type in ?"
		}
		args = append(args, libTypes)
	}

	if len(libClasses) > 0 {
		if query != "" {
			query += " and class in ?"
		} else {
			query += "class in ?"
		}
		args = append(args, libClasses)
	}

	if len(libRisks) > 0 {
		if query != "" {
			query += " and risk in ?"
		} else {
			query += "risk in ?"
		}
		args = append(args, libRisks)
	}

	if len(status) > 0 {
		if query != "" {
			query += " and status in ?"
		} else {
			query += "status in ?"
		}
		args = append(args, status)
	}

	if len(exploitImpact) != 0 {
		if query != "" {
			query += " and exploit_impact in ?"
		} else {
			query += "exploit_impact in ?"
		}
		args = append(args, exploitImpact)
	}

	if len(operateSystem) > 0 {
		if query != "" {
			query += " and operating_system in ?"
		} else {
			query += "operating_system in ?"
		}
		args = append(args, operateSystem)
	}

	if verifyType != 0 {
		if query != "" {
			query += " and verify_type = ?"
		} else {
			query += "verify_type = ?"
		}
		args = append(args, verifyType)
	}

	db = db.Debug().Where(query, args...)

	switch ptOrder {
	case "-published_time":
		db = db.Order("published_time DESC")
	case "published_time":
		db = db.Order("published_time ASC")
	}

	var vulIdList []int
	db.Pluck("id", &vulIdList)

	return vulIdList, nil
}

// 通过vul_id 获取vul_id
// exploitImpacts必选是[]int或[]string
func (v *VulLibraries) GetVulIDsById(ctx context.Context, ids []int, exploitImpacts any) []string {
	var (
		vulIDs []string
		db     = mysql.FromContext(ctx).Model(&VulLibraries{})
	)
	db.Select("vul_id").Where("id in ? and exploit_impact in ? and status = ?", ids, exploitImpacts, enums.VulLibrariesStatusSucess).Find(&vulIDs)
	return vulIDs
}

// GetVulLibrariesByPocname 通过pocname获取漏洞
func (v *VulLibraries) GetVulLibrariesByPocname(ctx context.Context, pocname string) (VulLibraries, error) {
	var (
		vulLibraries VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
	)
	curErr := db.Where("pocname = ?", pocname).First(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return vulLibraries, err
}

// GetVulLibrariesIdEqPocnames 通过pocnames获取漏洞
func (v *VulLibraries) GetVulLibrariesIdEqPocnames(ctx context.Context, pocnames []string) ([]VulLibraries, error) {
	var (
		vulLibraries []VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
	)

	curErr := db.Select("id").Where("pocname in ?", pocnames).Debug().Find(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return vulLibraries, err
}

// GetVulLibrariesIdByPocnames 通过pocnames模糊匹配获取漏洞
func (v *VulLibraries) GetVulLibrariesIdLikePocnames(ctx context.Context, pocnames []string) ([]VulLibraries, error) {
	var (
		vulLibraries []VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{}).Select("id")
	)
	where := make([]string, 0)
	for _, item := range pocnames {
		where = append(where, "pocname like '%"+item+"%'")
	}
	curErr := db.Select("id").Where(strings.Join(where, " or ")).Find(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return vulLibraries, err
}

// GetVulCount 获取漏洞数量
func (v *VulLibraries) GetVulCount(ctx context.Context) int64 {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&VulLibraries{})
	)
	db.Count(&count)
	return count
}

// GetVulLibrariesCanVerify 获取可以进行验证的脚本
func (v *VulLibraries) GetVulLibrariesCanVerify(ctx context.Context) ([]VulLibraries, error) {
	var (
		vulLibraries []VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
	)
	curErr := db.Where(`cve!="" and pocname !="" and vul_id not like "%-%"`).Find(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return vulLibraries, err
}

// 获取没有描述或修复建议的漏洞
func (v *VulLibraries) AllVulLibrariesByDataEmpty(ctx context.Context) ([]VulLibraries, error) {
	var (
		vulLibraries []VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
	)

	curErr := db.
		Where("description = '' or fix_suggest = '' or affect_range = ''").
		Find(&vulLibraries).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return vulLibraries, err
}

func (v *VulLibraries) UpdateByMap(ctx context.Context, id int, data map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&VulLibraries{})

	fmt.Println(id, data)
	if err := db.Where("id = ?", id).Updates(data).Debug().Error; err != nil {
		return err
	}

	return nil
}

// 获取所有漏洞数据  漏洞导出用
func (v *VulLibraries) AllVulLibraries(ctx context.Context, filter string) ([]VulLibraries, error) {
	var (
		vulLibraries []VulLibraries
		err          error
		db           = mysql.FromContext(ctx).Model(&VulLibraries{})
		curErr       error
	)
	if filter != "" {
		curErr = db.Where(filter).Find(&vulLibraries).Error
	} else {
		curErr = db.Find(&vulLibraries).Error
	}
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return vulLibraries, err
}

// AllVul 获取所有的漏洞数据
func (v *VulLibraries) AllVul(ctx context.Context) []VulLibraries {
	var (
		vulList []VulLibraries
		db      = mysql.FromContext(ctx).Model(&VulLibraries{})
	)

	db.Find(&vulList)

	return vulList
}

// AllVul 获取所有的漏洞数据
func (v *VulLibraries) GetVulLibByVulScripts(ctx context.Context, vulIdList []string) []VulLibraries {
	var (
		vulList []VulLibraries
		db      = mysql.FromContext(ctx).Model(&VulLibraries{})
	)
	db.Where("vul_id in ?", vulIdList).Find(&vulList)
	return vulList
}

// 删除不在名称列表中的漏洞
func (v *VulLibraries) DeleteVulLibrariesNotInNameList(ctx context.Context, nameList []string) error {
	var db = mysql.FromContext(ctx).Model(&VulLibraries{})
	v.UpdateTime = time.Now()
	if err := db.Where("vul_id not in ?", nameList).Delete(v).Error; err != nil {
		return err
	}
	return nil
}
