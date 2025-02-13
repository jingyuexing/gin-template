package core

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"template/i18n"
	"template/internal/builtin"
	builtinerrors "template/internal/builtinErrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	successResponse = gin.H{"code": 0, "message": i18n.I18N.T("Hints.Success")}
	validatorInstance = validator.New()
)

func Response(ctx *gin.Context, message string, data any, status int, code int) {
	ctx.JSON(status, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func ResponseError(ctx *gin.Context, err error, format ...map[string]any) {
	if customErr, ok := err.(*builtinerrors.Exception); ok {
		customErr.Format(format...)
		result := gin.H{
			"code":    customErr.Code(),
			"message": customErr.Error(),
		}

		if errMsg := customErr.ErrorMessage(); len(errMsg) > 0 {
			result["errors"] = errMsg
		}

		ctx.JSON(customErr.Status(), result)
		return
	}

	defaultErr := builtin.ErrInternalServer
	ctx.JSON(defaultErr.Status(), gin.H{
		"code":    defaultErr.Code(),
		"message": defaultErr.Error(),
	})
}

func SetHeaders(ctx *gin.Context, headers map[string]string) {
	for key := range headers {
		ctx.Header(key, headers[key])
	}
}

func ResponseData(ctx *gin.Context, data any) {
	if data == nil {
		ctx.JSON(http.StatusOK, successResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": successResponse["message"],
		"data":    data,
	})
}

func ValidateParams[T any](bind T) error {
	err := validatorInstance.Struct(bind)
	if err != nil {
		parmeterErr := builtinerrors.New("Errors.ParameterIsInvalid", http.StatusBadRequest, builtin.ParameterIsInvalidCode)
		if validErr, ok := err.(validator.ValidationErrors); ok {
			for _, val := range validErr {
				trans := make(map[string]any)
				path := ""
				trans["field"] = strings.ToLower(val.Field())
				trans["val"] = val.Param()
				trans["tag"] = val.Tag()
				switch val.Tag() {
				case "lte":
					path = "Validate.LessThanOrEqual"
				case "gte":
					path = "Validate.GreatThanOrEqual"
				case "lt":
					path = "Validate.LessThan"
				case "required":
					path = "Validate.Required"
				case "min":
					path = "Validate.Min"
				case "max":
					path = "Validate.Max"
				case "oneof":
					path = "Validate.Enum"
				case "email":
					path = "Validate.IsEmail"
				case "isbn":
					path = "Validate.ISBN"
				case "html":
					path = "Validate.HTML"
				case "uuid":
					path = "Validate.UUID"
				case "md4":
					path = "Validate.MD4"
				case "md5":
					path = "Validate.MD5"
				case "cve":
					path = "Validate.CVE"
				case "country_code":
					path = "Validate.Contry"
				case "boolean":
					path = "Validate.Boolean"
				case "number", "numeric", "float":
					path = "Validate.Numberic"
				case "alpha", "alphaunicode":
					path = "Validate.Alphabet"
				case "btc_addr":
					path = "Validate.BTC"
				case "eth_addr":
					path = "Validate.ETH"
				case "hexadecimal":
					path = "Validate.HEX"
				case "semver":
					path = "Validate.Semver"
				case "credit_card":
					path = "Validate.Credit"
				case "upper_required":
					path = "Validate.UpperCaseRequired"
				case "lower_required":
					path = "Validate.LowerCaseRequired"
				case "numberic_required":
					path = "Validate.NumberRequired"
				case "special_char":
					path = "Validate.SpecialChar"
				case "common_password":
					path = "Validate.CommonPassword"
				case "datetime":
					path = "Validate.Datetime"
				default:
					path = "Validate.Default"
				}
				if (val.Tag() == "min" || val.Tag() == "max") && (val.Kind() == reflect.Array || val.Kind() == reflect.String) {
					switch val.Tag() {
					case "min":
						path = "Validate.MinLength"
					case "max":
						path = "Validate.MaxLength"
					}
				}
				parmeterErr.SetErrorMessage(
					strings.ToLower(val.Field()),
					i18n.I18N.T(path, trans),
				)
			}
		}
		return parmeterErr
	}
	return nil
}


func GetTokenFromRequest(c *gin.Context) string {
	tokenAuth := c.GetHeader("Authorization")
	tokenQuery := c.DefaultQuery("token", "")
	xToken := c.GetHeader("X-Token")

	if xToken != "" {
		return xToken
	}
	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func JSONP(ctx *gin.Context, data any) {
	callback := ctx.DefaultQuery("callback", "")

	if callback != "" {
		ctx.Header("Content-Type", "application/javascript")

		jsonData, err := json.Marshal(data)
		if err != nil {
			ResponseError(ctx,builtin.ErrInternalServer)
			return
		}

		response := callback + "(" + string(jsonData) + ");"

		ctx.String(http.StatusOK, response)
	} else {
		ResponseData(ctx, data)
	}
}
