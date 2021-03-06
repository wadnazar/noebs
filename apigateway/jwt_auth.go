package gateway

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// VerifyJWT giving a jwt token and a secret it validates the token against a hard coded TokenClaims struct
func VerifyJWT(tokenString string, secret []byte) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})

	// a user might had submitted a non-jwt token
	// if err != nil {
	// 	return nil, err

	// }

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil

	} else {
		return nil, err
	}
}

func verifyWithClaim(tokenString string, secret []byte) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if token.Valid {
		fmt.Println("You look nice today")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return errors.New("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return errors.New("Timing is everything")
		} else {
			return errors.New("Couldn't handle this token:")
		}
	} else {
		return errors.New("Couldn't handle this token")
	}
	return nil
}

// GenerateJWT generates a JWT standard token with default values hardcoded. FIXME
func GenerateJWT(serviceID string, secret []byte) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	expiresAt := time.Now().Add(time.Hour * 1000).UTC().Unix()

	claims := TokenClaims{
		serviceID,
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "noebs",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	if tokenString, err := token.SignedString(secret); err == nil {
		fmt.Println(tokenString)
		return tokenString, nil
	} else {
		return "", err
	}
}

// GenerateJWTWithClaim generates a JWT standard token with default values hardcoded. FIXME
func GenerateJWTWithClaim(username string, secret []byte, tk TokenClaims) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	t := tk.Default(username)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)

	// Sign and get the complete encoded token as a string using the secret
	if tokenString, err := token.SignedString(secret); err == nil {
		fmt.Println(tokenString)
		return tokenString, nil
	} else {
		return "", err
	}
}

func generateClaims(iat, eat int64, issuer string) jwt.StandardClaims {
	claims := jwt.StandardClaims{
		IssuedAt:  iat,
		ExpiresAt: eat,
		Issuer:    issuer,
	}

	return claims
}

// TokenClaims noebs standard claim
type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Default populate token claims with default values
func (t TokenClaims) Default(username string) jwt.Claims {
	n := time.Now().Unix()
	n3h := time.Now().Add(3 * time.Hour).Unix()
	t.StandardClaims = generateClaims(n, n3h, username)
	t.Username = username
	return t
}

//secretFromClaims returns the claim's secret. in this case it is a user name
func secretFromClaims(token string, skipTime bool) (string, error) {
	claims, err := VerifyJWT(token, jwtKey)
	if e, ok := err.(*jwt.ValidationError); ok {
		if e.Errors&jwt.ValidationErrorExpired > 0 && skipTime {
			return claims.Username, nil
		} else {
			return "", errors.New("token is invalid")
		}
	} else {
		return "", errors.New("token is invalid")
	}
}
