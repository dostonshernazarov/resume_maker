package entity

type ListResume struct {
	Resumes    []*Resume
	TotalCount uint64
}

type ResumeContent struct {
	ResumeID    string
	UserID      string
	URL         string
	Filename    string
	Template    string
	Lang        string
	Salary      int64
	JobLocation string
}

type Resume struct {
	ID           string
	UserID       string
	URL          string
	Filename     string
	Basic        Basic
	Profiles     []*Profile
	Works        []*Work
	Projects     []*Project
	Educations   []*Education
	Certificates []*Certificate
	HardSkills   []*HardSkill
	SoftSkills   []*SoftSkill
	Languages    []*Language
	Interests    []*Interest
	Meta         Meta
	Salary       int64
	JobLocation  string
}

type Basic struct {
	Name        string
	JobTitle    string
	Image       string
	Email       string
	PhoneNumber string
	Website     string
	Summary     string
	LocationID  string
	City        string
	CountryCode string
	Region      string
}

type Profile struct {
	ProfileID string
	Network   string
	Username  string
	URL       string
}

type Work struct {
	WorkID    string
	Position  string
	Company   string
	StartDate string
	EndDate   string
	Location  string
	Summary   string
	Skills    []string
}

type Project struct {
	ProjectID   string
	Name        string
	URL         string
	Description string
}

type Education struct {
	EducationID string
	Institution string
	Area        string
	StudyType   string
	Location    string
	StartDate   string
	EndDate     string
	Score       string
	Courses     []string
}

type Certificate struct {
	CertificateID string
	Title         string
	Date          string
	Issuer        string
	Score         string
	URL           string
}

type HardSkill struct {
	HardSkillID string
	Name        string
	Level       string
}

type SoftSkill struct {
	SoftSkillID string
	Name        string
}

type Language struct {
	LanguageID string
	Language   string
	Fluency    string
}

type Interest struct {
	InterestID string
	Name       string
}

type Meta struct {
	Template string
	Lang     string
}
