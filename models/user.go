package models

type User struct {
	PkId        int    `gorm:"column:pkId" json:"pkId"`
	Id          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	AccessToken string `json:"accessToken"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}