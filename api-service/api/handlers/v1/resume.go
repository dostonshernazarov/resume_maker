package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dostonshernazarov/resume_maker/api-service/api/models"
	"github.com/dostonshernazarov/resume_maker/api-service/api/services"
	"github.com/dostonshernazarov/resume_maker/api-service/genproto/resume_service"
	l "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/logger"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/parser"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/pdf"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/template"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/utils"
	val "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/validation"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/utils/fs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
)

func createMultipartFileHeader(filePath string) *multipart.FileHeader {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	// create a buffer to hold the file in memory
	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// copy the content of the file to the form's file field
	if _, err := io.Copy(formPart, file); err != nil {
		log.Fatal(err)
		return nil
	}

	// close the form writer after the copying process is finished
	// I don't use defer in here to avoid unexpected EOF error
	formWriter.Close()

	// transform the bytes buffer into a form reader
	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		log.Fatal("multipart file not exists")
		return nil
	}

	return files[0]
}

// BasicReume
// BASIC RESUME  ...
// @Security BearerAuth
// @Router /v1/resume/basic [POST]
// @Summary BASIC RESUME
// @Description Api for post basic resume
// @Tags STEP-RESUME
// @Accept json
// @Produce json
// @Param BasicResumeData body models.Basics true "BasicResumeData"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
func (h HandlerV1) BasicResumeData(c *gin.Context) {
	var body models.Basics

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		},
		)
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	isEmail := val.IsValidEmail(body.Email)
	if !isEmail {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("Incorrect Email. Try again")
		return
	}

	redisId := uuid.New().String()

	// Connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-db:6379",
		Password: "",
		DB:       0,
	})
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			log.Println("failed to close redis connection", err)
			return
		}
	}(rdb)

	userByte, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to marshal body", l.Error(err))
		return
	}
	_, err = rdb.Set(context.Background(), redisId, userByte, time.Hour*5).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to set object to redis", l.Error(err))
		return
	}

	responseMessage := models.BacisResumeRes{
		Content:      "data has been saved",
		BasicRedisID: redisId,
		MainRedisID:  uuid.NewString(),
	}

	c.JSON(http.StatusOK, responseMessage)
}

// MainReume
// Main RESUME  ...
// @Security BearerAuth
// @Router /v1/resume/main [POST]
// @Summary Main RESUME
// @Description Api for post Main resume
// @Tags STEP-RESUME
// @Accept json
// @Produce json
// @Param MainResumeData body models.MainResumeReq true "MainResumeData"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
func (h HandlerV1) MainResumeData(c *gin.Context) {
	var body models.MainResumeReq

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		},
		)
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	// Connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-db:6379",
		Password: "",
		DB:       0,
	})
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			log.Println("failed to close redis connection", err)
			return
		}
	}(rdb)

	userByte, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to marshal body", l.Error(err))
		return
	}
	_, err = rdb.Set(context.Background(), body.MainRedisID, userByte, time.Hour*5).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Failed to set object to redis", l.Error(err))
		return
	}

	responseMessage := models.BacisResumeRes{
		Content:      "data has been saved",
		BasicRedisID: body.BasicRedisID,
		MainRedisID:  body.MainRedisID,
	}

	c.JSON(http.StatusOK, responseMessage)
}

