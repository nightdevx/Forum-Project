package main

type postData struct {
	UserData User
	PostData Post
}

type Post struct {
	PostTitle        string
	PostContent      string
	PostImage        string
	PostCategory     string
	PostLikeCount    int
	PostDislikeCount int
	PostCreatedAt    string
}

type User struct {
	Username     string
	Name         string
	Surname      string
	Email        string
	Password     string
	Biography    string
	CreationDate string
	BannerImage  string
	ProfileImage string
}

type Image struct {
	ImageData []byte
}
