package utils

import (
	"errors"
	"math/rand"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v4"
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	NickName    string
	AuthorityId uint
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT(key string) *JWT {
	return &JWT{
		[]byte(key),
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims, bufferTime time.Duration, expiresTime time.Duration, issuer string) CustomClaims {

	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bufferTime / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},                         // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),       // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresTime)), // 过期时间 7天  配置文件
			Issuer:    issuer,                                          // 签名的发行者
		},
	}
	return claims
}

// 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
	/*
		v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
			return j.CreateToken(claims)
		})
		return v.(string), err
	*/
	return "aaaaaaaaaa", nil
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}

// RandomString 生成指定长度和字符集的随机字符串
func RandomString(length int, numbers, letters, specials bool) string {
	rand.Seed(time.Now().UnixNano())

	var charSet []rune

	if numbers {
		charSet = append(charSet, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9')
	}
	if letters {
		charSet = append(charSet, 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z')
	}
	if specials {
		charSet = append(charSet, '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '+', '=', '{', '}', '[', ']', '|', '\\', ';', ':', ',', '.', '<', '>', '/', '?', '`', '~')
	}

	if len(charSet) == 0 {
		panic("At least one character type (numbers, letters, or specials) must be enabled.")
	}

	result := make([]rune, length)
	for i := range result {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}
