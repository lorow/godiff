package main

func GetProjects(limit, offset int) []Project {
	projects := []Project{
		{
			id:   1,
			name: "Project 1",
		},
		{
			id:   2,
			name: "Project 2",
		},
	}
	return projects
}
