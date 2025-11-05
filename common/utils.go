package common

import (
	"github.com/gin-gonic/gin"
)

const (
	WxOpenID  = "X-WX-OPENID"  // 用户ID
	WxAppID   = "X-WX-APPID"   // appID
	WxUnionID = "X-WX-UNIONID" // 用戶唯一ID
	WxEnv     = "X-WX-ENV"
)

const (
	UserID = "user_id"
	AppID  = "app_id"
	UUID   = "uuid"
	Env    = "local"
)

func HeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Request.Header.Get(WxOpenID)
		appID := c.Request.Header.Get(WxAppID)
		UUID := c.Request.Header.Get(WxUnionID)
		EnvValue := c.Request.Header.Get(WxEnv)

		c.Set(UserID, userID)
		c.Set(AppID, appID)
		c.Set(UUID, UUID)
		c.Set(Env, EnvValue)

		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	userID, exists := c.Get(UserID)
	if !exists {
		return ""
	}
	return userID.(string)
}
