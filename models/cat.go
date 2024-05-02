package models

import "github.com/lib/pq"

var CatRaces = [10]string{"Persian",
	"Maine Coon",
	"Siamese",
	"Ragdoll",
	"Bengal",
	"Sphynx",
	"British Shorthair",
	"Abyssinian",
	"Scottish Fold",
	"Birman"}

var Sex = [2]string{
	"male",
	"female"}

type Cat struct {
	PkId        int       `gorm:"column:pkId" json:"pkId"`
	Id          string    `gorm:"column:id" json:"id"`
	UserId      int       `gorm:"column:userId" json:"userId"`
	Name        string    `json:"name"`
	Race        string    `json:"race"`
	Sex         string    `json:"sex"`
	AgeInMonth  int       `gorm:"column:ageInMonth" json:"ageInMonth"`
	Description string    `json:"description"`
	HasMatched  bool      `json:"hasMatched"`
	ImageUrls   pq.StringArray `json:"imageUrls" gorm:"type:text[];column:imageUrls"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
}

type GetCat struct {
	Id          string   `gorm:"column:id"`
	Name        string   `json:"name"`
	Race        string   `json:"race"`
	Sex         string   `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	Description string   `json:"description"`
	HasMatched  bool     `json:"hasMatched"`
	ImageUrls    pq.StringArray `json:"imageUrls" gorm:"type:text[];column:imageUrls"`
	CreatedAt   string   `json:"created_at"`
}