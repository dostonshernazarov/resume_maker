# CV Maker

## Overview
CV Maker is a project developed in Golang with a microservices architecture. It allows users to generate resumes using three different templates in PDF format and sends the generated resume to a Telegram bot. The project uses RabbitMQ for message queuing, PostgreSQL for database management, Minio for file storage, and Docker Compose to run all microservices. The `bot-service` is written in Python using `aiogram` 2.25.

## Microservices
- **API-service**: Handles API requests from the frontend.
- **resume-service**: Generates resumes in PDF format.
- **user-service**: Manages user profiles and authentication.
- **bot-service**: Sends the generated resumes to a Telegram bot.

## Technologies and Libraries
- **Golang**: Used for the main microservices (`API-service`, `resume-service`, `user-service`).
- **Python**: Used for the `bot-service` with `aiogram` 2.25.
- **RabbitMQ**: Used for sending resumes to the Telegram bot.
- **PostgreSQL**: Used for storing user data.
- **Minio**: Used for storing user profile photos and resume files.
- **Docker Compose**: Used for orchestrating all microservices.

## Features
- Generate resumes in PDF format using three different templates.
- Send generated resumes to a Telegram bot.
- Manage user profiles and authentication.
- Store user profile photos and resumes.

## Getting Started

### Prerequisites
Make sure you have the following installed:
- Docker
- Docker Compose
- Go 1.15 or higher
- Python 3.7 or higher
- PostgreSQL

### Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/dostonshernazarov/resume_maker.git
    cd resume_maker
    ```

2. Set up environment variables:
    - Create a `.env` file in the root directory and add the necessary environment variables. Example:
      ```env
      SERVER_HOST=your_server_host
      SERVER_PORT=your_server_port
      REDIS_HOST=your_redis_host
      REDIS_PORT=your_redis_port
      RESUME_SERVICE_GRPC_HOST=resume_service_host
      RESUME_SERVICE_GRPC_PORT=resume_service_port
      USER_SERVICE_GRPC_HOST=user_service_host
      USER_SERVICE_GRPC_PORT=user_service_port
      POSTGRES_USER=your_postgres_user
      POSTGRES_PASSWORD=your_postgres_password
      POSTGRES_DATABASE=cv_maker_db
      AMQP_SERVER=your_rabbitmq_url
      QUEUE_NAME=rabbitmq_queue_name
      MINIO_HOST=your_minio_host
      MINIO_PORT=your_minio_port
      MINIO_ACCESS_KEY=your_minio_access_key
      MINIO_SECRET_KEY=your_minio_secret_key
      TELEGRAM_BOT_TOKEN=your_telegram_bot_token
      CHAT_ID=your_telegram_channel_id
      ```

3. Build and run the services using Docker Compose:
    ```bash
    docker-compose up --build
    ```

### Running the Project
1. Start the Docker containers:
    ```bash
    docker-compose up
    ```

2. The services will be available at the following endpoints:
    - API-service: `http://docker_container_name:9070`
    - User-service: `http://docker_container_name:9090`
    - Resume-service: `http://docker_container_name:9070`
    - Bot-service: Will be running in the background, listening for messages.

## Project Structure
- **api-service/**: Contains the API service code.
- **resume-service/**: Contains the resume generation service code.
- **user-service/**: Contains the user management service code.
- **bot-service/**: Contains the bot service code written in Python.
- **docker-compose.yml**: Docker Compose file to orchestrate the microservices.
- **.env**: Environment variables file.

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact
For any inquiries or questions, please contact [dostonshernazarov989@gmail.com].
