package user

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/william22913/example-dating-apps/repository"
)

func mockInsertUserData(dao UserDAO, mock sqlmock.Sqlmock) error {
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO \"user\" \( phone_number, password, first_name, middle_name, last_name, birth_date, gender \) VALUES \( \$1, \$2, \$3, \$4, \$5, \$6, \$7 \) RETURNING id`).
		WithArgs("12345", "password", "John", "Middle", "Doe", time.Now(), "Male").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Mock the rest of the queries (user preferences, passions, salt, etc.)
	mock.ExpectExec(`INSERT INTO \"user_preferences\" \( user_id, gender, min_age, max_age \) VALUES \( \$1, \$2, \$3, \$4 \)`).
		WithArgs(1, "Male", 19, 25).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO \"user_passions\" \( user_id, tags \) VALUES \( \$1, \$2 \)`).
		WithArgs(1, pq.Array([]string{"sports", "music"})).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO \"salt\" \( user_id, salt_key \) VALUES \( \$1, \$2 \)`).
		WithArgs(1, "salt123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Prepare the data
	data := repository.CompletedUserData{
		User: repository.User{
			PhoneNumber: sql.NullString{String: "12345"},
			Password:    sql.NullString{String: "password"},
			FirstName:   sql.NullString{String: "John"},
			MiddleName:  sql.NullString{String: "Middle"},
			LastName:    sql.NullString{String: "Doe"},
			BirthDate:   sql.NullTime{Time: time.Now()},
			Gender:      sql.NullString{String: "Male"},
		},
		UserPreferences: repository.UserPreferences{
			Gender: sql.NullString{String: "Male"},
			MinAge: sql.NullInt64{Int64: 19},
			MaxAge: sql.NullInt64{Int64: 25},
		},
		UserPassions: repository.UserPassions{
			Tags: []string{"sports", "music"},
		},
		Salt: repository.Salt{
			SaltKey: sql.NullString{String: "salt123"},
		},
	}

	return dao.InsertUserData(data)
}

func TestInsertUserData(t *testing.T) {
	db, mock, err := sqlmock.New() // Create a new mock DB
	require.NoError(t, err)
	defer db.Close()

	dao := &postgresqlUserDAO{db: db}

	err = mockInsertUserData(dao, mock)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %v", err)
	}
}

func TestCheckIsPhoneExist(t *testing.T) {
	db, mock, err := sqlmock.New() // Create a new mock DB
	require.NoError(t, err)
	defer db.Close()

	dao := &postgresqlUserDAO{db: db}

	// Prepare mock for the query (SELECT id FROM "user")
	mock.ExpectQuery(`SELECT id FROM \"user\" WHERE phone_number = \$1 AND deleted = FALSE`).
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Test the CheckIsPhoneExist method
	user := repository.User{PhoneNumber: sql.NullString{String: "12345"}}
	result, err := dao.CheckIsPhoneExist(user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID.Int64)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %v", err)
	}
}

func TestGetUserDataForLogin(t *testing.T) {
	// Set up sqlmock to mock the database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &postgresqlUserDAO{db: db}

	// Prepare mock query result for successful login query
	query := `
		SELECT 
			u.id, u.uuid_key, u.phone_number, 
			u.password, u.first_name, u.middle_name, 
			u.last_name, u.birth_date, u.gender, 
			s.salt_key, 
			up.user_id, up.purchase_at, up.price, 
			up.ended_at
		FROM \"user\" u
		LEFT JOIN \"salt\" s ON u.id = s.user_id
		LEFT JOIN \"user_premium\" up ON u.id = up.user_id AND up.ended_at > CURRENT_TIMESTAMP
		WHERE u.phone_number = \$1 AND u.deleted = FALSE
	`

	// Mock the expected query with the values
	mock.ExpectQuery(query).
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "uuid_key", "phone_number", "password", "first_name", "middle_name", "last_name",
			"birth_date", "gender", "salt_key", "user_id", "purchase_at", "price", "ended_at",
		}).
			AddRow(1, "uuid123", "12345", "password123", "John", "Middle", "Doe", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "M", "salt123", 1, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), 99.99, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))

	// Prepare input data for the function call
	data := repository.User{PhoneNumber: sql.NullString{String: "12345"}}

	// Test the GetUserDataForLogin function
	result, err := dao.GetUserDataForLogin(data)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.User.ID.Int64)
	assert.Equal(t, "uuid123", result.User.UUIDKey.String)
	assert.Equal(t, "12345", result.User.PhoneNumber.String)
	assert.Equal(t, "password123", result.User.Password.String)
	assert.Equal(t, "John", result.User.FirstName.String)
	assert.Equal(t, "Middle", result.User.MiddleName.String)
	assert.Equal(t, "Doe", result.User.LastName.String)
	assert.Equal(t, "1990-01-01", result.User.BirthDate.Time.Format("2006-01-02"))
	assert.Equal(t, "M", result.User.Gender.String)
	assert.Equal(t, "salt123", result.Salt.SaltKey.String)
	assert.Equal(t, "2024-01-01", result.UserPremium.PurchaseAt.Time.Format("2006-01-02"))
	assert.Equal(t, 99.99, result.UserPremium.Price.Float64)
	assert.Equal(t, "2025-01-01", result.UserPremium.EndedAt.Time.Format("2006-01-02"))

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %v", err)
	}
}
