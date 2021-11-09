package middleware

import (
	"app/api/constants"
	"app/api/infrastructure/lcontext"
	"app/api/infrastructure/lsession"
	"app/api/llog"
	"app/api/presentation/response"
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// var mode *string
var (
	allowHeaders = "*"
	allowOrigin  = "*"
)

type mode string

const (
	develop    mode = "develop"
	production mode = "production"
)

func init() {
	modeFlag := flag.String("mode", "production", "run mode. value=[develop, production]")
	flag.Parse()
	if *modeFlag == string(develop) {
		allowHeaders = "Content-Type"
		allowOrigin = "http://localhost:3000"
	}
	llog.Info(fmt.Sprintf("run mode is %s", *modeFlag))
}

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// CORS
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		token, err := lsession.GetSession(r)
		if err != nil {
			response.Unauthorized(w, errors.Wrap(err, "failed to get session"), "failed to authorization")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		ctx = lcontext.SetUserID(ctx, claims[constants.JWTUserIDClaimsKey].(string))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	//TODO AdminMidleware
	return nil
}
