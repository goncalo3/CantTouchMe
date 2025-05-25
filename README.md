# CantTouchMe Project 🚀

Este projeto consiste numa aplicação com backend em **Go** e frontend em **Vue.js**, com configurações separadas para desenvolvimento e produção.

🌐 **Aplicação em Produção:** [https://canttouchme.goncalo3.pt](https://canttouchme.goncalo3.pt)

## 📚 Contexto Académico

Projeto desenvolvido no âmbito da **Licenciatura em Engenharia Informática** da **Universidade da Beira Interior**, na unidade curricular de **Segurança Informática**.

## ⚠️ Aviso Importante

Trabalhar **apenas** na branch `dev`. Todas as alterações na branch `prod` serão automaticamente colocadas em produção.

## 🐳 Comandos Docker Compose

### Ambiente de Desenvolvimento

Para iniciar o ambiente de desenvolvimento (com hot-reload para o frontend e backend):

```bash
docker compose --profile dev up --build
````

Para parar o ambiente de desenvolvimento:

```bash
docker compose --profile dev down
```

### Ver Logs

```bash
docker compose logs -f
```

## 🔐 Configuração do Ambiente

Para correr a aplicação, é necessário criar um ficheiro `.env` com base no modelo:

```bash
cp .env.template .env
```

Para gerar uma chave segura para o JWT, usa:

```bash
openssl rand -base64 32
```

Copia o valor gerado para a variável JWT no teu `.env`.

## 📁 Estrutura do Projeto

* `backend/`: Servidor API em Go
* `frontend/`: Aplicação web em Vue.js
* `db/`: Scripts de inicialização da base de dados

## 👥 Autores

* Carolina Fernandes - Nº Aluno 50252
* Beatriz Laranjinha - Nº Aluno 50521
* Diogo Araújo       - Nº Aluno 49680
* Diogo Rodrigues    - Nº Aluno 49658
* Gonçalo Moreira    - Nº Aluno 49447
* Rodrigo Esteves    - Nº Aluno 49454