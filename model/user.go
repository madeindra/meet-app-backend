package model

type User struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

func CreateUser(data User) User {
	DB.Create(&data)
	return data
}
