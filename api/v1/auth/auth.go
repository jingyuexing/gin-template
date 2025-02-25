package auth

import (
	"template/core"
	"template/dto"
	"template/internal/builtin"
	"template/service"

	"github.com/gin-gonic/gin"
)


type AuthController struct{}


var authService = service.ServiceBoot.Auth
var accountService = service.ServiceBoot.Account

// LoginHandler 处理登录请求
//
// @Summary Refresh Token
// @Description refresh token
// @Param request body dto.UserLoginDTO true "login"
// @Tags auth
// @Accept json
// @Router /api/v1/auth/login [post]
func (auth AuthController) Login(ctx *gin.Context){
	param := new(dto.UserLoginDTO)
	if err := ctx.ShouldBindJSON(param); err != nil {
	core.ResponseError(ctx, builtin.ErrInvalidParams)
		return
	}
	err := core.ValidateParams(param)
	if err != nil {
		core.ResponseError(ctx, err)
		return
	}

	// 获取用户信息
    user, err := accountService.GetByEmail(param.Email)
    if err != nil {
        core.ResponseError(ctx, builtin.ErrUserNotFound)
        return
    }

    // 验证密码
    if user.Password != param.Password {
        core.ResponseError(ctx, builtin.ErrInvalidPassword)
        return
    }

    // 生成令牌
    accessToken, refreshToken, err := authService.Token(user.Email, user.GID, user.ID)
    if err != nil {
        core.ResponseError(ctx, err)
        return
    }
	ctx.Writer.Header().Set("Authorization", "Bearer "+accessToken)
	ctx.Writer.Header().Set("Refresh-Token", refreshToken)
    core.ResponseData(ctx, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "token_type":    "Bearer",
    })
}


// RefreshTokenHandler 处理刷新token请求
//
// @Summary Refresh Token
// @Description refresh token
// @Param request header string true "refresh token"
// @Tags auth
// @Accept json
// @Router /api/v1/auth/refresh_token [get]
func (auth AuthController) Refresh(ctx *gin.Context){
	refreshToken := core.GetTokenFromRequest(ctx)
    if refreshToken == "" {
        core.ResponseError(ctx, builtin.ErrTokenRequired)
        return
    }

    // 解析Token
    claims, err := authService.ParseToken(refreshToken)
    if err != nil {
        core.ResponseError(ctx, err)
        return
    }

    // 验证是否为刷新令牌
    if claims.Subject != "refresh" {
        core.ResponseError(ctx, builtin.ErrUnauthorized)
        return
    }

    // 生成新的访问令牌
    accessToken, err := authService.CreateAccessToken(claims.AccountId, claims.GID, claims.Uid)
    if err != nil {
        core.ResponseError(ctx, err)
        return
    }
	ctx.Writer.Header().Set("Authorization", "Bearer "+accessToken)
    core.ResponseData(ctx, gin.H{
        "access_token": accessToken,
        "token_type":   "Bearer",
    })
}

func (auth AuthController) AuthorizeClient(ctx *gin.Context){

}

func (auth AuthController) Oauth2Login(ctx *gin.Context){
	// TODO
}

func (auth AuthController) Oauth2Callback(ctx *gin.Context){
	// TODO
}

func (auth AuthController) Oauth2Logout(ctx *gin.Context){

}

func (auth AuthController) Oauth2Refresh(ctx *gin.Context){

}
