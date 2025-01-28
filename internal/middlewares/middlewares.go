package middlewares

import (
	"fmt"
	"net/http"
	"yoga-management/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthenticateMiddleware(ctx *gin.Context) {
	// Verify if token exists
	tokenStr, err := ctx.Cookie("tkn")
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		ctx.Abort()
		return
	}

	// Verify if token is valid
	token, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		ctx.Abort()
		return
	}

	fmt.Printf("Token validated succesfully. Claims: %+v\\n", token.Claims)
	ctx.Next()
}
