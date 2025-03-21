# Notify

Uma biblioteca Go para enviar notificações via gRPC para o serviço de notificações da AdSeleto.

## Instalação

```bash
go get github.com/AdSeleto/notify
```

## Uso Básico

Para enviar uma notificação:

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/AdSeleto/notify"
)

func main() {
	// Cria um cliente com configurações
	notifier, err := notify.NewClient(
		notify.WithServerAddress("notifications-service:50051"),
		notify.WithOrigin("nome-do-seu-servico"),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente: %v", err)
	}
	defer notifier.Close() // Sempre feche o cliente ao terminar

	// Prepara os parâmetros da notificação
	params := &notify.Data{
		ProjectID: "seu-projeto-id",
		Scope:     "SEU_ESCOPO", // SISTEMA, WARMUP, CAMPANHA, etc
		Type:      "TIPO_NOTIFICACAO", // SPAM, ENTREGABILIDADE, BLACKLIST, HIGH_BOUNCE, etc
		Severity:  notify.INFO, // INFO, WARNING, ERROR, CRITICAL
		Title:     "Título da notificação",
		Content:   "Conteúdo detalhado da notificação", // Conteúdo do e-mail
		Metadata: map[string]string{
			"chave": "valor",
		},
	}

	// Envia a notificação
	ctx := context.Background()
	if err := notifier.Notify(ctx, params); err != nil {
		log.Fatalf("Erro ao enviar notificação: %v", err)
	}

	log.Println("Notificação enviada com sucesso!")
}
```

## Configuração

A biblioteca usa valores padrão para a maioria das configurações, mas você pode personalizá-los:

| Opção           | Valor Padrão     | Descrição                              | Obrigatório |
|-----------------|------------------|----------------------------------------|-------------|
| Origin          | -                | Origem do serviço que envia a notificação | Sim        |
| ServerAddress   | -                | Endereço do servidor gRPC              | Sim         |
| Timeout         | 10 segundos      | Tempo máximo para cada requisição      | Não         |
| MaxRetries      | 3                | Número máximo de tentativas em caso de falha | Não     |
| RetryInterval   | 2 segundos       | Tempo entre tentativas de reconexão    | Não         |
| EnableTLS       | false            | Habilitar/desabilitar TLS              | Não         |

### Personalizando a configuração

```go
notifier, err := notify.NewClient(
    notify.WithServerAddress("notifications-service:50051"),
    notify.WithOrigin("meu-servico"),
    notify.WithTimeout(5 * time.Second),
    notify.WithMaxRetries(2),
)
```

## Níveis de Severidade

A biblioteca oferece constantes para os possíveis níveis de severidade:

```go
notify.INFO     // Para notificações informativas
notify.WARNING  // Para avisos
notify.ERROR    // Para erros
notify.CRITICAL // Para problemas críticos
```

## API de Referência

### Tipos

#### `notify.Client`

Interface principal do cliente:

```go
type Client interface {
    Notify(ctx context.Context, params *Data) error
    Close() error
}
```

#### `notify.Data`

Estrutura para os dados da notificação:

```go
type Data struct {
    ProjectID string            // ID do projeto
    Scope     string            // Escopo da notificação
    Type      string            // Tipo da notificação
    Title     string            // Título
    Content   string            // Conteúdo/mensagem
    Severity  string            // Nível de severidade (usar constantes)
    Metadata  map[string]string // Dados adicionais em formato chave-valor
}
```

### Funções

#### `notify.NewClient(opts ...Option) (Client, error)`

Cria uma nova instância do cliente de notificações.

### Opções de Configuração

- `notify.WithOrigin(origin string)`: Define a origem do serviço (obrigatório)
- `notify.WithServerAddress(address string)`: Define o endereço do servidor gRPC
- `notify.WithTimeout(timeout time.Duration)`: Define o timeout para requisições
- `notify.WithMaxRetries(retries int)`: Define o número máximo de tentativas
- `notify.WithRetryInterval(interval time.Duration)`: Define o intervalo entre tentativas
- `notify.WithTLS(certPath string)`: Habilita TLS com o certificado fornecido

## Contexto

A biblioteca respeita o padrão de `context.Context` de Go:

- Se você fornecer um contexto com deadline/timeout, ele será respeitado
- Se você não fornecer timeout no contexto, será usado o timeout padrão da biblioteca
- Se o contexto for cancelado, a operação será interrompida

## Tratamento de Erros

A biblioteca inclui retentativas automáticas em caso de falhas temporárias na comunicação. Depois de exceder o número máximo de tentativas, o erro da última tentativa é retornado.

Se o Sentry estiver configurado no ambiente, os erros também serão capturados automaticamente.

## Dicas de Uso

### Singleton com Inicialização Preguiçosa

Para aplicações de longa duração, considere usar um padrão singleton:

```go
package myapp

import (
	"context"
	"sync"

	"github.com/AdSeleto/notify"
)

var (
	notifier notify.Client
	once   sync.Once
	initErr error
)

func GetClient() (notify.Client, error) {
	once.Do(func() {
		notifier, initErr = notify.NewClient(
			notify.WithServerAddress("notifications-service:50051"),
			notify.WithOrigin("seu-servico"),
		)
	})
	return notifier, initErr
}

// SendNotification é um helper para enviar notificações facilmente
func SendNotification(title, content, severity string, metadata map[string]string) error {
	c, err := GetClient()
	if err != nil {
		return err
	}

	params := &notify.Data{
		ProjectID: "seu-projeto-id",
		Scope:     "SEU_ESCOPO",
		Type:      "NOTIFICACAO",
		Title:     title,
		Content:   content,
		Severity:  severity,
		Metadata:  metadata,
	}

	return c.Notify(context.Background(), params)
}
```

### Usando com HTTP Handlers

Ao usar com handlers HTTP, propague o contexto da requisição:

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Use o contexto da requisição
	err := notifier.Notify(r.Context(), params)
	// ...
}
```

### Integração com gRPC

Ao usar com serviços gRPC, utilize o contexto recebido:

```go
func (s *service) HandleSomething(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	// Usa o contexto do cliente gRPC
	err := notifier.Notify(ctx, params)
	// ...
}
```

## Licença

Copyright © 2025 AdSeleto.
