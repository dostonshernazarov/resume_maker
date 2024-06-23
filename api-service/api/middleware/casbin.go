package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/config"
	tokens "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/token"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type JwtRoleAuth struct {
	enforcer *casbin.Enforcer
	cfg      config.Config
}

func CheckCasbinPermission(casbin *casbin.Enforcer, cfg config.Config) gin.HandlerFunc {
	casbinHandler := &JwtRoleAuth{
		cfg:      cfg,
		enforcer: casbin,
	}

	return func(c *gin.Context) {
		allow, err := casbinHandler.CheckPermission(c)

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		if !allow {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Permission denied",
			})
		}
	}

}

func (casb *JwtRoleAuth) GetRole(c *gin.Context) (string, int) {
	var t string
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		t = strings.TrimPrefix(token, "Bearer ")
	} else {
		t = token
	}

	claims, err := tokens.ExtractClaim(t, []byte(casb.cfg.Token.SignInKey))
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}
	return cast.ToString(claims["role"]), 0
}

func (casb *JwtRoleAuth) CheckPermission(c *gin.Context) (bool, error) {

	method := c.Request.Method
	path := c.Request.URL.Path

	role, status := casb.GetRole(c)

	if role == "unauthorized" {
		allowed, err := casb.enforcer.Enforce(role, path, method)
		if err != nil {
			return false, err
		}
		return allowed, nil

	}

	if status != 0 {
		return false, errors.New(role)
	}

	allowed, err := casb.enforcer.Enforce(role, path, method)
	if err != nil {
		return false, err
	}

	return allowed, nil
}
