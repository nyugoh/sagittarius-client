package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	io "github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
	"github.com/nyugoh/sagittarius-client/app/api"
	"github.com/nyugoh/sagittarius-client/cmd/exporter"
	"github.com/nyugoh/sagittarius-client/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	r   *gin.Engine
	app api.App
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Unable to read .env file", err.Error())
	}

	// Set up logger
	if err := utils.InitLogger(); err != nil {
		utils.LogError(err.Error())
		utils.ExitApp(1)
	}

	utils.Log("Initializing app...")

	// Set app port
	app.Port = ":5000"                                 // Default value
	if len(strings.TrimSpace(os.Getenv("PORT"))) > 0 { // Check .env file
		utils.Log("Setting app port...")
		app.Port = ":" + os.Getenv("PORT")

	}
	utils.Log("Port set to", app.Port)

	// Set app name
	if len(strings.TrimSpace(os.Getenv("APP_NAME"))) == 0 {
		utils.Log("No app name set in .env file...")
		utils.Log("Setting default app name...")
		os.Setenv("APP_NAME", "sagittarius-client")
	}
	app.Name = os.Getenv("APP_NAME")
	utils.Log("App name is:", app.Name)

	utils.Log("Done initializing app...")

	appEnv := strings.ToLower(os.Getenv("APP_ENV"))
	if len(strings.TrimSpace(appEnv)) == 0 {
		utils.Log("App env has not been set...")
		os.Setenv("APP_ENV", "DEV")
		utils.Log("Setting to development...")
	}
	utils.Log("App is running in", strings.ToUpper(appEnv), "mode...")

	// Setting app mode to gin.ReleaseMode unless in DEV or LOCAL env
	if appEnv != "dev" && appEnv != "local" {
		gin.SetMode(gin.ReleaseMode)
	}
	utils.Log("Gin server is running in", gin.Mode(), "mode")

	// Connect to DB
	db, err := utils.DbConnect()
	if err != nil {
		utils.LogError("unable to connect to DB:", err.Error())
		utils.ExitApp(1)
	}
	app.DB = db

	// Init DB, update table structures
	utils.AutoMigrateDB(app.DB)

	if ok, err := strconv.ParseBool(os.Getenv("INIT_DB")); ok && err == nil {
		utils.SeedDB(app.DB) // Insert initial data required by the app to start
	}

}

// CORS Middleware
func CORS(c *gin.Context) {
	sagServer := os.Getenv("SAG_SERVER")
	c.Header("Access-Control-Allow-Origin", sagServer)
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func Run() {
	// Make a gin router
	r = gin.Default()

	r.Use(CORS)

	// Register all routes
	initRoutes()

	// Close DB on app close or panic
	defer app.DB.Close()

	// Socket.io server
	server, err := io.NewServer(nil)
	app.SocketServer = server
	if err != nil {
		utils.LogError("Unable to start socket.io realtime server")
		utils.ExitApp(1)
	}

	r.GET("/socket.io/*any", gin.WrapH(app.SocketServer))
	r.POST("/socket.io/*any", gin.WrapH(app.SocketServer))

	go app.IOServer()
	go app.SocketServer.Serve()
	defer app.SocketServer.Close()

	app.Login()
	app.Folders = exporter.StartMonitor()

	utils.Log("App folders:", app.Folders)

	utils.Log("Starting app...")
	utils.Log(fmt.Sprintf("Magic brewing on port %s", app.Port))
	if err := r.Run(app.Port); err != nil {
		utils.LogError("App terminated: ", err.Error())
	}
}

func initRoutes() {
	// Add middleware to monitor all request, /metric endpoint for analytics
	r.Use(utils.MetricsMonitor())

	r.GET("/", app.Index)
	r.GET("/logs", app.ListFolders)
	r.GET("/read", app.ReadLog)

	/*authRouter := r.Group("/auth")
	{
		authRouter.GET("/login", app.Login)
	}*/

}
