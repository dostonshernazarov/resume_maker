package models

import "github.com/dostonshernazarov/resume_maker/api-service/internal/utils/lang"

const (
	OutputDir      = "output"
	OutputHtmlFile = "output/resume.html"
	OutputPdfFile  = "output/resume.pdf"
)

const (
	ClassicTemplate = "classic"
)

const (
	EducationLabel   = "EducationLabel"
	ExperiencesLabel = "ExperiencesLabel"
	LanguagesLabel   = "LanguagesLabel"
	SkillsLabel      = "SkillsLabel"
	SoftSkillsLabel  = "SoftSkillsLabel"
	ProjectsLabel    = "ProjectsLabel"
	InterestsLabel   = "InterestsLabel"
	ProfileLabel     = "ProfileLabel"
	SinceLabel       = "SinceLabel"
)

type Filter struct {
	JobTitle    string `json:"job_title"`
	JobLocation string `json:"job_location" example:"offline"`
	JobType     string `json:"job_type" example:"full-time"`
	Salary      int64  `json:"salary"`
	Country     string `json:"country"`
	Experience  int64  `json:"experience"`
}

type Resume struct {
	Basics       Basics        `json:"basics"`
	Work         []Work        `json:"work"`
	Projects     []Project     `json:"projects"`
	Education    []Education   `json:"education"`
	Certificates []Certificate `json:"certificates"`
	Skills       []Skill       `json:"skills"`
	SoftSkills   []Skill       `json:"softSkills"`
	Languages    []Language    `json:"languages"`
	Interests    []Interest    `json:"interests"`
	Meta         Meta          `json:"meta"`
	Labels       ResumeLabels
}

type ResumeGenetare struct {
	Basics       Basics        `json:"basics"`
	Work         []Work        `json:"work"`
	Projects     []Project     `json:"projects"`
	Education    []Education   `json:"education"`
	Certificates []Certificate `json:"certificates"`
	Skills       []Skill       `json:"skills"`
	SoftSkills   []Skill       `json:"softSkills"`
	Languages    []Language    `json:"languages"`
	Interests    []Interest    `json:"interests"`
	Meta         Meta          `json:"meta"`
	Labels       ResumeLabels
	Salary       uint64 `json:"salary"`
	JobLocation  string `json:"job_location" example:"offline"`
}

type LastResumeReq struct {
	Certificates []Certificate `json:"certificates"`
	Skills       []Skill       `json:"skills"`
	Languages    []Language    `json:"languages"`
	Interests    []Interest    `json:"interests"`
	Meta         Meta          `json:"meta"`
	BasicRedisID string        `json:"basic_redis_id"`
	MainRedisID  string        `json:"main_redis_id"`
}

type Basics struct {
	Name           string    `json:"name"`
	Label          string    `json:"label"`
	Image          string    `json:"image"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Summary        string    `json:"summary"`
	Location       Location  `json:"location"`
	URL            string    `json:"url"`
	Profiles       []Profile `json:"profiles"`
	Salary         uint64    `json:"salary"`
	JobLocation    string    `json:"job_location"`
	JobType        string    `json:"job_type" example:"full-time"`
	ExperienceYear int32     `json:"experience_year"`
}

type BotProduce struct {
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phone_number"`
	JobTitle    string   `json:"job_title"`
	Resume      string   `json:"resume"`
	Links       []string `json:"links"`
	City        string   `json:"city"`
	Salary      uint64   `json:"salary"`
	Summary     string   `json:"summary"`
}

type Location struct {
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"region"`
}

type Profile struct {
	Network  string `json:"network"`
	Username string `json:"username"`
	URL      string `json:"url"`
}

type Work struct {
	Position     string   `json:"position"`
	Company      string   `json:"company"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate"`
	Summary      string   `json:"summary"`
	Location     string   `json:"location"`
	Skills       []string `json:"skills"`
	ContractType string   `json:"contract_type"`
}

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type Education struct {
	Institution string   `json:"institution"`
	Area        string   `json:"area"`
	StudyType   string   `json:"studyType"`
	Location    string   `json:"location"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Score       string   `json:"score"`
	Courses     []string `json:"courses"`
}

type MainResumeReq struct {
	Work         []Work      `json:"work"`
	Projects     []Project   `json:"projects"`
	Education    []Education `json:"education"`
	BasicRedisID string      `json:"basic_redis_id"`
	MainRedisID  string      `json:"main_redis_id"`
}

type Certificate struct {
	Title  string `json:"title"`
	Date   string `json:"date"`
	Issuer string `json:"issuer"`
	Score  string `json:"score"`
	URL    string `json:"url"`
}

type Skill struct {
	Name     string   `json:"name"`
	Level    string   `json:"level"`
	Keywords []string `json:"keywords"`
}

type Language struct {
	Language string `json:"language"`
	Fluency  string `json:"fluency"`
}

type Interest struct {
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
}

type Meta struct {
	Template string `json:"template"`
	Lang     string `json:"lang"`
}

type ResumeLabels struct {
	Education   string
	Experiences string
	Projects    string
	Skills      string
	SoftSkills  string
	Languages   string
	Interests   string
	Profile     string
	Since       string
}

type ResResume struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Filename    string `json:"filename"`
	JobTitle    string `json:"job_title"`
	City        string `json:"city"`
	Salary      uint64 `json:"salary"`
	JobLocation string `json:"job_location"`
	Experiance  int32  `json:"experiance_year"`
}

type ResResumeList struct {
	Resumes []ResResume `json:"resumes"`
	Count   uint64      `json:"count"`
}

func (r *Resume) GetEducationLabel() string {
	return lang.Translate(r.Meta.Lang, EducationLabel)
}

func (r *Resume) GetExperiencesLabel() string {
	return lang.Translate(r.Meta.Lang, ExperiencesLabel)
}

func (r *Resume) GetSkillsLabel() string {
	return lang.Translate(r.Meta.Lang, SkillsLabel)
}

func (r *Resume) GetSoftSkillsLabel() string {
	return lang.Translate(r.Meta.Lang, SoftSkillsLabel)
}

func (r *Resume) GetProjectsLabel() string {
	return lang.Translate(r.Meta.Lang, ProjectsLabel)
}

func (r *Resume) GetLanguagesLabel() string {
	return lang.Translate(r.Meta.Lang, LanguagesLabel)
}

func (r *Resume) GetInterestsLabel() string {
	return lang.Translate(r.Meta.Lang, InterestsLabel)
}

func (r *Resume) GetProfileLabel() string {
	return lang.Translate(r.Meta.Lang, ProfileLabel)
}

func (r *Resume) GetSinceLabel() string {
	return lang.Translate(r.Meta.Lang, SinceLabel)
}
