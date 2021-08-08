package template

import (
	"database/sql"
	"log"

	db "github.com/oci-hpc/slurm-cluster-api/src/database"
)

func convertRowsToTemplates(rows *sql.Rows, templates *[]SlurmTemplate) {
	if rows == nil {
		return
	}
	defer rows.Close()

	internalTemplates := []SlurmTemplate{}
	for rows.Next() {
		template := SlurmTemplate{}
		templateId := 0
		var key string
		err := rows.Scan(
			&template.Id,
			&template.Body,
			&template.Name,
			&key,
			&templateId,
		)
		if err != nil {
			log.Printf("WARN: convertRowsToTemplates: " + err.Error())
		}
		var lengthTemplates = len(internalTemplates)
		if lengthTemplates == 0 || internalTemplates[lengthTemplates-1].Id != template.Id {
			template.Keys = []string{key}
			internalTemplates = append(internalTemplates, template)
		} else {
			internalTemplates[lengthTemplates-1].Keys = append(internalTemplates[lengthTemplates-1].Keys, key)
		}
	}
	*templates = internalTemplates
}

func insertTemplate(template SlurmTemplate) {
	//TODO: this should be a transaction
	sqlString := `
		INSERT INTO t_template (
			m_body,
			m_name
		) VALUES (
			:m_body,
			:m_name
		)
	`
	db := db.GetDbConnection()
	defer db.Close()
	res, err := db.Exec(
		sqlString,
		sql.Named("m_body", template.Body),
		sql.Named("m_name", template.Name),
	)
	if err != nil {
		log.Printf("WARN: insertTemplate: " + err.Error())
		return
	}

	sqlString = `
		INSERT INTO t_template_keys (
			m_key,
			m_template_id
		) VALUES (
			:m_key,
			:m_template_id
		)
	`
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("WARN: insertTemplate can't get last insert ID: " + err.Error())
		return
	}
	for _, key := range template.Keys {
		_, err = db.Exec(
			sqlString,
			sql.Named("m_key", key),
			sql.Named("m_template_id", lastId),
		)
	}
	if err != nil {
		log.Printf("WARN: insertTemplate: " + err.Error())
		return
	}
}

func updateTemplate(template SlurmTemplate) {
	sqlString := `
		UPDATE t_template
		SET
			m_name = :m_name,
			m_body = :m_body
		WHERE id = :id
	`
	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_name", template.Name),
		sql.Named("m_body", template.Body),
		sql.Named("id", template.Id),
	)
	if err != nil {
		log.Printf("WARN: updateTemplate: " + err.Error())
	}
}

func deleteTemplateKey(templateId int, key string) {
	sqlString := `
		DELETE FROM t_template_keys WHERE m_template_id = :m_template_id AND m_key = :m_key
	`

	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_key", key),
		sql.Named("m_template_id", templateId),
	)
	if err != nil {
		log.Printf("WARN: deleteTemplateKey: " + err.Error())
	}
}

func insertTemplateKey(templateId int, key string) {
	sqlString := `
		INSERT INTO t_template_keys (
			m_key,
			m_template_id
		) VALUES (
			:m_key,
			:m_template_id
		)
	`
	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_key", key),
		sql.Named("m_template_id", templateId),
	)
	if err != nil {
		log.Printf("WARN: insertTemplateKey: " + err.Error())
	}
}

func queryAllTemplates() (templates []SlurmTemplate) {
	sqlString := `
		SELECT 
		  t_template.id,
			t_template.m_body,
			t_template.m_name,
			t_template_keys.m_key,
			t_template_keys.m_template_id
		FROM t_template
		JOIN t_template_keys ON t_template_keys.m_template_id = t_template.id
	`
	db := db.GetDbConnection()
	defer db.Close()
	rows, err := db.Query(sqlString)
	if err != nil {
		log.Printf("WARN: queryAllTemplates: " + err.Error())
	}
	convertRowsToTemplates(rows, &templates)
	return
}

func QueryTemplateById(id int) (template SlurmTemplate) {
	sqlString := `
		SELECT 
		  t_template.id,
			t_template.m_body,
			t_template.m_name,
			t_template_keys.m_key,
			t_template_keys.m_template_id
		FROM t_template
		JOIN t_template_keys ON t_template_keys.m_template_id = t_template.id
		WHERE t_template.id = :id
	`
	db := db.GetDbConnection()
	defer db.Close()
	rows, err := db.Query(sqlString, sql.Named("id", id))
	if err != nil {
		log.Printf("WARN: queryAllTemplates: " + err.Error())
	}
	var templates []SlurmTemplate
	convertRowsToTemplates(rows, &templates)
	template = templates[0]
	return
}

func queryTemplateKeys(templateId int) (keys []string) {
	sqlString := `
		SELECT 
		  m_key
		FROM t_template
		WHERE m_template_id = :m_template_id;
	`
	db := db.GetDbConnection()
	defer db.Close()
	rows, err := db.Query(sqlString, sql.Named("m_template_id", templateId))
	if err != nil {
		log.Printf("WARN: queryTemplateKeys: " + err.Error())
	}
	keys = []string{}
	for rows.Next() {
		var key string
		err := rows.Scan(
			&key,
		)
		if err != nil {
			log.Printf("WARN: convertRowsToTemplates: " + err.Error())
		}
		keys = append(keys, key)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("WARN: queryTemplateKeys: " + err.Error())
	}
	return
}
