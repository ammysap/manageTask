# **Task Management Service**

Task Management Service is a Go-based microservices application for managing tasks and users. It consists of two decoupled services: the Taskmanager (HTTP/REST) for task CRUD with pagination/filtering, and the User service (gRPC) for user details. Both services follow single-responsibility principles with clean directory layouts under internal/ (e.g. services, database, logging, dockerfiles) and communicate over gRPC internally  .

---

## **Architecture & Design Decisions**

- **Microservices:** We use a microservices architecture to isolate functionality. The Taskmanager handles task data, and the User service handles user data. This single-responsibility approach means each service can evolve and scale independently .
- **Service Structure:**
  - Taskmanager (HTTP API) is organized into layers: routes/handlers (defining REST endpoints), service (business logic), DAO (data access), and models (data structures).
  - User service (gRPC API) defines protobufs for user data and a corresponding service layer. Both have their own main.go and Dockerfile under internal/services/.
- **Docker Compose & Isolation:** Each service has its own Dockerfile and runs in a separate container. Docker Compose is used for orchestration, exposing Taskmanager on port 8080 and User gRPC on 50051. This setup isolates dependencies and ensures consistent environments, making it easy to scale and deploy .
- **Inter-service Communication:** We use gRPC (Protocol Buffers over HTTP/2) for efficient, type-safe calls between services. gRPC is high-performance and easy to implement, making it ideal for internal APIs .
- **Configuration:** Services support environment variables for ports and endpoints (configured in docker-compose.yml), allowing flexible deployment and easy linking between containers .

---

## **How to Run**

1. **Prerequisites:** Install Docker and Docker Compose.
2. **Start Services:** From the project root, run:

```json
docker compose -f internal/dockerfiles/docker-compose.yml up --build
```

1. **Ports:**
   - Taskmanager (HTTP API): localhost:8080
   - User service (gRPC): localhost:50051

Once running, you can access the Taskmanager REST endpoints and call the User service via gRPC.

---

## **API Reference**

Use the following example commands and JSON responses to interact with the APIs. Replace IDs and parameters as needed.

---

**Create Task**

```json
curl --location 'http://localhost:8080/createTask' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Finish Docker Setup 2",
    "description": "Set up Docker Compose for TaskManager with Postgres",
    "status": "in_progress",
    "priority": 0,
    "assignee_id": 2,
    "creator_id": 1,
    "due_at": "2025-08-15T17:00:00Z"
}
'
```

**Response (JSON):**

```json
{
  "message": "Task created successfully",
  "success": true
}
```

---

**Get Tasks (with pagination/filter)**

```json
curl --location 'http://localhost:8080/getTasks' \
--header 'Content-Type: application/json' \
--data '{
    "page": {
        "number": 0,
        "skip": 0,
        "rows": 10
    },
    "filter": {
        "field": "status",
        "type": "IN",
        "values": ["completed"]
    }
}'
```

**Response:**

```json
{
  "data": {
    "hits": [
      {
        "ID": 1,
        "CreatedAt": "2025-08-10T11:32:07.306741Z",
        "UpdatedAt": "2025-08-10T12:39:48.816058Z",
        "DeletedAt": null,
        "name": "Finish Docker Setup done",
        "description": "Set up Docker Compose for TaskManager with Postgres done",
        "status": "completed",
        "priority": 2,
        "assignee_id": 2,
        "creator_id": 1,
        "due_at": "2025-08-16T17:00:00Z"
      }
    ],
    "hasMore": false,
    "total": 1
  },
  "success": true
}
```

---

**Get Task by ID**

```json
curl --location 'http://localhost:8080/getTasks/2' \
--data ''
```

**Response:**

```json
{
  "ID": 2,
  "CreatedAt": "2025-08-10T11:41:27.746348Z",
  "UpdatedAt": "2025-08-10T11:41:27.746348Z",
  "DeletedAt": null,
  "name": "Finish Docker Setup",
  "description": "Set up Docker Compose for TaskManager with Postgres",
  "status": "in_progress",
  "priority": 0,
  "assignee_id": 2,
  "creator_id": 1,
  "due_at": "2025-08-15T17:00:00Z"
}
```

---

**Update Task**

```json
curl --location --request PUT 'http://localhost:8080/updateTask' \
--header 'Content-Type: application/json' \
--data '{
    "ID": 1,
    "CreatedAt": "2025-08-10T17:02:07.306741+05:30",
    "UpdatedAt": "2025-08-10T17:02:07.306741+05:30",
    "DeletedAt": null,
    "name": "Finish Docker Setup done",
    "description": "Set up Docker Compose for TaskManager with Postgres done",
    "status": "completed",
    "priority": 2,
    "assignee_id": 2,
    "creator_id": 1,
    "due_at": "2025-08-16T22:30:00+05:30"
}'
```

**Response:**

