package security

/*
Gestione token JWT
*/
import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"ricantares.com/ballroom/internal/domain"
	"ricantares.com/ballroom/internal/logger"
)

type JWTtoken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

const BEARER = "Bearer"

type JwtCustomClaims struct {
	Role domain.TipoRuolo `json:"role"`
	jwt.RegisteredClaims
}

// genera il token jwt
func GeneraToken(nomeutente string, ruolo domain.TipoRuolo) (JWTtoken, error) {
	exptime, _ := strconv.ParseInt(os.Getenv("JWTEXPTIME"), 10, 32)
	claims := JwtCustomClaims{
		ruolo,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(exptime))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ballroom",
			Subject:   nomeutente,
			ID:        "1",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWTSECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return JWTtoken{}, err
	}

	exp := exptime * 3600 // in secondi
	jwttoken := JWTtoken{
		AccessToken: tokenString,
		TokenType:   BEARER,
		ExpiresIn:   exp,
	}
	return jwttoken, nil
}

func GetTokenClaims(tokenString string) (*JwtCustomClaims, error) {
	// token parsing
	secretKey := os.Getenv("JWTSECRET")
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		// controllo segnatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.LogError(fmt.Sprintf("jwt metodo di segnatura inatteso: %v", token.Header["alg"]))
			return nil, errors.New("errore")
		}
		return []byte(secretKey), nil
	})
	// controllo validita' token
	if err != nil || !token.Valid {
		exp, _ := token.Claims.GetExpirationTime()
		iss, _ := token.Claims.GetIssuedAt()
		bef, _ := token.Claims.GetNotBefore()
		now := time.Now().Local()
		logger.LogError(fmt.Sprintf("jwt token non valido: now=%v, exp=%v, iss=%v, bef=%v", now, exp, iss, bef))
		return nil, errors.New("errore")
	}
	// recupero claims
	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		logger.LogError("jwt token claims non validi")
		return nil, errors.New("errore")
	}

	return claims, nil
}

// controllo validita' temporale del token (true se il token è scaduto)
func TokenScaduto(claims *JwtCustomClaims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}

// recupero ruolo
func GetTockenRole(claims *JwtCustomClaims) domain.TipoRuolo {
	return claims.Role
}
