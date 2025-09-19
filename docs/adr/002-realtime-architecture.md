# ADR-002: Arquitetura de Real-time (WebSocket + Server-Sent Events)

**Status**: Aceito  
**Data**: 2024-01-15  
**Autores**: Equipe Map Factions

## Contexto

O sistema requer atualizações em tempo real para:
- Alterações em polígonos/geometrias
- Novos incidentes/ocorrências
- Mudanças de status de assinatura
- Notificações administrativas

### Requisitos
- Baixa latência (< 100ms)
- Fallback para conexões instáveis
- Suporte a múltiplos clientes
- Escalabilidade horizontal
- Compatibilidade com diferentes navegadores/dispositivos

## Opções Consideradas

### 1. WebSocket Apenas
**Prós:**
- Bidirectional communication
- Baixa latência
- Overhead mínimo após handshake

**Contras:**
- Problemas com proxies/firewalls
- Reconexão complexa
- Não funciona em todos os ambientes

### 2. Server-Sent Events (SSE) Apenas
**Prós:**
- Unidirectional (suficiente para notificações)
- Reconexão automática
- Compatibilidade universal
- Simples de implementar

**Contras:**
- Apenas server-to-client
- Limitações em alguns proxies
- Overhead maior que WebSocket

### 3. Polling
**Prós:**
- Compatibilidade universal
- Simples de implementar
- Funciona em qualquer ambiente

**Contras:**
- Alto overhead
- Latência alta
- Desperdício de recursos

## Decisão

**Implementamos WebSocket como primário com SSE como fallback:**

### Estratégia Híbrida
1. **WebSocket First**: Tentativa inicial de conexão WebSocket
2. **SSE Fallback**: Se WebSocket falhar, usar SSE
3. **Graceful Degradation**: Detecção automática de falhas

## Arquitetura Técnica

### Backend (Go Fiber)
```go
// WebSocket handler
app.Get("/ws", websocket.New(func(c *websocket.Conn) {
    for {
        var msg Message
        if err := c.ReadJSON(&msg); err != nil {
            break
        }
        
        // Broadcast to Redis
        redisClient.Publish("updates", msg)
    }
}))

// SSE handler
app.Get("/events", func(c *fiber.Ctx) error {
    c.Set("Content-Type", "text/event-stream")
    c.Set("Cache-Control", "no-cache")
    c.Set("Connection", "keep-alive")
    
    // Subscribe to Redis
    pubsub := redisClient.Subscribe("updates")
    defer pubsub.Close()
    
    for msg := range pubsub.Channel() {
        fmt.Fprintf(c.Response().BodyWriter(), 
            "data: %s\n\n", msg.Payload)
        c.Response().SetBodyStreamWriter(
            fasthttp.StreamWriter(func(w *bufio.Writer) {
                w.Flush()
            }))
    }
    
    return nil
})
```

### Frontend (Next.js)
```typescript
class RealtimeService {
    private ws: WebSocket | null = null;
    private eventSource: EventSource | null = null;
    private fallbackToSSE = false;

    connect() {
        if (this.fallbackToSSE) {
            this.connectSSE();
        } else {
            this.connectWebSocket();
        }
    }

    private connectWebSocket() {
        this.ws = new WebSocket(WS_URL);
        
        this.ws.onopen = () => {
            console.log('WebSocket connected');
        };
        
        this.ws.onerror = () => {
            console.log('WebSocket failed, falling back to SSE');
            this.fallbackToSSE = true;
            this.connectSSE();
        };
        
        this.ws.onmessage = (event) => {
            this.handleMessage(JSON.parse(event.data));
        };
    }

    private connectSSE() {
        this.eventSource = new EventSource('/api/events');
        
        this.eventSource.onmessage = (event) => {
            this.handleMessage(JSON.parse(event.data));
        };
        
        this.eventSource.onerror = () => {
            // Reconnect automatically
            setTimeout(() => this.connectSSE(), 5000);
        };
    }
}
```

### Flutter App
```dart
class RealtimeService {
  IO.Socket? _socket;
  StreamSubscription? _sseSubscription;
  bool _useSSE = false;

  Future<void> connect() async {
    if (_useSSE) {
      await _connectSSE();
    } else {
      await _connectWebSocket();
    }
  }

  Future<void> _connectWebSocket() async {
    try {
      _socket = IO.io(Config.wsUrl, {
        'transports': ['websocket'],
        'autoConnect': false,
      });
      
      _socket!.on('connect', (_) {
        print('WebSocket connected');
      });
      
      _socket!.on('connect_error', (_) {
        print('WebSocket failed, falling back to SSE');
        _useSSE = true;
        _connectSSE();
      });
      
      _socket!.connect();
    } catch (e) {
      _useSSE = true;
      await _connectSSE();
    }
  }

  Future<void> _connectSSE() async {
    // Use dio for SSE
    final response = await dio.get(
      '/events',
      options: Options(responseType: ResponseType.stream),
    );
    
    _sseSubscription = response.data.stream
        .transform(utf8.decoder)
        .transform(const LineSplitter())
        .listen((line) {
      if (line.startsWith('data: ')) {
        final data = jsonDecode(line.substring(6));
        _handleMessage(data);
      }
    });
  }
}
```

## Distribuição via Redis

### Pub/Sub Pattern
```go
// Publisher (quando geometria é atualizada)
func (s *GeometryService) UpdateGeometry(geometry *Geometry) error {
    // Update database
    if err := s.db.Save(geometry).Error; err != nil {
        return err
    }
    
    // Publish update
    message := RealtimeMessage{
        Type: "geometry.updated",
        Data: geometry,
        Timestamp: time.Now(),
    }
    
    return s.redis.Publish("updates", message).Err()
}

// Subscriber (WebSocket/SSE handlers)
func (s *RealtimeService) Subscribe() {
    pubsub := s.redis.Subscribe("updates")
    defer pubsub.Close()
    
    for msg := range pubsub.Channel() {
        // Broadcast to all connected clients
        s.broadcastToClients(msg.Payload)
    }
}
```

## Consequências

### Positivas
- **Reliability**: Fallback garante funcionamento universal
- **Performance**: WebSocket oferece baixa latência quando disponível
- **Compatibility**: SSE funciona em todos os navegadores modernos
- **Scalability**: Redis Pub/Sub permite scaling horizontal

### Negativas
- **Complexity**: Duas implementações para manter
- **Testing**: Necessidade de testar ambos os cenários
- **Resources**: Duas conexões por cliente (potencial)

### Implementação

#### Phases
1. **Phase 1**: WebSocket básico
2. **Phase 2**: SSE fallback
3. **Phase 3**: Reconexão inteligente
4. **Phase 4**: Otimizações de performance

#### Monitoring
- Connection count por tipo
- Fallback rate
- Message delivery latency
- Reconnection frequency

## Referencias
- [WebSocket RFC 6455](https://tools.ietf.org/html/rfc6455)
- [Server-Sent Events W3C](https://www.w3.org/TR/eventsource/)
- [Redis Pub/Sub](https://redis.io/topics/pubsub)