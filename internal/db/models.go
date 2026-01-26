package db

type Editor struct {
	ID         int
	Name       string
	Content    string
	Project_id int
}

type Request struct {
	ID         int
	Name       string
	Url        string
	Method     string
	Headers    string
	Body       string
	Response   int
	Project_id int
}

type Project struct {
	ID   int
	Name string
}

func NewProject(id int, name string) Project {
	return Project{id, name}
}