// GenerateResume
// @Security 		BearerAuth
// @Summary 		Generate a Resume
// @Description 	This API for generate a resume
// @Tags 			STEP-RESUME
// @Accept			json
// @Produce 		json
// @Param 			data body models.LastResumeReq true "Resume Model"
// @Success 		200 {object} string "Resume URL"
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/resume/generate [POST]
func (h *HandlerV1) LastGenerateResume(c *gin.Context) {
	templateManager := template.NewTemplateManager("ui")
	htmlParser := parser.NewHTMLParser(models.OutputDir, models.OutputHtmlFile, templateManager)
	pdfGenerator := pdf.NewPDFGenerator()
	service := services.NewResumeService(htmlParser, pdfGenerator)

	var Lastbody models.LastResumeReq
	var resumeData models.Resume
	var bodyProd models.BotProduce

	if err := c.ShouldBindJSON(&Lastbody); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-db:6379",
		Password: "",
		DB:       0,
	})
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			log.Println("failed to close redis connection", err)
			return
		}
	}(rdb)

	basicValue, err := rdb.Get(context.Background(), Lastbody.BasicRedisID).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("Failed to get user from redis", l.Error(err))
		return
	}

	var basicDetail models.Basics
	if err := json.Unmarshal([]byte(basicValue), &basicDetail); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Error unmarshalling user detail", l.Error(err))
		return
	}

	MainValue, err := rdb.Get(context.Background(), Lastbody.MainRedisID).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("Failed to get user from redis", l.Error(err))
		return
	}

	var MainDetail models.MainResumeReq
	if err := json.Unmarshal([]byte(MainValue), &MainDetail); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("Error unmarshalling user detail", l.Error(err))
		return
	}

	resumeData.Basics = basicDetail
	resumeData.Work = MainDetail.Work
	resumeData.Projects = MainDetail.Projects
	resumeData.Education = MainDetail.Education
	resumeData.Certificates = Lastbody.Certificates
	resumeData.Skills = Lastbody.Skills
	resumeData.Languages = Lastbody.Languages
	resumeData.Interests = Lastbody.Interests
	resumeData.Meta = Lastbody.Meta

	htmlFile, err := service.Parser.ParseToHtml(resumeData)
	if err != nil {
		h.Logger.Error("ParseToHtml : " + err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "failed to parse HTML",
		})
		return
	}

	pdfData, err := service.Pdf.GenerateFromHTML(htmlFile)
	if err != nil {
		h.Logger.Error("GenerateFromHTML : " + err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "failed to generate PDF",
		})
		return
	}

	if err := fs.WriteFile(models.OutputPdfFile, pdfData); err != nil {
		h.Logger.Error("WriteFile : " + err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "failed to write PDF file",
		})
		return
	}

	multipartFile := createMultipartFileHeader(models.OutputPdfFile)

	endpoint := "3.76.217.224:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	bucketName := "resumes"
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "BucketAlreadyOwnedByYou" {
		} else {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: err.Error(),
			})
			log.Println(err.Error())
			return
		}
	}

	policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "AWS": ["*"]
                },
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::%s/*"]
            }
        ]
    }`, bucketName)

	err = minioClient.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	file := &models.File{
		File: multipartFile,
	}

	if file.File.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "File size cannot be larger than 10 MB",
		})
		return
	}

	ext := filepath.Ext(file.File.Filename)

	uploadDir := "./media"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			log.Println("Error creating media directory", err.Error())
			return
		}
	}

	userFullName := strings.Join(strings.Split(resumeData.Basics.Name, " "), "_") + "_CV_Maker"
	newFilename := userFullName + ext
	uploadPath := filepath.Join(uploadDir, newFilename)

	if err := c.SaveUploadedFile(file.File, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	objectName := newFilename
	contentType := "application/pdf"
	_, err = minioClient.FPutObject(context.Background(), bucketName, objectName, uploadPath, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	minioURL := fmt.Sprintf("https://media.cvmaker.uz/%s/%s", bucketName, objectName)

	userID, status := GetIdFromToken(c.Request, h.Config)
	if status == http.StatusUnauthorized {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: "Token is invalid",
		})
		h.Logger.Error(fmt.Sprintf("Token is invalid: %s", userID))
		return
	}

	timeNow := time.Now()

	var (
		profiles     []*resume_service.Profile
		works        []*resume_service.Work
		projects     []*resume_service.Project
		educations   []*resume_service.Education
		certificates []*resume_service.Certificate
		hardSkills   []*resume_service.HardSkill
		softSkills   []*resume_service.SoftSkill
		languages    []*resume_service.Language
		interests    []*resume_service.Interest
	)

	for _, profile := range resumeData.Basics.Profiles {
		profiles = append(profiles, &resume_service.Profile{
			Network:  profile.Network,
			Username: profile.Username,
			Url:      profile.URL,
		})
	}

	for _, work := range resumeData.Work {
		if work.StartDate == "" {
			work.StartDate = timeNow.String()
		}
		works = append(works, &resume_service.Work{
			Position:  work.Position,
			Company:   work.Company,
			StartDate: work.StartDate,
			EndDate:   work.EndDate,
			Location:  work.Location,
			Summary:   work.Summary,
		})
	}

	for _, project := range resumeData.Projects {
		projects = append(projects, &resume_service.Project{
			Name:        project.Name,
			Url:         project.URL,
			Description: project.Description,
		})
	}

	for _, education := range resumeData.Education {
		if education.StartDate == "" {
			education.StartDate = timeNow.String()
		}
		educations = append(educations, &resume_service.Education{
			EducationId: uuid.NewString(),
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

	for _, certificate := range resumeData.Certificates {

		if certificate.Date == "" {
			certificate.Date = timeNow.String()
		}
		certificates = append(certificates, &resume_service.Certificate{
			Title:  certificate.Title,
			Date:   certificate.Date,
			Issuer: certificate.Issuer,
			Score:  certificate.Score,
			Url:    certificate.URL,
		})
	}

	for _, hard := range resumeData.Skills {
		hardSkills = append(hardSkills, &resume_service.HardSkill{
			Name:  hard.Name,
			Level: hard.Level,
		})
	}

	for _, soft := range resumeData.SoftSkills {
		softSkills = append(softSkills, &resume_service.SoftSkill{
			Name: soft.Name,
		})
	}

	for _, lang := range resumeData.Languages {
		languages = append(languages, &resume_service.Language{
			Language: lang.Language,
			Fluency:  lang.Fluency,
		})
	}

	for _, interest := range resumeData.Interests {
		interests = append(interests, &resume_service.Interest{
			Name: interest.Name,
		})
	}

	_, err = h.Service.ResumeService().CreateResume(context.Background(), &resume_service.Resume{
		Id:          uuid.NewString(),
		UserId:      userID,
		Url:         resumeData.Basics.URL,
		Filename:    minioURL,
		Salary:      basicDetail.Salary,
		JobLocation: basicDetail.JobLocation,
		Basic: &resume_service.Basic{
			Name:        resumeData.Basics.Name,
			JobTitle:    resumeData.Basics.Label,
			Image:       resumeData.Basics.Image,
			Email:       resumeData.Basics.Email,
			PhoneNumber: resumeData.Basics.Phone,
			Website:     resumeData.Basics.URL,
			Summary:     resumeData.Basics.Summary,
			City:        resumeData.Basics.Location.City,
			CountryCode: resumeData.Basics.Location.CountryCode,
			Region:      resumeData.Basics.Location.Region,
		},
		Profiles:     profiles,
		Works:        works,
		Projects:     projects,
		Educations:   educations,
		Certificates: certificates,
		HardSkills:   hardSkills,
		SoftSkills:   softSkills,
		Languages:    languages,
		Interests:    interests,
		Meta: &resume_service.Meta{
			Template: resumeData.Meta.Template,
			Lang:     resumeData.Meta.Lang,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(fmt.Sprintf("failed to save resume content into service %v", err))
		return
	}

	// Message telegram

	bodyProd.FullName = resumeData.Basics.Name
	bodyProd.Email = resumeData.Basics.Email
	bodyProd.PhoneNumber = resumeData.Basics.Phone
	bodyProd.JobTitle = resumeData.Basics.Label
	bodyProd.City = resumeData.Basics.Location.City
	bodyProd.Resume = minioURL

	for _, v := range resumeData.Basics.Profiles {

		bodyProd.Links = append(bodyProd.Links, v.URL)
	}

	bodyProd.Salary = resumeData.Basics.Salary
	bodyProd.Summary = resumeData.Basics.Summary

	message := formatMessage(bodyProd)

	err = sendMessageToTelegram(h.Config.APIToken, h.Config.ChatID, message)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(fmt.Sprintf("failed to send resume to the telegram %v", err))
		return
	}

	// Message telegram

	c.JSON(http.StatusOK, minioURL)
}

// GenerateResume
// @Security 		BearerAuth
// @Summary 		Generate a Resume
// @Description 	This API for generate a resume
// @Tags 			RESUME
// @Accept			json
// @Produce 		json
// @Param 			data body models.ResumeGenetare true "Resume Model"
// @Success 		200 {object} string "Resume URL"
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/resume/generate-resume [POST]
func (h *HandlerV1) GenerateResume(c *gin.Context) {
	templateManager := template.NewTemplateManager("ui")
	htmlParser := parser.NewHTMLParser(models.OutputDir, models.OutputHtmlFile, templateManager)
	pdfGenerator := pdf.NewPDFGenerator()
	service := services.NewResumeService(htmlParser, pdfGenerator)

	var body models.ResumeGenetare
	var bodyProd models.BotProduce

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		return
	}

	var resumeData models.Resume
	resumeData.Basics = body.Basics
	resumeData.Certificates = body.Certificates
	resumeData.Work = body.Work
	resumeData.Projects = body.Projects
	resumeData.Education = body.Education
	resumeData.Skills = body.Skills
	resumeData.SoftSkills = body.SoftSkills
	resumeData.Languages = body.Languages
	resumeData.Interests = body.Interests
	resumeData.Meta = body.Meta
	resumeData.Labels = body.Labels

	htmlFile, err := service.Parser.ParseToHtml(resumeData)
	if err != nil {
		h.Logger.Error("ParseToHtml : " + err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "failed to parse HTML",
		})
		return
	}

	pdfData, err := service.Pdf.GenerateFromHTML(htmlFile)
	if err != nil {
		h.Logger.Error("GenerateFromHTML : " + err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "failed to generate PDF",
		})
		return
	}

	if err := fs.WriteFile(models.OutputPdfFile, pdfData); err != nil {
		h.Logger.Error("WriteFile : " + err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "failed to write PDF file",
		})
		return
	}

	multipartFile := createMultipartFileHeader(models.OutputPdfFile)

	endpoint := "3.76.217.224:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	bucketName := "resumes"
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "BucketAlreadyOwnedByYou" {
		} else {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: err.Error(),
			})
			log.Println(err.Error())
			return
		}
	}

	policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "AWS": ["*"]
                },
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::%s/*"]
            }
        ]
    }`, bucketName)

	err = minioClient.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	file := &models.File{
		File: multipartFile,
	}

	if file.File.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "File size cannot be larger than 10 MB",
		})
		return
	}

	ext := filepath.Ext(file.File.Filename)

	uploadDir := "./media"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			log.Println("Error creating media directory", err.Error())
			return
		}
	}

	cvuuid := uuid.NewString()
	userFullName := strings.Join(strings.Split(resumeData.Basics.Name, " "), "") + cvuuid[:5] + "CVMaker"
	newFilename := userFullName + ext
	uploadPath := filepath.Join(uploadDir, newFilename)

	if err := c.SaveUploadedFile(file.File, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	objectName := newFilename
	contentType := "application/pdf"
	_, err = minioClient.FPutObject(context.Background(), bucketName, objectName, uploadPath, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	minioURL := fmt.Sprintf("https://media.cvmaker.uz/%s/%s", bucketName, objectName)

	userID, status := GetIdFromToken(c.Request, h.Config)
	if status == http.StatusUnauthorized {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: "Token is invalid",
		})
		h.Logger.Error(fmt.Sprintf("Token is invalid: %s", userID))
		return
	}

	timeNow := time.Now()

	var (
		profiles     []*resume_service.Profile
		works        []*resume_service.Work
		projects     []*resume_service.Project
		educations   []*resume_service.Education
		certificates []*resume_service.Certificate
		hardSkills   []*resume_service.HardSkill
		softSkills   []*resume_service.SoftSkill
		languages    []*resume_service.Language
		interests    []*resume_service.Interest
	)

	for _, profile := range resumeData.Basics.Profiles {
		profiles = append(profiles, &resume_service.Profile{
			Network:  profile.Network,
			Username: profile.Username,
			Url:      profile.URL,
		})
	}

	for _, work := range resumeData.Work {
		if work.StartDate == "" {
			work.StartDate = timeNow.String()
		}
		works = append(works, &resume_service.Work{
			Position:  work.Position,
			Company:   work.Company,
			StartDate: work.StartDate,
			EndDate:   work.EndDate,
			Location:  work.Location,
			Summary:   work.Summary,
		})
	}

	for _, project := range resumeData.Projects {
		projects = append(projects, &resume_service.Project{
			Name:        project.Name,
			Url:         project.URL,
			Description: project.Description,
		})
	}

	for _, education := range resumeData.Education {
		if education.StartDate == "" {
			education.StartDate = timeNow.String()
		}
		educations = append(educations, &resume_service.Education{
			EducationId: uuid.NewString(),
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

	for _, certificate := range resumeData.Certificates {

		if certificate.Date == "" {
			certificate.Date = timeNow.String()
		}
		certificates = append(certificates, &resume_service.Certificate{
			Title:  certificate.Title,
			Date:   certificate.Date,
			Issuer: certificate.Issuer,
			Score:  certificate.Score,
			Url:    certificate.URL,
		})
	}

	for _, hard := range resumeData.Skills {
		hardSkills = append(hardSkills, &resume_service.HardSkill{
			Name:  hard.Name,
			Level: hard.Level,
		})
	}

	for _, soft := range resumeData.SoftSkills {
		softSkills = append(softSkills, &resume_service.SoftSkill{
			Name: soft.Name,
		})
	}

	for _, lang := range resumeData.Languages {
		languages = append(languages, &resume_service.Language{
			Language: lang.Language,
			Fluency:  lang.Fluency,
		})
	}

	for _, interest := range resumeData.Interests {
		interests = append(interests, &resume_service.Interest{
			Name: interest.Name,
		})
	}

	_, err = h.Service.ResumeService().CreateResume(context.Background(), &resume_service.Resume{
		Id:          uuid.NewString(),
		UserId:      userID,
		Url:         resumeData.Basics.URL,
		Filename:    minioURL,
		Salary:      body.Salary,
		JobLocation: body.JobLocation,
		Basic: &resume_service.Basic{
			Name:        resumeData.Basics.Name,
			JobTitle:    resumeData.Basics.Label,
			Image:       resumeData.Basics.Image,
			Email:       resumeData.Basics.Email,
			PhoneNumber: resumeData.Basics.Phone,
			Website:     resumeData.Basics.URL,
			Summary:     resumeData.Basics.Summary,
			City:        resumeData.Basics.Location.City,
			CountryCode: resumeData.Basics.Location.CountryCode,
			Region:      resumeData.Basics.Location.Region,
		},
		Profiles:     profiles,
		Works:        works,
		Projects:     projects,
		Educations:   educations,
		Certificates: certificates,
		HardSkills:   hardSkills,
		SoftSkills:   softSkills,
		Languages:    languages,
		Interests:    interests,
		Meta: &resume_service.Meta{
			Template: resumeData.Meta.Template,
			Lang:     resumeData.Meta.Lang,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(fmt.Sprintf("failed to save resume content into service %v", err))
		return
	}

	// Message telegram

	bodyProd.FullName = resumeData.Basics.Name
	bodyProd.Email = resumeData.Basics.Email
	bodyProd.PhoneNumber = resumeData.Basics.Phone
	bodyProd.JobTitle = resumeData.Basics.Label
	bodyProd.City = resumeData.Basics.Location.City
	bodyProd.Resume = minioURL

	for _, v := range resumeData.Basics.Profiles {

		bodyProd.Links = append(bodyProd.Links, v.URL)
	}

	bodyProd.Salary = resumeData.Basics.Salary
	bodyProd.Summary = resumeData.Basics.Summary

	message := formatMessage(bodyProd)

	err = sendMessageToTelegram(h.Config.APIToken, h.Config.ChatID, message)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(fmt.Sprintf("failed to send resume to the telegram %v", err))
		return
	}

	// Message telegram

	c.JSON(http.StatusOK, minioURL)
}

// UploadMedia
// @Summary     Upload Resume Photo
// @Security    BearerAuth
// @Description Through this api front-ent can upload resume photo and get the link to the resume.
// @Tags        MEDIA
// @Accept      json
// @Produce     json
// @Param       file formData file true "Image"
// @Success     200 {object} string
// @Failure     400 {object} models.Error
// @Failure     500 {object} models.Error
// @Router      /v1/resume/resume-photo [POST]
func (h *HandlerV1) UploadResumePhoto(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	endpoint := "3.76.217.224:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	bucketName := "resumes"
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "BucketAlreadyOwnedByYou" {
		} else {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: err.Error(),
			})
			log.Println(err.Error())
			return
		}
	}

	policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "AWS": ["*"]
                },
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::%s/*"]
            }
        ]
    }`, bucketName)

	err = minioClient.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	file := &models.File{}
	err = c.ShouldBind(&file)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	if file.File.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "File size cannot be larger than 10 MB",
		})
		return
	}

	ext := filepath.Ext(file.File.Filename)

	if ext != ".png" && ext != ".jpg" && ext != ".svg" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Only .jpg and .png format images are accepted",
		})
		return
	}

	uploadDir := "./media"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			log.Println("Error creating media directory", err.Error())
			return
		}
	}

	id := uuid.New().String()

	newFilename := id + ext
	uploadPath := filepath.Join(uploadDir, newFilename)

	if err := c.SaveUploadedFile(file.File, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	objectName := newFilename
	contentType := "image/jpeg"
	_, err = minioClient.FPutObject(context.Background(), bucketName, objectName, uploadPath, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	minioURL := fmt.Sprintf("https://media.cvmaker.uz/%s/%s", bucketName, objectName)

	c.JSON(http.StatusOK, minioURL)
}

// ListUsersResume
// @Summary LIST USER RESUME
// @Security BearerAuth
// @Description Api for ListUsersResume
// @Tags RESUME
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Success 200 {object} []models.ResResume
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/resume/list [get]
func (h *HandlerV1) ListUserResume(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		return
	}

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	userID, statusCode := GetIdFromToken(c.Request, h.Config)
	if statusCode == 401 {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: "Log In Again",
		})
		return
	}

	response, err := h.Service.ResumeService().GetUserResume(
		context.Background(), &resume_service.UserWithID{
			Page:   params.Page,
			Limit:  params.Limit,
			UserId: userID,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		l.Error(err)
		return
	}

	var resumes []*models.ResResume
	for _, val := range response.Resumes {
		var resRes models.ResResume
		resRes.ID = val.Id
		resRes.UserID = val.UserId
		resRes.Filename = val.Filename
		resRes.JobTitle = val.Basic.JobTitle
		resRes.Salary = val.Salary
		resRes.JobLocation = val.JobLocation

		resumes = append(resumes, &resRes)
	}

	c.JSON(http.StatusOK, resumes)
}

// ListResume
// @Summary LIST RESUME
// @Security BearerAuth
// @Description Api for ListREsume
// @Tags RESUME
// @Accept json
// @Produce json
// @Success 200 {object} models.ResResumeList
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/resume/list [get]
func (h *HandlerV1) ListResume(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		return
	}

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	response, err := h.Service.ResumeService().ListResume(
		context.Background(), &resume_service.ListRequest{
			Page:  params.Page,
			Limit: params.Limit,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		l.Error(err)
		return
	}

	var resumes models.ResResumeList
	for _, val := range response.Resumes {
		var resRes models.ResResume
		resRes.ID = val.Id
		resRes.UserID = val.UserId
		resRes.Filename = val.Filename
		resRes.JobTitle = val.Basic.JobTitle
		resRes.Salary = val.Salary
		resRes.JobLocation = val.JobLocation

		resumes.Resumes = append(resumes.Resumes, resRes)
	}
	resumes.Count = response.TotalCount

	c.JSON(http.StatusOK, resumes)
}

// DeleteResume
// @Summary DELETE
// @Security BearerAuth
// @Description Api for Delete Resume
// @Tags RESUME
// @Accept json
// @Produce json
// @Param id query string true "ID"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/resumes/{id} [delete]
func (h *HandlerV1) DeleteResume(c *gin.Context) {
	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	id := c.Query("id")

	user, err := h.Service.ResumeService().GetResumeByID(context.Background(), &resume_service.ResumeWithID{
		ResumeId: id,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to get resume in delete", l.Error(err))
		return
	}

	if user == nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to get resume in delete", l.Error(err))
		return
	}

	_, err = h.Service.ResumeService().DeleteResume(
		context.Background(), &resume_service.ResumeWithID{
			ResumeId: id,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to delete resume", l.Error(err))
		return
	}

	// if response != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Went wrong",
	// 	})
	// 	h.Logger.Error("failed to delete user", l.Error(err))
	// 	return
	// }

	c.JSON(http.StatusOK, &models.RegisterRes{
		Content: "Resume has been deleted",
	})
}

func formatMessage(botProduce models.BotProduce) string {
	links := ""
	for _, link := range botProduce.Links {
		links += fmt.Sprintf("- %s\n", link)
	}
	return fmt.Sprintf("NEW RESUME\n\n👱🏻‍♂️** Employee:** %s\n📧** Email:** %s\n📞** Phone Number:** %s\n📚** Job:** %s\n📄** Resume:** %s\n🔗** Links:**\n%s🏠** City:** %s\n💵** Salary:** %d\n🔎** Summary:** %s",
		botProduce.FullName, botProduce.Email, botProduce.PhoneNumber, botProduce.JobTitle, botProduce.Resume, links, botProduce.City, botProduce.Salary, botProduce.Summary)
}

func sendMessageToTelegram(apiToken string, chatID string, message string) error {
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return fmt.Errorf("failed to create a new bot: %w", err)
	}

	msg := tgbotapi.NewMessageToChannel(chatID, message)
	msg.ParseMode = "Markdown"

	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
