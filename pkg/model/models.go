package model

type People struct {
	ID         int    `json:"id" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic"`
}

type Car struct {
	ID     int    `json:"id"`
	RegNum string `json:"reg_num" validate:"required"`
	Mark   string `json:"mark" validate:"required"`
	Model  string `json:"model" validate:"required"`
	Year   int    `json:"year"`
	Owner  People `json:"owner" validate:"required"`
}
