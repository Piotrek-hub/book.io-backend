package db

type User struct {
	Login    string `form:"login"`
	Password string `form:"password"`
	UserKey  string `form:"user_key"`
}

type Book struct {
	Title         string `form:"title"`
	Author        string `form:"author"`
	Pages         uint16 `form:"pages"`
	DateCompleted string `form:"date-completed"`
	Status        string `form:"status"`
}
