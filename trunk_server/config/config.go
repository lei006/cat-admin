package config

import (
	"fmt"
	"go-admin/utils"
	"os"
	"path/filepath"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/spf13/viper"
)

var (
	Config StructConfig

	AppName       = "mapp-server"
	AppDesc       = "工作站服务器"
	IsDaemon      = true
	DefConfigFile = "config.yaml"

	RunAtVscode    bool   // 是否运行在vscode
	WorkPath       string // 工作路径
	ConfigFilePath string // 配置文件

)

type StructConfig struct {

	//
	Api ConfigApi `mapstructure:"api" json:"api" yaml:"api"`
	JWT JWT       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`

	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Logger  Logger  `mapstructure:"logger" json:"logger" yaml:"logger"`

	// gorm
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Sqlite Sqlite `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
}

func OnInit() error {

	///////////////////////////////////////////
	// 如果在vscode 则不是服务模式
	at_vscode, err := utils.RunAtVscode()
	if err != nil {
		zlog.Error("error: ")
		return err
	}

	RunAtVscode = at_vscode

	if !RunAtVscode {
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("os.Executable error " + err.Error())
		}
		WorkPath = filepath.Dir(exePath)
		ConfigFilePath = WorkPath + "/" + DefConfigFile
	} else {

		pwdPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("os.Getwd error " + err.Error())
		}
		WorkPath = pwdPath
		ConfigFilePath = WorkPath + "/" + DefConfigFile
	}

	// 加载配置-配置
	err = loadConfig()
	if err != nil {
		return err
	}

	PrintInfo()

	return nil
}

func loadConfig() error {

	///////////////////////////////////////////////
	// 1. 配置目录--支持 服务模式

	zlog.Debug("ReportCfg.RunAtVscode ", RunAtVscode)
	zlog.Debug("config： ", ConfigFilePath)

	v := viper.New()
	v.SetConfigFile(ConfigFilePath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error config file: "+err.Error(), ConfigFilePath)
	}

	if err = v.Unmarshal(&Config); err != nil {
		return fmt.Errorf("Fatal Unmarshal config file: " + err.Error())
	}

	if Config.Api.TokenExpired <= 0 {
		Config.Api.TokenExpired = 60
	}

	return nil
}

func PrintInfo() {
	zlog.Info("")
	zlog.Info("")
	zlog.Info("")
	zlog.Info("---------------------------------------------------------------------------------------------")
	zlog.Info("----                                                                                   -----")
	zlog.Info("---------------------------------------------------------------------------------------------")
	zlog.Info("--AppName", AppName)
	zlog.Info("--AppDesc", AppDesc)
	zlog.Info("--RunAtVscode", RunAtVscode)
	zlog.Info("--WorkPath", WorkPath)
	zlog.Info("--ConfigFilePath", ConfigFilePath)
	zlog.Info("--Logger.Path", Config.Logger.Path)
	zlog.Info("--Logger.SaveDay", Config.Logger.SaveDay)
	zlog.Info("--Api.AdminPassword", Config.Api.AdminPassword)
	zlog.Info("--Api.TokenExpired(秒):", Config.Api.TokenExpired)
	zlog.Info("--Api.OnlineMaxNum:", Config.Api.OnlineMaxNum)
	zlog.Info("")

}
