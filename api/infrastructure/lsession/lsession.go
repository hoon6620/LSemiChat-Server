package lsession

import (
	"app/api/constants"
	"app/api/infrastructure/lcontext"
	"app/api/infrastructure/nosql"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// StartSession トークンを発行して、cookieにつける
func StartSession(w http.ResponseWriter, userID string) (string, error) {
	token, err := createJWT(userID)
	if err != nil {
		return "", errors.Wrap(err, "failed to create jwt")
	}
	cookie := &http.Cookie{
		Name:  constants.SessionName,
		Value: token,
	}
	http.SetCookie(w, cookie)
	return token, nil
}

// RestartSession トークンを再発行してcookieにつける
func RestartSession(w http.ResponseWriter, r *http.Request, userID string) (string, error) {
	err := deleteCookie(w, r, constants.SessionName)
	if err != nil {
		return "", errors.Wrap(err, "failed to delete current session")
	}

	return StartSession(w, userID)
}

// EndSession cookieを消す
func EndSession(w http.ResponseWriter, r *http.Request) error {
	return deleteCookie(w, r, constants.SessionName)
}

// GetSession session tokenを取得
func GetSession(r *http.Request) (*jwt.Token, error) {
	cookie, err := r.Cookie(constants.SessionName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cookie")
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JWTSecret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse jwt")
	}

	err = nosql.CheckValidToken(cookie.Value)
	if err != nil {
		return nil, errors.Wrap(err, "auth is not valid more")
	}

	return token, nil
}

func deleteCookie(w http.ResponseWriter, r *http.Request, cookieName string) error {

	userid, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		return errors.Wrap(err, "failed to authentication")
	}
	nosql.DeleteAuth(userid)

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return errors.Wrap(err, "failed to delete cookie")
	}
	cookie.MaxAge = -1

	http.SetCookie(w, cookie)
	return nil
}

func createJWT(userID string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256")) //HS256 || RSA
	rTime := time.Now().Add(time.Hour * 1)

	token.Claims = jwt.MapClaims{
		constants.JWTUserIDClaimsKey: userID,
		"exp":                        rTime.Unix(),
	}

	secretKey := constants.JWTSecret
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.Wrap(err, "failed to get jwt string")
	}

	if err = nosql.CreateAuth(userID, rTime, tokenString); err != nil {
		return "", errors.Wrap(err, "failed to set jwt in nosql")
	}
	return tokenString, nil
}
