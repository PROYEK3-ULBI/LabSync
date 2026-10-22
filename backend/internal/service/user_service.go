package service

import (
	"context"
	"errors"
	"strings"

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
)

type UserService struct {
	userRepo  repository.UserRepository
	auditRepo repository.AuditRepository
}

func NewUserService(userRepo repository.UserRepository, auditRepo repository.AuditRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

func (s *UserService) Create(ctx context.Context, adminID string, req model.CreateUserRequest) (*model.UserResponse, error) {
	req.Identifier = strings.TrimSpace(req.Identifier)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.FullName = strings.TrimSpace(req.FullName)
	req.Role = strings.TrimSpace(strings.ToUpper(req.Role))

	if len(req.Identifier) < 3 {
		return nil, model.NewBadRequest("Identifier minimal 3 karakter")
	}
	if !strings.Contains(req.Email, "@") || len(req.Email) < 5 {
		return nil, model.NewBadRequest("Format email tidak valid")
	}
	if len(req.FullName) < 2 {
		return nil, model.NewBadRequest("Nama lengkap minimal 2 karakter")
	}
	if len(req.Password) < 8 {
		return nil, model.NewBadRequest("Password minimal 8 karakter")
	}
	if req.Role != "ADMIN" && req.Role != "STUDENT" {
		return nil, model.NewBadRequest("Role harus ADMIN atau STUDENT")
	}

	status := "ACTIVE"
	if req.Status != "" {
		st := strings.ToUpper(strings.TrimSpace(req.Status))
		if st != "ACTIVE" && st != "INACTIVE" {
			return nil, model.NewBadRequest("Status harus ACTIVE atau INACTIVE")
		}
		status = st
	}

	role, err := s.userRepo.GetRoleByCode(ctx, req.Role)
	if err != nil {
		return nil, model.NewBadRequest("Role tidak ditemukan dalam sistem")
	}

	// Cek duplikasi sebelum insert
	if existing, _ := s.userRepo.GetByIdentifier(ctx, req.Identifier); existing != nil {
		return nil, model.NewConflict("Identifier sudah terdaftar")
	}
	if existing, _ := s.userRepo.GetByEmail(ctx, req.Email); existing != nil {
		return nil, model.NewConflict("Email sudah terdaftar")
	}

	hashedPwd, err := HashPassword(req.Password)
	if err != nil {
		return nil, model.NewInternalError("Gagal memproses password", err)
	}

	newUser := &model.User{
		RoleID:       role.ID,
		RoleCode:     role.Code,
		Identifier:   req.Identifier,
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: hashedPwd,
		Status:       status,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, model.NewConflict("Identifier atau email sudah terdaftar")
		}
		return nil, model.NewInternalError("Gagal menyimpan pengguna baru", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "USER_CREATED",
		EntityType: "USER",
		EntityID:   newUser.ID,
	})

	return &model.UserResponse{
		ID:         newUser.ID,
		Identifier: newUser.Identifier,
		Email:      newUser.Email,
		FullName:   newUser.FullName,
		Role:       newUser.RoleCode,
		Status:     newUser.Status,
		CreatedAt:  &newUser.CreatedAt,
		UpdatedAt:  &newUser.UpdatedAt,
	}, nil
}

