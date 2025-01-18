package main

import "log"

const getProjectsQuery = "SELECT * FROM Project ORDER BY id LIMIT ? OFFSET ?"

func GetProjects(limit, offset int) []Project {
	db := GetDBConnection()
	projects := []Project{}

	stmt, err := db.Prepare(getProjectsQuery)
	if err != nil {
		log.Fatal("Preparing projects statement failed", err)
	}

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		log.Fatal("Querying projects statement failed", err)
	}

	for rows.Next() {
		var (
			id   int
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("Row is missing fields", err)
		}

		projects = append(projects, NewProject(id, name))
	}

	return projects
}
