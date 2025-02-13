package auth

import (
	"errors"
	"strings"
	parseduration "template/common/parseDuration"
	"template/common/utils"
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

type AuthService struct{
    tokenSign []byte
    accessTokenDuration  time.Duration
    refreshTokenDuration time.Duration
}

var BuiltInAuth = new(AuthService)

type OauthClaims struct {
	Uid        uint   `json:"uid"`
	AccountId  string `json:"accountId"`
	GID        string `json:"gid"`
	Permission string `json:"role"`
	ClientID   string `json:"client_id"`
	jwt.RegisteredClaims
}

// New creates a new AuthService instance
func New() *AuthService {
    auth := &AuthService{
        tokenSign: []byte(global.Config.System.User.Sign),
        accessTokenDuration:  25 * time.Minute,
        refreshTokenDuration: 7 * 24 * time.Hour,
    }

    // 初始化配置的过期时间
    if global.Config.System.Token.AccessTokenExpiration != "" {
        if duration, err := parseduration.ParseDuration(global.Config.System.Token.AccessTokenExpiration); err == nil {
            auth.accessTokenDuration = duration
        }
    }

    if global.Config.System.Token.RefreshTokenExpiration != "" {
        if duration, err := parseduration.ParseDuration(global.Config.System.Token.RefreshTokenExpiration); err == nil {
            auth.refreshTokenDuration = duration
        }
    }

    return auth
}

// GenerateToken optimization
func (auth *AuthService) GenerateToken(username, tokenType, gid string, id uint, duration time.Duration) (string, error) {
    claims := &OauthClaims{
        Uid:       id,
        AccountId: username,
        GID:      gid,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Subject:   tokenType,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(auth.tokenSign)
}

// CreateAccessToken optimization
func (auth *AuthService) CreateAccessToken(username, gid string, id uint) (string, error) {
    return auth.GenerateToken(username, AccessToken, gid, id, auth.accessTokenDuration)
}

// CreateRefreshToken optimization
func (auth *AuthService) CreateRefreshToken(username, gid string, id uint) (string, error) {
    return auth.GenerateToken(username, RefreshToken, gid, id, auth.refreshTokenDuration)
}

// Token optimization
func (auth *AuthService) Token(username, gid string, id uint) (accessToken, refreshToken string, err error) {
    if accessToken, err = auth.CreateAccessToken(username, gid, id); err != nil {
        return "", "", err
    }

    if refreshToken, err = auth.CreateRefreshToken(username, gid, id); err != nil {
        return "", "", err
    }

    return
}

// RemovePrefix checks if the input string has the specified prefix.
// If it does, it returns the string without the prefix; otherwise, it returns the original string.
func RemovePrefix(input, prefix string) string {
	if strings.HasPrefix(input, prefix) {
		return strings.TrimPrefix(input, prefix)
	}
	return input
}

// ParseToken optimization
func (auth *AuthService) ParseToken(tokenString string) (*OauthClaims, error) {
    if global.Config.System.Token.Type != "" {
        tokenString = RemovePrefix(tokenString, global.Config.System.Token.Type+" ")
    }

    token, err := jwt.ParseWithClaims(tokenString, &OauthClaims{}, func(token *jwt.Token) (interface{}, error) {
        return auth.tokenSign, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, builtin.ErrTokenExpired
        }
        return nil, builtin.ErrTokenInvalid
    }

    claims, ok := token.Claims.(*OauthClaims)
    return claims, utils.BoolToError(ok && token.Valid, builtin.ErrTokenInvalid)
}

// ValidateToken 验证token是否有效
func (auth AuthService) ValidateToken(tokenString string) bool {
	_, err := auth.ParseToken(tokenString)
	return err == nil
}
