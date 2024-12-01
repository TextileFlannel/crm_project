package models

type Contact struct {
	Name string  	 `gorm:"type:varchar(255)" json:"name"`
	Email string 	 `gorm:"type:varchar(255)" json:"email"`
	AccountID int    `gorm:"type:int" json:"account_id"`
	ID   int     	 `gorm:"primaryKey;autoIncrement" json:"id"`
}


type ContactChange struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	AccountID  int    `json:"account_id"`
	ID         int    `json:"id"`
	TypeChange string `json:"type_change"`
	TaskType   string `json:"task_type"`
}


type AmoContactsResponse struct {
	Embedded struct {
		Contacts []struct {
			ID                 int    `json:"id"`
			Name               string `json:"name"`
			CustomFieldsValues []struct {
				FieldName string `json:"field_name"`
				Values    []struct {
					Value string `json:"value"`
				} `json:"values"`
			} `json:"custom_fields_values"`
		} `json:"contacts"`
	} `json:"_embedded"`
}