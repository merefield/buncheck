package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserGroup struct {
	ID          uuid.UUID `bun:",unique,notnull"`
	CreatedAt   time.Time `bun:"default:now()"`
	UpdatedAt   time.Time
	DisplayName string `bun:",unique,notnull"`
	Description string
	OwnerID     int64
}

// User basic definition of a User and its meta
type User struct {
	ID             uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	CreatedAt      time.Time `bun:"default:now()"`
	UpdatedAt      time.Time
	Username       string       `bun:",notnull,unique"`
	FullName       string       `bun:",notnull"`
	Email          string       `bun:",unique,notnull"`
	FollowedGroups []uuid.UUID  `bun:",type:uuid[],array"`
	OwnerOfGroups  []*UserGroup `bun:"rel:has-many,join:id=owner_id""`
}
