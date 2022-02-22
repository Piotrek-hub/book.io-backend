package utils

type BookRequest struct {
	Username 	  string `form:"username"`
	UserKey       string `form:"user_key"`
	Title         string `form:"title"`
	Author        string `form:"author"`
	Pages         uint16 `form:"pages"`
	DateCompleted string `form:"date-completed"`
	Status        string `form:"status"`
}

type Resp struct {
	UserKey string `form:"user_key"`
}

type Config struct {
	Login string
	Password string
}

