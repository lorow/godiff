package main

type Editor struct {
	id         int
	name       string
	content    string
	project_id int
}

type Request struct {
	id         int
	name       string
	url        string
	method     string
	headers    string
	body       string
	response   int
	project_id int
}

type Project struct {
	id   int
	name string
}
