package v1

import (
	"net/http"
	"strings"

	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/config"
	tokens "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/token"

	"github.com/spf13/cast"
)

func GetIdFromToken(r *http.Request, cfg *config.Config) (string, int) {
	var softToken string
	token := r.Header.Get("Authorization")

	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		softToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		softToken = token
	}

	claims, err := tokens.ExtractClaim(softToken, []byte(cfg.Token.SignInKey))
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	resp := cast.ToString(claims["sub"])

	return resp, 200
}
