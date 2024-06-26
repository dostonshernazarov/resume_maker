package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dostonshernazarov/resume_maker/api-service/api/models"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/parser"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/pdf"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/utils/fs"
)

type ResumeService struct {
	Parser *parser.HTMLParser
	Pdf    *pdf.Generator
}

func NewResumeService(parser *parser.HTMLParser, pdf *pdf.Generator) *ResumeService {
	return &ResumeService{
		Parser: parser,
		Pdf:    pdf,
	}
}

func (s *ResumeService) GeneratePDF(resumeData models.Resume) error {
	htmlFile, err := s.Parser.ParseToHtml(resumeData)
	if err != nil {
		return err
	}

	pdfData, err := s.Pdf.GenerateFromHTML(htmlFile)
	if err != nil {
		return err
	}

	if err := fs.WriteFile(models.OutputPdfFile, pdfData); err != nil {
		return err
	}

	return nil
}

func (s *ResumeService) UnmarshalResume(data []byte) (models.Resume, error) {
	var resumeData models.Resume
	err := json.Unmarshal(data, &resumeData)
	if err != nil {
		return models.Resume{}, errors.New(fmt.Sprintf("failed to unmarshal JSON %v", err))
	}

	return resumeData, nil
}
