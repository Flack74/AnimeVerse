# 🌸 AnimeVerse API

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/Flack74/AnimeVerse) [![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE) [![Go Version](https://img.shields.io/badge/go-1.24+-blue)](https://golang.org/) [![Chi Router](https://img.shields.io/badge/router-chi-orange)](https://github.com/go-chi/chi)

Welcome to **AnimeVerse** – your ultimate RESTful API for managing and exploring your favorite anime collection! Built with **Go**, **MongoDB**, and **Chi Router**, this API provides a robust, scalable solution for anime enthusiasts. Whether you're a casual fan or a hardcore otaku, AnimeVerse API has got you covered! 🎉

---

## 🚀 Features

- **📝 CRUD Operations:** Create, Read, Update, and Delete anime records effortlessly
- **🔄 Partial Updates:** Send JSON payload with only the fields you need to update
- **🏛️ MongoDB Integration:** Secure and scalable NoSQL database storage
- **🌐 RESTful Design:** Clean, intuitive, and standardized API endpoints
- **🛡️ Chi Router:** Fast, lightweight HTTP router with middleware support
- **🔒 CORS Support:** Cross-Origin Resource Sharing enabled for web applications
- **📊 Request Logging:** Comprehensive logging with Chi middleware
- **⏱️ Timeout Protection:** Request timeout handling for better reliability
- **🗃️ Response Compression:** Automatic gzip compression for better performance
- **🚫 Duplicate Prevention:** Prevents duplicate anime entries by name
- **📊 Detailed Data:** Manage anime with name, type, score, progress, status, genre, and notes
- **🚑 Graceful Shutdown:** Proper server shutdown handling

---

## 📥 Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/Flack74/AnimeVerse.git
   cd AnimeVerse
   ```

2. **Install Dependencies:**

   Make sure you have [Go](https://golang.org/) installed. Then, run:

   ```bash
   go mod tidy
   ```

3. **Configure Environment Variables:**

   Create a `.env` file in the root directory with the following (adjust as needed):

   ```env
   ConnectionString=mongodb://<username>:<password>@localhost:27017
   DBName=anime
   CollectionName=watchlist
   ```

4. **Run the Application:**

   **Development Mode (with hot reload):**
   ```bash
   # Install Air for hot reloading (if not already installed)
   go install github.com/air-verse/air@latest
   
   # Run with hot reload
   air
   ```

   **Production Mode:**
   ```bash
   go run main.go
   ```

   The API will be available at [http://localhost:8000](http://localhost:8000) 🎉

5. **Docker Support (Optional):**

   **Development:**
   ```bash
   docker build --target dev -t animeverse-dev .
   docker run -p 8000:8000 -v $(pwd):/src animeverse-dev
   ```

   **Production:**
   ```bash
   docker build -t animeverse-prod .
   docker run -p 8000:8000 animeverse-prod
   ```

---

## ⚙️ Usage

### **Home Page**

Visit [http://localhost:8000](http://localhost:8000) to see a welcoming homepage.

### **API Requests**

- **Create an Anime:**

  ```bash
  curl -X POST http://localhost:8000/api/anime \
    -H "Content-Type: application/json" \
    -d '{
          "name": "My Hero Academia",
          "type": "TV",
          "score": 9,
          "progress": {
              "watched": 25,
              "total": 88
          },
          "status": "watching",
          "genre": ["action", "shounen", "superhero"],
          "notes": "Amazing superhero anime with great character development"
        }'
  ```

  **Response:**
  ```json
  {
    "success": true,
    "message": "Anime created successfully",
    "data": {
      "_id": "...",
      "name": "My Hero Academia",
      "type": "TV",
      "score": 9,
      "progress": {
        "watched": 25,
        "total": 88
      },
      "status": "watching",
      "genre": ["action", "shounen", "superhero"],
      "notes": "Amazing superhero anime with great character development"
    }
  }
  ```

- **Get All Anime:**

  ```bash
  curl http://localhost:8000/api/animes
  ```

  **Response:**
  ```json
  {
    "success": true,
    "message": "Animes retrieved successfully",
    "data": [
      {
        "_id": "...",
        "name": "Attack on Titan",
        "type": "TV",
        "score": 10,
        "status": "completed"
      }
    ]
  }
  ```

- **Get an Anime by Name:**

  ```bash
  curl http://localhost:8000/api/anime/attack-on-titan
  # or
  curl http://localhost:8000/api/anime/attack_on_titan
  ```

- **Update an Anime (Partial Update):**

  ```bash
  curl -X PUT http://localhost:8000/api/anime/{id} \
    -H "Content-Type: application/json" \
    -d '{
          "score": 10,
          "progress": {
              "watched": 88,
              "total": 88
          },
          "status": "completed",
          "notes": "Masterpiece! One of the best anime ever made."
        }'
  ```

- **Delete an Anime:**

  ```bash
  curl -X DELETE http://localhost:8000/api/anime/{id}
  ```

- **Delete All Anime:**

  ```bash
  curl -X DELETE http://localhost:8000/api/deleteallanime
  ```

---

## 🖼️ API Endpoints

| Method | Endpoint                  | Description                              |
| ------ | ------------------------- | ---------------------------------------- |
| GET    | `/`                       | Homepage with API information            |
| GET    | `/health`                 | Health check endpoint                    |
| GET    | `/api/animes`             | Retrieve all anime records               |
| GET    | `/api/anime/{animeName}`  | Retrieve a specific anime by name        |
| POST   | `/api/anime`              | Create a new anime record                |
| PUT    | `/api/anime/{id}`         | Update an anime record (partial update)  |
| DELETE | `/api/anime/{id}`         | Delete a specific anime record           |
| DELETE | `/api/deleteallanime`     | Delete all anime records                 |

📚 **For detailed API documentation with examples, see [API_DOCS.md](API_DOCS.md)**

---

## 📸 Screenshots

 ### Home Page ("/")
 ![Home Page](https://github.com/user-attachments/assets/6399dad4-a54a-4927-ad23-618b4d63f148)

### API JSON Output ("/api/animes")
 ![API JSON Output](https://github.com/user-attachments/assets/2c075413-dba0-4a5a-a813-838138547791)

---

## 🤝 Contributing

Contributions are always welcome! If you have ideas, suggestions, or improvements, please follow these steps:

1. **Fork the Repository**
2. **Create a Feature Branch:**  
   `git checkout -b feature/your-feature`
3. **Commit Your Changes:**  
   `git commit -am 'Add some feature'`
4. **Push to the Branch:**  
   `git push origin feature/your-feature`
5. **Open a Pull Request**

---

## 📝 Future Improvements for AnimeVerse API

### Bulk Anime Insertion (Planned Feature)

#### Overview
A potential enhancement for AnimeVerse API is adding support for bulk anime insertion. This would allow users to submit an array of anime data in a single API request, improving efficiency when adding multiple entries.

#### Proposed Route

- **Endpoint:** `POST /api/addmultipleanimes`
- **Request Body:** An array of anime JSON objects, e.g.,
  ```json
  [
    {"name": "Attack on Titan", "genre": "Action", "episodes": 75},
    {"name": "Demon Slayer", "genre": "Adventure", "episodes": 26}
  ]
  ```
- **Response:**
  - **Success:** Returns inserted anime details.
  - **Failure:** Returns an error message for invalid or duplicate entries.

#### Considerations
- **MongoDB Free Tier Limitations:** The free plan has a 512MB storage limit, so bulk insertion must be optimized.
- **Duplicate Handling:** Implement logic to prevent inserting duplicate anime based on name.
- **Abuse Prevention:** Consider adding rate limiting or authentication to prevent excessive API calls.
- **Performance Optimization:** Using MongoDB’s InsertMany function would be more efficient than inserting each entry individually.

---

## 🐳 Docker Support

AnimeVerse supports Docker for both **development** and **production** environments using multi-stage builds.

### For **Development**:

```bash
docker build --target dev -t animeverse-dev .
docker run -p 8000:8000 -v $(pwd):/app animeverse-dev
```

#### 🔹 For **Production**:

```bash
docker build -t animeverse-prod .
docker run -p 8000:8000 animeverse-prod
```

> Make sure you’ve configured your `.env` properly and MongoDB is accessible.

---

## 🌆 What's New in v2.0

### 🔄 Migration from Gorilla Mux to Chi Router

AnimeVerse API has been upgraded with significant improvements:

- **⚡ Performance:** Chi router provides better performance and lower memory footprint
- **🔧 Middleware:** Built-in middleware support for logging, recovery, CORS, and compression
- **📊 Better Responses:** Standardized JSON response format with success/error indicators
- **🚫 Input Validation:** Enhanced request validation and error handling
- **🚑 Graceful Shutdown:** Proper server shutdown handling
- **📝 Improved Logging:** Better request/response logging

### 🆕 Response Format

All API responses now follow a consistent format:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { /* actual data */ },
  "error": "" // only present when success is false
}
```

### 🔧 Technical Improvements

- **Chi Router:** Migrated from Gorilla Mux to Chi for better performance
- **Middleware Stack:** Request logging, recovery, timeout, and compression
- **CORS Support:** Cross-origin requests enabled for web applications
- **Graceful Shutdown:** Proper signal handling for clean server shutdown
- **Docker Multi-stage:** Optimized Docker builds for development and production

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

---

Made with ❤️ by Flack

