package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"ricantares.com/ballroom/src/internal/api"
	"ricantares.com/ballroom/src/internal/db"
	"ricantares.com/ballroom/src/internal/logger"
	"ricantares.com/ballroom/src/internal/security"
)

// main avvia l'applicazione web in base al parametro di lancio
// lancia il server web in modalita' normale se non specificato alcun parametro
// lancia il server web in modalita' di test se viene specificato il parametro "test"
func main() {
	// Caricamento environment
	err := loadEnvironment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore di caricamento dell'ambiente: %v", err)
		panic(err)
	}

	// set timezone
	loc, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore di impostazione TIMEZONE: %v", err)
		panic(err)
	}
	time.Local = loc

	// Inizializzazione logger
	outfile, err := os.Create(os.Getenv("LOGFILE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore di caricamento del file di log: %v", err)
		panic(err)
	}
	defer outfile.Close()
	logger.NewLog(outfile, os.Getenv("LOGLEVEL"))

	// Inizializzazione RBAC
	err = security.GetRbac(os.Getenv("RBACFILE"))
	if err != nil {
		logger.LogError(fmt.Sprintf("Errore RBAC %v", err))
		panic(err)
	}

	// Avvio applicazione
	args := os.Args
	op := "server"
	if len(args) > 1 {
		op = args[1]
		logger.LogInfo(fmt.Sprintf("Avvio applicazione: %v", op))
	}
	if err := run(op); err != nil {
		logger.LogError(fmt.Sprintf("Problemi nell'avvio del server. - err: %v", err))
	}

}

func loadEnvironment() error {
	candidates := []string{".env", filepath.Join("src", ".env")}
	var envFile string
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			envFile = candidate
			break
		}
	}
	if envFile == "" {
		return fmt.Errorf("file .env non trovato (cercati .env e src/.env)")
	}

	if err := godotenv.Load(envFile); err != nil {
		return err
	}

	baseDir := filepath.Dir(envFile)
	if baseDir == "." {
		return nil
	}
	return os.Chdir(baseDir)
}

/*************  ✨ Windsurf Command ⭐  *************/
// run starts the application according to the given op parameter.
// If op is "test", it performs a simple test using a mock db.
// Otherwise, it connects to the configured database, sets up the server,
// and starts the service.
// An error is returned if there is a problem with the database connection
// or with starting the server.
/*******  91723308-046f-483a-a1a9-cacc5e08f994  *******/
func run(op string) error {

	if op == "test" {
		/*
			mockDb := &mock.MockDb{}
			repo := db.NewRepository(mockDb)
			data, _ := repo.GetScuola()
			logger.LogDebug(fmt.Sprintf("mock scuola %v", data))
		*/
		return nil
	}

	// Connessione alla base dati
	var databaseUrl = os.Getenv("DBCONNECT")
	var ctxBg = context.Background()
	conn, err := db.NewDbConnection(ctxBg, databaseUrl)
	if err != nil {
		logger.LogError(fmt.Sprintf("Problemi nella connessione alla base dati: %v - err: %v", databaseUrl, err))
		return err
	}
	logger.LogDebug(fmt.Sprintf("Connesso al db %v", databaseUrl))
	repo := db.NewRepository(conn)
	defer conn.Close()

	// Impostazione server/routes/handlers(logger, security)
	gin.SetMode(os.Getenv("GIN_MODE"))
	//router := gin.Default()
	router := gin.New()
	router.SetTrustedProxies(nil)
	router.StaticFile("/favicon.ico", os.Getenv("FAVICON"))

	// Impostazione CORS
	// Deve essere impostata prima del routing altrimenti non funziona
	router.Use(CORS)

	routes := api.NewRouteHandler(*repo)
	api.LoggerWrapper(router)
	api.Routes(*routes, router)

	host := os.Getenv("HTTPHOST")
	if len(host) == 0 {
		host = getLocalIP()
	}
	port := os.Getenv("HTTPPORT")
	server := host + ":" + port
	router.Run(server)

	return nil
}

/*************  ✨ Windsurf Command ⭐  *************/
// getLocalIP trova l'IP locale non loopback
// Se non riesce, restituisce una stringa vuota
/*******  f658117f-e759-4a86-9301-d93eb72e82aa  *******/
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback then display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// Imposta la configurazione CORS
func CORS(c *gin.Context) {

	if os.Getenv("HDR_ALLOW_ORIGIN") == "same" {
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
	} else {
		c.Header("Access-Control-Allow-Origin", os.Getenv("HDR_ALLOW_ORIGIN"))
	}
	c.Header("Access-Control-Allow-Credentials", os.Getenv("HDR_ALLOW_CREDENTIALS"))
	c.Header("Access-Control-Allow-Headers", os.Getenv("HDR_ALLOW_HEADERS"))
	c.Header("Access-Control-Allow-Methods", os.Getenv("HDR_ALLOW_METHODS"))
	c.Header("Content-Type", os.Getenv("HDR_CONTENT_TYPE"))

	if c.Request.Method == "OPTIONS" {
		logger.LogDebug(fmt.Sprintf("Request: %v", c.Request.Method))
		logger.LogDebug(fmt.Sprintf("Request Header: %v", c.Request.Header))
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()

}
