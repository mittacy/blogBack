package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/mittacy/blogBack/pkg/store/cache"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
	"time"
)

var secret string // 服务开始后请勿更改密钥，否则会导致已经注册的token无法解压

var Token *token

type token struct {
	Expire    time.Duration
	Cache     cache.CustomRedis
	BlackName string // redis token黑名单集合键名
}

type TokenData struct {
	UserId int64 `json:"userId"`
	Role   int   `json:"role"`
	jwt.StandardClaims
}

// InitToken 初始化
// @param expireHours
func InitToken(customRedis cache.CustomRedis) {
	var expire time.Duration
	if err := viper.UnmarshalKey("jwt.expire", &expire); err != nil {
		panic(fmt.Sprintf("jwt init err: %s", err))
	}

	if err := viper.UnmarshalKey("jwt.secret", &secret); err != nil {
		panic(fmt.Sprintf("jwt init err: %s", err))
	}
	if secret == "" {
		panic("jwt.secret config cannot be empty")
	}

	Token = NewToken(expire, customRedis)
}

// NewToken 生成新的token配置
// @param expireHours token过期时间，单位：小时
// @param cache redis操作句柄
// @return *token token句柄
func NewToken(expire time.Duration, customRedis cache.CustomRedis) *token {
	expire = expire * time.Second * 3600

	var serverName string
	if err := viper.UnmarshalKey("server.name", &serverName); err != nil {
		panic(fmt.Sprintf("get redis err: %s", err))
	}

	return &token{
		Expire:    expire,
		Cache:     customRedis,
		BlackName: fmt.Sprintf("%s:token:blacklist", serverName),
	}
}

// Create 生成token
// @param userId 用户id
// @param role 用户角色
// @return string token字符串
// @return error
func (ctl *token) Create(userId int64, role int) (string, error) {
	claims := jwt.MapClaims{
		"id":     userId,
		"role":   role,
		"userId": userId,
		"exp":    time.Now().Add(ctl.Expire).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// IsValid 验证token是否有效(不过期且不在redis token黑名单中)
// @param tokenString token字符串
// @return bool 是否有效
func (ctl *token) IsValid(tokenString string) bool {
	// 1. 验证token有效期
	token, _ := jwt.ParseWithClaims(tokenString, &TokenData{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if token == nil || !token.Valid {
		return false
	}
	// 2. 验证是否在黑名单
	reply, _ := ctl.Cache.Do("zscore", ctl.BlackName, tokenString)
	if reply != nil {
		return false
	}
	return true
}

// Parse 解析token
// @param tokenString
// @return *CustomClaims
func (ctl *token) Parse(tokenString string) (*TokenData, error) {
	t, _ := ctl.Cache.Do("zscore", ctl.BlackName, tokenString)
	if t != nil {
		return nil, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenData{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if token != nil {
		if claims, ok := token.Claims.(*TokenData); ok && token.Valid {
			return claims, nil
		}
	}

	if err != nil && strings.Contains(err.Error(), "expire") {
		return nil, nil
	}

	if err != nil && err.Error() == "signature is invalid" {
		return nil, nil
	}

	return nil, err
}

// GetExpireTimestamp 获取过期时间戳
// @param tokenString
// @return int64 过期时间戳
func (ctl *token) GetExpireTimestamp(tokenString string) int64 {
	token, _ := jwt.ParseWithClaims(tokenString, &TokenData{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if token != nil {
		if claims, ok := token.Claims.(*TokenData); ok && token.Valid {
			return claims.ExpiresAt
		}
	}
	return 0
}

// JoinBlackList 加入黑名单
// @param token
// @return error
func (ctl *token) JoinBlackList(token string) error {
	// 清除已经过期的token，没必要留在黑名单
	nts := time.Now().Unix()

	_, err := ctl.Cache.Do("ZREMRANGEBYSCORE", ctl.BlackName, 0, nts)
	if err != nil {
		zap.S().Errorf("redis删除过期token错误: %s", err)
	}

	ts := ctl.GetExpireTimestamp(token)
	if ts == 0 {
		return nil
	}
	_, err = ctl.Cache.Do("ZADD", ctl.BlackName, ts, token)
	if err != nil {
		return errors.Wrap(err, "存储token黑名单出错")
	}
	return nil
}
