package utils

type JWT struct {
	Token   string
	Expires int64
}

type UserData struct {
	Id       int
	Username string
	Email    string
}

type TodoData struct {
	Id   int
	Text string
}

type TodosData struct {
	Todos    []TodoData
}

type SessionData struct {
	LoggedIn bool
	User     UserData
}

type IndexData struct {
	Session SessionData
	Todos    []TodoData
}

type SigninData struct {
	User  string
	Error string
}

type SignupData struct {
	Username string
	Email    string
	Error    string
}