func (s *UserService) List(ctx context.Context, q model.UserListQuery) (*model.UserListResponse, error) {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	users, total, err := s.userRepo.List(ctx, q)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil daftar pengguna", err)
	}

	resData := make([]model.UserResponse, len(users))
	for i, u := range users {
		faceEnrolled, _ := s.userRepo.HasFaceTemplate(ctx, u.ID)
		resData[i] = model.UserResponse{
			ID:           u.ID,
			Identifier:   u.Identifier,
			Email:        u.Email,
			FullName:     u.FullName,
			Role:         u.RoleCode,
			Status:       u.Status,
			FaceEnrolled: &faceEnrolled,
			CreatedAt:    &u.CreatedAt,
			UpdatedAt:    &u.UpdatedAt,
		}
	}

	return &model.UserListResponse{
		Data: resData,
		Meta: model.PaginationMeta{
			Page:  q.Page,
			Limit: q.Limit,
			Total: total,
		},
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, callerID, callerRole, targetID string) (*model.UserResponse, error) {
	if callerRole != "ADMIN" && callerID != targetID {
		return nil, model.NewForbidden("Akses ditolak: tidak dapat mengakses data pengguna lain")
	}

	user, err := s.userRepo.GetByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, model.NewNotFound("Pengguna tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data pengguna", err)
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

func (s *UserService) Update(ctx context.Context, callerID, callerRole, targetID string, req model.UpdateUserRequest) (*model.UserResponse, error) {
	if callerRole != "ADMIN" && callerID != targetID {
		return nil, model.NewForbidden("Akses ditolak: tidak dapat mengubah data pengguna lain")
	}

	user, err := s.userRepo.GetByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, model.NewNotFound("Pengguna tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal mengambil data pengguna", err)
	}

	// Mahasiswa tidak diizinkan mengubah role atau status miliknya
	if callerRole != "ADMIN" {
		if req.Role != nil && *req.Role != user.RoleCode {
			return nil, model.NewForbidden("Akses ditolak: hanya administrator yang dapat mengubah role")
		}
		if req.Status != nil && *req.Status != user.Status {
			return nil, model.NewForbidden("Akses ditolak: hanya administrator yang dapat mengubah status akun")
		}
	}

	if req.Identifier != "" && req.Identifier != user.Identifier {
		req.Identifier = strings.TrimSpace(req.Identifier)
		if len(req.Identifier) < 3 {
			return nil, model.NewBadRequest("Identifier minimal 3 karakter")
		}
		if existing, _ := s.userRepo.GetByIdentifier(ctx, req.Identifier); existing != nil && existing.ID != targetID {
			return nil, model.NewConflict("Identifier sudah digunakan oleh pengguna lain")
		}
		user.Identifier = req.Identifier
	}

	if req.Email != "" && req.Email != user.Email {
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if !strings.Contains(req.Email, "@") || len(req.Email) < 5 {
			return nil, model.NewBadRequest("Format email tidak valid")
		}
		if existing, _ := s.userRepo.GetByEmail(ctx, req.Email); existing != nil && existing.ID != targetID {
			return nil, model.NewConflict("Email sudah digunakan oleh pengguna lain")
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		req.FullName = strings.TrimSpace(req.FullName)
		if len(req.FullName) < 2 {
			return nil, model.NewBadRequest("Nama lengkap minimal 2 karakter")
		}
		user.FullName = req.FullName
	}

	if req.Password != nil && *req.Password != "" {
		if len(*req.Password) < 8 {
			return nil, model.NewBadRequest("Password minimal 8 karakter")
		}
		hashed, err := HashPassword(*req.Password)
		if err != nil {
			return nil, model.NewInternalError("Gagal memproses password baru", err)
		}
		user.PasswordHash = hashed
	}

	if callerRole == "ADMIN" && req.Role != nil && *req.Role != "" {
		newRoleCode := strings.ToUpper(strings.TrimSpace(*req.Role))
		if newRoleCode != user.RoleCode {
			role, err := s.userRepo.GetRoleByCode(ctx, newRoleCode)
			if err != nil {
				return nil, model.NewBadRequest("Role tidak valid")
			}
			user.RoleID = role.ID
			user.RoleCode = role.Code
		}
	}

	if callerRole == "ADMIN" && req.Status != nil && *req.Status != "" {
		newStatus := strings.ToUpper(strings.TrimSpace(*req.Status))
		if newStatus != "ACTIVE" && newStatus != "INACTIVE" {
			return nil, model.NewBadRequest("Status harus ACTIVE atau INACTIVE")
		}
		user.Status = newStatus
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, model.NewConflict("Identifier atau email sudah digunakan oleh pengguna lain")
		}
		return nil, model.NewInternalError("Gagal memperbarui data pengguna", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &callerID,
		Action:     "USER_UPDATED",
		EntityType: "USER",
		EntityID:   targetID,
	})

	return &model.UserResponse{
		ID:         user.ID,
		Identifier: user.Identifier,
		Email:      user.Email,
		FullName:   user.FullName,
		Role:       user.RoleCode,
		Status:     user.Status,
		CreatedAt:  &user.CreatedAt,
		UpdatedAt:  &user.UpdatedAt,
	}, nil
}

func (s *UserService) UpdateStatus(ctx context.Context, adminID, targetID string, req model.UpdateUserStatusRequest) (*model.UserResponse, error) {
	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status != "ACTIVE" && status != "INACTIVE" {
		return nil, model.NewBadRequest("Status harus ACTIVE atau INACTIVE")
	}

	if err := s.userRepo.UpdateStatus(ctx, targetID, status); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, model.NewNotFound("Pengguna tidak ditemukan")
		}
		return nil, model.NewInternalError("Gagal memperbarui status pengguna", err)
	}

	_ = s.auditRepo.Create(ctx, &model.AuditLog{
		UserID:     &adminID,
		Action:     "USER_STATUS_UPDATED",
		EntityType: "USER",
		EntityID:   targetID,
	})

	user, err := s.userRepo.GetByID(ctx, targetID)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil data pengguna setelah pembaruan", err)
	}

	return &model.UserResponse{
		ID:         user.ID,
		Identifier: user.Identifier,
		Email:      user.Email,
		FullName:   user.FullName,
		Role:       user.RoleCode,
		Status:     user.Status,
		CreatedAt:  &user.CreatedAt,
		UpdatedAt:  &user.UpdatedAt,
	}, nil
}
