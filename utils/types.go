package utils

type Todo struct {
	Id   int
	Text string
}

type Data struct {
	Todos []Todo
	LoggedIn bool
}

type JWT struct {
	Token string
	Expires int64
}

type Signin struct {
	User string
	Error string
}

type Signup struct {
	Username string
	Email string
	Error string
}
