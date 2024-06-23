package postgresql

import (
	"context"
	"database/sql"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/repository"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/pkg/postgres"
	"github.com/jackc/pgx/v4"
	"time"

	"github.com/Masterminds/squirrel"
)

const (
	usersTableName = "users"
)

type usersRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewUsersRepo(db *postgres.PostgresDB) repository.Users {
	return &usersRepo{
		tableName: usersTableName,
		db:        db,
	}
}

func (p *usersRepo) usersSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.
		Select(
			"id",
			"name",
			"image",
			"email",
			"phone_number",
			"refresh",
			"password",
			"role",
			"created_at",
			"updated_at",
		).From(p.tableName)
}

func (p usersRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	data := map[string]any{
		"id":           user.GUID,
		"name":         user.Name,
		"image":        user.Image,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"refresh":      user.Refresh,
		"password":     user.Password,
		"role":         user.Role,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, err
	}

	commandTag, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	return user, nil
}

func (p usersRepo) Update(ctx context.Context, users *entity.User) (*entity.User, error) {

	clauses := map[string]any{
		"name":         users.Name,
		"image":        users.Image,
		"email":        users.Email,
		"phone_number": users.PhoneNumber,
		"role":         users.Role,
		"updated_at":   users.UpdatedAt,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", users.GUID)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, err
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	user, err := p.Get(ctx, map[string]string{
		"id": users.GUID,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p usersRepo) Delete(ctx context.Context, guid string) error {
	clauses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", guid)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return err
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (p usersRepo) Get(ctx context.Context, params map[string]string) (*entity.User, error) {

	var (
		user entity.User
	)

	queryBuilder := p.usersSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		} else if key == "email" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		} else if key == "refresh" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
	}
	queryBuilder = queryBuilder.Where("deleted_at is null")
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var (
		nullImage       sql.NullString
		nullPhoneNumber sql.NullString
		nullRefresh     sql.NullString
		nullPassword    sql.NullString
	)
	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&user.GUID,
		&user.Name,
		&nullImage,
		&user.Email,
		&nullPhoneNumber,
		&nullRefresh,
		&nullPassword,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if nullImage.Valid {
		user.Image = nullImage.String
	}
	if nullPhoneNumber.Valid {
		user.PhoneNumber = nullPhoneNumber.String
	}
	if nullRefresh.Valid {
		user.Refresh = nullRefresh.String
	}
	if nullPassword.Valid {
		user.Password = nullPassword.String
	}

	return &user, nil
}

func (p usersRepo) List(ctx context.Context, limit uint64, offset uint64) (*entity.Users, error) {
	users := &entity.Users{}

	queryBuilder := p.usersSelectQueryPrefix()
	queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			user            entity.User
			nullImage       sql.NullString
			nullPhoneNumber sql.NullString
			nullRefresh     sql.NullString
			nullPassword    sql.NullString
		)
		if err = rows.Scan(
			&user.GUID,
			&user.Name,
			&nullImage,
			&user.Email,
			&nullPhoneNumber,
			&nullRefresh,
			&nullPassword,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if nullImage.Valid {
			user.Image = nullImage.String
		}
		if nullPhoneNumber.Valid {
			user.PhoneNumber = nullPhoneNumber.String
		}
		if nullRefresh.Valid {
			user.Refresh = nullRefresh.String
		}
		if nullPassword.Valid {
			user.Password = nullPassword.String
		}

		users.Users = append(users.Users, &user)
	}

	selectBuilder := p.db.Sq.Builder.Select("COUNT(*)")
	selectBuilder = selectBuilder.From(p.tableName)
	selectBuilder = selectBuilder.Where("deleted_at IS NULL")

	selectQuery, _, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var count uint64
	if err := p.db.QueryRow(ctx, selectQuery).Scan(&count); err != nil {
		return nil, err
	}
	users.Total = count

	return users, nil
}

func (p usersRepo) UniqueEmail(ctx context.Context, request *entity.IsUnique) (*entity.Response, error) {

	query := `SELECT COUNT(*) FROM users WHERE email = $1 and deleted_at is null`

	var count int
	err := p.db.QueryRow(ctx, query, request.Email).Scan(&count)
	if err != nil {
		return &entity.Response{Status: true}, err
	}
	if count != 0 {
		return &entity.Response{Status: true}, nil
	}

	return &entity.Response{Status: false}, nil
}

func (p usersRepo) UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error) {

	clauses := map[string]any{
		"refresh": request.RefreshToken,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", request.UserID)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return &entity.Response{Status: false}, err
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return &entity.Response{Status: false}, err
	}

	if commandTag.RowsAffected() == 0 {
		return &entity.Response{Status: false}, pgx.ErrNoRows
	}

	return &entity.Response{Status: true}, nil
}

func (p usersRepo) UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error) {

	clauses := map[string]any{
		"password": request.NewPassword,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", request.UserID)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return &entity.Response{Status: false}, err
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return &entity.Response{Status: false}, err
	}

	if commandTag.RowsAffected() == 0 {
		return &entity.Response{Status: false}, pgx.ErrNoRows
	}

	return &entity.Response{Status: true}, nil
}
