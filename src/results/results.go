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
	r.GET("/results/folders", getFolders)
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

	outputFileName := "slurm-" + strconv.Itoa(req.JobId) + ".out"
	outputFile := path.Join(req.WorkDir, outputFileName)
	dat, err := ioutil.ReadFile(outputFile)
	if err != nil {
		print(err.Error())
		cnx.JSON(500, "Error reading results file")
		return
	}

	files, _ := ioutil.ReadDir(req.WorkDir)
	var filesResp []ResultsFile
	for _, file := range files {
		var fileResp ResultsFile
		fileResp.FileName = file.Name()
		fileResp.Path = path.Join(req.WorkDir, file.Name())
		filesResp = append(filesResp, fileResp)
	}

	var resp ResultsResponse
	resp.JobId = req.JobId
	resp.Username = req.Username
	resp.AbsolutePath = outputFile
	resp.FileName = outputFileName
	resp.Body = string(dat)
	resp.Files = filesResp

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

	files, err := ioutil.ReadDir(req.WorkDir)
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
	var directoryTreeNode []DirectoryTreeNode
	outputDirPath := os.Getenv("OUTPUT_DIR")

	files, err := ioutil.ReadDir(outputDirPath)
	if err != nil {
		print("WARN: getFolders: " + err.Error())
		return
	}

	for _, file := range files {
		if file.IsDir() {
			var userFolder DirectoryTreeNode
			var userFolderPath = path.Join(outputDirPath, file.Name())
			recurseDirectory(userFolderPath, &userFolder)
			userFolder.Label = file.Name()
			userFolder.Path = userFolderPath
			userFolder.IsDir = true
			directoryTreeNode = append(directoryTreeNode, userFolder)
		}
	}

	cnx.JSON(200, directoryTreeNode)
}

func recurseDirectory(parentPath string, parent *DirectoryTreeNode) (err error) {
	files, err := ioutil.ReadDir(parentPath)
	if err != nil {
		print("WARN: recurseDirectory: " + err.Error())
		return
	}

	for _, file := range files {
		var child DirectoryTreeNode
		var childPath = path.Join(parentPath, file.Name())
		if file.IsDir() {
			recurseDirectory(childPath, &child)
			child.IsDir = true
		}
		child.Label = file.Name()
		child.Path = childPath
		parent.Children = append(parent.Children, child)
	}
	return
}

type ResultsRequest struct {
	JobId    int
	WorkDir  string
	Username string
}

type ResultsResponse struct {
	JobId        int
	FileName     string
	Username     string
	AbsolutePath string
	Body         string
	Files        []ResultsFile
}

type ResultsFile struct {
	FileName string
	Path     string
}

type DirectoryTreeNode struct {
	Label    string
	Path     string
	Children []DirectoryTreeNode
	IsDir    bool
}
