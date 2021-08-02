package results

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

func InitializeResultsEndpoint(r *gin.Engine) {
	r.POST("/results", getResults)
	r.POST("/results/files", getFiles)
	r.POST("/results/folders", getFolders)
}

func getResults(cnx *gin.Context) {
	jsonData, err := ioutil.ReadAll(cnx.Request.Body)

	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	var req ResultsRequest
	json.Unmarshal([]byte(jsonData), &req)
	if req.Username == "" {
		cnx.JSON(400, "Must provide username")
		return
	}

	if req.JobId == 0 {
		cnx.JSON(400, "Invalid job ID")
		return
	}

	outputDirPath := os.Getenv("OUTPUT_DIR")
	outputUserPath := path.Join(outputDirPath, req.Username)
	outputFileName := "slurm-" + strconv.Itoa(req.JobId) + ".out"
	outputFile := path.Join(outputUserPath + "/" + outputFileName)
	dat, err := ioutil.ReadFile(outputFile)
	if err != nil {
		print(err.Error())
		cnx.JSON(500, "Error reading results file")
		return
	}

	var resp ResultsResponse
	resp.JobId = req.JobId
	resp.Username = req.Username
	resp.AbsolutePath = outputFile
	resp.FileName = outputFileName
	resp.Body = string(dat)

	cnx.JSON(200, resp)
}

func getFiles(cnx *gin.Context) {
	jsonData, err := ioutil.ReadAll(cnx.Request.Body)

	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	var req ResultsRequest
	json.Unmarshal([]byte(jsonData), &req)

	outputDirPath := os.Getenv("OUTPUT_DIR")
	outputUserPath := path.Join(outputDirPath, req.Username)
	files, err := ioutil.ReadDir(outputUserPath)
	if err != nil {
		print(err.Error())
		cnx.JSON(500, "Error reading user directory")
		return
	}

	var res []string
	for _, file := range files {
		if !file.IsDir() {
			res = append(res, file.Name())
		}
	}

	cnx.JSON(200, res)
}

func getFolders(cnx *gin.Context) {
	outputDirPath := os.Getenv("OUTPUT_DIR")

	files, err := ioutil.ReadDir(outputDirPath)
	if err != nil {
		print(err.Error())
		cnx.JSON(500, "Error reading output directory")
		return
	}

	var folders []string
	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}

	cnx.JSON(200, folders)
}

type ResultsRequest struct {
	JobId    int
	Username string
}

type ResultsResponse struct {
	JobId        int
	FileName     string
	Username     string
	AbsolutePath string
	Body         string
}
