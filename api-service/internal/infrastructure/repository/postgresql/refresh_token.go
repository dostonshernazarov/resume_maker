package postgresql

import (
	"github.com/dostonshernazarov/resume_maker/internal/pkg/postgres"
	"github.com/dostonshernazarov/resume_maker/internal/usecase/refresh_token"
)

func NewRefreshTokenRepo(db *postgres.PostgresDB) refresh_token.RefreshTokenRepo {
	return nil
}
