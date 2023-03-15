package instance

import (
	"github.com/alexperezortuno/go-webshell/tools/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		v := session.Get("instance")
		c := common.RandomStr(10)

		if v == nil {
			session.Set("instance", c)
			err := session.Save()
			if err != nil {
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": c,
		})
	}
}
