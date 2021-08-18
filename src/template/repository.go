package template

import (
	"database/sql"
	"log"
	"strings"

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
		//templateKeys := []TemplateKey{}
		templateKey := TemplateKey{}
		var templateId sql.NullInt32
		var templateKeyId sql.NullInt32
		var templateKeyKey sql.NullString
		var templateKeyType sql.NullString
		var templateKeyDescription sql.NullString
		err := rows.Scan(
			&template.Id,
			&template.Header,
			&template.Body,
			&template.Name,
			&template.Description,
			&template.Version,
			&template.IsPublished,
			&template.OriginalId,
			&templateKeyId,
			&templateKeyKey,
			&templateKeyType,
			&templateKeyDescription,
			&templateId,
		)
		if templateKeyId.Valid {
			templateKey.Id = int(templateKeyId.Int32)
			templateKey.Key = templateKeyKey.String
			templateKey.Type = templateKeyType.String
			templateKey.Description = templateKeyDescription.String

			if err != nil {
				log.Printf("WARN: convertRowsToTemplates: " + err.Error())
			}
			//TODO: Doing this here is probably bad form
			if strings.ToUpper(templateKey.Type) == "PICKLIST" {
				templateKey.Picklist = []string{}
				pickSql := `
					SELECT 
						m_value
					FROM t_template_keys_picklist 
					WHERE m_template_keys_id = :m_template_keys_id
				`
				db := db.GetDbConnection()
				defer db.Close()
				keyRows, err := db.Query(
					pickSql,
					sql.Named("m_template_keys_id", templateKey.Id),
				)
				var pickValue string
				for keyRows.Next() {
					err := keyRows.Scan(
						&pickValue,
					)
					if err != nil {
						log.Printf("WARN: convertRowsToTemplates: " + err.Error())
					}
					templateKey.Picklist = append(templateKey.Picklist, pickValue)
				}
				if err != nil {
					log.Printf("WARN: convertRowsToTemplates: " + err.Error())
				}
			}
		}
		var lengthTemplates = len(internalTemplates)
		if lengthTemplates == 0 || internalTemplates[lengthTemplates-1].Id != template.Id {
			if templateKeyId.Valid {
				template.Keys = []TemplateKey{templateKey}
			} else {
				template.Keys = []TemplateKey{}
			}

			internalTemplates = append(internalTemplates, template)
		} else {
			internalTemplates[lengthTemplates-1].Keys = append(internalTemplates[lengthTemplates-1].Keys, templateKey)
		}
	}
	*templates = internalTemplates
}

func insertTemplate(template SlurmTemplate) {
	//TODO: this should be a transaction
	sqlString := `
		INSERT INTO t_template (
			m_header,
			m_body,
			m_name,
			m_description,
			m_version,
			m_is_published,
			m_original_id
		) VALUES (
			:m_header,
			:m_body,
			:m_name,
			:m_description,
			:m_version,
			:m_is_published,
			:m_original_id
		)
	`
	db := db.GetDbConnection()
	defer db.Close()
	res, err := db.Exec(
		sqlString,
		sql.Named("m_header", template.Header),
		sql.Named("m_body", template.Body),
		sql.Named("m_name", template.Name),
		sql.Named("m_description", template.Description),
		sql.Named("m_version", template.Version),
		sql.Named("m_is_published", template.IsPublished),
		sql.Named("m_original_id", template.OriginalId),
	)
	if err != nil {
		log.Printf("WARN: insertTemplate: " + err.Error())
		return
	}

	sqlString = `
		INSERT INTO t_template_keys (
			m_key,
			m_type,
			m_description,
			m_template_id
		) VALUES (
			:m_key,
			:m_type,
			:m_description,
			:m_template_id
		)
	`
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("WARN: insertTemplate can't get last insert ID: " + err.Error())
		return
	}
	for _, key := range template.Keys {
		keyRes, err := db.Exec(
			sqlString,
			sql.Named("m_key", key.Key),
			sql.Named("m_type", key.Type),
			sql.Named("m_description", key.Description),
			sql.Named("m_template_id", lastId),
		)
		if err != nil {
			log.Printf("WARN: insertTemplate: " + err.Error())
			return
		}

		if len(key.Picklist) > 0 {
			keyLastId, err := keyRes.LastInsertId()
			if err != nil {
				log.Printf("WARN: insertTemplate: " + err.Error())
				return
			}
			keySqlString := `
				INSERT INTO t_template_keys_picklist (
					m_template_keys_id,
					m_value
				) VALUES (
					:m_template_keys_id,
					:m_value
				)
			`
			for _, value := range key.Picklist {
				_, err := db.Exec(
					keySqlString,
					sql.Named("m_template_keys_id", keyLastId),
					sql.Named("m_value", value),
				)
				if err != nil {
					log.Printf("WARN: insertTemplate: " + err.Error())
					return
				}
			}
		}
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
			m_description = :m_description
		WHERE id = :id
	`
	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_name", template.Name),
		sql.Named("m_description", template.Description),
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

func deleteTemplate(templateId int) {
	sqlString := `
		DELETE FROM t_template WHERE id = :id
	`

	db := db.GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("id", templateId),
	)
	if err != nil {
		log.Printf("WARN: deleteTemplate: " + err.Error())
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
			t_template.m_header,
			t_template.m_body,
			t_template.m_name,
			t_template.m_description,
			t_template.m_version,
			t_template.m_is_published,
			t_template.m_original_id,
			t_template_keys.id,
			t_template_keys.m_key,
			t_template_keys.m_type,
			t_template_keys.m_description,
			t_template_keys.m_template_id
		FROM t_template
		LEFT JOIN t_template_keys ON t_template_keys.m_template_id = t_template.id
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
			t_template.m_header,
			t_template.m_body,
			t_template.m_name,
			t_template.m_description,
			t_template.m_version,
			t_template.m_is_published,
			t_template.m_original_id,
			t_template_keys.id,
			t_template_keys.m_key,
			t_template_keys.m_type,
			t_template_keys.m_description,
			t_template_keys.m_template_id
		FROM t_template
		LEFT JOIN t_template_keys ON t_template_keys.m_template_id = t_template.id
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

/*func queryTemplateKeys(templateId int) (keys []string) {
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
}*/
