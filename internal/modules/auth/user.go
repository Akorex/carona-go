package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gender string
type Role string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID                   uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirstName            string     `gorm:"not null" json:"first_name"`
	LastName             string     `gorm:"not null" json:"last_name"`
	Email                string     `gorm:"not null;uniqueIndex" json:"email"`
	Password             string     `gorm:"not null" json:"-"`
	Gender               Gender     `gorm:"not null" json:"gender"`
	IsVerified           bool       `gorm:"default:false" json:"is_verified"`
	Role                 Role       `gorm:"type:varchar(20);default:'user'" json:"role"`
	PasswordResetToken   *string    `gorm:"default:null" json:"-"`
	PasswordResetExpires *time.Time `gorm:"default:null" json:"-"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
