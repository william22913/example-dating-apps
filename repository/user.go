package repository

import (
	"database/sql"

	"github.com/lib/pq"
)

type CompletedUserData struct {
	User            User
	UserPreferences UserPreferences
	UserPassions    UserPassions
	Salt            Salt
	UserPremium     UserPremium
}

type User struct {
	ID          sql.NullInt64
	UUIDKey     sql.NullString
	PhoneNumber sql.NullString
	Password    sql.NullString
	FirstName   sql.NullString
	MiddleName  sql.NullString
	LastName    sql.NullString
	BirthDate   sql.NullTime
	Gender      sql.NullString
	MaxSwipe    sql.NullInt64
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
	CreatedBy   sql.NullInt64
	UpdatedBy   sql.NullInt64
	Deleted     sql.NullBool
}

type UserPreferences struct {
	ID        sql.NullInt64
	UUIDKey   sql.NullString
	UserID    sql.NullInt64
	Gender    sql.NullString
	MinAge    sql.NullInt64
	MaxAge    sql.NullInt64
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	CreatedBy sql.NullInt64
	UpdatedBy sql.NullInt64
	Deleted   sql.NullBool
}

type UserPassions struct {
	ID        sql.NullInt64
	UUIDKey   sql.NullString
	UserID    sql.NullInt64
	Tags      pq.StringArray
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	CreatedBy sql.NullInt64
	UpdatedBy sql.NullInt64
	Deleted   sql.NullBool
}

type UserPremium struct {
	ID         sql.NullInt64
	UUIDKey    sql.NullString
	UserID     sql.NullInt64
	PurchaseAt sql.NullTime
	Price      sql.NullFloat64
	EndedAt    sql.NullTime
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
	CreatedBy  sql.NullInt64
	UpdatedBy  sql.NullInt64
	Deleted    sql.NullBool
}

type Salt struct {
	ID        sql.NullInt64
	UUIDKey   sql.NullString
	UserID    sql.NullInt64
	SaltKey   sql.NullString
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	CreatedBy sql.NullInt64
	UpdatedBy sql.NullInt64
	Deleted   sql.NullBool
}
