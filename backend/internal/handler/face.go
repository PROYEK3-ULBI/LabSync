package handler

import (
	"io"
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type FaceHandler struct {
	faceService *service.FaceService
}

func NewFaceHandler(faceService *service.FaceService) *FaceHandler {
	return &FaceHandler{faceService: faceService}
}

// Enroll menangani POST /api/v1/face/enroll (Mahasiswa)
func (h *FaceHandler) Enroll(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	consentAgreedStr := c.FormValue("consent_agreed")
	consentAgreed := consentAgreedStr == "true" || consentAgreedStr == "1"

	reEnrollStr := c.FormValue("re_enroll")
	reEnroll := reEnrollStr == "true" || reEnrollStr == "1"

	frontFile, err := c.FormFile("front_image")
	if err != nil {
		return model.NewBadRequest("Berkas front_image wajib diunggah")
	}

	leftFile, err := c.FormFile("left_image")
	if err != nil {
		return model.NewBadRequest("Berkas left_image wajib diunggah")
	}

	rightFile, err := c.FormFile("right_image")
	if err != nil {
		return model.NewBadRequest("Berkas right_image wajib diunggah")
	}

	frontBytes, frontMime, err := readUploadedFile(frontFile)
	if err != nil {
		return model.NewBadRequest("Gagal membaca berkas front_image", err)
	}

	leftBytes, leftMime, err := readUploadedFile(leftFile)
	if err != nil {
		return model.NewBadRequest("Gagal membaca berkas left_image", err)
	}

	rightBytes, rightMime, err := readUploadedFile(rightFile)
	if err != nil {
		return model.NewBadRequest("Gagal membaca berkas right_image", err)
	}

	reqID := c.GetRespHeader(fiber.HeaderXRequestID)
	if reqID == "" {
		reqID = c.Get(fiber.HeaderXRequestID)
	}

	in := service.EnrollInput{
		UserID:        claims.UserID,
		RequestID:     reqID,
		FrontImg:      frontBytes,
		LeftImg:       leftBytes,
		RightImg:      rightBytes,
		FrontName:     frontFile.Filename,
		LeftName:      leftFile.Filename,
		RightName:     rightFile.Filename,
		FrontMime:     frontMime,
		LeftMime:      leftMime,
		RightMime:     rightMime,
		ConsentAgreed: consentAgreed,
		ReEnroll:      reEnroll,
	}

	res, err := h.faceService.Enroll(c.Context(), in)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetStatus menangani GET /api/v1/face/status
func (h *FaceHandler) GetStatus(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	targetID := claims.UserID
	// Jika admin menyediakan query param user_id, izinkan admin melihat status user tersebut
	if claims.Role == "ADMIN" && c.Query("user_id") != "" {
		targetID = c.Query("user_id")
	}

	res, err := h.faceService.GetStatus(c.Context(), claims.UserID, claims.Role, targetID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Revoke menangani DELETE /api/v1/face/template/:user_id (Admin only)
func (h *FaceHandler) Revoke(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	targetUserID := c.Params("user_id")
	if targetUserID == "" {
		return model.NewBadRequest("Parameter user_id tidak boleh kosong")
	}

	if err := h.faceService.Revoke(c.Context(), claims.UserID, targetUserID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Template biometrik wajah berhasil dicabut",
	})
}

func readUploadedFile(fh *multipart.FileHeader) ([]byte, string, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	mime := fh.Header.Get("Content-Type")
	if mime == "" || mime == "application/octet-stream" {
		lowerName := strings.ToLower(fh.Filename)
		if strings.HasSuffix(lowerName, ".jpg") || strings.HasSuffix(lowerName, ".jpeg") {
			mime = "image/jpeg"
		} else if strings.HasSuffix(lowerName, ".png") {
			mime = "image/png"
		} else if strings.HasSuffix(lowerName, ".webp") {
			mime = "image/webp"
		}
	}

	return data, mime, nil
}
