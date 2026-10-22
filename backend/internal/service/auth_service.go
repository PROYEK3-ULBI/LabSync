package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
)

type AuthService struct {
	userRepo  repository.UserRepository
	auditRepo repository.AuditRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	auditRepo repository.AuditRepository,
	jwtSecret string,
	jwtExpiryHours int,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		auditRepo: auditRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: time.Duration(jwtExpiryHours) * time.Hour,
	}
}

// HashPassword membuat bcrypt hash dari plaintext password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("gagal hashing password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword memeriksa kecocokan password dengan hash
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login memvalidasi kredensial pengguna dan mengembalikan profil serta JWT token
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.UserResponse, string, error) {
	if req.Identifier == "" || req.Password == "" {
		return nil, "", model.NewBadRequest("Identifier dan password wajib diisi")
	}

	user, err := s.userRepo.GetByIdentifier(ctx, req.Identifier)
	if err != nil && strings.Contains(req.Identifier, "@") {
		user, err = s.userRepo.GetByEmail(ctx, req.Identifier)
	}
	if err != nil {
		// Log audit percobaan gagal tanpa membocorkan identifier terdaftar atau tidak
		_ = s.auditRepo.Create(ctx, &model.AuditLog{
			Action:     "LOGIN_FAILED",
			EntityType: "AUTH",
			EntityID:   "00000000-0000-0000-0000-000000000000",
		})
		return nil, "", model.NewUnauthorized("Identifier atau password salah")
	}

	// Periksa kecocokan hash password
	if !CheckPassword(user.PasswordHash, req.Password) {
		_ = s.auditRepo.Create(ctx, &model.AuditLog{
			UserID:     &user.ID,
			Action:     "LOGIN_FAILED",
			EntityType: "AUTH",
			EntityID:   user.ID,
		})
		return nil, "", model.NewUnauthorized("Identifier atau password salah")
	}

	// Cegah login akun nonaktif
	if user.Status != "ACTIVE" {
		_ = s.auditRepo.Create(ctx, &model.AuditLog{
			UserID:     &user.ID,
			Action:     "LOGIN_FAILED_INACTIVE",
			EntityType: "AUTH",
			EntityID:   user.ID,
		})
		return nil, "", model.NewForbidden("Akun tidak aktif, silakan hubungi administrator")
	}

	// Terbitkan token JWT
	tokenString, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", model.NewInternalError("Gagal membuat sesi login", err)
	}

	// Cek status pendaftaran wajah
	faceEnrolled, _ := s.userRepo.HasFaceTemplate(ctx, user.ID)

	// Catat audit login sukses
	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &user.ID,
		Action:     "LOGIN_SUCCESS",
		EntityType: "AUTH",
		EntityID:   user.ID,
	})

	res := &model.UserResponse{
		ID:           user.ID,
		Identifier:   user.Identifier,
		Email:        user.Email,
		FullName:     user.FullName,
		Role:         user.RoleCode,
		Status:       user.Status,
		FaceEnrolled: &faceEnrolled,
		CreatedAt:    &user.CreatedAt,
		UpdatedAt:    &user.UpdatedAt,
	}

	return res, tokenString, nil
}

// Logout mencatat log audit keluar
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	if userID != "" {
		_ = s.auditRepo.Create(ctx, &model.AuditLog{
			UserID:     &userID,
			Action:     "LOGOUT",
			EntityType: "AUTH",
			EntityID:   userID,
		})
	}
	return nil
}

// Me mengembalikan data profil pengguna yang sedang login
func (s *AuthService) Me(ctx context.Context, userID string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, model.NewUnauthorized("Pengguna tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data profil", err)
	}

	if user.Status != "ACTIVE" {
		return nil, model.NewForbidden("Akun tidak aktif")
	}

	faceEnrolled, _ := s.userRepo.HasFaceTemplate(ctx, user.ID)

	return &model.UserResponse{
		ID:           user.ID,
		Identifier:   user.Identifier,
		Email:        user.Email,
		FullName:     user.FullName,
		Role:         user.RoleCode,
		Status:       user.Status,
		FaceEnrolled: &faceEnrolled,
		CreatedAt:    &user.CreatedAt,
		UpdatedAt:    &user.UpdatedAt,
	}, nil
}

func (s *AuthService) GenerateToken(u *model.User) (string, error) {
	now := time.Now()
	claims := model.Claims{
		UserID:     u.ID,
		Identifier: u.Identifier,
		Email:      u.Email,
		FullName:   u.FullName,
		Role:       u.RoleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   u.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak diharapkan: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token tidak valid")
}
