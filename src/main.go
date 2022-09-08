package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"istorage/config"
	"istorage/controllers"
	"istorage/logger"
	"istorage/middleware"
	"istorage/s3"
	"istorage/services"
)

var envConfiguration = "server.cfg"

func init() {
	godotenv.Load()

	endpoint, _ := os.LookupEnv("FIREWALL_ENDPOINT")
	if endpoint == "" {
		logger.Infof("Upload Guard turned OFF")
	}

	services.RunMigrations()

	initConfig()
	SentryInit()
	GraylogInit()
	InitS3()
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	dbEngine := services.InitDb()

	attachmentController := &controllers.AttachmentController{dbEngine}
	fileController := &controllers.FileController{dbEngine}
	fileProwlerController := &controllers.FileProwlerController{dbEngine}

	router.GET("/__healthcheck", controllers.HealthCheck)
	router.GET("/info/:uuid", fileController.FileInfo)

	router.POST("/download", middleware.AuthGuard, fileProwlerController.ProwlFile)

	file := router.Group("/files")
	{
		file.GET("/:uuid/:profile", fileController.ReadFile, fileController.ReadFileWithResize)
		file.GET("/:uuid", fileController.ReadFile)
		file.DELETE("/:uuid", middleware.AuthGuard, fileController.DeleteFile)

		fileUpl := file.Group("/upload", middleware.AuthGuard)
		{
			fileUpl.POST("/", attachmentController.StoreAttachment)
			fileUpl.POST("/:context", attachmentController.StoreAttachmentWithContext)
		}
	}

	runServer(router)
}

func runServer(router *gin.Engine) {
	server := &http.Server{
		Addr:           config.Config.Server.Port,
		Handler:        router,
		ReadTimeout:    config.Config.Server.ReadTimeout,
		WriteTimeout:   config.Config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Infof("Start server on %s%s", config.Config.Server.Host, config.Config.Server.Port)

	server.ListenAndServe()
}

func initConfig() {
	envName := *flag.String("c", envConfiguration, "Environment config name")

	err := config.LoadConfig(envName)
	if err != nil {
		logger.Fatal(err)
	}
}

func SentryInit() {
	dsn, exists := os.LookupEnv("SENTRY_DSN")
	if !exists {
		panic("Not found SENTRY_DSN environment variable")
	}

	raven.SetDSN(dsn)
	raven.SetEnvironment(os.Getenv("APP_ENV_NAME"))
}

func GraylogInit() {
	host, exists := os.LookupEnv("GRAYLOG_HOST")
	if !exists {
		panic("Not found GRAYLOG_HOST environment variable")
	}

	port, exists := os.LookupEnv("GRAYLOG_PORT")
	if !exists {
		panic("Not found GRAYLOG_PORT environment variable")
	}

	hook := graylog.NewAsyncGraylogHook(host+":"+port, map[string]interface{}{})
	defer hook.Flush()

	logrus.AddHook(hook)
}

func InitS3() {
	bucket, exists := os.LookupEnv("AWS_S3_BUCKET")
	if !exists {
		panic("Not found AWS_S3_BUCKET environment variable")
	}

	s3host := os.Getenv("AWS_S3_HOST")

	s3.SetConfig(s3.Config{
		S3host:   s3host,
		S3bucket: bucket,
	})
}
