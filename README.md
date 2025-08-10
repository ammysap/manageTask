Got it — here’s your README.md exactly as you provided, without any single word changed:

Task management Service
yes the service will run like command docker compose -f internal/dockerfiles/docker-compose.yml up --build from root folder
yes
dont mention

Task Management Service

Task Management Service is a Go-based microservices application for managing tasks and users. It consists of two decoupled services: the Taskmanager (HTTP/REST) for task CRUD with pagination/filtering, and the User service (gRPC) for user details. Both services follow single-responsibility principles with clean directory layouts under internal/ (e.g. services, database, logging, dockerfiles) and communicate over gRPC internally ￼ ￼.

Architecture & Design Decisions
• Microservices: We use a microservices architecture to isolate functionality. The Taskmanager handles task data, and the User service handles user data. This single-responsibility approach means each service can evolve and scale independently ￼ ￼.
• Service Structure:
• Taskmanager (HTTP API) is organized into layers: routes/handlers (defining REST endpoints), service (business logic), DAO (data access), and models (data structures).
• User service (gRPC API) defines protobufs for user data and a corresponding service layer. Both have their own main.go and Dockerfile under internal/services/.
• Docker Compose & Isolation: Each service has its own Dockerfile and runs in a separate container. Docker Compose is used for orchestration, exposing Taskmanager on port 8080 and User gRPC on 50051. This setup isolates dependencies and ensures consistent environments, making it easy to scale and deploy ￼ ￼.
• Inter-service Communication: We use gRPC (Protocol Buffers over HTTP/2) for efficient, type-safe calls between services. gRPC is high-performance and easy to implement, making it ideal for internal APIs ￼.
• Configuration: Services support environment variables for ports and endpoints (configured in docker-compose.yml), allowing flexible deployment and easy linking between containers ￼.

How to Run 1. Prerequisites: Install Docker and Docker Compose. 2. Start Services: From the project root, run:

docker compose -f internal/dockerfiles/docker-compose.yml up --build

    3.	Ports:
    •	Taskmanager (HTTP API): localhost:8080
    •	User service (gRPC): localhost:50051

Once running, you can access the Taskmanager REST endpoints and call the User service via gRPC.

API Reference

Use the following example commands and JSON responses to interact with the APIs. Replace IDs and parameters as needed.
• Create Task:

curl -X POST -H "Content-Type: application/json" -d '{"title":"Task 1","description":"Sample task"}' http://localhost:8080/tasks

Response (JSON):

{
"id": 1,
"title": "Task 1",
"description": "Sample task",
"status": "pending"
}

    •	Get Tasks (with pagination/filter):

curl "http://localhost:8080/tasks?page=1&limit=10&status=pending"

Response:

{
"tasks": [
{"id":1,"title":"Task 1","description":"Sample task","status":"pending"}
],
"page": 1,
"limit": 10,
"total": 1
}

    •	Get Task by ID:

curl http://localhost:8080/tasks/1

Response:

{
"id": 1,
"title": "Task 1",
"description": "Sample task",
"status": "pending"
}

    •	Update Task:

curl -X PUT -H "Content-Type: application/json" -d '{"title":"Updated Task 1","status":"done"}' http://localhost:8080/tasks/1

Response:

{
"id": 1,
"title": "Updated Task 1",
"description": "Sample task",
"status": "done"
}

    •	Delete Task:

curl -X DELETE http://localhost:8080/tasks/1

Response:

{
"message": "Task deleted"
}

    •	Get User by ID (gRPC): Use grpcurl or a gRPC client. For example:

grpcurl -plaintext -d '{"id":1}' localhost:50051 user.UserService/GetUser

Response:

{
"id": 1,
"name": "Alice",
"email": "alice@example.com"
}

Microservices Demonstration

This project exemplifies microservice best practices:
• Decoupled services: Task and User logic are separated into independent services, each handling one domain (tasks or users) ￼ ￼.
• Containerization: Every service has its own Dockerfile and runs in its own container. Docker Compose links them, enabling independent scaling and deployment ￼ ￼.
• Efficient communication: Services communicate over gRPC (Protocol Buffers/HTTP2), which is performant and well-suited for inter-service RPC calls ￼.
• Configurable: Ports, URLs, and other settings are managed via environment variables and Compose files, making the system flexible across environments ￼.

These choices ensure a scalable, maintainable system where each service can be developed, tested, and deployed separately, while still working together seamlessly.

Do you want me to also create the README.md file and save it to your project folder? That way you can directly commit it.
