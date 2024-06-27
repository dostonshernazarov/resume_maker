package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/repository"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	resumesTableName      = "resumes"
	usersTableName        = "users"
	locationTableName     = "locations"
	profileTableName      = "profiles"
	worksTableName        = "works"
	projectsTableName     = "projects"
	educationsTableName   = "educations"
	coursesTableName      = "courses"
	certificatesTableName = "certificates"
	hardSkillsTableName   = "hard_skills"
	softSkillsTableName   = "soft_skills"
	languagesTableName    = "languages"
	interestsTableName    = "interests"
)

type resumeRepo struct {
	resumeTableName       string
	usersTableName        string
	locationsTableName    string
	profileTableName      string
	worksTableName        string
	projectsTableName     string
	educationsTableName   string
	coursesTableName      string
	certificatesTableName string
	hardSkillsTableName   string
	softSkillsTableName   string
	languagesTableName    string
	interestsTableName    string
	db                    *postgres.PostgresDB
}

func NewResumeRepo(db *postgres.PostgresDB) repository.Resumes {
	return &resumeRepo{
		resumeTableName:       resumesTableName,
		usersTableName:        usersTableName,
		locationsTableName:    locationTableName,
		profileTableName:      profileTableName,
		worksTableName:        worksTableName,
		projectsTableName:     projectsTableName,
		educationsTableName:   educationsTableName,
		coursesTableName:      coursesTableName,
		certificatesTableName: certificatesTableName,
		hardSkillsTableName:   hardSkillsTableName,
		softSkillsTableName:   softSkillsTableName,
		languagesTableName:    languagesTableName,
		interestsTableName:    interestsTableName,
		db:                    db,
	}
}

