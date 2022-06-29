package users

import (
	"database/sql"
	"log"
	"time"

	db "github.com/oci-hpc/database"
)

// addRefreshToken creates a new refresh token in the DB
func addRefreshToken(refreshToken string, userName string) error {
	// add row to refreshToken table
	// (userid, refreshToken, active)
	// (1, "hash", true)

	currentTime := time.Now()
	sqlString := `
		INSERT INTO t_user_token (
			m_username,
			m_refresh_token,
			m_active
			m_created
		) VALUES (
			:m_username,
			:m_refresh_token,
			:m_active
			:m_created
		)
	`
	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_username", userName),
		sql.Named("m_refresh_token", refreshToken),
		sql.Named("m_active", true),
		sql.Named("m_created", currentTime),
	)
	if err != nil {
		log.Printf("WARN: addRefreshToken: " + err.Error())
		return err
	}
	return nil
}

// revokeRefreshToken deactivates a refresh token in the DB
func revokeRefreshToken(userName string, refreshToken string) error {
	// lookup refreshToken
	// set active to false
	sqlString := `
		UPDATE t_user_token 
		SET m_active=false
		WHERE m_username=:m_username
		AND m_refresh_token=:m_refresh_token
	`
	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_username", userName),
		sql.Named("m_refresh_token", refreshToken),
	)
	if err != nil {
		log.Printf("WARN: revokeRefreshToken: " + err.Error())
		return err
	}
	return nil
}

// validateRefreshToken checks whether a refresh token in the DB and active
func validateRefreshToken(refreshToken string) (bool, error) {
	// lookup refresh token
	sqlString := `
	SELECT count(*)
	FROM t_user_token 
	WHERE m_active=true
	AND m_refresh_token=:m_refresh_token
`
	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_refresh_token", refreshToken),
	)
	if err != nil {
		log.Printf("WARN: validateRefreshToken: " + err.Error())
		return false, err
	}
	// if active:
	return true, nil
}
