# ADR-001: Escolha do Framework Backend (Go Fiber vs Gin)

**Status**: Aceito  
**Data**: 2024-01-15  
**Autores**: Equipe Map Factions

## Contexto

Precisamos escolher um framework Go para o backend da aplicação que atenda aos seguintes requisitos:
- Performance elevada para operações geoespaciais
- Suporte nativo a WebSockets para tempo real
- Facilidade de desenvolvimento e manutenção
- Ecosystem maduro
- Documentação abrangente

### Opções Consideradas

#### 1. Gin (github.com/gin-gonic/gin)
**Prós:**
- Framework mais popular no ecosystem Go
- Documentação extensa e comunidade ativa
- Performance excelente
- Middleware ecosystem rico
- Sintaxe clean e familiar

**Contras:**
- Baseado em net/http (limitações de performance)
- WebSocket requer bibliotecas externas
- Não otimizado para alta concorrência

#### 2. Fiber (github.com/gofiber/fiber)
**Prós:**
- Baseado em fasthttp (performance superior)
- WebSocket support nativo
- API similar ao Express.js (familiar)
- Menor consumo de memória
- Otimizado para alta concorrência
- Middleware oficial abrangente

**Contras:**
- Comunidade menor que Gin
- Menos recursos de terceiros
- Curva de aprendizado para fasthttp

## Decisão

**Escolhemos Fiber** pelos seguintes motivos:

### 1. Performance Superior
- Baseado no fasthttp, oferece 3-4x melhor performance que frameworks baseados em net/http
- Essencial para operações geoespaciais que exigem baixa latência
- Menor overhead de memória

### 2. WebSocket Nativo
- Suporte built-in para WebSockets
- Essencial para atualizações em tempo real do mapa
- Não requer bibliotecas externas

### 3. Concorrência Otimizada
- Melhor handling de conexões simultâneas
- Pool de workers eficiente
- Ideal para aplicações real-time

### 4. Ecosystem Adequado
- Middleware oficial para todas nossas necessidades:
  - JWT authentication
  - CORS
  - Rate limiting
  - Prometheus metrics
  - WebSocket

## Consequências

### Positivas
- Melhor performance em operações geoespaciais
- WebSocket integrado facilita tempo real
- Menor latência para consultas de mapa
- Melhor utilização de recursos do servidor

### Negativas
- Ecosystem menor que Gin
- Menos exemplos na comunidade
- Time precisará aprender particularidades do fasthttp

### Riscos Mitigados
- Documentação oficial é completa
- API familiar (Express-like)
- Suporte oficial ativo
- Pode migrar para Gin se necessário (API similar)

## Implementação

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/gofiber/cors/v2"
)

func main() {
    app := fiber.New(fiber.Config{
        ServerHeader: "Map Factions API",
        AppName:      "Map Area Factions v1.0.0",
    })

    // Middleware
    app.Use(cors.New())
    
    // WebSocket endpoint
    app.Use("/ws", func(c *fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })
    
    app.Get("/ws", websocket.New(handleWebSocket))
    
    app.Listen(":8080")
}
```

## Referências
- [Fiber Documentation](https://docs.gofiber.io/)
- [Fasthttp Performance](https://github.com/valyala/fasthttp)
- [Go Web Framework Benchmarks](https://github.com/smallnest/go-web-framework-benchmark)