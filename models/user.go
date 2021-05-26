package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserGroup struct {
	ID          int64     `bun:",unique,notnull"`
	CreatedAt   time.Time `bun:"default:now()"`
	UpdatedAt   time.Time
	DisplayName string `bun:",unique,notnull"`
	Description string
}

// User basic definition of a User and its meta
type User struct {
	ID             uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	CreatedAt      time.Time `bun:"default:now()"`
	UpdatedAt      time.Time
	Username       string      `bun:",notnull,unique"`
	FullName       string      `bun:",notnull"`
	Email          string      `bun:",unique,notnull"`
	FollowedGroups []uuid.UUID `bun:",type:uuid[],array"`
	OwnerOfGroups  []UserGroup `bun:"rel:has-many"`
}
