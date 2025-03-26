package notify

import (
	"fmt"

	"github.com/AdSeleto/notify/internal/infrastructure/grpc/notifications"
)

// Constantes para Severity
const (
	INFO     = "INFO"
	WARNING  = "WARNING"
	ERROR    = "ERROR"
	CRITICAL = "CRITICAL"
)

// Constantes para Scope
const (
	CAMPAIGN = "CAMPAIGN"
	PROJECT  = "PROJECT"
	SYSTEM   = "SYSTEM"
	WARMUP   = "WARMUP"
)

// Constantes para Type
const (
	BLACKLIST           = "BLACKLIST"
	HIGH_BOUNCE         = "HIGH_BOUNCE"
	DELIVERABILITY_DROP = "DELIVERABILITY_DROP"
	COMPLETED           = "COMPLETED"
	FAILED              = "FAILED"
	ISSUES              = "ISSUES"
	IMPORT_COMPLETED    = "IMPORT_COMPLETED"
	STATE_CHANGE        = "STATE_CHANGE"
	DAILY_SUMMARY       = "DAILY_SUMMARY"
	PAUSED              = "PAUSED"
	BOUNCE              = "BOUNCE"
	SPAM_COMPLAINTS     = "SPAM_COMPLAINTS"
)

// Data representa os parâmetros para criar uma notificação
type Data struct {
	ProjectID string            `json:"project_id"`
	Scope     string            `json:"scope"`
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Severity  string            `json:"severity"`
	Metadata  map[string]string `json:"metadata"`
}

// Valida se a severity está entre os valores permitidos
func (np *Data) validateSeverity() error {
	switch np.Severity {
	case INFO, WARNING, ERROR, CRITICAL:
		return nil
	default:
		return fmt.Errorf("severity inválida: %s. Use uma das constantes: INFO, WARNING, ERROR, CRITICAL", np.Severity)
	}
}

// Valida se o scope está entre os valores permitidos
func (np *Data) validateScope() error {
	switch np.Scope {
	case CAMPAIGN, PROJECT, SYSTEM, WARMUP:
		return nil
	default:
		return fmt.Errorf("invalid scope: %s. Use one of the constants: CAMPAIGN, PROJECT, SYSTEM, WARMUP", np.Scope)
	}
}

// Valida se o type está entre os valores permitidos
func (np *Data) validateType() error {
	switch np.Type {
	case BLACKLIST, HIGH_BOUNCE, DELIVERABILITY_DROP, COMPLETED, FAILED, ISSUES, IMPORT_COMPLETED, STATE_CHANGE, DAILY_SUMMARY, PAUSED, BOUNCE, SPAM_COMPLAINTS:
		return nil
	default:
		return fmt.Errorf("invalid type: %s. Use one of the constants: BLACKLIST, HIGH_BOUNCE, DELIVERABILITY_DROP, COMPLETED, FAILED, ISSUES, IMPORT_COMPLETED, STATE_CHANGE, DAILY_SUMMARY, PAUSED, BOUNCE, SPAM_COMPLAINTS", np.Type)
	}
}

// Converte os parâmetros de notificação para uma request gRPC
func (np *Data) toGRPCRequest(origin string) (*notifications.NotifyRequest, error) {
	if err := np.validateSeverity(); err != nil {
		return nil, err
	}
	if err := np.validateScope(); err != nil {
		return nil, err
	}
	if err := np.validateType(); err != nil {
		return nil, err
	}

	return &notifications.NotifyRequest{
		ProjectId: np.ProjectID,
		Scope:     np.Scope,
		Type:      np.Type,
		Title:     np.Title,
		Content:   np.Content,
		Severity:  np.Severity,
		Origin:    origin,
		Metadata:  np.Metadata,
	}, nil
}
