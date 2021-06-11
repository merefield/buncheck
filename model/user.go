package model

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserGroup struct {
	ID          uuid.UUID `bun:",unique,notnull"`
	CreatedAt   time.Time `bun:"default:now()"`
	UpdatedAt   time.Time
	DisplayName string `bun:",unique,notnull"`
	Description string
	OwnerID     uuid.UUID
	Active      bool `bun:",default:false"`
}

// User basic definition of a User and its meta
type User struct {
	ID             uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	CreatedAt      time.Time `bun:"default:now()"`
	UpdatedAt      time.Time
	Username       string         `bun:",notnull,unique"`
	FullName       string         `bun:",notnull"`
	Email          string         `bun:",unique,notnull"`
	DisplayName    sql.NullString `bun:"type:varchar(200)"`
	FollowedGroups []uuid.UUID    `bun:",type:uuid[],array"`
	OwnerOfGroups  []*UserGroup   `bun:"rel:has-many,join:id=owner_id""`
}
