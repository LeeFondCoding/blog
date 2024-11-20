package middlewares

import (
	"blog/controller"
	"blog/pkg/jwt"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// 基于JWT的中间件认证
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseErrorWithMsg(c, controller.CodeInvalidToken,
				"请求头缺少Auth Token")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseErrorWithMsg(c, controller.CodeInvalidToken,
				"请求头token格式有误")
			c.Abort()
			return
		}

		myClaims, err := jwt.ParseToken(parts[1])
		if err != nil {
			fmt.Println(err)
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}

		c.Set(controller.ContextUserIDKey, myClaims.UserID)
		c.Next()
	}
}
