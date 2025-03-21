package notifications

import (
	"fmt"

	"github.com/AdSeleto/notifications-client/internal/infrastructure/grpc/notifications"
)

// Constantes para Severity
const (
	SEVERITY_INFO     = "INFO"
	SEVERITY_WARNING  = "WARNING"
	SEVERITY_ERROR    = "ERROR"
	SEVERITY_CRITICAL = "CRITICAL"
)

// NotificationParams representa os parâmetros para criar uma notificação
type NotificationParams struct {
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
func (np *NotificationParams) validateSeverity() error {
	switch np.Severity {
	case SEVERITY_INFO, SEVERITY_WARNING, SEVERITY_ERROR, SEVERITY_CRITICAL:
		return nil
	default:
		return fmt.Errorf("severity inválida: %s. Use uma das constantes: SeverityInfo, SeverityWarning, SeverityError, SeverityCritical", np.Severity)
	}
}

// Converte os parâmetros de notificação para uma request gRPC
func (np *NotificationParams) toGRPCRequest() (*notifications.NotifyRequest, error) {
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
