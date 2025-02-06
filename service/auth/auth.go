package auth

import (
	"errors"
	"strings"
	parseduration "template/common/parseDuration"
	"template/global"
	"template/internal/builtin"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (

	// Token类型
	AccessToken  = "access"
	RefreshToken = "refresh"

	// Token过期时间
	// AccessTokenExpiry  = time.Hour * 2     // 访问令牌2小时过期
	// RefreshTokenExpiry = time.Hour * 24 * 7 // 刷新令牌7天过期
)

type AuthService struct{}

type OauthClaims struct {
	Uid        uint   `json:"uid"`
	AccountId  string `json:"accountId"`
	GID        string `json:"gid"`
	Permission string `json:"role"`
	ClientID   string `json:"client_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成指定类型的token
func (auth AuthService) GenerateToken(username, tokenType, gid string, id uint, duration time.Duration) (string, error) {
	claims := &OauthClaims{
		Uid:       uint(id),
		AccountId: username,
		GID:       gid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   tokenType,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(global.Config.System.User.Sign))
	if err != nil {
		return "", builtin.ErrInternalServer
	}
	return signedToken, nil
}

// CreateAccessToken 创建访问令牌
func (auth AuthService) CreateAccessToken(username, gid string, id uint) (string, error) {
	dura := 25 * time.Minute // 访问令牌有效期
	if global.Config.System.Token.AccessTokenExpiration != "" {
		result, err := parseduration.ParseDuration(
			global.Config.System.Token.AccessTokenExpiration,
		) // 访问令牌有效期
		if err != nil {
			global.Logger.Error("访问令牌解析出现错误")

		}
		if err == nil || result != 0 {
			dura = result
		}
	}
	return auth.GenerateToken(username, AccessToken, gid, id, dura)
}

// CreateRefreshToken 创建刷新令牌
func (auth AuthService) CreateRefreshToken(username, gid string, id uint) (string, error) {
	dura := 7 * 24 * time.Hour // 刷新令牌有效期
	if global.Config.System.Token.RefreshTokenExpiration != "" {
		result, err := parseduration.ParseDuration(
			global.Config.System.Token.RefreshTokenExpiration,
		) // 刷新令牌有效期
		if err != nil {
			global.Logger.Error("刷新令牌解析出现错误")
		}
		if err == nil || result != 0 {
			dura = result
		}
	}
	return auth.GenerateToken(username, RefreshToken, gid, id, dura)
}

// Token 创建访问令牌和刷新令牌
func (auth AuthService) Token(username, gid string, id uint) (accessToken, refreshToken string, err error) {
	accessToken, err = auth.CreateAccessToken(username, gid, id)
	if err != nil {
		return "", "", builtin.ErrInternalServer
	}

	refreshToken, err = auth.CreateRefreshToken(username, gid, id)
	if err != nil {
		return "", "", builtin.ErrInternalServer
	}

	return accessToken, refreshToken, nil
}

// RemovePrefix checks if the input string has the specified prefix.
// If it does, it returns the string without the prefix; otherwise, it returns the original string.
func RemovePrefix(input, prefix string) string {
	if strings.HasPrefix(input, prefix) {
		return strings.TrimPrefix(input, prefix)
	}
	return input
}

// ParseToken 解析token
func (auth AuthService) ParseToken(tokenString string) (*OauthClaims, error) {
	if global.Config.System.Token.Type != "" {
		tokenString = RemovePrefix(tokenString, global.Config.System.Token.Type+" ")
	}
	token, err := jwt.ParseWithClaims(tokenString, &OauthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.System.User.Sign), nil
	})

	if err != nil {
		if errors.Is(err,jwt.ErrTokenExpired) {
			return nil, builtin.ErrTokenExpired
		}
		return nil, builtin.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*OauthClaims)
	if !ok || !token.Valid {
		return nil, builtin.ErrTokenInvalid
	}

	return claims, nil
}

// ValidateToken 验证token是否有效
func (auth AuthService) ValidateToken(tokenString string) bool {
	_, err := auth.ParseToken(tokenString)
	return err == nil
}