# Quiz App Backend Documentation

## Base URL

The backend is hosted at: [https://quiz-app-go-backend.onrender.com](https://quiz-app-go-backend.onrender.com)

## Project Setup

### Prerequisites

- Go (version 1.17 or later)

### Clone the Repository

```bash
git clone https://github.com/your-repository/quiz-backend.git
cd quiz-backend
```

### Install Dependencies

Navigate to the project directory and run:

```bash
go mod download
```
### Run the Application

To start the backend server, run:

```bash
go run main.go
```

The server will start on port 8080. You can test it using the provided endpoints.

## API Endpoints

### 1. `GET /users`

**Description:** Fetches a list of all users.

**Request:**

```http
GET /users HTTP/1.1
Host: quiz-app-go-backend.onrender.com
```

**Response:**

- **200 OK**
  
  ```json
  [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "points": 100
    },
    ...
  ]
  ```

### 2. `POST /users`

**Description:** Registers a new user.

**Request:**

```http
POST /users HTTP/1.1
Host: quiz-app-go-backend.onrender.com
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "password123"
}
```

**Response:**

- **201 Created**

  ```json
  {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "points": 0
  }
  ```

- **400 Bad Request** (if password is too short or request is invalid)

  ```json
  {
    "error": "Register: password must be at least 6 characters long"
  }
  ```

- **409 Conflict** (if user already exists)

  ```json
  {
    "error": "Register: user already exists"
  }
  ```

### 3. `POST /login`

**Description:** Authenticates a user and returns user details.

**Request:**

```http
POST /login HTTP/1.1
Host: quiz-app-go-backend.onrender.com
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

**Response:**

- **200 OK**

  ```json
  {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "points": 100
  }
  ```

- **401 Unauthorized** (if invalid email or password)

  ```json
  {
    "error": "Login: invalid email or password"
  }
  ```

### 4. `POST /games`

**Description:** Creates a new game and returns it along with sample questions.

**Request:**

```http
POST /games HTTP/1.1
Host: quiz-app-go-backend.onrender.com
Content-Type: application/json

{
  "name": "Stock Market Quiz",
  "description": "A quiz about stock market basics."
}
```

**Response:**

- **200 OK**

  ```json
  {
    "game": {
      "id": 1,
      "name": "Stock Market Quiz",
      "description": "A quiz about stock market basics."
    },
    "questions": [
      {
        "questionText": "What is a stock?",
        "options": ["A type of bond", "A share in the ownership of a company", "A loan given to the government", "A type of commodity"],
        "correctAnswerIndex": 1,
        "explanation": "A stock represents a share in the ownership of a company..."
      },
      ...
    ]
  }
  ```

### 5. `GET /games`

**Description:** Fetches games available for the user.

**Request:**

```http
GET /games HTTP/1.1
Host: quiz-app-go-backend.onrender.com
```

**Response:**

- **200 OK**

  ```json
  [
    {
      "id": 1,
      "name": "Stock Market Quiz",
      "description": "A quiz about stock market basics."
    },
    ...
  ]
  ```

### 6. `POST /user-games`

**Description:** Creates a new user game entry.

**Request:**

```http
POST /user-games HTTP/1.1
Host: quiz-app-go-backend.onrender.com
Content-Type: application/json

{
  "userID": 1,
  "gameID": 1
}
```

**Response:**

- **200 OK**

  ```json
  {
    "userID": 1,
    "gameID": 1
  }
  ```

- **500 Internal Server Error** (if game already exists for the user)

  ```json
  {
    "error": "CreateUserGame: Game already exists"
  }
  ```

### 7. `GET /users/user`

**Description:** Fetches user details by ID.

**Request:**

```http
GET /users/user?id=1 HTTP/1.1
Host: quiz-app-go-backend.onrender.com
```

**Response:**

- **200 OK**

  ```json
  {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "points": 100
  }
  ```

- **400 Bad Request** (if ID is missing or invalid)

  ```json
  {
    "error": "Invalid user ID"
  }
  ```

- **404 Not Found** (if user not found)

  ```json
  {
    "error": "User not found"
  }
  ```

### 8. `POST /game-end`

**Description:** Ends a game for a user and updates their score.

**Request:**

```http
POST /game-end HTTP/1.1
Host: quiz-app-go-backend.onrender.com
Content-Type: application/json

{
  "user_id": 1,
  "game_id": 1,
  "point": 20,
  "gamePoint": 100
}
```

**Response:**

- **200 OK**

  ```json
  {
    "user_id": 1,
    "game_id": 1,
    "point": 20,
    "gamePoint": 100
  }
  ```

- **500 Internal Server Error** (if game or user game entry not found)

  ```json
  {
    "error": "EndGame: failed to get game"
  }
  ```
