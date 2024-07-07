package entity

type ListResume struct {
	Resumes    []*Resume
	TotalCount uint64
}

type Resume struct {
	ID          string
	UserID      string
	URL         string
	Salary      int64
	JobTitle    string
	Region      string
	JobLocation string
	JobType     string
	Experience  int64
	Template    string
}

type ListRequest struct {
	Page        int64
	Limit       int64
	JobTitle    string
	JobLocation string
	JobType     string
	Salary      int64
	Region      string
	Experience  int64
}
