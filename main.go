package main

import (
	"fmt"

	ws "github.com/liornabat/golibs/webservice"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"runtime"

)

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

func authenticator(userId string, password string, c *gin.Context) (string, bool) {
	if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
		return userId, true
	}
	fmt.Println(userId, password)
	return userId, false
}

func authorizator(userId string, c *gin.Context) bool {
	if userId == "admin" {
		return true
	}

	return false
}

func main() {

	j := ws.NewJwtAuth().
		AddRoute(ws.AuthGET, "/hello", helloHandler).
		SetAuthenticatorFunc(authenticator).
		SetKey("asfukasdasasd-awdasda-234dasdad").
		SetAuthorizatorFunc(authorizator)
	s := ws.NewServer("8080").
		SetJwtAuth(j)

	s.Run()

	runtime.Goexit()
}
