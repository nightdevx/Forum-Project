package main

type postData struct {
	UserData User
	PostData Post
}

type Post struct {
	PostTitle        string
	PostContent      string
	PostCategory     string
	PostLikeCount    int
	PostDislikeCount int
	PostCreatedAt    string
}

type Image struct {
	ImageData []byte
}

type User struct {
	Username     string
	Name         string
	Surname      string
	Email        string
	Password     string
	Biography    string
	CreationDate string
	Image        Image
}

type ChangeMessage struct {
	Message   string
	IsChanged bool
}
