# AnimeVerse

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/Flack74/AnimeApi) [![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)


**AnimeVerse** is a RESTful API for managing and exploring your anime collection. Built with **Go**, **MongoDB**, and **Gorilla Mux**, it provides clean CRUD endpoints, partial updates, and robust duplicate prevention—perfect for casual watchers and hardcore otaku alike! 🎉


## 🚀 Features

- **CRUD Operations**: Create, Read, Update, Delete anime records.  
- **Partial Updates**: Send only the fields you need to change.  
- **MongoDB Integration**: Secure, scalable storage.  
- **RESTful Endpoints**: Clean URL design.  
- **Rich Data Model**: Fields for `name`, `type`, `score`, `progress`, `status`, and `tags`.  
- **Duplicate Prevention**: No repeated entries by name.

---

## 📥 Installation

1. **Clone the repo**  
   ```bash
   git clone https://github.com/Flack74/AnimeVerse.git
   cd AnimeVerse
   ```

2. **Install Go modules**  
   ```bash
   go mod tidy
   ```

3. **Configure environment** (see below).  
4. **Run locally**  
   ```bash
   go run main.go
   ```
   The API will be available at `http://localhost:8000`.

---

## ⚙️ Configuration

Create a `.env` file in the project root:

```env
ConnectionString=mongodb://<username>:<password>@localhost:27017
DBName=anime
CollectionName=watchlist
```

> **Tip:** Never commit your `.env` to version control.

---

## ⚙️ Usage

### Home Page

Open your browser to:

```
http://localhost:8000
```

You’ll see a friendly welcome message.

### API Endpoints

| Method | Endpoint                  | Description                              |
| :----- | :------------------------ | :--------------------------------------- |
| **GET**    | `/api/animes`             | List all anime records                   |
| **GET**    | `/api/anime/{name}`       | Get one anime by **name**                |
| **POST**   | `/api/anime`              | Create a new anime                       |
| **PUT**    | `/api/anime/{id}`         | Partially update an anime by **ID**      |
| **DELETE** | `/api/anime/{id}`         | Delete one anime by **ID**               |
| **DELETE** | `/api/deleteallanime`     | Delete _all_ anime records               |

#### Examples

- **Create**  
  ```bash
  curl -X POST http://localhost:8000/api/anime \
    -H "Content-Type: application/json" \
    -d '{
      "name": "My Hero Academia",
      "type": "TV",
      "score": 9,
      "progress": {"watched": 25, "total": 88},
      "status": "watching",
      "tags": ["action","shounen","superhero"]
    }'
  ```

- **Read All**  
  ```bash
  curl http://localhost:8000/api/animes
  ```

- **Update**  
  ```bash
  curl -X PUT http://localhost:8000/api/anime/60c72b2f5f1b2c0015e4d3a7 \
    -H "Content-Type: application/json" \
    -d '{"score": 10, "status": "completed"}'
  ```

- **Delete**  
  ```bash
  curl -X DELETE http://localhost:8000/api/anime/60c72b2f5f1b2c0015e4d3a7
  ```

---

## 📸 Screenshots

Home Page  
![Home Page](https://github.com/user-attachments/assets/6399dad4-a54a-4927-ad23-618b4d63f148)

API JSON Output  
![API JSON Output](https://github.com/user-attachments/assets/2c075413-dba0-4a5a-a813-838138547791)

---

## 🤝 Contributing

1. **Fork** the repo  
2. **Create** a feature branch  
   ```bash
   git checkout -b feature/YourFeature
   ```
3. **Commit** your changes  
   ```bash
   git commit -m "Add YourFeature"
   ```
4. **Push** to your branch  
   ```bash
   git push origin feature/YourFeature
   ```
5. **Open** a Pull Request  

Please follow conventional commits and maintain test coverage.

---

## 📝 Future Improvements

- **Bulk Insertion**: `POST /api/addmultipleanimes` to insert arrays of anime.  
- **Authentication**: JWT or OAuth2 support.  
- **Rate Limiting**: Prevent abuse.  
- **Metrics & Logging**: Integrate Prometheus & structured logs.

---

## 🐳 Docker Support

We provide two Docker targets—**development** and **production**—using multi‑stage builds.

### Development

```bash
docker build --target dev -t animeverse-dev .
```
```bash
docker run --rm \
  -p 8000:8000 \       # API
  -p 40000:40000 \     # Delve debugger
  -v "$(pwd)":/app \   # Live reload
  animeverse-dev
```

- **Hot‑Reload:** via [Air](https://github.com/cosmtrek/air)  
- **Debugging:** with [Delve](https://github.com/go-delve/delve)  
- **Full Go Toolchain** inside container  

### Production

```bash
docker build -t animeverse-prod .
```
```bash
docker run --rm \
  --env-file .env \          # Secrets
  -v "$(pwd)/.env:/.env:ro" \# (Optional) read‑only
  -p 8000:8000 \             # API
  animeverse-prod
```

- **Scratch Base:** ~5–10 MB image  
- **Non‑Root User:** Enhanced security  
- **Immutable Binary:** No package manager or shell  

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).  

---

Made with ❤️ by **Flack**
