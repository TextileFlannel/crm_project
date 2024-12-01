package models

type Task struct {
	AccountID int    `json:"account_id"`
	UnisenderKey string `json:"unisender_key"`
	TaskType     string `json:"task_type"`
}