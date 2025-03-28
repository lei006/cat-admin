package model

import (
	"time"

	"github.com/sohaha/zlsgo/zlog"
)

type SysAuthInfo struct {
	BASE_MODEL
	Username string `json:"username" gorm:"index;comment:用户名"`              // 用户登录名
	Token    string `json:"token" gorm:"index;default:token;comment:token"` // token
	Index    uint   `json:"index" gorm:"default:0"`                         // 主键ID
	CurNum   int64  `json:"curnum" gorm:"default:0"`                        // curnum
	MaxNum   int64  `json:"maxnum" gorm:"default:1"`                        // maxnum

}

func (model *SysAuthInfo) CreateAuthInfo(username string, newtoken string) (*SysAuthInfo, error) {

	val := SysAuthInfo{
		Username: username,
		Token:    newtoken,
	}

	err := g_db.Create(&val).Error
	if err != nil {
		zlog.Error(err)
		return nil, err
	}

	return &val, nil
}

func (model *SysAuthInfo) DeleteAuthInfo(username string) error {

	//删除username的token
	err := g_db.Unscoped().Where("username = ?", username).Delete(&SysAuthInfo{}).Error
	if err != nil {
		zlog.Error(err)
		return err
	}

	return nil
}

// 根据更新时间删除token
func (model *SysAuthInfo) DeleteAuthInfoByUpdatedTime(delete_time int) error {

	tmp_Val := time.Duration(delete_time) * time.Second

	//删除，更新时间超过指定时间的数据
	//err := g_db.Unscoped().Debug().Where("updated_at < ?", time.Now().Add(-tmp_Val)).Delete(&AuthInfo{}).Error
	err := g_db.Unscoped().Where("updated_at < ?", time.Now().Add(-tmp_Val)).Delete(&SysAuthInfo{}).Error

	return err
}

func (model *SysAuthInfo) GetOneByToken(token string) (retVal *SysAuthInfo, err error) {
	retVal = &SysAuthInfo{}
	err = g_db.Where("token = ?", token).First(retVal).Error
	return
}

func (model *SysAuthInfo) GetOneByUsername(username string) (retVal *SysAuthInfo, err error) {
	retVal = &SysAuthInfo{}
	err = g_db.Where("username = ?", username).First(retVal).Error
	return
}

func (model *SysAuthInfo) UpdateOneByToken(token string) (retVal *SysAuthInfo, err error) {

	//更新 Index为 Index + 1
	retVal = &SysAuthInfo{}
	err = g_db.Where("token = ?", token).First(retVal).Error
	if err != nil {
		zlog.Error(err)
		return nil, err
	}
	retVal.Index = retVal.Index + 1
	err = g_db.Save(retVal).Error

	return retVal, err
}

// 取得数据数量
func (model *SysAuthInfo) GetNum() (int64, error) {

	var num int64
	err := g_db.Model(&SysAuthInfo{}).Count(&num).Error
	if err != nil {
		zlog.Error(err)
		return 0, err
	}

	return num, nil
}
