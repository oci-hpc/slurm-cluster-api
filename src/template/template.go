package template

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitializeTemplateEndpoint(r *gin.Engine) {
	r.POST("/template", createTemplate)
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
	for _, str := range slurmTemplate.Variables {
		sb.WriteString(strings.ToUpper(str))
		sb.WriteString("=")
		sb.WriteString("{{." + strings.ToLower(str) + "}}")
	}
	res := sb.String()

	//t := template.Must(template.New("t2").
	//	Parse(sb.String()))

	//for _, v := range (slurm)
	//t.Execute(os.Stdout, card2)

	cnx.JSON(200, res)
}

type SlurmTemplate struct {
	Body      string
	Variables []string
}
