package database

type User struct {
	IdUser int `"json:idUser"`
	Username string `"json:username"`
	Biography string `"json:biography"`
}

type Image struct {
	IdImage int64 `"json:idImage"`
	IdOwner int `"json:idOwner"`
	DateTime string `"json:dateTime"`
	Url	string `"json:url"`
}