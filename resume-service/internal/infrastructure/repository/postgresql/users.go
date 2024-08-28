package postgresql

import (
	"context"
	"time"

	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/repository"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

const (
	resumesTableName = "resumes"
	usersTableName   = "users"
)

type resumeRepo struct {
	resumeTableName string
	usersTableName  string
	db              *postgres.PostgresDB
}

func NewResumeRepo(db *postgres.PostgresDB) repository.Resumes {
	return &resumeRepo{
		resumeTableName: resumesTableName,
		usersTableName:  usersTableName,
		db:              db,
	}
}

func (r resumeRepo) CreateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error) {
	resumeBuilder := r.db.Sq.Builder.Insert(r.resumeTableName)
	resumeBuilder = resumeBuilder.SetMap(map[string]interface{}{
		"id":           resume.ID,
		"user_id":      resume.UserID,
		"url":          resume.URL,
		"salary":       resume.Salary,
		"job_title":    resume.JobTitle,
		"region":       resume.Region,
		"job_location": resume.JobLocation,
		"job_type":     resume.JobType,
		"experience":   resume.Experience,
		"template":     resume.Template,
	})

	resumeQuery, resumeArgs, err := resumeBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	commandTag, err := r.db.Exec(ctx, resumeQuery, resumeArgs...)
	if err != nil {
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	return resume, nil
}

func (r resumeRepo) UpdateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error) {
	resumeBuilder := r.db.Sq.Builder.Update(r.resumeTableName)
	resumeBuilder = resumeBuilder.Where(r.db.Sq.Equal("id", resume.ID))
	resumeBuilder = resumeBuilder.SetMap(map[string]interface{}{
		"id":           resume.ID,
		"user_id":      resume.UserID,
		"url":          resume.URL,
		"salary":       resume.Salary,
		"job_title":    resume.JobTitle,
		"region":       resume.Region,
		"job_location": resume.JobLocation,
		"job_type":     resume.JobType,
		"experience":   resume.Experience,
		"template":     resume.Template,
		"updated_at":   time.Now().Format(time.RFC3339),
	})

	resumeQuery, resumeArgs, err := resumeBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	commandTag, err := r.db.Exec(ctx, resumeQuery, resumeArgs...)
	if err != nil {
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	return resume, nil
}

func (r resumeRepo) DeleteResume(ctx context.Context, resumeID string) error {
	clauses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	builder := r.db.Sq.Builder.Update(r.resumeTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("id", resumeID))
	builder = builder.SetMap(clauses)

	resumeQuery, resumeArgs, err := builder.ToSql()
	if err != nil {
		return err
	}

	commandTag, err := r.db.Exec(ctx, resumeQuery, resumeArgs...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r resumeRepo) DeleteUserResumes(ctx context.Context, userID string) error {
	clauses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	builder := r.db.Sq.Builder.Update(r.resumeTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("user_id", userID))
	builder = builder.SetMap(clauses)

	resumeQuery, resumeArgs, err := builder.ToSql()
	if err != nil {
		return err
	}

	commandTag, err := r.db.Exec(ctx, resumeQuery, resumeArgs...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r resumeRepo) GetResumeByID(ctx context.Context, resumeID string) (*entity.Resume, error) {
	var response entity.Resume

	builder := r.db.Sq.Builder.Select("id, user_id, url, salary, job_title, region, job_location, job_type, experience, template")
	builder = builder.From(resumesTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("id", resumeID))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&response.ID,
		&response.UserID,
		&response.URL,
		&response.Salary,
		&response.JobTitle,
		&response.Region,
		&response.JobLocation,
		&response.JobType,
		&response.Experience,
		&response.Template,
	)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r resumeRepo) GetUserResume(ctx context.Context, userID string, limit, offset uint64) (*entity.ListResume, error) {
	var (
		ids      []string
		response entity.ListResume
	)

	builder := r.db.Sq.Builder.Select("id")
	builder = builder.From(r.resumeTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("user_id", userID))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	for _, id := range ids {
		resume, err := r.GetResumeByID(ctx, id)
		if err != nil {
			return nil, err
		}

		response.Resumes = append(response.Resumes, resume)
	}

	return &response, nil
}

func (r resumeRepo) ListResume(ctx context.Context, request *entity.ListRequest) (*entity.ListResume, error) {
	var (
		total    int64
		response entity.ListResume
	)
	offset := request.Limit * (request.Page - 1)

	builder := r.db.Sq.Builder.Select("id, user_id, url, salary, job_title, region, job_location, job_type, experience, template")
	builder = builder.From(r.resumeTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.ILike("job_title", "%"+request.JobTitle+"%"))
	builder = builder.Where(r.db.Sq.ILike("job_location", "%"+request.JobLocation+"%"))
	builder = builder.Where(r.db.Sq.ILike("job_type", "%"+request.JobType+"%"))
	builder = builder.Where(r.db.Sq.ILike("region", "%"+request.Region+"%"))
	builder = builder.Limit(uint64(request.Limit))
	builder = builder.Offset(uint64(offset))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var resume entity.Resume
		err = rows.Scan(
			&resume.ID,
			&resume.UserID,
			&resume.URL,
			&resume.Salary,
			&resume.JobTitle,
			&resume.Region,
			&resume.JobLocation,
			&resume.JobType,
			&resume.Experience,
			&resume.Template,
		)
		if err != nil {
			return nil, err
		}

		response.Resumes = append(response.Resumes, &resume)
	}

	totalBuilder := r.db.Sq.Builder.Select("COUNT(*)")
	totalBuilder = totalBuilder.From(r.resumeTableName)
	totalBuilder = totalBuilder.Where("deleted_at IS NULL")
	totalBuilder = totalBuilder.Where(r.db.Sq.ILike("job_title", "%"+request.JobTitle+"%"))
	totalBuilder = totalBuilder.Where(r.db.Sq.ILike("job_location", "%"+request.JobLocation+"%"))
	totalBuilder = totalBuilder.Where(r.db.Sq.ILike("job_type", "%"+request.JobType+"%"))
	builder = builder.Where(r.db.Sq.ILike("region", "%"+request.Region+"%"))
	totalQuery, args, err := totalBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(ctx, totalQuery, args...).Scan(&total); err != nil {
		return nil, err
	}
	response.TotalCount = uint64(total)

	return &response, nil
}
