package model

import "time"

type FaceTemplate struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	ModelName      string    `json:"model_name"`
	ModelVersion   string    `json:"model_version"`
	ConsentVersion string    `json:"consent_version"`
	ConsentedAt    time.Time `json:"consented_at"`
	Status         string    `json:"status"` // ACTIVE, REVOKED
	EnrolledAt     time.Time `json:"enrolled_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type FaceTemplateSample struct {
	ID                  string    `json:"id"`
	FaceTemplateID      string    `json:"face_template_id"`
	Pose                string    `json:"pose"` // FRONT, LEFT, RIGHT
	EmbeddingCiphertext []byte    `json:"-"`    // Jangan pernah ekspos ke JSON
	QualityStatus       string    `json:"quality_status"`
	CreatedAt           time.Time `json:"created_at"`
}

type FaceEnrollResponse struct {
	Message      string `json:"message"`
	TemplateID   string `json:"template_id"`
	ModelName    string `json:"model_name"`
	ModelVersion string `json:"model_version"`
	SampleCount  int    `json:"sample_count"`
	EnrolledAt   string `json:"enrolled_at"`
}

type FaceStatusResponse struct {
	Enrolled       bool     `json:"enrolled"`
	Status         string   `json:"status"`
	ModelName      string   `json:"model_name,omitempty"`
	ModelVersion   string   `json:"model_version,omitempty"`
	ConsentVersion string   `json:"consent_version,omitempty"`
	EnrolledAt     string   `json:"enrolled_at,omitempty"`
	SampleCount    int      `json:"sample_count"`
	Poses          []string `json:"poses"`
}

type PythonExtractResponse struct {
	ModelName      string    `json:"model_name"`
	ModelVersion   string    `json:"model_version"`
	Pose           string    `json:"pose"`
	Embedding      []float64 `json:"embedding"`
	Dim            int       `json:"dim"`
	FaceConfidence float64   `json:"face_confidence"`
}

type EnrolledSamplePayload struct {
	Pose      string    `json:"pose"`
	Embedding []float64 `json:"embedding"`
}

type PythonVerifyResponse struct {
	Verified        bool               `json:"verified"`
	Distance        float64            `json:"distance"`
	Threshold       float64            `json:"threshold"`
	Strategy        string             `json:"strategy"`
	BestMatchedPose string             `json:"best_matched_pose"`
	DistancesByPose map[string]float64 `json:"distances_by_pose"`
	IsCalibrated    bool               `json:"is_calibrated"`
	DecisionRule    string             `json:"decision_rule"`
	ModelName       string             `json:"model_name"`
}

type PythonErrorResponse struct {
	Error struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		RequestID string `json:"request_id"`
	} `json:"error"`
}
