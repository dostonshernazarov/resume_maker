package repository

import (
	"context"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
)

type Resumes interface {
	CreateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error)
	UpdateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error)
	DeleteResume(ctx context.Context, resumeID string) error
	DeleteUserResumes(ctx context.Context, userID string) error
	GetResumeByID(ctx context.Context, resumeID string) (*entity.Resume, error)
	GetUserResume(ctx context.Context, userID string, limit, offset uint64) (*entity.ListResume, error)
	ListResume(ctx context.Context, limit, offset uint64) (*entity.ListResume, error)
	GetContent(ctx context.Context, params map[string]string) (*entity.ResumeContent, error)
	GetBasic(ctx context.Context, resumeID string) (*entity.Basic, error)
	GetProfiles(ctx context.Context, resumeID string) ([]*entity.Profile, error)
	GetWorks(ctx context.Context, resumeID string) ([]*entity.Work, error)
	GetProjects(ctx context.Context, resumeID string) ([]*entity.Project, error)
	GetEducations(ctx context.Context, resumeID string) ([]*entity.Education, error)
	GetCertificates(ctx context.Context, resumeID string) ([]*entity.Certificate, error)
	GetHardSkills(ctx context.Context, resumeID string) ([]*entity.HardSkill, error)
	GetSoftSkills(ctx context.Context, resumeID string) ([]*entity.SoftSkill, error)
	GetLanguages(ctx context.Context, resumeID string) ([]*entity.Language, error)
	GetInterests(ctx context.Context, resumeID string) ([]*entity.Interest, error)
}
