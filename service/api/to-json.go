package api

func UserToJson(id int, username string) User {
	user := User{
		IdUser:   id,
		Username: username,
	}
	return user
}
