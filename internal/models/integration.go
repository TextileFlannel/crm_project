package models

type Integration struct {
    ID                 int    `gorm:"primaryKey;autoIncrement json:"id"`
    SecretKey          string `gorm:"type:text" json:"secret_key"`
    ClientID           string `gorm:"type:varchar(255)" json:"client_id"`
    RedirectURL        string `gorm:"type:text" json:"redirect_url"` 
    AuthenticationCode string `gorm:"type:text" json:"authentication_code"`
    Account_ID         int    `gorm:"type:int" json:"account_id"`
}
