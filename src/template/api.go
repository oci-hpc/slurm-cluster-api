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
	return
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

	var sb strings.Builder
	sb.WriteString("#!/bin/bash\n")
	for _, key := range slurmTemplate.Keys {
		sb.WriteString(strings.ToUpper(key))
		sb.WriteString("=")
		sb.WriteString("{{." + strings.ToLower(key) + "}}")
		sb.WriteString("\n")
	}
	sb.WriteString(slurmTemplate.Body)
	res := sb.String()
	slurmTemplate.Body = res
	//TODO: Make sure errors return to the front end
	insertTemplate(slurmTemplate)

	//t := template.Must(template.New("t2").Parse(res))
	//Pass in a map[string]string with keys in lowercase
	//t.Execute(os.Stdout, passIn)

	cnx.JSON(200, res)
}
