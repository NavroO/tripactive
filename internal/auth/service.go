package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NavroO/tripactive/internal/shared"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	secretKey  = []byte(os.Getenv("JWT_SECRET"))
	usersCache = shared.NewCache(5 * time.Minute)
)

type TokenClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Service interface {
	Login(ctx context.Context, body LoginRequest) (RefreshToken, error)
	Register(ctx context.Context, body RegisterRequest) (RefreshToken, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func CreateToken(userID int64) (string, error) {
	if len(secretKey) == 0 {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	claims := TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
}

func RespondWithToken(w http.ResponseWriter, token string) {
	shared.RespondJSON(w, http.StatusOK, TokenResponse{Token: token})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	fmt.Printf("token claims: %+v\n", claims)
	return claims, nil
}

type contextKey string

const ClaimsKey contextKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Fields(authHeader)
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			fmt.Println("Invalid authorization header format", authHeader)
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(bearerToken[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetSecretKeyForTests(secret string) {
	secretKey = []byte(secret)
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*TokenClaims)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}

func (s *service) Register(ctx context.Context, body RegisterRequest) (RefreshToken, error) {
	// if _, found := usersCache.Get(body.Email); found {
	// 	return "", fmt.Errorf("user already exists")
	// }

	// passwordHash, err := auth.HashPassword(body.Password)
	// if err != nil {
	// 	return RefreshToken{}, fmt.Errorf("failed to hash password: %w", err)
	// }

	// user := &User{
	// 	Username:     username,
	// 	Email:        email,
	// 	PasswordHash: passwordHash,
	// }

	// if _, err := s.repo.Register(ctx, username, email, passwordHash); err != nil {
	// 	return "", fmt.Errorf("failed to register user: %w", err)
	// }

	// usersCache.Set(email, user)

	// token, err := auth.CreateToken(user.ID)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to create token: %w", err)
	// }

	// return token, nil
}

func (s *service) Login(ctx context.Context, body LoginRequest) (RefreshToken, error) {
	// user, err := s.repo.Login(ctx, body.Email)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return RefreshToken{}, fmt.Errorf("invalid credentials")
	// 	}
	// 	return RefreshToken{}, fmt.Errorf("login error: %w", err)
	// }

	// if !CheckPasswordHash(body.Password, user.PasswordHash) {
	// 	return RefreshToken{}, fmt.Errorf("invalid credentials")
	// }

	// token, err := CreateToken(user.ID)
	// if err != nil {
	// 	return RefreshToken{}, fmt.Errorf("failed to create token: %w", err)
	// }

	// return RefreshToken{Token: token}, nil
}
