package template

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitializeTemplateEndpoint(r *gin.Engine) {
	r.GET("/templates", getTemplatesEndpoint)
	r.POST("/templates", createTemplate)
	r.PUT("/templates", updateTemplateEndpoint)
	r.DELETE("/templates", deleteTemplateEndpoint)
}

func getTemplatesEndpoint(cnx *gin.Context) {
	templates := queryAllTemplates()
	cnx.JSON(200, templates)
	/*templateId, err := strconv.Atoi(templateIdString)
	if err != nil {
		log.Println("WARN: getTemplatesEndpoint: Invalid templateID " + err.Error())
		cnx.JSON(400, "Invalid templateId")
	}
	*/

}

func deleteTemplateEndpoint(cnx *gin.Context) {
	query := cnx.Request.URL.Query()
	val, ok := query["templateId"]
	if !ok || len(val) == 0 {
		cnx.JSON(400, "Invalid templateId")
		return
	}
	templateId, err := strconv.Atoi(val[0])
	if err != nil {
		log.Println("WARN: deleteTemplate: Invalid templateID " + err.Error())
		cnx.JSON(400, "Invalid templateId")
		return
	}
	deleteTemplate(templateId)
	cnx.JSON(200, "Template deleted")
}

func createTemplate(cnx *gin.Context) {
	slurmTemplate, err := getSlurmTemplateFromContext(cnx)
	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	slurmTemplate.Version = 1
	slurmTemplate.Header = buildTemplateHeader(slurmTemplate)
	//TODO: Make sure errors return to the front end
	insertTemplate(slurmTemplate)

	cnx.JSON(200, "Template created")
}

func updateTemplateEndpoint(cnx *gin.Context) {
	slurmTemplate, err := getSlurmTemplateFromContext(cnx)
	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	if slurmTemplate.OriginalId == 0 {
		slurmTemplate.OriginalId = slurmTemplate.Id
	}
	slurmTemplate.Version = slurmTemplate.Version + 1
	slurmTemplate.Header = buildTemplateHeader(slurmTemplate)
	//TODO: Make sure errors return to the front end
	insertTemplate(slurmTemplate)

	cnx.JSON(200, "Template updated")
}

func getSlurmTemplateFromContext(cnx *gin.Context) (slurmTemplate SlurmTemplate, err error) {
	jsonData, err := ioutil.ReadAll(cnx.Request.Body)

	if err != nil {
		log.Println("WARN: getSlurmTemplateFromContext: " + err.Error())
		return
	}

	json.Unmarshal([]byte(jsonData), &slurmTemplate)
	return
}

func buildTemplateHeader(slurmTemplate SlurmTemplate) (body string) {
	var sb strings.Builder
	sb.WriteString("#!/bin/bash\n")
	for _, key := range slurmTemplate.Keys {
		if key.Type == "Number" {
			sb.WriteString(strings.ToUpper(key.Key))
			sb.WriteString("=")
			sb.WriteString("{{." + strings.ToLower(key.Key) + "}}")
			sb.WriteString("\n")
		} else {
			sb.WriteString(strings.ToUpper(key.Key))
			sb.WriteString("=\"")
			sb.WriteString("{{." + strings.ToLower(key.Key) + "}}")
			sb.WriteString("\"\n")
		}
	}
	return sb.String()
}
