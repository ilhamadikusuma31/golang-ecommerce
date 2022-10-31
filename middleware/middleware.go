package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhamadikusuma31/golang-ecommerce/tokens"
)

func Autentikasi() gin.HandlerFunc {
	return func(c *gin.Context){
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ga ada izin"})
			c.Abort()
			return 
		}
		klaims,err := tokens.ValidasiToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"gagal validasi token"})
			c.Abort()
			return
		}

		//set context
		c.Set("email", klaims.Email)
		c.Set("uid", klaims.Uid)

		//ke route yang dituju
		c.Next()
	}
}