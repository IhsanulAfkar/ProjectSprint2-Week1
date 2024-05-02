package models

type Matcher struct {
	PkId        int    `gorm:"column:pkId" json:"pkId"`
	Id          string `json:"id"`
	UserId      int    `gorm:"column:userId" json:"userId"`
	MatchUserId int    `gorm:"column:matchUserId" json:"matchUserId"`
	UserCatId   int    `gorm:"column:userCatId" json:"userCatId"`
	MatchCatId  int    `gorm:"column:matchCatId" json:"matchCatId"`
	Message     string `json:"message"`
	IsApproved  bool   `json:"isApproved" gorm:"column:isApproved"`
	IsValid     bool   `json:"isValid" gorm:"column:isValid"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}