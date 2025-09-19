# map-area-factions

🎯 Prompt para o Copilot Agent — “Mapa de Áreas por Facção (web + mobile + desktop)”

Quero um sistema end-to-end que inclui: backend (Go), painel admin web, app Flutter (Android/iOS) e Flutter desktop (Windows/macOS/Linux), com mapas e polígonos (polylines/polygons) para marcar áreas dominadas por facções, oferecer busca por bairro/cidade, assinaturas pagas e atualização em tempo real. Produza tudo que for necessário para colocar em produção.

0) Diretrizes gerais
	•	Linguagem do backend: Go (priorize performance e concorrência).
	•	Framework: Fiber ou Gin (sugira e justifique no README).
	•	Banco: PostgreSQL + PostGIS para geoespacial.
	•	Cache/tempo real: Redis (Pub/Sub) + WebSockets; disponibilize também SSE como fallback.
	•	Geocodificação/tiles: abstrair provedor (Mapbox, MapTiler, Google, OpenStreetMap).
	•	Admin web: Next.js 14 (App Router) com Mapbox GL JS + Mapbox GL Draw (ou alternativa compatível) para desenhar/editar polígonos/linhas.
	•	App do cidadão: Flutter (um único código para mobile e desktop).
	•	Observabilidade: OpenAPI/Swagger, métricas Prometheus, logs estruturados, tracing OpenTelemetry.
	•	Infra: Docker + docker-compose para dev, IaC para prod (Terraform ou CDKTF), deploy em AWS (RDS Postgres com PostGIS, ElastiCache Redis, ECS Fargate/EKS, S3+CloudFront).
	•	Segurança: RBAC, JWT + Refresh, roteamento zero-trust no backend, auditoria de alterações geoespaciais.
	•	Sem dados sensíveis pessoais. Incluir disclaimers legais e mecanismo de report de correções.

1) Personas e fluxos
	1.	Admin (gestor): cria/edita bairros/zonas via mapa, define facção dominante (ex.: TDN, CV etc.), cores e rótulos, períodos de validade, nível de risco, notas, fonte/autor, status (rascunho/revisão/publicado).
	2.	Revisor/Moderador: revisa submissões da comunidade, aprova/recusa, registra histórico.
	3.	Cidadão (assinante ou gratuito limitado): pesquisa por bairro/cidade/CEP, visualiza polígonos coloridos, rótulos grandes (como na imagem), legenda por facção, aviso legal e atualizações em tempo real.
	4.	Colaborador (opcional): sugere novos polígonos/eventos (modo crowd), com fila de moderação.

2) Funcionalidades obrigatórias
	•	Desenho/edição geoespacial no Admin: criar, mover, dividir, mesclar, excluir polígonos e polilinhas; snap a vias; importar/exportar GeoJSON.
	•	Camadas e estilos: por facção, por nível de risco, por data de atualização; rótulos grandes e configuráveis.
	•	Busca: por bairro/cidade, geo-autocomplete, “perto de mim” (GPS).
	•	Assinaturas: planos 1 mês (R$14,90), 3 meses (R$29,90), vitalício (R$79,90). Integração com Mercado Pago e Pix; controle de acesso por feature flags e paywall.
	•	Tempo real: quando um polígono mudar de facção/risco/status, push para clientes conectados (WebSocket/SSE).
	•	Histórico/Auditoria: versionamento de cada polígono (quem alterou, quando, dif GeoJSON).
	•	Incidentes (opcional, mas planeje): pontos/quadrículas com validade curta (ex.: bloqueios, ocorrências), feed ao vivo, expiração automática.
	•	Offline básico no app: cache de mosaico/GeoJSON por cidade, com TTL e aviso de dados possivelmente desatualizados.
	•	Acessibilidade: contraste/cores daltônicas, legendas descritivas, tamanho de fonte.
	•	LGPD/Legal: tela de termos/privacidade, denúncia de informação incorreta, contato.

3) Modelo de dados (alto nível)
	•	factions (id, nome, sigla, cor_hex, prioridade_exibicao, ativo)
	•	regions (id, nome, tipo: bairro/zona/custom, municipio, estado, centroide, bounds, status: rascunho/revisão/publicado)
	•	geometries (id, region_id, geojson, tipo: poligono/polilinha, nivel_risco: baixo/médio/alto, faccao_id, inicio_vigencia, fim_vigencia, fonte, anotacoes)
	•	geometry_versions (id, geometry_id, geojson, diff, autor_id, created_at)
	•	incidents (id, tipo, descricao, geom ponto/área, inicio, fim, severidade, status)
	•	users (id, nome, email, hash, role: admin/revisor/colaborador/cidadao, status)
	•	subscriptions (id, user_id, plano, status, provider, external_id, dt_inicio, dt_fim, nota_fiscal)
	•	audit_logs (id, user_id, entidade, entidade_id, acao, diff, ip, user_agent, created_at)
	•	reports (id, user_id?, geometry_id?, incident_id?, texto, anexos, status)

Inclua migrations com PostGIS (índices GIST, SRID padrão 4326), constraints de integridade e views para consultas rápidas por município.

4) Contratos de API (REST + eventos)

