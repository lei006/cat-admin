package model

import (
	"errors"

	"github.com/sohaha/zlsgo/zlog"
	"gorm.io/gorm"
)

type SysUser struct {
	BASE_MODEL
	Username string `json:"username" gorm:"uniqueIndex;comment:用户登录名"` // 用户登录名
	Password string `json:"-"  gorm:"comment:用户登录密码"`                  // 用户登录密码
	Nickname string `json:"nickname" gorm:"default:系统用户;comment:用户昵称"` // 用户昵称
	Token    string `json:"token" gorm:"-"`                            // token
	//Headerimg string `json:"headerimg" gorm:"type:mediumtext; comment:用户头像"` // 用户头像
	UserSign string `json:"user_sign" gorm:"type:mediumtext;comment:用户签名"` // 用户签名
	IsAdmin  bool   `json:"is_admin" gorm:"default:0;comment:是否管理员"`       // 用户角色ID
	Phone    string `json:"phone"  gorm:"comment:用户手机号"`                   // 用户手机号
	//Email     string `json:"email"  gorm:"comment:用户邮箱"`                          // 用户邮箱
	IsDisable bool `json:"is_disable" gorm:"default:0;comment:用户是否被冻结 0正常 1冻结"` //用户是否被冻结 0正常 1冻结
}

func (SysUser) TableName() string {
	return "sys_user"
}

// CreateReportStruct 创建报告的数据结构记录
// Author [piexlmax](https://github.com/piexlmax)
func (model *SysUser) Create(newVal *SysUser) (err error) {
	err = g_db.Create(newVal).Error
	return err
}

// DeleteReportStruct 删除报告的数据结构记录
// Author [piexlmax](https://github.com/piexlmax)
func (model *SysUser) DeleteOne(id string) (err error) {
	err = g_db.Unscoped().Delete(&SysUser{}, id).Error
	return err
}

// DeleteReportStructByIds 批量删除报告的数据结构记录
// Author [piexlmax](https://github.com/piexlmax)
func (model *SysUser) DeleteMany(ids []uint) (err error) {
	err = g_db.Unscoped().Delete(&[]SysUser{}, "id in ?", ids).Error
	return err
}

func (model *SysUser) UpdateOne(val *SysUser) (err error) {
	err = g_db.Model(val).Select("nickname", "phone", "user_sign").Save(val).Error
	return err
}

func (model *SysUser) PatchOne(id string, field string, data interface{}) error {

	result := g_db.Model(&SysUser{}).Where("id = ?", id).Update(field, data)
	if result.RowsAffected == 0 {
		return errors.New("No rows were updated")
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (model *SysUser) GetField(id string, field string) (*SysUser, error) {
	//err = g_db.Where("id = ?", id).First(&retVal).Error
	retVal := &SysUser{}
	err := g_db.Where("id = ?", id).Select(field).First(retVal).Error
	return retVal, err
}

func (model *SysUser) GetOne(id string) (retVal SysUser, err error) {
	err = g_db.Where("id = ?", id).First(&retVal).Error
	return
}

func (model *SysUser) GetOneByUsername(username string) (retVal *SysUser, err error) {
	// 通过username 取得一行
	retVal = &SysUser{}
	err = g_db.Where("username = ?", username).First(retVal).Error
	return
}

func (model *SysUser) GetOneByToken(token string) (retVal *SysUser, err error) {
	// 通过username 取得一行
	retVal = &SysUser{}
	err = g_db.Where("token = ?", token).First(retVal).Error
	return
}

// GetReportStructInfoList 分页获取报告的数据结构记录
// Author [piexlmax](https://github.com/piexlmax)
func (model *SysUser) GetList() (list []SysUser, total int64, err error) {

	// 创建db
	db := g_db.Model(&SysUser{})
	var reportItems []SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&reportItems).Error
	return reportItems, total, err
}

// 只能取一个
func (model *SysUser) GetListOne(id string) (list []SysUser, total int64, err error) {

	// 创建db
	db := g_db.Model(&SysUser{})
	db = db.Where("id = ?", id)

	var reportItems []SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&reportItems).Error
	return reportItems, total, err
}

func (model *SysUser) NoFoundCreate(username string, password string, nickname string, isadmin bool, phone string, isDisable bool) error {

	val := SysUser{
		Username:  username,
		Password:  password,
		Nickname:  nickname,
		IsAdmin:   isadmin,
		Phone:     phone,
		IsDisable: isDisable,
	}
	err := g_db.Where("username = ?", username).First(&val).Error
	if err != gorm.ErrRecordNotFound { // 如果不是未找到错误，则说明用户已存在
		zlog.Error(err)
		return err
	}

	err = g_db.Create(&val).Error
	if err != nil {
		zlog.Error(err)
		return err
	}
	return nil
}
