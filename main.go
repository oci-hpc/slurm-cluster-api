package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	db "github.com/oci-hpc/database"
	download "github.com/oci-hpc/download"
	jobs "github.com/oci-hpc/jobs"
	nodes "github.com/oci-hpc/nodes"
	results "github.com/oci-hpc/results"
	template "github.com/oci-hpc/template"
	users "github.com/oci-hpc/users"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	initResultsFolder()

	r := gin.Default()
	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}
	config.AllowCredentials = true
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}

	r.Use(cors.New(config))

	initialize(r)
	r.Run("0.0.0.0:8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func initialize(r *gin.Engine) {
	nodes.InitializeNodesEndpoint(r)
	InitializeConfEndpoint(r)
	jobs.InitializeJobsEndpoint(r)
	results.InitializeResultsEndpoint(r)
	template.InitializeTemplateEndpoint(r)
	download.InitializeDownloadEndpoint(r)
	users.InitializeUsersEndpoint(r)
	db.InitDatabase()
	go nodes.RunNodeMonitor()
	go jobs.RunJobMonitor()
}

func initResultsFolder() {
	outputDirPath := os.Getenv("OUTPUT_DIR")
	if outputDirPath == "" {
		print("WARN: Output directory (OUTPUT_DIR) not set in .env, using /home/opc/slurm-api-output")
		outputDirPath = "/home/opc/slurm-api-output"
	}

	if _, err := os.Stat(outputDirPath); os.IsNotExist(err) {
		err = os.Mkdir(outputDirPath, 0755)
		if err != nil {
			print("ERROR: initResultsFolder: Cannot create results directory")
			log.Fatal(err.Error())
		}
	}

	/*outputJobsDir := path.Join(outputDirPath, "/jobs")
	if _, err := os.Stat(outputJobsDir); os.IsNotExist(err) {
		err = os.Mkdir(outputJobsDir, 0755)
		if err != nil {
			print("ERROR: initResultsFolder: Cannot create jobs directory")
			log.Fatal(err.Error())
		}
	}*/
}
