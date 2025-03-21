package main

import (
	"context"
	"log"
	"time"

	"github.com/AdSeleto/notifications-client/pkg/notifications"
)

func main() {
	// Cria um novo cliente com opções personalizadas
	notifier, err := notifications.NewClient(
		notifications.WithServerAddress("notifications-service:50051"),
		notifications.WithTimeout(5*time.Second),
		notifications.WithMaxRetries(2),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente de notificações: %v", err)
	}
	defer notifier.Close()

	// Prepara os parâmetros da notificação
	params := &notifications.NotificationParams{
		ProjectID: "dc5e5aa5-ccdd-4cc7-be00-ccfa5ec37058",
		Scope:     "WARMUP",
		Type:      "HIGH_BOUNCE_RATE",
		Title:     "Taxa de bounces elevada",
		Content:   "A taxa de bounce está acima do limite aceitável.",
		Severity:  notifications.SEVERITY_INFO,
		Origin:    "go-warmups",
		Metadata: map[string]string{
			"bounce_rate": "0.5",
		},
	}
	// Envia a notificação
	ctx := context.Background()
	if err := notifier.Notify(ctx, params); err != nil {
		log.Fatalf("Erro ao enviar notificação: %v", err)
	}

	log.Println("Notificação enviada com sucesso!")
}
