# lms-be

## Introduction

`lms-be` is an API service designed to support the Learning Management System frontend (`lms-fe`). It provides a robust set of endpoints for managing courses, students, instructors, and more, offering a seamless integration with the frontend application to create a comprehensive learning management experience.

## Getting Started

These instructions will get your copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Before you begin, ensure you have the following installed:
- Go (version 1.19 or later recommended)
- Any other dependencies your project might need (e.g., PostgreSQL, Docker if applicable)

### Setting Up Environment Variables

Before running the service, you need to set up the required environment variables. The project includes a `.env.example` file with all the necessary environment variable declarations.

To set up your environment variables:

1. Copy the `.env.example` file to a new file named `.env` in the same directory:

    ```bash
    cp .env.example .env
    ```

2. Open the `.env` file in your preferred text editor and modify the environment variables to fit your development environment.

    Make sure to replace placeholder values with actual data where necessary.

### Running Migrations

Before starting the service, you need to run the database migrations to set up the necessary tables and schema. The project uses `migrate` for managing database migrations.

```bash
    go run main.go migrate:up
```
Or if you want to run with seeder you can run this command

```bash
    go run main.go migrate:fresh
```
But it will be down your migration first

### Running the Service

To run `lms-be` on your local machine, you have two options:

#### Using Make

If you have `make` installed, you can start the service using the following command:

```bash
make serve-http
```

This command is configured in your Makefile to set up necessary environment variables, compile the code, and run the service, making it ready to accept requests from lms-fe.

#### Using Go Directly
Alternatively, you can run the service directly using Go with the following command:
```bash
go run main.go http
```

