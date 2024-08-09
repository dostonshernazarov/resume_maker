package v1

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/dostonshernazarov/resume_maker/api-service/api/docs"
	"github.com/dostonshernazarov/resume_maker/api-service/api/models"
	pbu "github.com/dostonshernazarov/resume_maker/api-service/genproto/user_service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// UploadMedia
// @Summary     Upload User photo
// @Security    BearerAuth
// @Description Through this api front-ent can upload user photo and get the link to the media.
// @Tags        MEDIA
// @Accept      json
// @Produce     json
// @Param       file formData file true "Image"
// @Success     200 {object} string
// @Failure     400 {object} models.Error
// @Failure     500 {object} models.Error
// @Router      /v1/media/user-photo [POST]
func (h *HandlerV1) UploadMedia(c *gin.Context) {
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

	endpoint := h.Config.Minio.Host + h.Config.Minio.Port
	accessKeyID := h.Config.Minio.AccessKey
	secretAccessKey := h.Config.Minio.SecretKey
	bucketName := h.Config.Minio.BucketName
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

	userID, statusCode := GetIdFromToken(c.Request, h.Config)
	if statusCode == 401 {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: "Log In Again",
		})
		return
	}

	user, err := h.Service.UserService().GetUser(ctx, &pbu.Filter{
		Filter: map[string]string{
			"id": userID,
		},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		return
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

	minioURL := fmt.Sprintf("https://media.cvmaker.uz/%s/%s", endpoint, bucketName, objectName)

	user.Image = minioURL
	_, err = h.Service.UserService().UpdateUser(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		return
	}

	c.JSON(http.StatusOK, minioURL)
}