func (r resumeRepo) CreateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// insert data to resume table

	resumeBuilder := r.db.Sq.Builder.Insert(r.resumeTableName)
	resumeBuilder = resumeBuilder.SetMap(map[string]interface{}{
		"id":            resume.ID,
		"user_id":       resume.UserID,
		"url":           resume.URL,
		"filename":      resume.Filename,
		"full_name":     resume.Basic.Name,
		"job_title":     resume.Basic.JobTitle,
		"summary":       resume.Basic.Summary,
		"salary":        resume.Salary,
		"job_location":  resume.JobLocation,
		"website":       resume.Basic.Website,
		"profile_image": resume.Basic.Image,
		"email":         resume.Basic.Email,
		"phone_number":  resume.Basic.PhoneNumber,
		"template":      resume.Meta.Template,
		"lang":          resume.Meta.Lang,
	})

	resumeQuery, resumeArgs, err := resumeBuilder.ToSql()
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	commandTag, err := r.db.Exec(ctx, resumeQuery, resumeArgs...)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, pgx.ErrNoRows
	}

	// insert data to locations table

	locationBuilder := r.db.Sq.Builder.Insert(r.locationsTableName)
	locationBuilder = locationBuilder.SetMap(map[string]interface{}{
		"user_id":      resume.UserID,
		"resume_id":    resume.ID,
		"city":         resume.Basic.City,
		"country_code": resume.Basic.CountryCode,
		"region":       resume.Basic.Region,
	})

	locationQuery, locationArgs, err := locationBuilder.ToSql()
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	commandTag, err = r.db.Exec(ctx, locationQuery, locationArgs...)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, pgx.ErrNoRows
	}

	// insert profiles data to profiles table

	for _, profile := range resume.Profiles {
		profileBuilder := r.db.Sq.Builder.Insert(r.profileTableName)
		profileBuilder = profileBuilder.SetMap(map[string]interface{}{
			"user_id":   resume.UserID,
			"resume_id": resume.ID,
			"network":   profile.Network,
			"username":  profile.Username,
			"url":       profile.URL,
		})

		profileQuery, profileArgs, err := profileBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, profileQuery, profileArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// insert work experience data to works table

	for _, work := range resume.Works {
		worksBuilder := r.db.Sq.Builder.Insert(r.worksTableName)
		worksBuilder = worksBuilder.SetMap(map[string]interface{}{
			"user_id":    resume.UserID,
			"resume_id":  resume.ID,
			"position":   work.Position,
			"company":    work.Company,
			"start_date": work.StartDate,
			"end_date":   work.EndDate,
			"location":   work.Location,
			"summary":    work.Summary,
			"skills":     strings.Join(work.Skills, ","),
		})

		workQuery, workArgs, err := worksBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, workQuery, workArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// insert projects data to project table

	for _, project := range resume.Projects {
		projectsBuilder := r.db.Sq.Builder.Insert(r.projectsTableName)
		projectsBuilder = projectsBuilder.SetMap(map[string]interface{}{
			"user_id":     resume.UserID,
			"resume_id":   resume.ID,
			"name":        project.Name,
			"url":         project.URL,
			"description": project.Description,
		})

		projectsQuery, projectsArgs, err := projectsBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, projectsQuery, projectsArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// insert educations data to educations table

	for _, education := range resume.Educations {
		educationsBuilder := r.db.Sq.Builder.Insert(r.educationsTableName)
		educationsBuilder = educationsBuilder.SetMap(map[string]interface{}{
			"id":          education.EducationID,
			"user_id":     resume.UserID,
			"resume_id":   resume.ID,
			"institution": education.Institution,
			"area":        education.Area,
			"location":    education.Location,
			"study_type":  education.StudyType,
			"start_date":  education.StartDate,
			"end_date":    education.EndDate,
			"score":       education.Score,
		})

		educationsQuery, educationsArgs, err := educationsBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, educationsQuery, educationsArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		// insert educations courses to courses table

		for _, course := range education.Courses {
			coursesBuilder := r.db.Sq.Builder.Insert(r.coursesTableName)
			coursesBuilder = coursesBuilder.SetMap(map[string]interface{}{
				"id":           uuid.New().String(),
				"education_id": education.EducationID,
				"course_name":  course,
			})

			coursesQuery, coursesArgs, err := coursesBuilder.ToSql()
			if err != nil {
				if err := tx.Rollback(ctx); err != nil {
					return nil, err
				}
				return nil, err
			}

			commandTag, err = r.db.Exec(ctx, coursesQuery, coursesArgs...)
			if err != nil {
				if err := tx.Rollback(ctx); err != nil {
					return nil, err
				}
				return nil, err
			}
		}
	}

	// insert certificates data to certificates table

	for _, cert := range resume.Certificates {
		certificatesBuilder := r.db.Sq.Builder.Insert(r.certificatesTableName)
		certificatesBuilder = certificatesBuilder.SetMap(map[string]interface{}{
			"user_id":   resume.UserID,
			"resume_id": resume.ID,
			"title":     cert.Title,
			"date":      cert.Date,
			"issuer":    cert.Issuer,
			"score":     cert.Score,
			"url":       cert.URL,
		})

		certificatesQuery, certificatesArgs, err := certificatesBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, certificatesQuery, certificatesArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// insert into hard skills data to hard_skills table

	for _, hard := range resume.HardSkills {
		hardSkillsBuilder := r.db.Sq.Builder.Insert(r.hardSkillsTableName)
		hardSkillsBuilder = hardSkillsBuilder.SetMap(map[string]interface{}{
			"user_id":   resume.UserID,
			"resume_id": resume.ID,
			"name":      hard.Name,
			"level":     hard.Level,
		})

		hardSkillsQuery, hardSkillsArgs, err := hardSkillsBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, hardSkillsQuery, hardSkillsArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// insert soft skills data to soft_skills table

	for _, soft := range resume.SoftSkills {
		softSkillsBuilder := r.db.Sq.Builder.Insert(r.softSkillsTableName)
		softSkillsBuilder = softSkillsBuilder.SetMap(map[string]interface{}{
			"user_id":   resume.UserID,
			"resume_id": resume.ID,
			"name":      soft.Name,
		})

		softSkillsQuery, softSkillsArgs, err := softSkillsBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, softSkillsQuery, softSkillsArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// insert languages data to languages table

	for _, lang := range resume.Languages {
		languagesBuilder := r.db.Sq.Builder.Insert(r.languagesTableName)
		languagesBuilder = languagesBuilder.SetMap(map[string]interface{}{
			"user_id":   resume.UserID,
			"resume_id": resume.ID,
			"language":  lang.Language,
			"fluency":   lang.Fluency,
		})

		languageQuery, languageArgs, err := languagesBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, languageQuery, languageArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// insert interests data to interests table

	for _, interest := range resume.Interests {
		interestBuilder := r.db.Sq.Builder.Insert(r.interestsTableName)
		interestBuilder = interestBuilder.SetMap(map[string]interface{}{
			"user_id":   resume.UserID,
			"resume_id": resume.ID,
			"name":      interest.Name,
		})

		interestQuery, interestArgs, err := interestBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, interestQuery, interestArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	return resume, nil
}

func (r resumeRepo) UpdateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// update data to resume table

	resumeBuilder := r.db.Sq.Builder.Update(r.resumeTableName)
	resumeBuilder = resumeBuilder.Where(r.db.Sq.Equal("id", resume.ID))
	resumeBuilder = resumeBuilder.SetMap(map[string]interface{}{
		"id":            resume.ID,
		"user_id":       resume.UserID,
		"url":           resume.URL,
		"filename":      resume.Filename,
		"full_name":     resume.Basic.Name,
		"job_title":     resume.Basic.JobTitle,
		"summary":       resume.Basic.Summary,
		"website":       resume.Basic.Website,
		"profile_image": resume.Basic.Image,
		"email":         resume.Basic.Email,
		"phone_number":  resume.Basic.PhoneNumber,
		"template":      resume.Meta.Template,
		"lang":          resume.Meta.Lang,
		"updated_at":    time.Now().Format(time.RFC3339),
	})

	resumeQuery, resumeArgs, err := resumeBuilder.ToSql()
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	commandTag, err := r.db.Exec(ctx, resumeQuery, resumeArgs...)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, pgx.ErrNoRows
	}

	// update data to locations table

	locationBuilder := r.db.Sq.Builder.Update(r.locationsTableName)
	locationBuilder = locationBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
	locationBuilder = locationBuilder.SetMap(map[string]interface{}{
		"user_id":      resume.UserID,
		"city":         resume.Basic.City,
		"country_code": resume.Basic.CountryCode,
		"region":       resume.Basic.Region,
		"updated_at":   time.Now().Format(time.RFC3339),
	})

	locationQuery, locationArgs, err := locationBuilder.ToSql()
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	commandTag, err = r.db.Exec(ctx, locationQuery, locationArgs...)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, pgx.ErrNoRows
	}

	// update profiles data to profiles table

	for _, profile := range resume.Profiles {
		profileBuilder := r.db.Sq.Builder.Update(r.profileTableName)
		profileBuilder = profileBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
		profileBuilder = profileBuilder.SetMap(map[string]interface{}{
			"user_id":    resume.UserID,
			"resume_id":  resume.ID,
			"network":    profile.Network,
			"username":   profile.Username,
			"url":        profile.URL,
			"updated_at": time.Now().Format(time.RFC3339),
		})

		profileQuery, profileArgs, err := profileBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, profileQuery, profileArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// update work experience data to works table

	for _, work := range resume.Works {
		worksBuilder := r.db.Sq.Builder.Update(r.worksTableName)
		worksBuilder = worksBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
		worksBuilder = worksBuilder.SetMap(map[string]interface{}{
			"user_id":    resume.UserID,
			"resume_id":  resume.ID,
			"position":   work.Position,
			"company":    work.Company,
			"start_date": work.StartDate,
			"end_date":   work.EndDate,
			"location":   work.Location,
			"summary":    work.Summary,
			"skills":     strings.Join(work.Skills, ","),
			"updated_at": time.Now().Format(time.RFC3339),
		})

		workQuery, workArgs, err := worksBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, workQuery, workArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// update projects data to project table

	for _, project := range resume.Projects {
		projectsBuilder := r.db.Sq.Builder.Update(r.projectsTableName)
		projectsBuilder = projectsBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
		projectsBuilder = projectsBuilder.SetMap(map[string]interface{}{
			"user_id":     resume.UserID,
			"resume_id":   resume.ID,
			"name":        project.Name,
			"url":         project.URL,
			"description": project.Description,
			"updated_at":  time.Now().Format(time.RFC3339),
		})

		projectsQuery, projectsArgs, err := projectsBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, projectsQuery, projectsArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// update educations data to educations table

	for _, education := range resume.Educations {
		educationsBuilder := r.db.Sq.Builder.Update(r.educationsTableName)
		educationsBuilder = educationsBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
		educationsBuilder = educationsBuilder.SetMap(map[string]interface{}{
			"id":          education.EducationID,
			"user_id":     resume.UserID,
			"resume_id":   resume.ID,
			"institution": education.Institution,
			"area":        education.Area,
			"location":    education.Location,
			"study_type":  education.StudyType,
			"start_date":  education.StartDate,
			"end_date":    education.EndDate,
			"score":       education.Score,
			"updated_at":  time.Now().Format(time.RFC3339),
		})

		educationsQuery, educationsArgs, err := educationsBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, educationsQuery, educationsArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		// update educations courses to courses table

		for _, course := range education.Courses {
			coursesBuilder := r.db.Sq.Builder.Update(r.coursesTableName)
			coursesBuilder = coursesBuilder.SetMap(map[string]interface{}{
				"id":           uuid.New().String(),
				"education_id": education.EducationID,
				"course_name":  course,
				"updated_at":   time.Now().Format(time.RFC3339),
			})
			coursesBuilder = coursesBuilder.Where(r.db.Sq.Equal("education_id", education.EducationID))

			coursesQuery, coursesArgs, err := coursesBuilder.ToSql()
			if err != nil {
				if err := tx.Rollback(ctx); err != nil {
					return nil, err
				}
				return nil, err
			}

			commandTag, err = r.db.Exec(ctx, coursesQuery, coursesArgs...)
			if err != nil {
				if err := tx.Rollback(ctx); err != nil {
					return nil, err
				}
				return nil, err
			}
		}
	}

	// update certificates data to certificates table

	for _, cert := range resume.Certificates {
		certificatesBuilder := r.db.Sq.Builder.Update(r.certificatesTableName)
		certificatesBuilder = certificatesBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
		certificatesBuilder = certificatesBuilder.SetMap(map[string]interface{}{
			"user_id":    resume.UserID,
			"resume_id":  resume.ID,
			"title":      cert.Title,
			"date":       cert.Date,
			"issuer":     cert.Issuer,
			"score":      cert.Score,
			"url":        cert.URL,
			"updated_at": time.Now().Format(time.RFC3339),
		})

		certificatesQuery, certificatesArgs, err := certificatesBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, certificatesQuery, certificatesArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// update into hard skills data to hard_skills table

	for _, hard := range resume.HardSkills {
		hardSkillsBuilder := r.db.Sq.Builder.Update(r.hardSkillsTableName)
		hardSkillsBuilder = hardSkillsBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
		hardSkillsBuilder = hardSkillsBuilder.SetMap(map[string]interface{}{
			"user_id":    resume.UserID,
			"resume_id":  resume.ID,
			"name":       hard.Name,
			"level":      hard.Level,
			"updated_at": time.Now().Format(time.RFC3339),
		})

		hardSkillsQuery, hardSkillsArgs, err := hardSkillsBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, hardSkillsQuery, hardSkillsArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// update soft skills data to soft_skills table

	for _, soft := range resume.SoftSkills {
		softSkillsBuilder := r.db.Sq.Builder.Update(r.softSkillsTableName)
		softSkillsBuilder = softSkillsBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
		softSkillsBuilder = softSkillsBuilder.SetMap(map[string]interface{}{
			"user_id":    resume.UserID,
			"resume_id":  resume.ID,
			"name":       soft.Name,
			"updated_at": time.Now().Format(time.RFC3339),
		})

		softSkillsQuery, softSkillsArgs, err := softSkillsBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, softSkillsQuery, softSkillsArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// update languages data to languages table

	for _, lang := range resume.Languages {
		languagesBuilder := r.db.Sq.Builder.Update(r.languagesTableName)
		languagesBuilder = languagesBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
		languagesBuilder = languagesBuilder.SetMap(map[string]interface{}{
			"user_id":    resume.UserID,
			"resume_id":  resume.ID,
			"language":   lang.Language,
			"fluency":    lang.Fluency,
			"updated_at": time.Now().Format(time.RFC3339),
		})

		languageQuery, languageArgs, err := languagesBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, languageQuery, languageArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// update interests data to interests table

	for _, interest := range resume.Interests {
		interestBuilder := r.db.Sq.Builder.Update(r.interestsTableName)
		interestBuilder = interestBuilder.Where(r.db.Sq.Equal("resume_id", resume.ID))
		interestBuilder = interestBuilder.SetMap(map[string]interface{}{
			"user_id":    resume.UserID,
			"resume_id":  resume.ID,
			"name":       interest.Name,
			"updated_at": time.Now().Format(time.RFC3339),
		})

		interestQuery, interestArgs, err := interestBuilder.ToSql()
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}

		commandTag, err = r.db.Exec(ctx, interestQuery, interestArgs...)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
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

	content, err := r.GetContent(ctx, map[string]string{
		"id": resumeID,
	})
	if err != nil {
		return nil, err
	}

	response.ID = resumeID
	response.UserID = content.UserID
	response.URL = content.URL
	response.Filename = content.Filename
	response.Meta.Template = content.Template
	response.Meta.Lang = content.Lang
	response.Salary = content.Salary
	response.JobLocation = content.JobLocation

	basic, err := r.GetBasic(ctx, resumeID)
	if err != nil {
		return nil, err
	}

	response.Basic.Name = basic.Name
	response.Basic.JobTitle = basic.JobTitle
	response.Basic.Summary = basic.Summary
	response.Basic.Website = basic.Website
	response.Basic.Image = basic.Image
	response.Basic.Email = basic.Email
	response.Basic.PhoneNumber = basic.PhoneNumber
	response.Basic.LocationID = basic.LocationID
	response.Basic.City = basic.City
	response.Basic.CountryCode = basic.CountryCode
	response.Basic.Region = basic.Region

	profiles, err := r.GetProfiles(ctx, resumeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	for _, profile := range profiles {
		response.Profiles = append(response.Profiles, &entity.Profile{
			ProfileID: profile.ProfileID,
			Network:   profile.Network,
			Username:  profile.Username,
			URL:       profile.URL,
		})
	}

	works, err := r.GetWorks(ctx, resumeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	for _, work := range works {
		response.Works = append(response.Works, &entity.Work{
			WorkID:    work.WorkID,
			Position:  work.Position,
			Company:   work.Company,
			StartDate: work.StartDate,
			EndDate:   work.EndDate,
			Location:  work.Location,
			Summary:   work.Summary,
			Skills:    work.Skills,
		})
	}

	projects, err := r.GetProjects(ctx, resumeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	for _, project := range projects {
		response.Projects = append(response.Projects, &entity.Project{
			ProjectID:   project.ProjectID,
			Name:        project.Name,
			URL:         project.URL,
			Description: project.Description,
		})
	}

	educations, err := r.GetEducations(ctx, resumeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	for _, education := range educations {
		response.Educations = append(response.Educations, &entity.Education{
			EducationID: education.EducationID,
			Institution: education.Institution,
			Area:        education.Area,
			StudyType:   education.StudyType,
			Location:    education.Location,
			StartDate:   education.StartDate,
			EndDate:     education.EndDate,
			Score:       education.Score,
			Courses:     education.Courses,
		})
	}

	certificates, err := r.GetCertificates(ctx, resumeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	for _, certificate := range certificates {
		response.Certificates = append(response.Certificates, &entity.Certificate{
			CertificateID: certificate.CertificateID,
			Title:         certificate.Title,
			Date:          certificate.Date,
			Issuer:        certificate.Issuer,
			Score:         certificate.Score,
			URL:           certificate.URL,
		})
	}

	hardSkills, err := r.GetHardSkills(ctx, resumeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	for _, hardSkill := range hardSkills {
		response.HardSkills = append(response.HardSkills, &entity.HardSkill{
			HardSkillID: hardSkill.HardSkillID,
			Name:        hardSkill.Name,
			Level:       hardSkill.Level,
		})
	}

	softSkills, err := r.GetSoftSkills(ctx, resumeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	for _, softSkill := range softSkills {
		response.SoftSkills = append(response.SoftSkills, &entity.SoftSkill{
			SoftSkillID: softSkill.SoftSkillID,
			Name:        softSkill.Name,
		})
	}

	languages, err := r.GetLanguages(ctx, resumeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	for _, language := range languages {
		response.Languages = append(response.Languages, &entity.Language{
			LanguageID: language.LanguageID,
			Language:   language.Language,
			Fluency:    language.Fluency,
		})
	}

	interests, err := r.GetInterests(ctx, resumeID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	for _, interest := range interests {
		response.Interests = append(response.Interests, &entity.Interest{
			InterestID: interest.InterestID,
			Name:       interest.Name,
		})
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

func (r resumeRepo) ListResume(ctx context.Context, limit, offset uint64) (*entity.ListResume, error) {
	var (
		ids      []string
		response entity.ListResume
	)

	builder := r.db.Sq.Builder.Select("id")
	builder = builder.From(r.resumeTableName)
	builder = builder.Where("deleted_at IS NULL")
	//builder = builder.Limit(limit).Offset(offset)

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

	//countBuilder := r.db.Sq.Builder.Select("COUNT(*)")
	//countBuilder = countBuilder.From(r.resumeTableName)
	//countBuilder = countBuilder.Where("deleted_at IS NULL")
	//
	//query, args, err = countBuilder.ToSql()
	//if err != nil {
	//	return nil, err
	//}
	//if err := r.db.QueryRow(ctx, query, args...).Scan(&response.TotalCount); err != nil {
	//	return nil, err
	//}

	return &response, nil
}

func (r resumeRepo) GetContent(ctx context.Context, params map[string]string) (*entity.ResumeContent, error) {
	var content entity.ResumeContent

	builder := r.db.Sq.Builder.Select("id, user_id, url, filename, template, lang, salary, job_location")
	builder = builder.From(r.resumeTableName)
	builder = builder.Where("deleted_at IS NULL")

	for k, v := range params {
		builder = builder.Where(r.db.Sq.Equal(k, v))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&content.ResumeID,
		&content.UserID,
		&content.URL,
		&content.Filename,
		&content.Template,
		&content.Lang,
		&content.Salary,
		&content.JobLocation,
	)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func (r resumeRepo) GetBasic(ctx context.Context, resumeID string) (*entity.Basic, error) {
	var (
		summary      sql.NullString
		website      sql.NullString
		profileImage sql.NullString
		response     entity.Basic
	)

	builder := r.db.Sq.Builder.Select("full_name, job_title, summary, website, profile_image, email, phone_number")
	builder = builder.From(r.resumeTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("id", resumeID))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	err = row.Scan(
		&response.Name,
		&response.JobTitle,
		&summary,
		&website,
		&profileImage,
		&response.Email,
		&response.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}

	if summary.Valid {
		response.Summary = summary.String
	}
	if website.Valid {
		response.Website = website.String
	}
	if profileImage.Valid {
		response.Image = profileImage.String
	}

	selectLocationBuilder := r.db.Sq.Builder.Select("id, city, country_code, region")
	selectLocationBuilder = selectLocationBuilder.From(r.locationsTableName)
	selectLocationBuilder = selectLocationBuilder.Where("deleted_at IS NULL")
	selectLocationBuilder = selectLocationBuilder.Where(r.db.Sq.Equal("resume_id", resumeID))

	query, args, err = selectLocationBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		city        sql.NullString
		countryCode sql.NullString
		region      sql.NullString
	)
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&response.LocationID,
		&city,
		&countryCode,
		&region,
	)
	if err != nil {
		return nil, err
	}

	if city.Valid {
		response.City = city.String
	}
	if countryCode.Valid {
		response.CountryCode = countryCode.String
	}
	if region.Valid {
		response.Region = region.String
	}

	return &response, nil
}

func (r resumeRepo) GetProfiles(ctx context.Context, resumeID string) ([]*entity.Profile, error) {
	var response []*entity.Profile

	builder := r.db.Sq.Builder.Select("id, network, username, url")
	builder = builder.From(r.profileTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("resume_id", resumeID))

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
		var profile entity.Profile

		err := rows.Scan(
			&profile.ProfileID,
			&profile.Network,
			&profile.Username,
			&profile.URL,
		)
		if err != nil {
			return nil, err
		}

		response = append(response, &profile)
	}

	return response, nil
}

func (r resumeRepo) GetWorks(ctx context.Context, resumeID string) ([]*entity.Work, error) {
	var response []*entity.Work

	builder := r.db.Sq.Builder.Select("id, position, company, start_date, end_date, location, summary, skills")
	builder = builder.From(r.worksTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("resume_id", resumeID))

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
		var (
			summary   sql.NullString
			skills    sql.NullString
			startDate sql.NullString
			endDate   sql.NullString
			work      entity.Work
		)

		err := rows.Scan(
			&work.WorkID,
			&work.Position,
			&work.Company,
			&startDate,
			&endDate,
			&work.Location,
			&summary,
			&skills,
		)
		if err != nil {
			return nil, err
		}
		if skills.Valid {
			work.Skills = strings.Split(skills.String, ",")
		}
		if summary.Valid {
			work.Summary = summary.String
		}
		if startDate.Valid {
			work.StartDate = startDate.String
		}
		if endDate.Valid {
			work.EndDate = endDate.String
		}

		response = append(response, &work)
	}

	return response, nil
}

func (r resumeRepo) GetProjects(ctx context.Context, resumeID string) ([]*entity.Project, error) {
	var response []*entity.Project

	builder := r.db.Sq.Builder.Select("id, name, url, description")
	builder = builder.From(r.projectsTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("resume_id", resumeID))

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
		var (
			url         sql.NullString
			description sql.NullString
			project     entity.Project
		)

		err := rows.Scan(&project.ProjectID, &project.Name, &url, &description)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			project.Description = description.String
		}
		if url.Valid {
			project.URL = url.String
		}

		response = append(response, &project)
	}

	return response, nil
}

func (r resumeRepo) GetEducations(ctx context.Context, resumeID string) ([]*entity.Education, error) {
	var response []*entity.Education

	builder := r.db.Sq.Builder.Select("id, institution, area, location, study_type, start_date, end_date, score")
	builder = builder.From(r.educationsTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("resume_id", resumeID))

	educationQuery, educationArgs, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, educationQuery, educationArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			startDate sql.NullString
			endDate   sql.NullString
			education entity.Education
		)
		err := rows.Scan(
			&education.EducationID,
			&education.Institution,
			&education.Area,
			&education.Location,
			&education.StudyType,
			&startDate,
			&endDate,
			&education.Score,
		)
		if err != nil {
			return nil, err
		}
		if endDate.Valid {
			education.EndDate = endDate.String
		}
		if startDate.Valid {
			education.StartDate = startDate.String
		}

		courseBuilder := r.db.Sq.Builder.Select("course_name")
		courseBuilder = courseBuilder.From(r.coursesTableName)
		courseBuilder = courseBuilder.Where("deleted_at IS NULL")
		courseBuilder = courseBuilder.Where(r.db.Sq.Equal("education_id", education.EducationID))

		query, args, err := courseBuilder.ToSql()
		if err != nil {
			return nil, err
		}

		rows, err := r.db.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var courseName string
			if err = rows.Scan(&courseName); err != nil {
				return nil, err
			}

			education.Courses = append(education.Courses, courseName)
		}
		rows.Close()

		response = append(response, &education)
	}

	return response, nil
}

func (r resumeRepo) GetCertificates(ctx context.Context, resumeID string) ([]*entity.Certificate, error) {
	var response []*entity.Certificate

	builder := r.db.Sq.Builder.Select("id, title, date, issuer, score, url")
	builder = builder.From(r.certificatesTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("resume_id", resumeID))

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
		var (
			date  sql.NullString
			score sql.NullString
			url   sql.NullString
			cert  entity.Certificate
		)

		err := rows.Scan(
			&cert.CertificateID,
			&cert.Title,
			&date,
			&cert.Issuer,
			&score,
			&url,
		)
		if err != nil {
			return nil, err
		}
		if score.Valid {
			cert.Score = score.String
		}
		if url.Valid {
			cert.URL = url.String
		}
		if date.Valid {
			cert.Date = date.String
		}

		response = append(response, &cert)
	}

	return response, nil
}

func (r resumeRepo) GetHardSkills(ctx context.Context, resumeID string) ([]*entity.HardSkill, error) {
	var response []*entity.HardSkill

	builder := r.db.Sq.Builder.Select("id, name, level")
	builder = builder.From(r.hardSkillsTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("resume_id", resumeID))

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
		var (
			level sql.NullString
			hard  entity.HardSkill
		)

		err := rows.Scan(&hard.HardSkillID, &hard.Name, &level)
		if err != nil {
			return nil, err
		}
		if level.Valid {
			hard.Level = level.String
		}

		response = append(response, &hard)
	}

	return response, nil
}

func (r resumeRepo) GetSoftSkills(ctx context.Context, resumeID string) ([]*entity.SoftSkill, error) {
	var response []*entity.SoftSkill

	builder := r.db.Sq.Builder.Select("id, name")
	builder = builder.From(r.softSkillsTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("resume_id", resumeID))

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
		var soft entity.SoftSkill

		err := rows.Scan(&soft.SoftSkillID, &soft.Name)
		if err != nil {
			return nil, err
		}

		response = append(response, &soft)
	}

	return response, nil
}

func (r resumeRepo) GetLanguages(ctx context.Context, resumeID string) ([]*entity.Language, error) {
	var response []*entity.Language

	builder := r.db.Sq.Builder.Select("id, language, fluency")
	builder = builder.From(r.languagesTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("resume_id", resumeID))

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
		var (
			fluency sql.NullString
			lang    entity.Language
		)

		err := rows.Scan(&lang.LanguageID, &lang.Language, &fluency)
		if err != nil {
			return nil, err
		}
		if fluency.Valid {
			lang.Fluency = fluency.String
		}

		response = append(response, &lang)
	}

	return response, nil
}

func (r resumeRepo) GetInterests(ctx context.Context, resumeID string) ([]*entity.Interest, error) {
	var response []*entity.Interest

	builder := r.db.Sq.Builder.Select("id, name")
	builder = builder.From(r.interestsTableName)
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(r.db.Sq.Equal("resume_id", resumeID))

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
		var interest entity.Interest

		err := rows.Scan(&interest.InterestID, &interest.Name)
		if err != nil {
			return nil, err
		}

		response = append(response, &interest)
	}

	return response, nil
}
