package api

func UserToJson(id int, username string, biography string) User{
	user := User {
		IdUser: id,
		Username: username,
		Biography: biography,
	}
	return user
}