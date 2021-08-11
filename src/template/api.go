package template

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitializeTemplateEndpoint(r *gin.Engine) {
	r.GET("/templates", getTemplatesEndpoint)
	r.POST("/templates", createTemplate)
}

func getTemplatesEndpoint(cnx *gin.Context) {
	query := cnx.Request.URL.Query()
	if val, ok := query["templateId"]; ok {
		print(val)
	}
	templates := queryAllTemplates()
	cnx.JSON(200, templates)
	/*templateId, err := strconv.Atoi(templateIdString)
	if err != nil {
		log.Println("WARN: getTemplatesEndpoint: Invalid templateID " + err.Error())
		cnx.JSON(400, "Invalid templateId")
	}
	*/

}

func createTemplate(cnx *gin.Context) {
	jsonData, err := ioutil.ReadAll(cnx.Request.Body)

	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	var slurmTemplate SlurmTemplate
	json.Unmarshal([]byte(jsonData), &slurmTemplate)
	if slurmTemplate.Id != 0 && slurmTemplate.Version != 0 {
		slurmTemplate.Version = slurmTemplate.Version + 1
	}

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
	sb.WriteString(slurmTemplate.Body)
	res := sb.String()
	slurmTemplate.Body = res
	//TODO: Make sure errors return to the front end
	insertTemplate(slurmTemplate)

	cnx.JSON(200, res)
}
