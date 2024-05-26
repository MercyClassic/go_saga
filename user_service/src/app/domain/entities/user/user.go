package userEntities

type User struct {
	Id       int
	Name     string
	Username string
	Balance  float32
}

func NewUser(
	name string,
	username string,
) *User {
	return &User{
		Name:     name,
		Username: username,
	}
}
