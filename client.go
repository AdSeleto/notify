package notify

import (
	"context"
	"fmt"
	"time"

	"github.com/AdSeleto/notify/pb/notifications"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// NotifyClient é a estrutura concreta para o cliente de notificações
type NotifyClient struct {
	conn    *grpc.ClientConn
	client  notifications.NotificationsServiceClient
	options *ClientOptions
}

// NewClient cria uma nova instância do cliente de notificações
// Retorna um ponteiro para NotifyClient
func NewClient(opts ...Option) (*NotifyClient, error) {
	// Aplica as opções padrão
	options := DefaultOptions()

	// Aplica as opções personalizadas
	for _, opt := range opts {
		opt(options)
	}

	// Verifica se origin foi configurado
	if options.Origin == "" {
		return nil, fmt.Errorf("a origem (Origin) do serviço deve ser configurada usando WithOrigin()")
	}

	// Verifica se ServerAddress foi configurado
	if options.ServerAddress == "" {
		return nil, fmt.Errorf("o endereço do servidor (ServerAddress) deve ser configurado explicitamente usando WithServerAddress()")
	}

	// Estabelece a conexão gRPC
	conn, err := createConnection(options)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar conexão gRPC: %w", err)
	}

	// Cria o cliente gRPC
	client := notifications.NewNotificationsServiceClient(conn)

	return &NotifyClient{
		conn:    conn,
		client:  client,
		options: options,
	}, nil
}

// Notify envia uma notificação através do serviço gRPC
func (c *NotifyClient) Notify(ctx context.Context, params *Data) error {
	if params == nil {
		return fmt.Errorf("parâmetros de notificação não podem ser nulos")
	}

	// Converte os parâmetros para o formato gRPC, incluindo validação e adicionando origin
	req, err := params.toGRPCRequest(c.options.Origin)
	if err != nil {
		return fmt.Errorf("parâmetros inválidos: %w", err)
	}

	// Adiciona timeout ao contexto se não houver um
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.Timeout)
		defer cancel()
	}

	// Tenta enviar a notificação com retentativas
	var lastErr error
	for attempt := 0; attempt <= c.options.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(c.options.RetryInterval)
		}

		_, err := c.client.Notify(ctx, req)
		if err == nil {
			return nil
		}

		lastErr = err
		// Captura erro no Sentry, se configurado
		sentry.CaptureException(fmt.Errorf("tentativa %d falhou ao enviar notificação: %w", attempt+1, err))
	}

	return fmt.Errorf("falha ao enviar notificação após %d tentativas: %w", c.options.MaxRetries+1, lastErr)
}

// Close fecha a conexão gRPC
func (c *NotifyClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// createConnection estabelece uma conexão gRPC com as opções configuradas
func createConnection(options *ClientOptions) (*grpc.ClientConn, error) {
	// Define as opções de dial
	dialOpts := []grpc.DialOption{}

	// Configura TLS se habilitado
	if options.EnableTLS {
		creds, err := credentials.NewClientTLSFromFile(options.TLSCertPath, "")
		if err != nil {
			return nil, fmt.Errorf("falha ao carregar certificados TLS: %w", err)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Estabelece a conexão com timeout
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, options.ServerAddress, dialOpts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
