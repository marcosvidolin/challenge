package entity

import "time"

type User struct {
	ID           int64      `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	MergedAt     *time.Time `json:"merged_at"`
	ParentUserID *int64     `json:"parent_user_id"`
}
