package users

import (
	"database/sql"
	"log"

	db "github.com/oci-hpc/database"
)

func addRefreshToken(refreshToken string, userId string) error {
	// add row to refreshToken table
	// (userid, refreshToken, active)
	// (1, "hash", true)

	sqlString := `
		INSERT INTO t_user_token (
			m_user_id,
			m_refresh_token,
			m_active
		) VALUES (
			:m_user_id,
			:m_refresh_token,
			:m_active
		)
	`
	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_user_id", userId),
		sql.Named("m_refresh_token", refreshToken),
		sql.Named("m_active", true),
	)
	if err != nil {
		log.Printf("WARN: addRefreshToken: " + err.Error())
		return err
	}
	return nil
}

func revokeRefreshToken(userId string, refreshToken string) error {
	// lookup refreshToken
	// set active to false
	sqlString := `
		UPDATE t_user_token 
		SET m_active=false
		WHERE m_user_id=:m_user_id
		AND m_refresh_token=:m_refresh_token
	`
	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_user_id", userId),
		sql.Named("m_refresh_token", refreshToken),
	)
	if err != nil {
		log.Printf("WARN: revokeRefreshToken: " + err.Error())
		return err
	}
	return nil
}

func validateRefreshToken(refreshToken string) (bool, error) {
	// lookup refresh token
	sqlString := `
	SELECT count(*)
	FROM t_user_token 
	WHERE m_user_id=:m_user_id
	AND m_active=true
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
