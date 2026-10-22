package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"labsync-backend/internal/model"
)

type FaceClient interface {
	Extract(ctx context.Context, imgBytes []byte, filename, contentType, pose, requestID string) (*model.PythonExtractResponse, error)
	Verify(ctx context.Context, imgBytes []byte, filename, contentType string, enrolledSamples []model.EnrolledSamplePayload, threshold *float64, strategy, requestID string) (*model.PythonVerifyResponse, error)
}

type HTTPFaceClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewHTTPFaceClient(baseURL, apiKey string, timeoutSec int) *HTTPFaceClient {
	return &HTTPFaceClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
	}
}

func (c *HTTPFaceClient) Extract(ctx context.Context, imgBytes []byte, filename, contentType, pose, requestID string) (*model.PythonExtractResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := createImageFormFile(writer, "file", filename, contentType)
	if err != nil {
		return nil, model.NewInternalError("Gagal mempersiapkan payload citra", err)
	}
	if _, err := part.Write(imgBytes); err != nil {
		return nil, model.NewInternalError("Gagal menulis payload citra", err)
	}

	if err := writer.WriteField("pose", pose); err != nil {
		return nil, model.NewInternalError("Gagal menulis form pose", err)
	}

	if err := writer.Close(); err != nil {
		return nil, model.NewInternalError("Gagal menutup multipart payload", err)
	}

	url := fmt.Sprintf("%s/api/v1/internal/face/extract", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, model.NewInternalError("Gagal membuat HTTP request ke face service", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Internal-API-Key", c.apiKey)
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) || strings.Contains(err.Error(), "Client.Timeout") {
			return nil, &model.AppError{
				StatusCode: http.StatusGatewayTimeout,
				Code:       "FACE_SERVICE_TIMEOUT",
				Message:    "Waktu tunggu pemrosesan biometrik wajah terlampaui",
			}
		}
		return nil, &model.AppError{
			StatusCode: http.StatusServiceUnavailable,
			Code:       "FACE_SERVICE_UNAVAILABLE",
			Message:    "Layanan internal biometrik wajah tidak dapat dihubungi",
			Internal:   err,
		}
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, model.NewInternalError("Gagal membaca respons dari layanan biometrik", err)
	}

	if resp.StatusCode == http.StatusOK {
		var result model.PythonExtractResponse
		if err := json.Unmarshal(respBytes, &result); err != nil {
			return nil, model.NewInternalError("Gagal deserialisasi respons embedding biometrik", err)
		}
		return &result, nil
	}

	// Parsing format error terstruktur dari face service
	var pyErr model.PythonErrorResponse
	if err := json.Unmarshal(respBytes, &pyErr); err == nil && pyErr.Error.Code != "" {
		code := pyErr.Error.Code
		msg := pyErr.Error.Message
		statusCode := resp.StatusCode
		if statusCode == http.StatusUnprocessableEntity {
			return nil, &model.AppError{
				StatusCode: http.StatusUnprocessableEntity,
				Code:       code,
				Message:    msg,
			}
		}
		if statusCode == http.StatusBadRequest {
			return nil, model.NewBadRequest(msg)
		}
		if statusCode >= 500 {
			return nil, &model.AppError{
				StatusCode: http.StatusServiceUnavailable,
				Code:       code,
				Message:    msg,
			}
		}
	}

	return nil, &model.AppError{
		StatusCode: http.StatusServiceUnavailable,
		Code:       "FACE_SERVICE_ERROR",
		Message:    fmt.Sprintf("Layanan biometrik mengembalikan status %d", resp.StatusCode),
	}
}

func (c *HTTPFaceClient) Verify(ctx context.Context, imgBytes []byte, filename, contentType string, enrolledSamples []model.EnrolledSamplePayload, threshold *float64, strategy, requestID string) (*model.PythonVerifyResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := createImageFormFile(writer, "file", filename, contentType)
	if err != nil {
		return nil, model.NewInternalError("Gagal mempersiapkan payload citra", err)
	}
	if _, err := part.Write(imgBytes); err != nil {
		return nil, model.NewInternalError("Gagal menulis payload citra", err)
	}

	samplesBytes, err := json.Marshal(enrolledSamples)
	if err != nil {
		return nil, model.NewInternalError("Gagal serialisasi sampel enrollment", err)
	}

	if err := writer.WriteField("enrolled_samples", string(samplesBytes)); err != nil {
		return nil, model.NewInternalError("Gagal menulis form enrolled_samples", err)
	}

	if strategy != "" {
		if err := writer.WriteField("strategy", strategy); err != nil {
			return nil, model.NewInternalError("Gagal menulis form strategy", err)
		}
	}

	if threshold != nil {
		if err := writer.WriteField("threshold", fmt.Sprintf("%.4f", *threshold)); err != nil {
			return nil, model.NewInternalError("Gagal menulis form threshold", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, model.NewInternalError("Gagal menutup multipart payload", err)
	}

	url := fmt.Sprintf("%s/api/v1/internal/face/verify", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, model.NewInternalError("Gagal membuat HTTP request ke face service", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Internal-API-Key", c.apiKey)
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) || strings.Contains(err.Error(), "Client.Timeout") {
			return nil, &model.AppError{
				StatusCode: http.StatusGatewayTimeout,
				Code:       "FACE_SERVICE_TIMEOUT",
				Message:    "Waktu tunggu pemrosesan verifikasi wajah terlampaui",
			}
		}
		return nil, &model.AppError{
			StatusCode: http.StatusServiceUnavailable,
			Code:       "FACE_SERVICE_UNAVAILABLE",
			Message:    "Layanan internal biometrik wajah tidak dapat dihubungi",
			Internal:   err,
		}
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, model.NewInternalError("Gagal membaca respons dari layanan biometrik", err)
	}

	if resp.StatusCode == http.StatusOK {
		var result model.PythonVerifyResponse
		if err := json.Unmarshal(respBytes, &result); err != nil {
			return nil, model.NewInternalError("Gagal deserialisasi respons verifikasi biometrik", err)
		}
		return &result, nil
	}

	// Parsing format error terstruktur dari face service
	var pyErr model.PythonErrorResponse
	if err := json.Unmarshal(respBytes, &pyErr); err == nil && pyErr.Error.Code != "" {
		code := pyErr.Error.Code
		msg := pyErr.Error.Message
		statusCode := resp.StatusCode
		if statusCode == http.StatusUnprocessableEntity {
			return nil, &model.AppError{
				StatusCode: http.StatusUnprocessableEntity,
				Code:       code,
				Message:    msg,
			}
		}
		if statusCode == http.StatusBadRequest {
			return nil, model.NewBadRequest(msg)
		}
		if statusCode >= 500 {
			return nil, &model.AppError{
				StatusCode: http.StatusServiceUnavailable,
				Code:       code,
				Message:    msg,
			}
		}
	}

	return nil, &model.AppError{
		StatusCode: http.StatusServiceUnavailable,
		Code:       "FACE_SERVICE_ERROR",
		Message:    fmt.Sprintf("Layanan biometrik mengembalikan status %d", resp.StatusCode),
	}
}

func createImageFormFile(writer *multipart.Writer, fieldname, filename, contentType string) (io.Writer, error) {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}
	header.Set("Content-Type", contentType)
	return writer.CreatePart(header)
}
