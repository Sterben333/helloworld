package model

// Kratos hello kratos.
type Kratos struct {
	Id    int
	Title string
}

type Article struct {
	ID      int64
	Content string
	Author  string
}

type User struct {
	Uid      int64
	Nickname string
	Age      int32
	Uptime   int32
	Addtime  int32
}
