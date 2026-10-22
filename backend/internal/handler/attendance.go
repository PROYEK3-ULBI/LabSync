package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type AttendanceHandler struct {
	attendanceService *service.AttendanceService
}

func NewAttendanceHandler(attendanceService *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{attendanceService: attendanceService}
}

// CheckIn menangani POST /api/v1/attendance/check-in (Mahasiswa)
func (h *AttendanceHandler) CheckIn(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	sessionID := strings.TrimSpace(c.FormValue("session_id"))
	if sessionID == "" {
		return model.NewBadRequest("Parameter session_id wajib diisi")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return model.NewBadRequest("Berkas citra wajah live (image) wajib diunggah")
	}

	imgBytes, mimeType, err := readUploadedFile(file)
	if err != nil {
		return model.NewBadRequest("Gagal membaca berkas citra live", err)
	}

	reqID := c.GetRespHeader(fiber.HeaderXRequestID)
	if reqID == "" {
		reqID = c.Get(fiber.HeaderXRequestID)
	}

	in := service.CheckInInput{
		UserID:    claims.UserID,
		SessionID: sessionID,
		LiveImg:   imgBytes,
		Filename:  file.Filename,
		MimeType:  mimeType,
		RequestID: reqID,
	}

	att, err := h.attendanceService.CheckIn(c.Context(), in)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.CheckInResponse{
		Message:    "Presensi masuk (check-in) berhasil",
		Attendance: *att,
	})
}

// CheckOut menangani POST /api/v1/attendance/check-out (Mahasiswa)
func (h *AttendanceHandler) CheckOut(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	var sessionID string
	var req model.CheckOutRequest
	if err := c.BodyParser(&req); err == nil && req.SessionID != "" {
		sessionID = req.SessionID
	} else {
		sessionID = c.FormValue("session_id")
	}

	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return model.NewBadRequest("Parameter session_id wajib diisi")
	}

	att, err := h.attendanceService.CheckOut(c.Context(), claims.UserID, sessionID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.CheckOutResponse{
		Message:    "Presensi keluar (check-out) berhasil",
		Attendance: *att,
	})
}

// History menangani GET /api/v1/attendance/history (Mahasiswa)
func (h *AttendanceHandler) History(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, err := h.attendanceService.GetHistory(c.Context(), claims.UserID, page, limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// List menangani GET /api/v1/attendance (Admin only)
func (h *AttendanceHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	q := model.AttendanceListQuery{
		Page:      page,
		Limit:     limit,
		SessionID: c.Query("session_id"),
		UserID:    c.Query("user_id"),
		Date:      c.Query("date"),
	}

	res, err := h.attendanceService.List(c.Context(), q)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetByID menangani GET /api/v1/attendance/:id (Admin atau Pemilik)
func (h *AttendanceHandler) GetByID(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := c.Params("id")
	if id == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	res, err := h.attendanceService.GetByID(c.Context(), claims.UserID, claims.Role, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
