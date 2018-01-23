package webservice

import (
	"time"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"errors"
)
type AuthRouteType int
const (
	AuthUndefined AuthRouteType = 0
	AuthGET       AuthRouteType = 1
	AuthPOST      AuthRouteType = 2
	AuthPUT       AuthRouteType = 3
	AuthDELETE    AuthRouteType = 4
)

type authRoute struct {
	kind AuthRouteType
	path string
	f    func(c *gin.Context)
}
type JwtAuth struct {
	realm string
	key []byte
	timeout time.Duration
	maxRefresh time.Duration
	authenticatorFunc func(userId,password string,c *gin.Context) (string,bool)
	authorizatorFunc func(userId string,c *gin.Context) bool
	unauthorizedFunc func(c *gin.Context, code int, message string)
	loginPath string
	authPath string
	refreshPath string
	tokenLookUp string
	tokenHeadName string
	timeFunc func() time.Time
	routes map[string]*authRoute
}

func NewJwtAuth() *JwtAuth{
	j:=&JwtAuth{
		realm: "test zone",
		key:[]byte(uuid.New().String()),
		timeout: time.Hour,
		maxRefresh: time.Hour,
		unauthorizedFunc: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		loginPath: "/login",
		authPath: "/auth",
		refreshPath:"/refresh_token",
		tokenLookUp: "header:Authorization",
		tokenHeadName:"Bearer",
		timeFunc: time.Now,
		routes: make (map[string]*authRoute),

	}
	return j
	}



// Set Realm value, default is "test zone"
func (j *JwtAuth) SetRealm(value string) *JwtAuth {
	j.realm= value
	return j
}

// Set Key value, default will be generated automatically
func (j *JwtAuth) SetKey(value string) *JwtAuth {
	j.key= []byte(value)
	return j
}

// Set Timeout, default will one hour
func (j *JwtAuth) SetTimeout(value time.Duration) *JwtAuth {
	j.timeout= value
	return j
}

// Set Max Refresh Time, default will one hour
func (j *JwtAuth) SetMaxRefresh(value time.Duration) *JwtAuth {
	j.maxRefresh= value
	return j
}

// Set Authentication func. should be func(userId,password string,c *gin.Context) (string,bool) format
func (j *JwtAuth) SetAuthenticatorFunc(f func(userId,password string,c *gin.Context) (string,bool)) *JwtAuth {
	j.authenticatorFunc = f
	return j
}

// Set Authorizator func. should be func(userId string,c *gin.Context) (bool) format
func (j *JwtAuth) SetAuthorizatorFunc(f func(userId string,c *gin.Context) (bool)) *JwtAuth {
	j.authorizatorFunc = f
	return j
}

// Set Unauthorize func. should be func(c *gin.Context , code int, message string) format
func (j *JwtAuth) SetUnauthorizedFunc(f func(c *gin.Context , code int, message string)) *JwtAuth {
	j.unauthorizedFunc = f
	return j
}

// Set Login path, default is /login
func (j *JwtAuth) SetLoginPath(value string) *JwtAuth {
	j.loginPath= value
	return j
}
// Set Group Auth path path, default is /auth
func (j *JwtAuth) SetGroupAuthPath(value string) *JwtAuth {
	j.authPath= value
	return j
}

// Set Refresh token path, default is /refresh_token
func (j *JwtAuth) SetRefreshTokenPath(value string) *JwtAuth {
	j.refreshPath= value
	return j
}

// Set Token lookup. Default "header:Authorization"
// Options:
// TokenLookup is a string in the form of "<source>:<name>" that is used
// to extract token from the request.
// Optional. Default value "header:Authorization".
// Possible values:
// - "header:<name>"
// - "query:<name>"
// - "cookie:<name>"
func (j *JwtAuth) SetTokenLookup(value string) *JwtAuth {
	j.tokenLookUp= value
	return j
}

// Set TokenHeadName, default value is "Bearer"
// TokenHeadName is a string in the header.
func (j *JwtAuth) SetTokenHeadName(value string) *JwtAuth {
	j.tokenHeadName= value
	return j
}

// Set Time Function, default is time.Now
func (j *JwtAuth) SetTimeFunc(value func () time.Time) *JwtAuth {
	j.timeFunc= value
	return j
}
// Add route to auth group.
func (j *JwtAuth) AddRoute(kind AuthRouteType, path string, f func(c *gin.Context)) *JwtAuth {
	j.routes[path] = &authRoute{kind: kind, path: path, f: f}
	return j
}

func (j *JwtAuth) set (r *gin.Engine) (error) {
	if j.authenticatorFunc ==nil {
		return errors.New("authenitcator function is missing")
	}

	if j.authorizatorFunc ==nil {
		return errors.New("autherizator function is missing")
	}

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      j.realm,
		Key:        j.key,
		Timeout:    j.timeout,
		MaxRefresh: j.maxRefresh,
		Authenticator: j.authenticatorFunc,
		Authorizator: j.authorizatorFunc,
		Unauthorized: j.unauthorizedFunc,
		TokenLookup: j.tokenLookUp,
		TokenHeadName: j.tokenHeadName,
		TimeFunc: j.timeFunc,
	}

	r.POST(j.loginPath,authMiddleware.LoginHandler)
	auth := r.Group(j.authPath)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		for _, value := range j.routes {
			switch value.kind {
			case AuthGET:
				auth.GET(value.path, value.f)
			case AuthPOST:
				auth.POST(value.path, value.f)
			case AuthPUT:
				auth.PUT(value.path, value.f)
			case AuthDELETE:
				auth.DELETE(value.path, value.f)
			}
		}
		auth.GET(j.refreshPath, authMiddleware.RefreshHandler)
	}

	return nil
}