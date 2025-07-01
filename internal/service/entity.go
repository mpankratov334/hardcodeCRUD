package service

type PostRequest struct {
	Title  string `json:"title" validate:"required"`
	Data   string `json:"data"`
	Status string `json:"status"`
	UserID string `json:"user_id" validate:"required"`
}

type RequestWithId struct {
	ID string `validate:"required,intString,min=1"`
}

type UpdateRequest struct {
	Status string `json:"status" validate:"required"`
	ID     string `validate:"required,intString,min=1"`
}

type PostUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RequestWithUserName struct {
	Name string `json:"name" validate:"required"`
}
