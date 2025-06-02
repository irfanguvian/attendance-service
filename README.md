# Attendance Service API

This is a backend API service for managing employee attendance, built with Go. It features a clean architecture design and uses JWT for authentication with access and refresh tokens.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Authentication](#authentication)
- [Data Models](#data-models)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Installation and Setup](#installation-and-setup)
  - [Running Locally](#running-locally)
  - [Running with Docker](#running-with-docker)
  - [Database Migration](#database-migration)
- [API Endpoints](#api-endpoints)
- [Makefile Commands](#makefile-commands)

## Features

- User authentication (Signup, Login, Signout)
- JWT-based authorization with access and refresh tokens
- Employee management (Create, Update, Delete, Get)
- Attendance tracking (Clock-in)
- List today's attendance

## Architecture

The project follows the principles of Clean Architecture to ensure a separation of concerns, making the codebase modular, scalable, and testable. The main layers are:

-   **Models**: Defines the data structures (database schema).
-   **Repositories**: Handles data access logic, interacting directly with the database.
-   **Services**: Contains the business logic of the application.
-   **Controllers**: Handles incoming HTTP requests, validates input, and calls the appropriate services.
-   **Routes**: Defines the API endpoints and maps them to controller actions.
-   **Middlewares**: Handles cross-cutting concerns like authentication and CORS.

## Authentication

Authentication is handled using JSON Web Tokens (JWT).
-   **Access Token**: A short-lived token used to access protected resources.
-   **Refresh Token**: A long-lived token used to obtain a new access token when the current one expires.

The `protect_middleware.go` ensures that protected routes are only accessible with a valid access token. The `/auth/exchange-token` endpoint allows users to get a new access token using their refresh token.

## Data Models

The service uses the following GORM models:

1.  **`User`**: (`models/user.go`)
    *   `ID`: uint (Primary Key)
    *   `Email`: string (Unique)
    *   `Password`: string
    *   `CreatedAt`: time.Time
    *   `UpdatedAt`: time.Time

2.  **`AccessToken`**: (`models/access_token.go`)
    *   `ID`: string (Primary Key, UUID)
    *   `UserID`: uint (Foreign Key to `User`)
    *   `ExpiredAt`: time.Time
    *   `User`: `User` (Belongs to relationship)

3.  **`RefreshToken`**: (`models/refresh_token.go`)
    *   `ID`: string (Primary Key, UUID)
    *   `AccessTokenID`: string (Foreign Key to `AccessToken`, CASCADE on delete)
    *   `ExpiredAt`: time.Time
    *   `AccessToken`: `AccessToken` (Belongs to relationship)

4.  **`Employees`**: (`models/employee.go`)
    *   `ID`: uint (Primary Key)
    *   `EmpID`: string (Unique, Employee ID)
    *   `Fullname`: string
    *   `CreatedAt`: time.Time
    *   `UpdatedAt`: time.Time
    *   `DeletedAt`: *time.Time (For soft deletes)
    *   `Attandance`: *[]`Attendance` (Has many relationship)

5.  **`Attendance`**: (`models/attandance.go`)
    *   `ID`: uint (Primary Key)
    *   `EmployeeID`: uint (Foreign Key to `Employees`)
    *   `ClockIn`: time.Time
    *   `CreatedAt`: time.Time
    *   `UpdatedAt`: time.Time
    *   `Employee`: `Employees` (Belongs to relationship)

## Technologies Used

-   **Go**: Programming language
-   **Gin**: HTTP web framework
-   **GORM**: ORM library for database interaction
-   **PostgreSQL**: SQL database
-   **JWT (golang-jwt/jwt)**: For token-based authentication
-   **godotenv**: For managing environment variables
-   **Docker**: For containerization
-   **CompileDaemon**: For live reloading during development

## Prerequisites

-   Go (version 1.22.3 or higher recommended)
-   PostgreSQL
-   Docker (optional, for containerized setup)
-   Make (for using Makefile commands)
-   CompileDaemon (for live reload `make serve`)
    ```bash
    go install github.com/githubnemo/CompileDaemon@latest
    ```

## Installation and Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/irfanguvian/attendance-service.git
    cd attendance-service
    ```

2.  **Set up environment variables:**
    Create a `.env` file in the root directory by copying the example or creating a new one:
    ```env
    DATA_BASE_URL="host=localhost user=postgres password='' dbname=attendance-service port=5432 TimeZone=UTC"
    PORT=your_db_user
    JWT_SECRET_KEY=attendance_db
    ```
    Ensure your PostgreSQL server is running and the specified database exists.

3.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

### Running Locally

1.  **Run database migrations:**
    This will create the necessary tables in your database.
    ```bash
    make run-migration
    ```

2.  **Start the server:**
    This command uses CompileDaemon for live reloading.
    ```bash
    make serve
    ```
    The API will be available at `http://localhost:3000` (or the port specified in your `.env` or config if different). The default port in `docker-compose.yml` is 3000.

### Running with Docker

1.  **Ensure Docker and Docker Compose are installed.**

2.  **Build and run the application using Docker Compose:**
    This will build the Go application image and start the `app` service.
    ```bash
    docker-compose up --build
    ```
    The API will be available at `http://localhost:3000`.
    *Note: The provided `docker-compose.yml` only defines the `app` service. You might need to add a PostgreSQL service to it for a complete containerized development environment or ensure your local PostgreSQL is accessible from the Docker container (e.g., using `host.docker.internal` for the `DB_HOST`).*

    To run migrations when using Docker, you can execute the command inside the running container:
    ```bash
    docker-compose exec app go run ./migrate/migrate.go
    ```
    Or, modify the Dockerfile or docker-compose entrypoint/command to run migrations on startup.

### Database Migration

To set up or update the database schema, run the migration command:
```bash
make run-migration
```
This executes the `migrate/migrate.go` file, which uses GORM's `AutoMigrate` feature to create or update tables based on the defined models.

## API Endpoints

**Auth Routes (`/auth`)**
-   `POST /login`: User login.
    -   Request Body: `dto.LoginBody` (`{ "email": "user@example.com", "password": "password123" }`)
-   `POST /signup`: User registration.
    -   Request Body: `dto.SignupBody` (`{ "email": "user@example.com", "password": "password123" }`)
-   `POST /signout`: User signout (requires Authorization token).
-   `POST /exchange-token`: Exchange a refresh token for a new access token.
    -   Request Body: `{ "refresh_token": "your_refresh_token" }`

**Employee Routes (`/employee`)** (All require Authorization token)
-   `POST /create`: Create a new employee.
    -   Request Body: `dto.CreateEmployeeBody` (`{ "fullname": "John Doe" }`)
-   `PUT /update`: Update an existing employee.
    -   Request Body: `dto.UpdateEmployeeBody` (`{ "id": 1, "fullname": "John Doe Updated" }`)
-   `DELETE /:employeeID`: Delete an employee by ID.
    -   Path Param: `employeeID` (uint)
-   `GET /:employeeID`: Get an employee by ID.
    -   Path Param: `employeeID` (uint)
-   `GET /list`: Get a paginated list of all employees.
    -   Query Params: `page` (int8), `limit` (int8)

**Attendance Routes (`/attendance`)** (All require Authorization token)
-   `POST /create`: Create a new attendance record (clock-in).
    -   Request Body: `dto.CreateAttendanceBody` (`{ "employee_id": 1 }`)
-   `GET /list-today`: Get a list of today's attendance records.
    -   Query Params: `page` (int8), `limit` (int8)
-   `GET /salaries`: Get a list of salaries employee records on a month.
    -   Query Params: `page` (int8), `limit` (int8), `start_date` (yyyy-mm-dd), `end_date` (yyyy-mm-dd)
-   `GET /summary-today`: Get a summary of today data.
-   `GET /daily-trends`: Get a daily trends data.
    -   Query Params: `start_date` (yyy-mm-dd), `end_date` (yyy-mm-dd)
-   `GET /monthly-trends`: Get a monthly trends data.
    -   Query Params: `start_date` (yyy-mm-dd), `end_date` (yyy-mm-dd)

**Health Check**
-   `GET /ping`: Returns a pong response.
    ```json
    {
        "message": "pong"
    }
    ```

## Makefile Commands

-   `make serve`: Runs the application with live reload using CompileDaemon.
-   `make run-migration`: Executes the database migrations.

---

This README provides a comprehensive overview of the Attendance Service API.
