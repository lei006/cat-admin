package config

type ConfigApi struct {
	Enable          bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	Debug           bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	Addr            string `mapstructure:"addr" json:"addr" yaml:"addr"`                                        // 端口值
	UploadPrintPath string `mapstructure:"upload-print-path" json:"upload-print-path" yaml:"upload-print-path"` // 上传文件路径
	TokenExpired    int    `mapstructure:"token-expired" json:"token-expired" yaml:"token-expired"`             // token过期时间(秒)
	AdminPassword   string `mapstructure:"admin-password" json:"admin-password" yaml:"admin-password"`          // admin密码
	OnlineMaxNum    int    `mapstructure:"online-max-num" json:"online-max-num" yaml:"online-max-num"`          // 在线最大数量
}

type Captcha struct {
	Enable             bool `mapstructure:"enable" json:"enable" yaml:"enable"`
	KeyLong            int  `mapstructure:"key-long" json:"key-long" yaml:"key-long"`                                     // 验证码长度
	ImgWidth           int  `mapstructure:"img-width" json:"img-width" yaml:"img-width"`                                  // 验证码宽度
	ImgHeight          int  `mapstructure:"img-height" json:"img-height" yaml:"img-height"`                               // 验证码高度
	OpenCaptcha        int  `mapstructure:"open-captcha" json:"open-captcha" yaml:"open-captcha"`                         // 防爆破验证码开启此数，0代表每次登录都需要验证码，其他数字代表错误密码此数，如3代表错误三次后出现验证码
	OpenCaptchaTimeOut int  `mapstructure:"open-captcha-timeout" json:"open-captcha-timeout" yaml:"open-captcha-timeout"` // 防爆破验证码超时时间，单位：s(秒)
}