Crie documentação OpenAPI 3 completa (sem código), cobrindo:
	•	Auth: login, refresh, rotas públicas mínimas, scopes por role.
	•	Factions: CRUD + ordenação/cores.
	•	Regions: listar por cidade/UF, detalhe, busca, estatísticas.
	•	Geometries: CRUD, importar GeoJSON, publicar, versionar, histórico.
	•	Incidents: CRUD, proximidade (raio), expiração.
	•	Subscriptions: webhooks do provedor (Pix/Mercado Pago), status, revalidação.
	•	Search: ?q=bairro | ?lat=…&lng=…&radius=… | sugestão por município.
	•	Realtime: canal WS/SSE (geometry.updated, incident.created, subscription.changed).
	•	Reports: abrir ticket/correção, anexos (S3), moderação.
Defina códigos de erro, padrão de paginação, rate limit headers e exemplos de request/response na spec.

5) Front-end Admin (Next.js)
	•	Páginas: Login, Dashboard, Mapa (CRUD de polígonos/linhas), Lista por cidade, Moderar Sugestões, Assinaturas, Audit Log, Configurações de estilos/legendas.
	•	Ferramentas do mapa: desenhar, editar, desfazer/refazer, cortar/mesclar, importar/exportar GeoJSON, “pintar” facção e risco, modo calor (opcional).
	•	UX: legendas fixas (cores por facção), toggle camadas, filtros por status e data, preview rótulo grande (como a arte “CENTRAL CEARA”).
	•	Integração WS para refletir alterações em tempo real em quem estiver no mesmo município.

6) App Flutter (mobile + desktop)
	•	Telas: Onboarding (termos/privacidade), Login/Assinatura, Mapa, Busca (bairro/cidade), Detalhe da área, Legenda, Configurações, Reportar erro/ocorrência.
	•	Funcionalidades:
	•	Mapa com renderização de polígonos por facção/risco, rótulos legíveis.
	•	Busca e “perto de mim”.
	•	Sincronização via WS/SSE, com reconexão e fila offline.
	•	Paywall + gestão de assinatura (planos, recibos).
	•	Notificações push (quando área próxima muda de nível/facção).
	•	Cache offline por cidade (limite de tamanho configurável).

7) Assinaturas e monetização
	•	Planos: 1 mês R$14,90; 3 meses R$29,90; Vitalício R$79,90.
	•	Integrar Mercado Pago com Pix e cartão; webhooks seguros (validação de assinatura).
	•	Gate no backend: middleware verifica escopo/validade antes de servir dados detalhados; modo gratuito mostra resumo simplificado.

8) Qualidade, segurança e compliance
	•	Testes: unitários + integração (APIs + geoespacial), testes de carga em consultas geográficas.
	•	Segurança: headers, CORS restrito, rate limit, anti-abuso em endpoints de busca, hashing moderno, rotação de secrets, backup automatizado do banco, WAF no edge.
	•	Compliance: LGPD (minimização de dados), auditoria completa, retenção e descarte de reports/incident após TTL.
	•	Conteúdo sensível: disclaimer claro (“dados colaborativos, podem conter imprecisões; não substituem orientação oficial; use por sua conta e risco”), botão “Reportar imprecisão”.

9) DevEx e entrega
	•	Monorepo com pastas: /backend, /admin, /app, /infra, /docs.
	•	Scripts make para dev/up/lint/test.
	•	CI/CD: build, testes, lint, scan SCA, gerar imagem, rodar e2e, publicar, migrações automáticas com lock.
	•	Seed de dados (cidades de exemplo, facções TDN/CV, 10 bairros com polígonos exemplo).
	•	Documentação completa (/docs):
	•	Arquitetura, decisões ADR, como rodar local, como publicar, guias de moderação, como importar GeoJSON, política de cores/legendagem, política de privacidade/termos.
	•	Comparativo dos provedores de mapas e como trocar a chave.

10) Critérios de aceite
	•	Eu consigo:
	1.	Desenhar um polígono no admin, marcar facção/risco e publicar.
	2.	Ver esse polígono no app após poucos segundos (WS/SSE).
	3.	Pesquisar “Messejana, Fortaleza-CE” e ver camadas por facção (cores e rótulos grandes).
	4.	Comprar assinatura (Pix em sandbox), ter acesso completo; expirar assinatura e perder o acesso premium.
	5.	Ver histórico de versões do polígono e restauração (“rollback”).
	6.	Exportar GeoJSON de uma cidade.
	7.	Receber push quando uma área próxima muda de risco.
	8.	Rodar tudo em docker-compose local e ter manifestos para deploy em AWS.

11) Extras (se couber no escopo)
	•	Tile server próprio para estilização e custo baixo (Vector tiles a partir de GeoJSON).
	•	Quad-tree/R-tree e tesselação para escalabilidade de render.
	•	Heatmap por incidentes recentes (com TTL).
	•	Suporte a múltiplas cidades/estados com sharding lógico.
	•	Modo seguro: ocultar rótulos em zoom < X para evitar “marcação” indevida.

Entregue: estrutura de repositório, READMEs detalhados, especificação OpenAPI, diagramas (C4/Seq), scripts de setup, seeds, pipelines CI, manifestos IaC, políticas de segurança, termos/privacidade e arquivos de configuração para chaves de mapas. Não gerar código agora — apenas todos os artefatos de especificação, planejamento e implantação necessários para que o time parta para implementação imediata.
