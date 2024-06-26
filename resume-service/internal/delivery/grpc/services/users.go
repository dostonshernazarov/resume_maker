package services

import (
	"context"
	"github.com/dostonshernazarov/resume_maker/user-service/genproto/resume_service"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/entity"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/infrastructure/grpc_service_clients"
	"github.com/dostonshernazarov/resume_maker/user-service/internal/usecase"
	"go.uber.org/zap"
)

type resumeRPC struct {
	logger        *zap.Logger
	resumeUseCase usecase.Resume
	client        grpc_service_clients.ServiceClients
}

func NewRPC(logger *zap.Logger, resumeUseCase usecase.Resume, client *grpc_service_clients.ServiceClients) resume_service.ResumeServiceServer {
	return &resumeRPC{
		logger:        logger,
		resumeUseCase: resumeUseCase,
		client:        *client,
	}
}

func (s resumeRPC) CreateResume(ctx context.Context, in *resume_service.Resume) (*resume_service.ResumeWithID, error) {
	var (
		entityProfiles     []*entity.Profile
		entityWorks        []*entity.Work
		entityProjects     []*entity.Project
		entityEducations   []*entity.Education
		entityCertificates []*entity.Certificate
		entityHardSkills   []*entity.HardSkill
		entitySoftSkills   []*entity.SoftSkill
		entityLanguages    []*entity.Language
		entityInterests    []*entity.Interest
	)

	for _, profile := range in.Profiles {
		entityProfiles = append(entityProfiles, &entity.Profile{
			ProfileID: profile.ProfileId,
			Network:   profile.Network,
			Username:  profile.Username,
			URL:       profile.Url,
		})
	}

	for _, work := range in.Works {
		entityWorks = append(entityWorks, &entity.Work{
			WorkID:    work.WorkId,
			Position:  work.Position,
			Company:   work.Company,
			StartDate: work.StartDate,
			EndDate:   work.EndDate,
			Location:  work.Location,
			Summary:   work.Summary,
			Skills:    work.Skills,
		})
	}

	for _, project := range in.Projects {
		entityProjects = append(entityProjects, &entity.Project{
			ProjectID:   project.ProjectId,
			Name:        project.Name,
			URL:         project.Url,
			Description: project.Description,
		})
	}

	for _, education := range in.Educations {
		entityEducations = append(entityEducations, &entity.Education{
			EducationID: education.EducationId,
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

	for _, cert := range in.Certificates {
		entityCertificates = append(entityCertificates, &entity.Certificate{
			CertificateID: cert.CertificateId,
			Title:         cert.Title,
			Date:          cert.Date,
			Issuer:        cert.Issuer,
			Score:         cert.Score,
			URL:           cert.Url,
		})
	}

	for _, hard := range in.HardSkills {
		entityHardSkills = append(entityHardSkills, &entity.HardSkill{
			HardSkillID: hard.HardSkillId,
			Name:        hard.Name,
			Level:       hard.Level,
		})
	}

	for _, soft := range in.SoftSkills {
		entitySoftSkills = append(entitySoftSkills, &entity.SoftSkill{
			SoftSkillID: soft.SoftSkillId,
			Name:        soft.Name,
		})
	}

	for _, language := range in.Languages {
		entityLanguages = append(entityLanguages, &entity.Language{
			LanguageID: language.LanguageId,
			Language:   language.Language,
			Fluency:    language.Fluency,
		})
	}

	for _, interest := range in.Interests {
		entityInterests = append(entityInterests, &entity.Interest{
			InterestID: interest.InterestId,
			Name:       interest.Name,
		})
	}

	resume, err := s.resumeUseCase.CreateResume(ctx, &entity.Resume{
		ID:          in.Id,
		UserID:      in.UserId,
		URL:         in.Url,
		Filename:    in.Filename,
		Salary:      int64(in.Salary),
		JobLocation: in.JobLocation,
		Basic: entity.Basic{
			Name:        in.Basic.Name,
			JobTitle:    in.Basic.JobTitle,
			Image:       in.Basic.Image,
			Email:       in.Basic.Email,
			PhoneNumber: in.Basic.PhoneNumber,
			Website:     in.Basic.Website,
			Summary:     in.Basic.Summary,
			LocationID:  in.Basic.LocationId,
			City:        in.Basic.City,
			CountryCode: in.Basic.CountryCode,
			Region:      in.Basic.Region,
		},
		Profiles:     entityProfiles,
		Works:        entityWorks,
		Projects:     entityProjects,
		Educations:   entityEducations,
		Certificates: entityCertificates,
		HardSkills:   entityHardSkills,
		SoftSkills:   entitySoftSkills,
		Languages:    entityLanguages,
		Interests:    entityInterests,
		Meta: entity.Meta{
			Template: in.Meta.Template,
			Lang:     in.Meta.Lang,
		},
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &resume_service.ResumeWithID{
		ResumeId: resume.ID,
	}, nil
}

func (s resumeRPC) UpdateResume(ctx context.Context, in *resume_service.Resume) (*resume_service.Resume, error) {
	entityResume := &entity.Resume{}

	for _, profile := range in.Profiles {
		entityResume.Profiles = append(entityResume.Profiles, &entity.Profile{
			ProfileID: profile.ProfileId,
			Network:   profile.Network,
			Username:  profile.Username,
			URL:       profile.Url,
		})
	}

	for _, work := range in.Works {
		entityResume.Works = append(entityResume.Works, &entity.Work{
			WorkID:    work.WorkId,
			Position:  work.Position,
			Company:   work.Company,
			StartDate: work.StartDate,
			EndDate:   work.EndDate,
			Location:  work.Location,
			Summary:   work.Summary,
			Skills:    work.Skills,
		})
	}

	for _, project := range in.Projects {
		entityResume.Projects = append(entityResume.Projects, &entity.Project{
			ProjectID:   project.ProjectId,
			Name:        project.Name,
			URL:         project.Url,
			Description: project.Description,
		})
	}

	for _, education := range in.Educations {
		entityResume.Educations = append(entityResume.Educations, &entity.Education{
			EducationID: education.EducationId,
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

	for _, cert := range in.Certificates {
		entityResume.Certificates = append(entityResume.Certificates, &entity.Certificate{
			CertificateID: cert.CertificateId,
			Title:         cert.Title,
			Date:          cert.Date,
			Issuer:        cert.Issuer,
			Score:         cert.Score,
			URL:           cert.Url,
		})
	}

	for _, hard := range in.HardSkills {
		entityResume.HardSkills = append(entityResume.HardSkills, &entity.HardSkill{
			HardSkillID: hard.HardSkillId,
			Name:        hard.Name,
			Level:       hard.Level,
		})
	}

	for _, soft := range in.SoftSkills {
		entityResume.SoftSkills = append(entityResume.SoftSkills, &entity.SoftSkill{
			SoftSkillID: soft.SoftSkillId,
			Name:        soft.Name,
		})
	}

	for _, language := range in.Languages {
		entityResume.Languages = append(entityResume.Languages, &entity.Language{
			LanguageID: language.LanguageId,
			Language:   language.Language,
			Fluency:    language.Fluency,
		})
	}

	for _, interest := range in.Interests {
		entityResume.Interests = append(entityResume.Interests, &entity.Interest{
			InterestID: interest.InterestId,
			Name:       interest.Name,
		})
	}

	resume, err := s.resumeUseCase.UpdateResume(ctx, &entity.Resume{
		ID:       in.Id,
		UserID:   in.UserId,
		URL:      in.Url,
		Filename: in.Filename,
		Basic: entity.Basic{
			Name:        in.Basic.Name,
			JobTitle:    in.Basic.JobTitle,
			Image:       in.Basic.Image,
			Email:       in.Basic.Email,
			PhoneNumber: in.Basic.PhoneNumber,
			Website:     in.Basic.Website,
			Summary:     in.Basic.Summary,
			LocationID:  in.Basic.LocationId,
			City:        in.Basic.City,
			CountryCode: in.Basic.CountryCode,
			Region:      in.Basic.Region,
		},
		Profiles:     entityResume.Profiles,
		Works:        entityResume.Works,
		Projects:     entityResume.Projects,
		Educations:   entityResume.Educations,
		Certificates: entityResume.Certificates,
		HardSkills:   entityResume.HardSkills,
		SoftSkills:   entityResume.SoftSkills,
		Languages:    entityResume.Languages,
		Interests:    entityResume.Interests,
		Meta: entity.Meta{
			Template: in.Meta.Template,
			Lang:     in.Meta.Lang,
		},
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	protoResume := &resume_service.Resume{}

	for _, profile := range resume.Profiles {
		protoResume.Profiles = append(protoResume.Profiles, &resume_service.Profile{
			ProfileId: profile.ProfileID,
			Network:   profile.Network,
			Username:  profile.Username,
			Url:       profile.URL,
		})
	}

	for _, work := range resume.Works {
		protoResume.Works = append(protoResume.Works, &resume_service.Work{
			WorkId:    work.WorkID,
			Position:  work.Position,
			Company:   work.Company,
			StartDate: work.StartDate,
			EndDate:   work.EndDate,
			Location:  work.Location,
			Summary:   work.Summary,
			Skills:    work.Skills,
		})
	}

	for _, project := range resume.Projects {
		protoResume.Projects = append(protoResume.Projects, &resume_service.Project{
			ProjectId:   project.ProjectID,
			Name:        project.Name,
			Url:         project.URL,
			Description: project.Description,
		})
	}

	for _, education := range resume.Educations {
		protoResume.Educations = append(protoResume.Educations, &resume_service.Education{
			EducationId: education.EducationID,
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

	for _, cert := range resume.Certificates {
		protoResume.Certificates = append(protoResume.Certificates, &resume_service.Certificate{
			CertificateId: cert.CertificateID,
			Title:         cert.Title,
			Date:          cert.Date,
			Issuer:        cert.Issuer,
			Score:         cert.Score,
			Url:           cert.URL,
		})
	}

	for _, hard := range resume.HardSkills {
		protoResume.HardSkills = append(protoResume.HardSkills, &resume_service.HardSkill{
			HardSkillId: hard.HardSkillID,
			Name:        hard.Name,
			Level:       hard.Level,
		})
	}

	for _, soft := range resume.SoftSkills {
		protoResume.SoftSkills = append(protoResume.SoftSkills, &resume_service.SoftSkill{
			SoftSkillId: soft.SoftSkillID,
			Name:        soft.Name,
		})
	}

	for _, language := range resume.Languages {
		protoResume.Languages = append(protoResume.Languages, &resume_service.Language{
			LanguageId: language.LanguageID,
			Language:   language.Language,
			Fluency:    language.Fluency,
		})
	}

	for _, interest := range resume.Interests {
		protoResume.Interests = append(protoResume.Interests, &resume_service.Interest{
			InterestId: interest.InterestID,
			Name:       interest.Name,
		})
	}

	return &resume_service.Resume{
		Id:       resume.ID,
		UserId:   resume.UserID,
		Url:      resume.URL,
		Filename: resume.Filename,
		Basic: &resume_service.Basic{
			Name:        resume.Basic.Name,
			JobTitle:    resume.Basic.JobTitle,
			Image:       resume.Basic.Image,
			Email:       resume.Basic.Email,
			PhoneNumber: resume.Basic.PhoneNumber,
			Website:     resume.Basic.Website,
			Summary:     resume.Basic.Summary,
			LocationId:  resume.Basic.LocationID,
			City:        resume.Basic.City,
			CountryCode: resume.Basic.CountryCode,
			Region:      resume.Basic.Region,
		},
		Profiles:     protoResume.Profiles,
		Works:        protoResume.Works,
		Projects:     protoResume.Projects,
		Educations:   protoResume.Educations,
		Certificates: protoResume.Certificates,
		HardSkills:   protoResume.HardSkills,
		SoftSkills:   protoResume.SoftSkills,
		Languages:    protoResume.Languages,
		Interests:    protoResume.Interests,
		Meta: &resume_service.Meta{
			Template: resume.Meta.Template,
			Lang:     resume.Meta.Lang,
		},
	}, nil
}

func (s resumeRPC) DeleteResume(ctx context.Context, in *resume_service.ResumeWithID) (*resume_service.Status, error) {
	if err := s.resumeUseCase.DeleteResume(ctx, in.ResumeId); err != nil {
		s.logger.Error(err.Error())
		return &resume_service.Status{Action: false}, err
	}

	return &resume_service.Status{Action: true}, nil
}

func (s resumeRPC) DeleteUserResume(ctx context.Context, in *resume_service.UserWithID) (*resume_service.Status, error) {
	if err := s.resumeUseCase.DeleteUserResumes(ctx, in.UserId); err != nil {
		s.logger.Error(err.Error())
		return &resume_service.Status{Action: false}, err
	}

	return &resume_service.Status{Action: true}, nil
}

func (s resumeRPC) GetResumeByID(ctx context.Context, in *resume_service.ResumeWithID) (*resume_service.Resume, error) {
	resume, err := s.resumeUseCase.GetResumeByID(ctx, in.ResumeId)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	protoResume := &resume_service.Resume{}

	for _, profile := range resume.Profiles {
		protoResume.Profiles = append(protoResume.Profiles, &resume_service.Profile{
			ProfileId: profile.ProfileID,
			Network:   profile.Network,
			Username:  profile.Username,
			Url:       profile.URL,
		})
	}

	for _, work := range resume.Works {
		protoResume.Works = append(protoResume.Works, &resume_service.Work{
			WorkId:    work.WorkID,
			Position:  work.Position,
			Company:   work.Company,
			StartDate: work.StartDate,
			EndDate:   work.EndDate,
			Location:  work.Location,
			Summary:   work.Summary,
			Skills:    work.Skills,
		})
	}

	for _, project := range resume.Projects {
		protoResume.Projects = append(protoResume.Projects, &resume_service.Project{
			ProjectId:   project.ProjectID,
			Name:        project.Name,
			Url:         project.URL,
			Description: project.Description,
		})
	}

	for _, education := range resume.Educations {
		protoResume.Educations = append(protoResume.Educations, &resume_service.Education{
			EducationId: education.EducationID,
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

	for _, cert := range resume.Certificates {
		protoResume.Certificates = append(protoResume.Certificates, &resume_service.Certificate{
			CertificateId: cert.CertificateID,
			Title:         cert.Title,
			Date:          cert.Date,
			Issuer:        cert.Issuer,
			Score:         cert.Score,
			Url:           cert.URL,
		})
	}

	for _, hard := range resume.HardSkills {
		protoResume.HardSkills = append(protoResume.HardSkills, &resume_service.HardSkill{
			HardSkillId: hard.HardSkillID,
			Name:        hard.Name,
			Level:       hard.Level,
		})
	}

	for _, soft := range resume.SoftSkills {
		protoResume.SoftSkills = append(protoResume.SoftSkills, &resume_service.SoftSkill{
			SoftSkillId: soft.SoftSkillID,
			Name:        soft.Name,
		})
	}

	for _, language := range resume.Languages {
		protoResume.Languages = append(protoResume.Languages, &resume_service.Language{
			LanguageId: language.LanguageID,
			Language:   language.Language,
			Fluency:    language.Fluency,
		})
	}

	for _, interest := range resume.Interests {
		protoResume.Interests = append(protoResume.Interests, &resume_service.Interest{
			InterestId: interest.InterestID,
			Name:       interest.Name,
		})
	}

	return &resume_service.Resume{
		Id:          resume.ID,
		UserId:      resume.UserID,
		Url:         resume.URL,
		Filename:    resume.Filename,
		Salary:      uint64(resume.Salary),
		JobLocation: resume.JobLocation,
		Basic: &resume_service.Basic{
			Name:        resume.Basic.Name,
			JobTitle:    resume.Basic.JobTitle,
			Image:       resume.Basic.Image,
			Email:       resume.Basic.Email,
			PhoneNumber: resume.Basic.PhoneNumber,
			Website:     resume.Basic.Website,
			Summary:     resume.Basic.Summary,
			LocationId:  resume.Basic.LocationID,
			City:        resume.Basic.City,
			CountryCode: resume.Basic.CountryCode,
			Region:      resume.Basic.Region,
		},
		Profiles:     protoResume.Profiles,
		Works:        protoResume.Works,
		Projects:     protoResume.Projects,
		Educations:   protoResume.Educations,
		Certificates: protoResume.Certificates,
		HardSkills:   protoResume.HardSkills,
		SoftSkills:   protoResume.SoftSkills,
		Languages:    protoResume.Languages,
		Interests:    protoResume.Interests,
		Meta: &resume_service.Meta{
			Template: resume.Meta.Template,
			Lang:     resume.Meta.Lang,
		},
	}, nil
}

func (s resumeRPC) GetUserResume(ctx context.Context, in *resume_service.UserWithID) (*resume_service.ListResumeResponse, error) {
	offset := in.Limit * (in.Page - 1)
	resumes, err := s.resumeUseCase.GetUserResume(ctx, in.UserId, in.Limit, offset)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	response := &resume_service.ListResumeResponse{}

	for _, resume := range resumes.Resumes {
		protoResume := &resume_service.Resume{}

		for _, profile := range resume.Profiles {
			protoResume.Profiles = append(protoResume.Profiles, &resume_service.Profile{
				ProfileId: profile.ProfileID,
				Network:   profile.Network,
				Username:  profile.Username,
				Url:       profile.URL,
			})
		}

		for _, work := range resume.Works {
			protoResume.Works = append(protoResume.Works, &resume_service.Work{
				WorkId:    work.WorkID,
				Position:  work.Position,
				Company:   work.Company,
				StartDate: work.StartDate,
				EndDate:   work.EndDate,
				Location:  work.Location,
				Summary:   work.Summary,
				Skills:    work.Skills,
			})
		}

		for _, project := range resume.Projects {
			protoResume.Projects = append(protoResume.Projects, &resume_service.Project{
				ProjectId:   project.ProjectID,
				Name:        project.Name,
				Url:         project.URL,
				Description: project.Description,
			})
		}

		for _, education := range resume.Educations {
			protoResume.Educations = append(protoResume.Educations, &resume_service.Education{
				EducationId: education.EducationID,
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

		for _, cert := range resume.Certificates {
			protoResume.Certificates = append(protoResume.Certificates, &resume_service.Certificate{
				CertificateId: cert.CertificateID,
				Title:         cert.Title,
				Date:          cert.Date,
				Issuer:        cert.Issuer,
				Score:         cert.Score,
				Url:           cert.URL,
			})
		}

		for _, hard := range resume.HardSkills {
			protoResume.HardSkills = append(protoResume.HardSkills, &resume_service.HardSkill{
				HardSkillId: hard.HardSkillID,
				Name:        hard.Name,
				Level:       hard.Level,
			})
		}

		for _, soft := range resume.SoftSkills {
			protoResume.SoftSkills = append(protoResume.SoftSkills, &resume_service.SoftSkill{
				SoftSkillId: soft.SoftSkillID,
				Name:        soft.Name,
			})
		}

		for _, language := range resume.Languages {
			protoResume.Languages = append(protoResume.Languages, &resume_service.Language{
				LanguageId: language.LanguageID,
				Language:   language.Language,
				Fluency:    language.Fluency,
			})
		}

		for _, interest := range resume.Interests {
			protoResume.Interests = append(protoResume.Interests, &resume_service.Interest{
				InterestId: interest.InterestID,
				Name:       interest.Name,
			})
		}

		response.Resumes = append(response.Resumes, &resume_service.Resume{
			Id:          resume.ID,
			UserId:      resume.UserID,
			Url:         resume.URL,
			Filename:    resume.Filename,
			Salary:      uint64(resume.Salary),
			JobLocation: resume.JobLocation,
			Basic: &resume_service.Basic{
				Name:        resume.Basic.Name,
				JobTitle:    resume.Basic.JobTitle,
				Image:       resume.Basic.Image,
				Email:       resume.Basic.Email,
				PhoneNumber: resume.Basic.PhoneNumber,
				Website:     resume.Basic.Website,
				Summary:     resume.Basic.Summary,
				LocationId:  resume.Basic.LocationID,
				City:        resume.Basic.City,
				CountryCode: resume.Basic.CountryCode,
				Region:      resume.Basic.Region,
			},
			Profiles:     protoResume.Profiles,
			Works:        protoResume.Works,
			Projects:     protoResume.Projects,
			Educations:   protoResume.Educations,
			Certificates: protoResume.Certificates,
			HardSkills:   protoResume.HardSkills,
			SoftSkills:   protoResume.SoftSkills,
			Languages:    protoResume.Languages,
			Interests:    protoResume.Interests,
			Meta: &resume_service.Meta{
				Template: resume.Meta.Template,
				Lang:     resume.Meta.Lang,
			},
		})
	}
	response.TotalCount = resumes.TotalCount

	return response, nil
}

func (s resumeRPC) ListResume(ctx context.Context, in *resume_service.ListRequest) (*resume_service.ListResumeResponse, error) {
	offset := in.Limit * (in.Page - 1)
	resumes, err := s.resumeUseCase.ListResume(ctx, in.Limit, offset)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	response := &resume_service.ListResumeResponse{}

	for _, resume := range resumes.Resumes {
		protoResume := &resume_service.Resume{}

		for _, profile := range resume.Profiles {
			protoResume.Profiles = append(protoResume.Profiles, &resume_service.Profile{
				ProfileId: profile.ProfileID,
				Network:   profile.Network,
				Username:  profile.Username,
				Url:       profile.URL,
			})
		}

		for _, work := range resume.Works {
			protoResume.Works = append(protoResume.Works, &resume_service.Work{
				WorkId:    work.WorkID,
				Position:  work.Position,
				Company:   work.Company,
				StartDate: work.StartDate,
				EndDate:   work.EndDate,
				Location:  work.Location,
				Summary:   work.Summary,
				Skills:    work.Skills,
			})
		}

		for _, project := range resume.Projects {
			protoResume.Projects = append(protoResume.Projects, &resume_service.Project{
				ProjectId:   project.ProjectID,
				Name:        project.Name,
				Url:         project.URL,
				Description: project.Description,
			})
		}

		for _, education := range resume.Educations {
			protoResume.Educations = append(protoResume.Educations, &resume_service.Education{
				EducationId: education.EducationID,
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

		for _, cert := range resume.Certificates {
			protoResume.Certificates = append(protoResume.Certificates, &resume_service.Certificate{
				CertificateId: cert.CertificateID,
				Title:         cert.Title,
				Date:          cert.Date,
				Issuer:        cert.Issuer,
				Score:         cert.Score,
				Url:           cert.URL,
			})
		}

		for _, hard := range resume.HardSkills {
			protoResume.HardSkills = append(protoResume.HardSkills, &resume_service.HardSkill{
				HardSkillId: hard.HardSkillID,
				Name:        hard.Name,
				Level:       hard.Level,
			})
		}

		for _, soft := range resume.SoftSkills {
			protoResume.SoftSkills = append(protoResume.SoftSkills, &resume_service.SoftSkill{
				SoftSkillId: soft.SoftSkillID,
				Name:        soft.Name,
			})
		}

		for _, language := range resume.Languages {
			protoResume.Languages = append(protoResume.Languages, &resume_service.Language{
				LanguageId: language.LanguageID,
				Language:   language.Language,
				Fluency:    language.Fluency,
			})
		}

		for _, interest := range resume.Interests {
			protoResume.Interests = append(protoResume.Interests, &resume_service.Interest{
				InterestId: interest.InterestID,
				Name:       interest.Name,
			})
		}

		response.Resumes = append(response.Resumes, &resume_service.Resume{
			Id:          resume.ID,
			UserId:      resume.UserID,
			Url:         resume.URL,
			Filename:    resume.Filename,
			Salary:      uint64(resume.Salary),
			JobLocation: resume.JobLocation,
			Basic: &resume_service.Basic{
				Name:        resume.Basic.Name,
				JobTitle:    resume.Basic.JobTitle,
				Image:       resume.Basic.Image,
				Email:       resume.Basic.Email,
				PhoneNumber: resume.Basic.PhoneNumber,
				Website:     resume.Basic.Website,
				Summary:     resume.Basic.Summary,
				LocationId:  resume.Basic.LocationID,
				City:        resume.Basic.City,
				CountryCode: resume.Basic.CountryCode,
				Region:      resume.Basic.Region,
			},
			Profiles:     protoResume.Profiles,
			Works:        protoResume.Works,
			Projects:     protoResume.Projects,
			Educations:   protoResume.Educations,
			Certificates: protoResume.Certificates,
			HardSkills:   protoResume.HardSkills,
			SoftSkills:   protoResume.SoftSkills,
			Languages:    protoResume.Languages,
			Interests:    protoResume.Interests,
			Meta: &resume_service.Meta{
				Template: resume.Meta.Template,
				Lang:     resume.Meta.Lang,
			},
		})
	}
	response.TotalCount = resumes.TotalCount

	return response, nil
}
