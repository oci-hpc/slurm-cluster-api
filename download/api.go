package download

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

func InitializeDownloadEndpoint(r *gin.Engine) {
	r.POST("/download", FileDownload)
}

func FileDownload(cnx *gin.Context) {
	jsonData, err := ioutil.ReadAll(cnx.Request.Body)

	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	var req DownloadRequest
	json.Unmarshal([]byte(jsonData), &req)

	file, err := getFile(req.Path)
	if err != nil {
		cnx.JSON(400, "Cannot open file: "+req.Path)
		return
	}
	fileStat, err := file.Stat()
	if err != nil {
		cnx.JSON(500, "Cannot get file size: "+req.Path)
		return
	}
	cnx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name())) //fmt.Sprintf("attachment; filename=%s", filename) Downloaded file renamed
	cnx.Writer.Header().Add("Content-Length", fmt.Sprintf("%d", fileStat.Size()))
	cnx.Writer.Header().Add("Content-Type", "application/octet-stream")
	file.Close()
	cnx.File(req.Path)
}

func getFile(path string) (file *os.File, err error) {
	file, err = os.Open(path)
	if err != nil {
		return
	}
	return
}

type DownloadRequest struct {
	Path string
}
