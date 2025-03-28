package config

type Logger struct {
	Path    string `mapstructure:"path" json:"path" yaml:"path"`             // 环境值
	SaveDay int    `mapstructure:"save_day" json:"save_day" yaml:"save_day"` // 保存日期
	Enable  bool   `mapstructure:"enable" json:"enable" yaml:"enable"`       // 使用redis
}
