package notify

import (
	"time"
)

// ClientOptions contém todas as opções configuráveis para o cliente de notificações
type ClientOptions struct {
	// Endereço do servidor gRPC
	ServerAddress string

	// Timeout para conexões gRPC
	Timeout time.Duration

	// Tentativas máximas de reconexão
	MaxRetries int

	// Intervalo entre tentativas de reconexão
	RetryInterval time.Duration

	// Habilitar TLS
	EnableTLS bool

	// Certificado TLS (se EnableTLS for true)
	TLSCertPath string

	// Origin identifica o serviço que está enviando a notificação
	Origin string
}

// DefaultOptions retorna as opções padrão para o cliente
func DefaultOptions() *ClientOptions {
	return &ClientOptions{
		ServerAddress: "", // ServerAddress deve ser configurado explicitamente
		Timeout:       time.Second * 10,
		MaxRetries:    3,
		RetryInterval: time.Second * 2,
		EnableTLS:     false,
		Origin:        "",
	}
}

// Option é um tipo para funções de configuração
type Option func(*ClientOptions)

// WithServerAddress define o endereço do servidor
func WithServerAddress(address string) Option {
	return func(o *ClientOptions) {
		o.ServerAddress = address
	}
}

// WithTimeout define o timeout para requisições
func WithTimeout(timeout time.Duration) Option {
	return func(o *ClientOptions) {
		o.Timeout = timeout
	}
}

// WithMaxRetries define o número máximo de tentativas de reconexão
func WithMaxRetries(retries int) Option {
	return func(o *ClientOptions) {
		o.MaxRetries = retries
	}
}

// WithRetryInterval define o intervalo entre tentativas
func WithRetryInterval(interval time.Duration) Option {
	return func(o *ClientOptions) {
		o.RetryInterval = interval
	}
}

// WithTLS habilita TLS com o certificado fornecido
func WithTLS(certPath string) Option {
	return func(o *ClientOptions) {
		o.EnableTLS = true
		o.TLSCertPath = certPath
	}
}

// WithOrigin define a origem/serviço que está enviando a notificação
func WithOrigin(origin string) Option {
	return func(o *ClientOptions) {
		o.Origin = origin
	}
}
