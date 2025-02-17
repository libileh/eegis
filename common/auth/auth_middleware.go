package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"

	"net/http"
)

func (aut *JWTAuthenticator) BasicAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				aut.HttpError.UnauthorizedErrorResponse(w, r, "authorization header is missing")
				return
			}

			//parse it -> get the base64
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Basic" {
				aut.HttpError.UnauthorizedErrorResponse(w, r, "authorization header is malformed")
				return
			}

			//decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				aut.HttpError.UnauthorizedErrorResponse(w, r, err.Error())
				return
			}

			// check the credentials
			username := aut.Properties.Username
			password := aut.Properties.Password
			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				aut.HttpError.UnauthorizedErrorResponse(w, r, fmt.Sprintf("invalid credentials: %v", err.Error()))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ExtractUserFromToken extracts the user ID and role from the JWT token in the Authorization header.
func (aut *JWTAuthenticator) ExtractUserFromToken(r *http.Request) (*CtxUser, error) {
	// Read the auth header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}

	// Split: authorization: Bearer <token>
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("authorization header is malformed")
	}

	// Get token
	tokenString := parts[1]
	jwtToken, err := aut.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Parse userId as UUID
	userIdStr, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID format")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	ctxRole, err := ExtractRole(claims)
	if err != nil {
		return nil, err
	}

	// Return the CtxUser
	return &CtxUser{
		ID:          userId,
		ContextRole: ctxRole,
	}, nil
}

// ExtractRole safely extracts a CtxRole from JWT claims.
// It expects the "role" claim to be a map with a "Value" field.
func ExtractRole(claims jwt.MapClaims) (CtxRole, error) {
	roleClaim, exists := claims["role"]
	if !exists {
		return CtxRole{}, fmt.Errorf("role claim missing")
	}
	// The role claim is expected to be a map[string]interface{}.
	roleMap, ok := roleClaim.(map[string]interface{})
	if !ok {
		return CtxRole{}, fmt.Errorf("role claim not in expected format, got: %T", roleClaim)
	}
	roleValue, ok := roleMap["Value"].(float64)
	if !ok {
		return CtxRole{}, fmt.Errorf("role value is not a number, got: %T", roleMap["Value"])
	}
	return MapToCtxRole(roleValue)
}

// AuthMiddleware returns the HTTP handler for the AuthTokenMiddleware.
func (aut *JWTAuthenticator) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from the token
		ctxUser, err := aut.ExtractUserFromToken(r)
		if err != nil {
			aut.HttpError.UnauthorizedTokenResponse(w, r, err.Error())
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), "authedUser", ctxUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
