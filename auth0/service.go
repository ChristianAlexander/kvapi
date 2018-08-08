package auth0

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"github.com/auth0-community/auth0"
	jose "gopkg.in/square/go-jose.v2"
)

// Service provides request validation via Auth0.
type Service struct {
	audience []string
	issuer   string
}

// NewService returns a new Auth0 Service.
func NewService(audience []string, issuer string) *Service {
	return &Service{
		audience,
		issuer,
	}
}

func (s *Service) validator() *auth0.JWTValidator {
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: fmt.Sprintf("%s.well-known/jwks.json", s.issuer)}, nil)
	configuration := auth0.NewConfiguration(client, s.audience, s.issuer, jose.RS256)

	return auth0.NewValidator(configuration, nil)
}

// Middleware produces a mux.MiddlewareFunc that requires the request to be authenticated.
func (s *Service) Middleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			validator := s.validator()
			token, err := validator.ValidateRequest(r)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Token is not valid:", token)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			claims := map[string]interface{}{}
			err = validator.Claims(r, token, &claims)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized, invalid claims"))
				fmt.Println("Invalid claims:", err)
				return
			}

			context.Set(r, "uid", claims["sub"])
			next.ServeHTTP(w, r)
		})
	}
}
