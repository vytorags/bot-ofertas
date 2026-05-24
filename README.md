# ORVIT

<p>
    BOT para o Whatsapp que recebe ofertas de Shopee, Mercado Livre e AliExpress e as encaminha para grupos.
</p>

## Funções
- **Encaminhamento de Ofertas**: Recebe ofertas via HTTP e envia para grupos do WhatsApp
- **Multi-loja**: Suporte a Shopee, Mercado Livre e AliExpress
- **Multi-grupo**: Envia para um ou vários grupos configurados
- **Banco de Dados**: Persistência de sessão com SQLite

## Pré-requisitos
- [Go](https://go.dev/)
- [Git](https://git-scm.com/)

## Instalação/Configuração

#### Clone o repositório
```bash
git clone https://github.com/viitorags/orvit
cd orvit
```

#### Configure as variáveis de ambiente
```bash
cp .env_example .env
```

#### Edite o arquivo .env

| Variável    | Descrição                                                                 |
|:------------|:--------------------------------------------------------------------------|
| `DB_PATH`   | Caminho para o banco SQLite (ex: `file:data/orvit.db?_foreign_keys=on`)  |
| `GROUP_JID` | JID(s) do(s) grupo(s) de destino, separados por vírgula                  |
| `HTTP_PORT` | Porta do servidor HTTP (padrão: `8080`)                                   |

**Exemplo de `.env`:**
```env
DB_PATH="file:data/orvit.db?_foreign_keys=on"
GROUP_JID=111111111-111111@g.us,222222222-222222@g.us
HTTP_PORT=8080
```

> **Como obter o GROUP_JID:** Inicie o bot, envie uma mensagem no grupo desejado e o JID aparecerá nos logs.

#### Baixando as dependências do Go
```bash
go mod download
```

## Executar

O projeto tem um [Makefile](./Makefile) para facilitar a execução:

```bash
make build   # compilar
make run     # executar
make clean   # limpar binário
```

Na primeira execução, um QR Code será exibido no terminal para autenticação no WhatsApp.

## API HTTP

O bot expõe um endpoint para receber ofertas do código Python:

### `POST /offer`

**Body (JSON):**

| Campo           | Tipo   | Obrigatório | Descrição                                             |
|:----------------|:-------|:-----------:|:------------------------------------------------------|
| `store`         | string | Sim         | Loja: `shopee`, `ml`, `mercadolivre`, `aliexpress`    |
| `title`         | string | Sim         | Nome do produto                                       |
| `price`         | string | Sim         | Preço atual                                           |
| `original_price`| string | Não         | Preço original (exibe riscado)                        |
| `discount`      | string | Não         | Percentual de desconto (ex: `50%`)                    |
| `link`          | string | Não         | Link da oferta                                        |
| `group_jid`     | string | Não         | Envia apenas para este grupo (ignora o `.env`)        |

**Exemplo — enviar para todos os grupos configurados:**
```json
{
  "store": "shopee",
  "title": "Fone Bluetooth XYZ",
  "price": "R$ 49,90",
  "original_price": "R$ 99,90",
  "discount": "50%",
  "link": "https://shopee.com.br/..."
}
```

**Exemplo — enviar para um grupo específico:**
```json
{
  "store": "ml",
  "title": "Tênis Esportivo ABC",
  "price": "R$ 199,90",
  "link": "https://mercadolivre.com.br/...",
  "group_jid": "111111111-111111@g.us"
}
```

**Exemplo com Python:**
```python
import requests

requests.post("http://localhost:8080/offer", json={
    "store": "aliexpress",
    "title": "Cabo USB-C 2m",
    "price": "R$ 12,90",
    "original_price": "R$ 30,00",
    "discount": "57%",
    "link": "https://aliexpress.com/..."
})
```

## Estrutura
```
.
├── cmd/bot/        # Ponto de entrada (main.go)
├── data/           # Banco de dados SQLite
├── internal/
│   └── bot/        # Conexão WhatsApp, servidor HTTP e envio de ofertas
└── Makefile        # Scripts de automação
```

## Licença
Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](./LICENSE) para mais detalhes.
