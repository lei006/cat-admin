package model

import (
	"time"

	"github.com/songzhibin97/gkit/cache/local_cache"
)

type Captcha struct {
}

func (model *Captcha) Generate(key string, width int, height int, long int, timeout time.Duration) {

}

func (model *Captcha) Verify(key string, captucha string) {
	local_cache.NewCache()
}