```json
{
  "data": {
    "ID": 1,
    "CreatedAt": "2025-08-10T17:02:07.306741+05:30",
    "UpdatedAt": "2025-08-10T19:02:03.715536301Z",
    "DeletedAt": null,
    "name": "Finish Docker Setup done",
    "description": "Set up Docker Compose for TaskManager with Postgres done",
    "status": "completed",
    "priority": 2,
    "assignee_id": 2,
    "creator_id": 1,
    "due_at": "2025-08-16T22:30:00+05:30"
  },
  "success": true
}
```

---

**Delete Task**

```json
curl --location --request DELETE 'http://localhost:8080/deleteTask/3' \
--data ''
```

**Response:**

```json
{
  "message": "Task deleted successfully",
  "success": true
}
```

---

**Get User by ID (gRPC Inter Service Communication)**

Example:

```json
curl --location 'http://localhost:8080/getUser/1' \
--data ''
```

**Response:**

```json
{
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "name": "John Doe"
  }
}
```

---

## **Supported Filters in Pagination API**

The pagination API supports a rich set of filters that can be combined to perform complex queries. Each filter applies to a specific field in the dataset and can be combined using logical operators like AND and OR.

Below is a description of each filter type from the FilterType enum:

---

### **Equality & Comparison**

| **Filter Type**                  | **Description**                                                     | **Example**        |
| -------------------------------- | ------------------------------------------------------------------- | ------------------ |
| **EQUALS**                       | Matches values exactly.                                             | status = "active"  |
| **NOT_EQUALS**                   | Matches values that are not equal.                                  | status != "active" |
| **GT** (Greater Than)            | Matches values greater than a specified number or date.             | age > 18           |
| **GTE** (Greater Than or Equals) | Matches values greater than or equal to a specified number or date. | age >= 18          |
| **LT** (Less Than)               | Matches values less than a specified number or date.                | price < 100        |
| **LTE** (Less Than or Equals)    | Matches values less than or equal to a specified number or date.    | price <= 100       |

---

### **String Matching**

| **Filter Type** | **Description**                                                                                 | **Example**                  |
| --------------- | ----------------------------------------------------------------------------------------------- | ---------------------------- |
| **CONTAINS**    | Matches values that contain the given substring (case-insensitive unless configured otherwise). | name contains "john"         |
| **LIKE**        | SQL-style pattern matching using % for wildcards.                                               | email LIKE "%@gmail.com"     |
| **STARTS_WITH** | Matches values starting with the given prefix.                                                  | username starts with "admin" |
| **ENDS_WITH**   | Matches values ending with the given suffix.                                                    | filename ends with ".pdf"    |

---

### **Set Membership**

| **Filter Type** | **Description**                                  | **Example**                            |
| --------------- | ------------------------------------------------ | -------------------------------------- |
| **IN**          | Matches values present in the provided list.     | category IN ["tech", "finance"]        |
| **NOT_IN**      | Matches values not present in the provided list. | category NOT IN ["banned", "archived"] |

---

### **Logical Operators**

| **Filter Type** | **Description**                                     | **Example**                          |
| --------------- | --------------------------------------------------- | ------------------------------------ |
| **AND**         | Combines multiple filters; all must match.          | status = "active" AND role = "admin" |
| **OR**          | Combines multiple filters; at least one must match. | role = "admin" OR role = "moderator" |
| **NOT**         | Negates a filter.                                   | NOT (status = "inactive")            |

---

### **Range Filtering**

A Range object can be provided for numeric or date-based filtering.

- **from**: The lower bound (inclusive).
- **to**: The upper bound (inclusive).

Example:

```json
{
  "field": "created_at",
  "type": "GTE",
  "range": { "from": "2024-01-01", "to": "2024-12-31" }
}
```

---

### **Value Types**

The ValueType defines the expected data type for the filter value(s):

- **string** – for text fields (e.g., names, categories)
- **numeric** – for numbers (e.g., price, quantity, age)
- **boolean** – for true/false values

---

### **Example PaginatedRequest**

```json
{
  "filter": {
    "field": "status",
    "type": "EQUALS",
    "values": ["active"],
    "valueType": "string"
  },
  "sorts": [{ "field": "created_at", "order": "desc" }],
  "page": {
    "number": 1,
    "skip": 0,
    "rows": 20
  }
}
```

---

## **Microservices Demonstration**

This project exemplifies microservice best practices:

- **Decoupled services:** Task and User logic are separated into independent services, each handling one domain (tasks or users) .
- **Containerization:** Every service has its own Dockerfile and runs in its own container. Docker Compose links them, enabling independent scaling and deployment .
- **Efficient communication:** Services communicate over gRPC (Protocol Buffers/HTTP2), which is performant and well-suited for inter-service RPC calls .
- **Configurable:** Ports, URLs, and other settings are managed via environment variables and Compose files, making the system flexible across environments .

These choices ensure a scalable, maintainable system where each service can be developed, tested, and deployed separately, while still working together seamlessly.

---

this is my documentaion for the assignemt for notion now i need the same for my readme.md please write the readme.md file contetnt with the above content only striclty dont change any single word and design. also please give me the readme.md file so that it also looks beautifull
