package database

type User struct {
	IdUser   int    `json:"idUser"`
	Username string `json:"username"`
}

type Image struct {
	IdImage       int              `json:"idImage"`
	IdOwner       int              `json:"idOwner"`
	DateTime      string           `json:"dateTime"`
	File          []byte           `json:"file"`
	LikesCount    int              `json:"likesCount"`
	CommentsCount int              `json:"commentsCount"`
	LikeStatus    bool             `json:"likeStatus"`
	Comments      []CommentWrapper `json:"comments"`
}

type CommentWrapper struct {
	CommentData Commenting `json:"commentData"`
	Username    string     `json:"username"`
}

type DoubleIdUser struct {
	IdUser  int `json:"idUser"`
	IdUser2 int `json:"idUser2"`
}

type Commenting struct {
	IdComment    int    `json:"idComment"`
	IdImage      int    `json:"idImage"`
	IdOwner      int    `json:"idOwner"`
	IdUserWriter int    `json:"idUserWriter"`
	Text         string `json:"text"`
}

type Liking struct {
	IdImage int `json:"idImage"`
	IdOwner int `json:"idOwner"`
	IdLiker int `json:"idComment"`
}

type UserProfile struct {
	User          User    `json:"user"`
	FollowCount   int     `json:"followCount"`
	FollowerCount int     `json:"followerCount"`
	Image         []Image `json:"images"`
	FollowStatus  bool    `json:"followStatus"`
}

type Stream2 struct {
	Image    Image  `json:"image"`
	Username string `json:"username"`
}
