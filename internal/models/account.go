package models

type Account struct {
    ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
    AccessToken     string `gorm:"type:text;not null" json:"access_token"`
    RefreshToken    string `gorm:"type:text;not null" json:"refresh_token"`
    Expires         int    `gorm:"type:int" json:"expires"`
    Subdomain       string `gorm:"type:varchar(100);not null" json:"subdomain"`
    AccountID       int    `gorm:"type:int" json:"account_id"`
	UnisenderKey    string `gorm:"type:text" json:"unisender_key"`
}