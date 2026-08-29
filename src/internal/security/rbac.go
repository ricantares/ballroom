package security

/*
	Gestione delle autorizzazioni "Role Based Access Control"

	Il file di configurazione "rbac.json" definisce gli accessi. Di seguito un estratto

	{
	"rbac": [
        {...},
        {...},
        {
            "pattern": "/scuola/sala/\\d*",
            "permissions": {
                "can_read": ["Admin", "Direzione", "Iscritto", "Maestro", "Staff"],
                "can_write": ["Admin"],
                "can_delete": ["Admin"]
            }
        }
    	]
	}

	"pattern" contiene una regex che corrisponde alla risorsa richiesta
	"permissions" contiene i ruoli autorizzati ad accedere alla risorsa raggruppati per metodo:
		- can_read: GET
		- can_write: POST, PUT
		- can_delete: DELETE
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	dom "ricantares.com/ballroom/src/internal/domain"
	"ricantares.com/ballroom/src/internal/logger"
)

type Utente struct {
	dom.Model
}

type Rbac struct {
	Resources []RbacResource `json:"rbac,omitempty"`
}
type RbacResource struct {
	Pattern     string     `json:"pattern,omitempty"`
	Permissions Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Can_read   []string `json:"can_read,omitempty"`
	Can_write  []string `json:"can_write,omitempty"`
	Can_delete []string `json:"can_delete,omitempty"`
}

var rbacData Rbac

func GetRbac(rbacjson string) error {
	configFile, err := os.Open(rbacjson)
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&rbacData); err != nil {
		return err
	}

	logger.LogDebug(fmt.Sprintf("RBAC: %v", rbacData))
	return nil
}

func AccessGranted(c *gin.Context) bool {
	header := c.GetHeader("Authorization")
	if header == "" {
		logger.LogError(fmt.Sprintln("Header 'Authorization' non fornita"))
		return false
	}
	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 || jwtToken[0] != "Bearer" {
		logger.LogError(fmt.Sprintln("Header 'Authorization' malformata"))
		return false
	}

	tokenString := jwtToken[1]
	claims, err := GetTokenClaims(tokenString)
	if err != nil {
		return false
	}
	role := claims.Role.String()
	method := c.Request.Method
	url := c.Request.URL

	logger.LogDebug(fmt.Sprintf("RBAC: %v, %v, %v", role, method, url))

	if isGranted(role, method, url.String()) {
		c.Set("Role", role)
		return true
	}
	return false
}

func isGranted(role string, method string, url string) bool {
	for _, rbac := range rbacData.Resources {
		matched, _ := regexp.MatchString(rbac.Pattern, url)
		if matched {
			roles := roleGranted(method, rbac)
			for _, r := range roles {
				if r == role {
					return true
				}
			}
		}
	}

	return false
}

func roleGranted(method string, rbac RbacResource) []string {
	switch method {
	case http.MethodGet:
		return rbac.Permissions.Can_read
	case http.MethodPost, http.MethodPut:
		return rbac.Permissions.Can_write
	case http.MethodDelete:
		return rbac.Permissions.Can_delete
	}

	return nil
}
