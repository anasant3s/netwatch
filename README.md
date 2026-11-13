
# NetWatch

**Sistema de gerenciamento, monitoramento e automação de dispositivos de rede.**

## Sobre o projeto

O NetWatch é um projeto backend desenvolvido em Go (Golang), com o objetivo de solucionar uma necessidade identificada no ambiente corporativo: facilitar a visualização, o gerenciamento e o acompanhamento dos dispositivos presentes em uma infraestrutura de rede.

A proposta é desenvolver uma solução centralizada que permita consultar informações dos equipamentos, verificar sua conectividade e, futuramente, integrar esses dados a ferramentas de automação, como o n8n.

O projeto utiliza uma arquitetura baseada em API REST, com foco em organização do código, escalabilidade, integração entre sistemas e monitoramento de infraestrutura.

## Objetivos

- Centralizar informações de dispositivos de rede em uma única aplicação.
- Facilitar a consulta e o gerenciamento de equipamentos.
- Monitorar a conectividade dos dispositivos.
- Identificar possíveis indisponibilidades na infraestrutura.
- Permitir integração com ferramentas externas por meio de APIs.
- Automatizar processos e notificações utilizando n8n.
- Criar uma base para futuras funcionalidades de análise e monitoramento inteligente.

## Funcionalidades

### Implementadas / em desenvolvimento

- [x] Cadastro de dispositivos de rede.
- [x] Listagem de dispositivos cadastrados.
- [x] Consulta individual por identificador.
- [x] Atualização das informações dos dispositivos.
- [x] Exclusão de dispositivos.
- [x] Validação de endereços IP e portas.
- [x] Verificação de conectividade por meio de conexões TCP.
- [x] Identificação de estados de conectividade (`online`, `offline` e `unknown`).
- [x] Respostas HTTP estruturadas em JSON.
- [x] Controle de concorrência no armazenamento dos dados.

### Planejadas

- [ ] Integração com n8n para criação de fluxos automatizados.
- [ ] Notificações automáticas sobre indisponibilidade de dispositivos.
- [ ] Armazenamento persistente em banco de dados.
- [ ] Histórico de verificações e alterações de status.
- [ ] Dashboard para visualização dos dispositivos e seus estados.
- [ ] Monitoramento periódico e automatizado.
- [ ] Exploração de recursos de inteligência artificial para análise de informações e apoio ao diagnóstico.

## Tecnologias e ferramentas

| Tecnologia | Aplicação |
|---|---|
| Go (Golang) | Desenvolvimento do backend |
| API REST | Comunicação entre sistemas |
| HTTP | Comunicação e gerenciamento de requisições |
| JSON | Formato de troca de informações |
| TCP Networking | Verificação de conectividade |
| Postman | Testes e validação dos endpoints |
| Git | Controle de versão |
| GitHub | Versionamento e documentação |
| n8n | Integração e automação de fluxos (planejado) |

### Tecnologias em estudo e expansão

- **PostgreSQL:** persistência e organização dos dados.
- **Docker:** containerização e organização do ambiente.
- **Frontend:** desenvolvimento de uma interface para visualização e gerenciamento dos dispositivos.
- **Inteligência Artificial:** exploração de recursos para análise de dados e auxílio na identificação de possíveis problemas.

## Arquitetura do projeto

A aplicação está sendo organizada em camadas, buscando separar as responsabilidades e facilitar a manutenção e expansão do código.

```text
netwatch-api/
│
├── internal/
│   ├── models/
│   │   └── device.go
│   │
│   ├── repository/
│   │   └── device_repository.go
│   │
│   ├── handlers/
│   │
│   ├── services/
│   │
│   └── database/
│
├── main.go
├── main_test.go
└── go.mod
```

**Responsabilidades das camadas:**

- `models`: definição das estruturas e entidades utilizadas pela aplicação.
- `repository`: gerenciamento e acesso aos dados dos dispositivos.
- `handlers`: tratamento das requisições HTTP.
- `services`: implementação das regras de negócio.
- `database`: futura integração e configuração do banco de dados.
- `main.go`: inicialização da aplicação e configuração do servidor.

A estrutura está em evolução e será ampliada conforme novas funcionalidades forem implementadas.

## API Endpoints

A API utiliza o prefixo `/api`.

| Método | Endpoint | Descrição |
|---|---|---|
| GET | `/api/health` | Verifica a disponibilidade da API |
| GET | `/api/devices` | Lista os dispositivos |
| GET | `/api/devices/{id}` | Consulta um dispositivo |
| POST | `/api/devices` | Cadastra um dispositivo |
| PUT | `/api/devices/{id}` | Atualiza um dispositivo |
| DELETE | `/api/devices/{id}` | Exclui um dispositivo |
| GET | `/api/devices/{id}/status` | Verifica a conectividade de um dispositivo |

## Monitoramento de conectividade

O projeto utiliza conexões TCP para verificar se uma porta específica de determinado dispositivo aceita conexões.

A verificação considera o endereço IP e a porta configurada para cada equipamento.

Os estados utilizados são:

- `online`: conexão TCP estabelecida.
- `offline`: conexão TCP não estabelecida.
- `unknown`: conectividade ainda não verificada.

**Observação:** uma falha na conexão TCP não significa necessariamente que o equipamento esteja desligado. A porta pode estar bloqueada, filtrada ou com algum serviço indisponível.

## Integração planejada com n8n

Uma das propostas de expansão do NetWatch API é sua integração com o n8n, plataforma de automação de workflows.

A integração permitirá utilizar os dados disponibilizados pela API para construir processos automatizados.

### Exemplo de fluxo futuro

1. O n8n inicia uma execução em um intervalo configurado.
2. Realiza uma requisição HTTP ao NetWatch API.
3. Consulta os dispositivos cadastrados e seus estados.
4. Identifica dispositivos que apresentem indisponibilidade.
5. Executa uma ação automatizada, como enviar uma notificação.

Essa integração poderá reduzir a necessidade de consultas manuais e facilitar o acompanhamento da infraestrutura.

## Como executar o projeto

### Pré-requisitos

- Go instalado.
- Git instalado.
- Postman (opcional, para testar a API).

### Clonar o repositório

```bash
git clone URL_DO_REPOSITORIO
```

### Acessar a pasta

```bash
cd netwatch-api
```

### Executar a aplicação

```bash
go run .
```

A API estará disponível em:

```text
http://localhost:8080
```

### Verificar se está funcionando

```http
GET http://localhost:8080/api/health
```

Resposta esperada:

```json
{
  "status": "ok"
}
```

## Testes

Os endpoints podem ser testados utilizando o Postman, permitindo enviar requisições HTTP, configurar parâmetros, encaminhar dados JSON e verificar os códigos de resposta da API.

Os testes automatizados são desenvolvidos utilizando o pacote `testing` da própria linguagem Go.

## Status do projeto

**Em desenvolvimento.**

O projeto está sendo construído de forma incremental, com implementação das funcionalidades principais, organização da arquitetura e planejamento das integrações futuras.


