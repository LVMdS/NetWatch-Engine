# 📡 NetWatch Pro - Enterprise Network Monitor

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![SQLite](https://img.shields.io/badge/SQLite-Native%20(CGO--Free)-003B57?style=for-the-badge&logo=sqlite)
![TailwindCSS](https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css)
![Alpine.js](https://img.shields.io/badge/Alpine.js-8BC0D0?style=for-the-badge&logo=alpine.js)

O **NetWatch Pro** é uma plataforma de monitorização de infraestruturas de rede e observabilidade corporativa. Concebido para ser **100% autossuficiente e portátil**, roda a partir de um único executável nativo, incorporando um banco de dados SQLite local (sem dependência de CGO) e uma interface Web moderna em tempo real.

---

## ✨ Funcionalidades Principais

- **Monitorização Mista (ICMP & TCP):** Validação de disponibilidade de hosts via Ping tradicional ou testes diretos a portas (ex: Web, Banco de Dados, SSH).
- **Descoberta Automática de Rede (Scanner ARP):** Mapeamento de sub-redes em tempo real para encontrar hosts ativos e resolver Endereços MAC fisicamente.
- **Alertas Ricos via Discord:** Integração nativa com Webhooks do Discord, enviando Embeds (Cartões Visuais) formatados imediatamente na queda ou restabelecimento de serviços.
- **Lógica Anti-Flapping:** Filtro inteligente que exige falhas consecutivas antes de declarar um host como `OFFLINE`, evitando alertas falsos causados por micro-quedas ou picos de rede.
- **Gestão de SLA e Uptime:** Cálculo dinâmico da disponibilidade de cada nó da rede com exportação de inventário completa para arquivo CSV com um clique.
- **Dashboard NOC (Network Operations Center):** Interface reativa com alternância de visualização entre *Grid* (Cartões) e *Table* (Alta Densidade), sincronizada instantaneamente via WebSockets.

---

## 🚀 Como Compilar e Executar

A arquitetura do NetWatch Pro utiliza o pacote `glebarez/sqlite`, o que significa que **não requer compiladores C (GCC/MinGW)** instalados na sua máquina. A compilação é limpa e feita puramente em Go.

Certifique-se de ter o [Go instalado](https://go.dev/dl/) na sua máquina antes de prosseguir.

### 🪟 Windows
Para gerar um executável "invisível" (que roda de forma oculta em segundo plano, sem a janela do CMD) e otimizado:
(Opcional)
```powershell
# 1. Sincronizar dependências
go mod tidy

# 2. Compilar o binário otimizado e silencioso
$env:CGO_ENABLED="0"
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -ldflags="-H windowsgui -w -s" -o netwatch_engine.exe ./cmd/server/main.go
```

Para iniciar, dê um duplo clique em netwatch_engine.exe. A aplicação registrará a si mesma na inicialização e abrirá no seu navegador padrão.

#🐧 Linux
Para compilar um binário otimizado para servidores ou desktops Linux:

```Bash
## 1. Sincronizar dependências
go mod tidy

## 2. Compilar
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o netwatch_engine ./cmd/server/main.go

# 3. Dar permissão de execução e rodar
chmod +x netwatch_engine
./netwatch_engine
```

# 🍎 macOS
Para compilar no macOS (escolha a arquitetura correta do seu processador):

```Bash
## 1. Sincronizar dependências
go mod tidy

## 2A. Para Macs com Apple Silicon (M1/M2/M3):
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o netwatch_engine_mac ./cmd/server/main.go

## 2B. Para Macs com Processadores Intel:
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o netwatch_engine_mac ./cmd/server/main.go

## 3. Executar
./netwatch_engine_mac
```

⚙️ Configuração Inicial
Ao executar o NetWatch pela primeira vez, o sistema gerará automaticamente um arquivo netwatch.db (banco de dados SQLite) na mesma pasta do executável.

Acesse http://localhost:8080.

O sistema detectará que é o primeiro acesso e o redirecionará para a página de Registro.

Crie a sua conta de Administrador (você pode inserir a URL do seu Webhook do Discord neste momento).

O login será efetuado automaticamente e o painel de controle estará pronto para uso.

🌍 Acesso Remoto Seguro (Zero Trust VPN via ZeroTier)
Tratando-se de dados sensíveis de infraestrutura, não é recomendado expor a dashboard do NetWatch por meio de URLs públicas ou encaminhamento de portas (Port Forwarding) no seu roteador.

A solução definitiva para acessar o sistema de qualquer lugar com segurança de nível militar é utilizar o ZeroTier, que cria uma rede (SD-WAN) com criptografia ponta a ponta. O NetWatch continuará escutando na porta 8080, mas apenas os seus dispositivos aprovados manualmente conseguirão acessá-lo.

Passo a Passo da Configuração:
Criar a Rede Virtual:

Crie uma conta gratuita no ZeroTier.

Clique em Create a Network e copie o Network ID gerado (código alfanumérico).

Conectar o Servidor (Onde o NetWatch roda):

Instale o cliente ZeroTier na máquina do servidor.

Clique em Join Network e cole o Network ID.

Conectar o Dispositivo Remoto (Seu Celular ou Notebook):

Instale o aplicativo ZeroTier no seu dispositivo.

Adicione a rede colando o mesmo Network ID.

Autorizar o Acesso (A Chave de Segurança):

Retorne ao painel web do ZeroTier e role até a seção Members.

Localize os dois dispositivos que acabaram de entrar e marque a caixa Auth? para ambos.

O painel atribuirá um IP interno exclusivo (ex: 10.147.20.15) ao seu Servidor.

Acessando a Dashboard Remotamente:

No navegador do seu dispositivo remoto (certifique-se de que a conexão ZeroTier está ativada), digite a URL usando o IP gerado: http://10.147.20.15:8080.

NetWatch/
 ├── cmd/server/main.go       # Ponto de entrada, trava de diretório e rotas HTTP
 ├── database/                # Conexão GORM e AutoMigrate configurado para SQLite
 ├── frontend/                # Interface SPA Completa (HTML, Alpine.js, Tailwind)
 ├── internal/
 │    ├── domain/             # Modelos de Dados (User, Device, Group)
 │    ├── handlers/           # Controladores de API HTTP
 │    ├── monitoring/         # Motor de ICMP/TCP, Anti-Flapping e Discord Embeds
 │    ├── repositories/       # Queries SQL e Hard Deletes para flexibilidade
 │    ├── services/           # Lógica de Autenticação e Segurança (Bcrypt)
 │    └── websocket/          # Hub de sincronização em Tempo Real via WebSocket
 ├── go.mod                   # Gestão de dependências do Go
 └── netwatch.db              # (Criado automaticamente durante a execução)


 O NetWatch Engine foi desenhado como uma solução "Drop & Run", oferecendo observabilidade profunda e infraestrutura profissional sem a dor de cabeça de configurações complexas.
