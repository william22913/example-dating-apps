package user

import (
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/william22913/example-dating-apps/repository"
)

func NewPostgresqlUserDAO(
	db *sql.DB,
) UserDAO {
	return &postgresqlUserDAO{
		db: db,
	}
}

type postgresqlUserDAO struct {
	db *sql.DB
}

func (p *postgresqlUserDAO) InsertUserData(
	data repository.CompletedUserData,
) error {
	// Start the transaction
	var err error

	tx, err := p.db.Begin()
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Error starting transaction")
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query := `
		INSERT INTO "user" (
			phone_number, password, first_name, 
			middle_name, last_name, birth_date,
			gender
		) VALUES (
			$1, $2, $3, 
			$4, $5, $6, 
			$7
		)
		RETURNING id
	`
	param := []interface{}{
		data.User.PhoneNumber.String, data.User.Password.String, data.User.FirstName.String,
		data.User.MiddleName.String, data.User.LastName.String, data.User.BirthDate.Time,
		data.User.Gender.String,
	}

	var userID int64

	err = tx.QueryRow(query, param...).Scan(&userID)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Error inserting user data")
		return err
	}

	queryPreferences := `
		INSERT INTO "user_preferences" (
			user_id, gender, min_age, 
			max_age
		) VALUES (
			$1, $2, $3, 
			$4
		)
	`

	param = []interface{}{
		userID, data.User.Gender.String, data.UserPreferences.MinAge.Int64,
		data.UserPreferences.MaxAge.Int64,
	}

	_, err = tx.Exec(queryPreferences, param...)

	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Error inserting user preferences")
		return err
	}

	queryPassions := `
		INSERT INTO "user_passions" (
			user_id, tags
		) VALUES (
			$1, $2
		)
	`

	param = []interface{}{
		userID, data.UserPassions.Tags,
	}

	_, err = tx.Exec(queryPassions, param...)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Error inserting user passions")
		return err
	}

	querySalt := `
		INSERT INTO "salt" (
			user_id, salt_key
		) VALUES (
			$1, $2
		)
	`
	param = []interface{}{
		userID, data.Salt.SaltKey.String,
	}

	_, err = tx.Exec(querySalt, param...)

	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Error inserting salt")
		return err
	}

	return nil
}

func (p *postgresqlUserDAO) CheckIsPhoneExist(
	data repository.User,
) (
	repository.User,
	error,
) {
	result := repository.User{}

	query := `
		SELECT id 
		FROM "user" 
		WHERE phone_number = $1 AND deleted = FALSE
	`

	err := p.db.QueryRow(query, data.PhoneNumber.String).Scan(&result.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}

		return result, err
	}

	return result, nil
}

func (p *postgresqlUserDAO) GetUserDataForLogin(
	data repository.User,
) (
	repository.CompletedUserData,
	error,
) {
	result := repository.CompletedUserData{}

	// Query to get user, salt, and user premium data for login
	query := `
		SELECT 
			u.id, u.uuid_key, u.phone_number, 
			u.password, u.first_name, u.middle_name, 
			u.last_name, u.birth_date, u.gender, 
			s.salt_key, up.user_id, up.purchase_at, 
			up.price, up.ended_at
		FROM "user" u
		LEFT JOIN "salt" s ON u.id = s.user_id
		LEFT JOIN "user_premium" up ON u.id = up.user_id AND up.ended_at > CURRENT_TIMESTAMP
		WHERE u.phone_number = $1 AND u.deleted = FALSE
	`

	err := p.db.QueryRow(query, data.PhoneNumber.String).Scan(
		&result.User.ID, &result.User.UUIDKey, &result.User.PhoneNumber,
		&result.User.Password, &result.User.FirstName, &result.User.MiddleName,
		&result.User.LastName, &result.User.BirthDate, &result.User.Gender,
		&result.Salt.SaltKey, &result.UserPremium.UserID, &result.UserPremium.PurchaseAt,
		&result.UserPremium.Price, &result.UserPremium.EndedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}

		log.Error().
			Err(err).
			Caller().
			Msg("Error getting user data for login")
		return result, err
	}

	return result, nil
}
