package notify

import (
	"fmt"

	"github.com/AdSeleto/notify/internal/infrastructure/grpc/notifications"
)

// Constantes para Severity
const (
	SeverityInfo     = "INFO"
	SeverityWarning  = "WARNING"
	SeverityError    = "ERROR"
	SeverityCritical = "CRITICAL"
)

// Data representa os parâmetros para criar uma notificação
type Data struct {
	ProjectID string            `json:"project_id"`
	Scope     string            `json:"scope"`
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Severity  string            `json:"severity"`
	Origin    string            `json:"origin"`
	Metadata  map[string]string `json:"metadata"`
}

// Valida se a severity está entre os valores permitidos
func (np *Data) validateSeverity() error {
	switch np.Severity {
	case SeverityInfo, SeverityWarning, SeverityError, SeverityCritical:
		return nil
	default:
		return fmt.Errorf("severity inválida: %s. Use uma das constantes: SeverityInfo, SeverityWarning, SeverityError, SeverityCritical", np.Severity)
	}
}

// Converte os parâmetros de notificação para uma request gRPC
func (np *Data) toGRPCRequest() (*notifications.NotifyRequest, error) {
	if err := np.validateSeverity(); err != nil {
		return nil, err
	}

	return &notifications.NotifyRequest{
		ProjectId: np.ProjectID,
		Scope:     np.Scope,
		Type:      np.Type,
		Title:     np.Title,
		Content:   np.Content,
		Severity:  np.Severity,
		Origin:    np.Origin,
		Metadata:  np.Metadata,
	}, nil
}
