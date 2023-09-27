package utils

type Todo struct {
	Id   int
	Text string
}

type Data struct {
	Todos []Todo
}

type JWT struct {
	Token string
	Expires int64
}
