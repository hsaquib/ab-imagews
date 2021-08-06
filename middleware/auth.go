package middleware

import (
	"encoding/json"
	"fmt"
	rest_error "github.com/hsaquib/ab-imagews/error"
	"github.com/hsaquib/ab-imagews/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func AuthenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get(utils.AuthorizationKey)
		if jwtTkn == "" {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing access token"))
			return
		}

		jwtTkn = stripBearerFromToken(jwtTkn)

		//todo: verify token from auth

		//r.Header.Set(utils.UsernameKey, claims.Username)

		next.ServeHTTP(w, r)
	})
}

func AuthenticatedMerchantOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get(utils.AuthorizationKey)
		if jwtTkn == "" {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing access token"))
			return
		}

		jwtTkn = stripBearerFromToken(jwtTkn)

		jwtClaim := newJWTClaim()

		if !jwtClaim.TokenVerified(jwtTkn, utils.UserTypeMerchant) {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Invalid merchant token"))
			return
		}

		r.Header.Set(utils.UsernameKey, jwtClaim.Username)

		next.ServeHTTP(w, r)
	})
}

func AuthenticatedAdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get(utils.AuthorizationKey)
		if jwtTkn == "" {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing access token"))
			return
		}

		jwtTkn = stripBearerFromToken(jwtTkn)

		jwtClaim := newJWTClaim()

		if !jwtClaim.TokenVerified(jwtTkn, utils.UserTypeAdmin) {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Invalid admin token"))
			return
		}

		r.Header.Set(utils.UsernameKey, jwtClaim.Username)

		//if claims.UserType != utils.UserTypeAdmin {
		//	utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusForbidden, "Not a merchant"))
		//	return
		//}
		//r.Header.Set(utils.UsernameKey, claims.Username)

		next.ServeHTTP(w, r)
	})
}

func stripBearerFromToken(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		//log.Println("token contains Bearer")
	}

	if strings.HasPrefix(token, "bearer ") {
		token = strings.TrimPrefix(token, "bearer ")
		//log.Println("token contains bearer")
	}

	return token
}

// ClaimRes example
type ClaimRes struct {
	Success   bool     `json:"success" example:"false"`
	Status    string   `json:"status" example:"OK"`
	Message   string   `json:"message" example:"success message"`
	Timestamp string   `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      JWTClaim `json:"data"`
}

type JWTClaim struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	UserType string `json:"user_type"`
}

func newJWTClaim() *JWTClaim {
	return &JWTClaim{}
}

func (j *JWTClaim) TokenVerified(jwtTkn string, userType string) bool {
	client := http.Client{Timeout: time.Minute * 2}

	authHost := os.Getenv("AUTH_HOST")
	authUrl, err := url.Parse(authHost)
	if err != nil {
		return false
	}
	authUrl.Path = path.Join(authUrl.Path, "/api/v1/private/", fmt.Sprintf("%ss", userType), "verify-token")

	req, err := http.NewRequest(http.MethodGet, authUrl.String(), nil)
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Set("Authorization", jwtTkn)

	//log.Println("sending query:")
	// send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	//log.Println("query sent: reading resp: status:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	err = resp.Body.Close()
	if err != nil {
		log.Println(err)
		return false
	}
	//
	//log.Println("decoding body:", string(body))

	var reply ClaimRes

	err = json.Unmarshal(body, &reply)

	//log.Println(string(body))

	j.Username = reply.Data.Username
	j.UserType = reply.Data.UserType
	j.Role = reply.Data.Role

	return true
}
