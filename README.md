# Go URL Shortener

A high-performance, lightweight URL shortener service built with Go and MongoDB. This service provides a RESTful API to shorten URLs, redirect to original links, list user-specific shortened URLs, and remove existing shortcuts.

## 🚀 Features

- **Quick Shortening**: Generate short, unique codes for any URL.
- **Dynamic Redirection**: Instantly redirect from short codes to original long URLs.
- **User Management**: Track and manage URLs based on `userID`.
- **API-First Design**: Clean RESTful endpoints for integration.
- **MongoDB Backend**: Persistent and scalable storage using MongoDB.
- **Logging Middleware**: Built-in request logging for monitoring.

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Database**: MongoDB (v2 Driver)
- **Environment**: Godotenv for configuration
- **Routing**: standard `net/http` Mux

## 📋 Prerequisites

- [Go](https://golang.org/doc/install) 1.22+ installed
- [MongoDB](https://www.mongodb.com/try/download/community) instance running locally or on the cloud (Atlas)

## ⚙️ Project Structure

```text
url-shortner/
├── config/         # Database connection logic
├── handlers/       # API route handlers
├── middlware/      # HTTP middleware (Logger)
├── models/         # MongoDB data models
├── .env            # Environment configuration
├── main.go         # Application entry point
└── go.mod          # Go module definitions
```

## 🚀 Getting Started

### 1. Clone the repository
```bash
git clone <repository-url>
cd url-shortner
```

### 2. Configure Environment Variables
Create a `.env` file in the root directory:
```env
PORT=3000
MONGO_URI=mongodb://localhost:27017
SERVER=http://localhost:3000/
```

### 3. Install dependencies
```bash
go mod tidy
```

### 4. Run the application
```bash
go run main.go
```

## 🔌 API Endpoints

### 1. Shorten a URL
- **Endpoint**: `POST /shorten`
- **Body**:
  ```json
  {
    "url": "https://www.google.com",
    "userid": "user123"
  }
  ```
- **Response**:
  ```json
  {
    "short_url": "http://localhost:3000/abc123"
  }
  ```

### 2. Redirect to Original URL
- **Endpoint**: `GET /{shortCode}`
- **Behavior**: Redirects with `302 Found` status to the long URL.

### 3. List User URLs
- **Endpoint**: `POST /list`
- **Body**:
  ```json
  {
    "userid": "user123"
  }
  ```

### 4. Remove a URL
- **Endpoint**: `POST /remove`
- **Body**:
  ```json
  {
    "code": "abc123",
    "userid": "user123"
  }
  ```

## 📝 License

This project is open-source and available under the MIT License.
