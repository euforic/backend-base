package jwt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cristalhq/jwt/v3"
	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
)

// type for CTX jwt key
type ContextKey string

// CTX JWT key
const ContextJWTKey ContextKey = "jwt"

// Decoder ...
type Decoder struct {
	KeyServerURL string
	DevMode      bool
}

// JSONWebToken ...
type JSONWebToken struct {
	Version   int      `json:"ver"`
	ID        string   `json:"jti"`
	Issuer    string   `json:"iss"`
	Audience  string   `json:"aud"`
	IssuedAt  int      `json:"iat"`
	ExpiresAt int      `json:"exp"`
	ClientID  string   `json:"cid"`
	UserID    string   `json:"uid"`
	Scopes    []string `json:"scp"`
	Subject   string   `json:"sub"`
	Groups    []string `json:"groups"`
	IsValid   bool     `json:"is_valid"`
}

// New decoder
func New() *Decoder {
	return &Decoder{}
}

// Decode ...
func (d Decoder) Decode(t string) *JSONWebToken {
	var parsedToken JSONWebToken
	token, err := jwt.Parse([]byte(t))
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(token.RawClaims(), &parsedToken); err != nil {
		os.Stderr.WriteString(err.Error())
		return nil
	}
	return &parsedToken
}

// Verify ...
func (d Decoder) Verify(t JSONWebToken, raw string) JSONWebToken {
	var token *jwtverifier.Jwt
	validated := JSONWebToken{}
	toValidate := map[string]string{}
	claims := make(map[string]interface{})

	toValidate["aud"] = t.Audience
	toValidate["cid"] = t.ClientID

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer:           t.Issuer,
		ClaimsToValidate: toValidate,
	}

	verifier := jwtVerifierSetup.New()
	token, err := verifier.VerifyAccessToken(raw)
	if token == nil {
		token = &jwtverifier.Jwt{Claims: claims}
	}

	if err != nil {
		fmt.Print(err)
		token.Claims["is_valid"] = false
		decode(token.Claims, &validated)
		return validated
	}
	token.Claims["is_valid"] = true
	decode(token.Claims, &validated)
	return validated
}

// FromContext ...
func FromContext(ctx context.Context) (*JSONWebToken, error) {
	/*
		token := ctx.Value(ContextJWTKey)
		v, ok := token.(JSONWebToken)
		if !ok {
			return nil, errors.New("auth token missing or invalid")
		}
		return &v, nil
	*/

	return nil, nil

}

// AuthMiddleware ...
func AuthMiddleware(decoder *Decoder, next http.Handler) http.Handler {
	//TODO(euforic): possibly generate anon JWT for non logged in users
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*

			var token string
			var verifiedToken JSONWebToken

			auth := r.Header.Get("Authorization")
			tokenArr := strings.Split(auth, "Bearer ")

			if len(tokenArr) < 2 {
				// invalid or missing token so return early
				next.ServeHTTP(w, r)
				return
			}

			token = tokenArr[1]
			decodedToken := decoder.Decode(token)
			if decodedToken == nil {
				// no token so return early
				next.ServeHTTP(w, r)
				return
			}

			verifiedToken = decoder.Verify(*decodedToken, token)
			if !verifiedToken.IsValid || !sliceutil.ContainsString(verifiedToken.Groups, "group-name") {
				// token is invalid so return early
				next.ServeHTTP(w, r)
				return
			}

			// token is good so add it to the context and hit the next route
			ctx := context.WithValue(r.Context(), ContextJWTKey, verifiedToken)
			r = r.WithContext(ctx)
		*/

		next.ServeHTTP(w, r)
	})
}

func decode(in, out interface{}) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(in); err != nil {
		panic(err)
	}
	if err := json.NewDecoder(&b).Decode(out); err != nil {
		panic(err)
	}
}
