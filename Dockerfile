# ==========================================
# ETAPA 1: Ambiente de Construção (Builder)
# ==========================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Instala dependências do sistema operacional para baixar pacotes
RUN apk add --no-cache git

# Copia os gerenciadores de dependência e baixa tudo (aproveitando o cache do Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o restante do código
COPY . .

# Compila o binário otimizado
# CGO_ENABLED=0 garante um binário 100% estático sem dependências de bibliotecas C
# -ldflags="-w -s" remove informações de debug, reduzindo drasticamente o tamanho do arquivo
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o netwatch_app ./cmd/server/main.go

# ==========================================
# ETAPA 2: Ambiente de Produção (Final)
# ==========================================
FROM alpine:latest

WORKDIR /app

# Instala certificados de segurança (CA) e configura o fuso horário corporativo
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=America/Sao_Paulo

# Copia APENAS o binário limpo da Etapa 1
COPY --from=builder /app/netwatch_app .
# Copia a pasta do frontend estático para o painel funcionar
COPY --from=builder /app/frontend ./frontend

# Expõe a porta interna do contêiner
EXPOSE 8080

# Comando de inicialização
CMD ["./netwatch_app"]