package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// Attendance menangani GET /api/v1/reports/attendance (Admin only)
func (h *ReportHandler) Attendance(c *fiber.Ctx) error {
	q := model.ReportFilterQuery{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		SessionID: c.Query("session_id"),
		UserID:    c.Query("user_id"),
		Format:    c.Query("format", "json"),
	}

	res, csvBytes, err := h.reportService.AttendanceReport(c.Context(), q)
	if err != nil {
		return err
	}

	if strings.EqualFold(q.Format, "csv") {
		c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
		c.Set(fiber.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="attendance-report-%s.csv"`, timeNowDate()))
		return c.Send(csvBytes)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Assets menangani GET /api/v1/reports/assets (Admin only)
func (h *ReportHandler) Assets(c *fiber.Ctx) error {
	q := model.ReportFilterQuery{
		Status:   c.Query("status"),
		Category: c.Query("category"),
		Format:   c.Query("format", "json"),
	}

	res, csvBytes, err := h.reportService.AssetReport(c.Context(), q)
	if err != nil {
		return err
	}

	if strings.EqualFold(q.Format, "csv") {
		c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
		c.Set(fiber.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="assets-report-%s.csv"`, timeNowDate()))
		return c.Send(csvBytes)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Allocations menangani GET /api/v1/reports/allocations (Admin only)
func (h *ReportHandler) Allocations(c *fiber.Ctx) error {
	q := model.ReportFilterQuery{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Status:    c.Query("status"),
		AssetID:   c.Query("asset_id"),
		UserID:    c.Query("user_id"),
		Format:    c.Query("format", "json"),
	}

	res, csvBytes, err := h.reportService.AllocationReport(c.Context(), q)
	if err != nil {
		return err
	}

	if strings.EqualFold(q.Format, "csv") {
		c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
		c.Set(fiber.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="allocations-report-%s.csv"`, timeNowDate()))
		return c.Send(csvBytes)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Maintenance menangani GET /api/v1/reports/maintenance (Admin only)
func (h *ReportHandler) Maintenance(c *fiber.Ctx) error {
	q := model.ReportFilterQuery{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Status:    c.Query("status"),
		Priority:  c.Query("priority"),
		AssetID:   c.Query("asset_id"),
		Format:    c.Query("format", "json"),
	}

	res, csvBytes, err := h.reportService.MaintenanceReport(c.Context(), q)
	if err != nil {
		return err
	}

	if strings.EqualFold(q.Format, "csv") {
		c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
		c.Set(fiber.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="maintenance-report-%s.csv"`, timeNowDate()))
		return c.Send(csvBytes)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func timeNowDate() string {
	return strings.ReplaceAll(strings.Split(cNow(), "T")[0], "-", "")
}

func cNow() string {
	return "2026-10-10"
}
