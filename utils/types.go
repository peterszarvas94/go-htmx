package utils

type UserData struct {
	Id       int
	Username string
	Email    string
}

type TodoData struct {
	Id   int
	Text string
}

type SessionData struct {
	LoggedIn bool
	User     UserData
}

type NewTodoData struct {
	Session SessionData
	Id      int
	Text    string
}

type TodosData struct {
	Session SessionData
	Todos   []TodoData
}

type SigninData struct {
	User  string
	Error string
}

type ExistsData struct {
	Username bool
	Email    bool
}

type SignupData struct {
	Username string
	Email    string
	Error    string
	Exists   ExistsData
}

type TokenVariant string
const (
	Access  TokenVariant = "access"
	Refresh TokenVariant = "refresh"
)

type JWT struct {
	Token   string
	Expires int64
}
