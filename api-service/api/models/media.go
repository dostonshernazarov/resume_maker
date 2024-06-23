package models

import "mime/multipart"

type (
	File struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
	}

	UploadPhotoRes struct {
		URL string `json:"photo_url"`
	}

	MediaResponse struct {
		ErrorCode    int            `json:"error_code"`
		ErrorMessage string         `json:"error_message"`
		Body         UploadPhotoRes `json:"body"`
	}

	Media struct {
		UserId string `json:",omitempty"`
		ProfileImg  string `json:"image_url,omitempty"`
		FileName  string `json:"file_name"`
	}

	ProductImages struct {
		Images []*Media `json:"images,omitempty"`
	}

	EstablishmentImageRespons struct {
		ImageURL string `json:"image_url"`
		Message string `json:"message"`
	}

)

