# CantTouchMe Project 🚀

This project consists of a backend (Go) and frontend (Vue.js) application with development and production configurations.

## Warning ⚠️

**Important:** Only work on the `dev` branch. All changes to the `prod` branch will be automatically built and deployed.

## Docker Compose Commands

### Development Environment

To start the development environment (includes hot-reloading for the frontend):

```bash
docker compose --profile dev up --build
```

To stop the development environment:

```bash
docker compose --profile dev down
```

## Additional Commands 🛠️

To view logs:

```bash
docker compose logs -f
```

## Environment Configuration

This project uses environment variables which should be defined in a `.env` file at the root of the project. The following variables are used:

- `API_PORT`: Port for the backend API
- `MYSQL_ROOT_PASSWORD`: Root password for MySQL
- `MYSQL_DATABASE`: MySQL database name
- `MYSQL_USER`: MySQL username
- `MYSQL_PASSWORD`: MySQL password
- `LOG_LEVEL`: Logging level for the backend

## Project Structure 📂

- `backend/`: Go API server
- `frontend/`: Vue.js web application
- `db/`: Database initialization scripts
