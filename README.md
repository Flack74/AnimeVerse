# 🎬 AnimeVerse

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/Flack74/AnimeApi) [![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

Welcome to **AnimeVerse** – your one-stop RESTful API for managing and exploring your favorite anime collection! Built with **Go**, **MongoDB**, and **Gorilla Mux**, this API lets you create, read, update, and delete anime records with ease. Whether you're a casual fan or a hardcore otaku, Anime API has got you covered! 🎉

---

## 🚀 Features

- **CRUD Operations:** Create, Read, Update, and Delete anime records effortlessly.
- **Partial Updates:** Send a JSON payload with only the fields you need to update.
- **MongoDB Integration:** Secure and scalable storage with MongoDB.
- **RESTful Design:** Clean and intuitive endpoints.
- **Detailed Data:** Manage anime with fields like name, type, score, progress, status, and genre.

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

   ```bash
   go run main.go
   ```

   The API will be available at [http://localhost:8000](http://localhost:8000) 🎉

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
          "tags": ["action", "shounen", "superhero"]
        }'
  ```

- **Get All Anime:**

  ```bash
  curl http://localhost:8000/api/animes
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
          "status": "completed"
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

## 🛠 API Endpoints

| Method | Endpoint                  | Description                              |
| ------ | ------------------------- | ---------------------------------------- |
| GET    | `/api/animes`             | Retrieve all anime records               |
| POST   | `/api/anime`              | Create a new anime record                |
| PUT    | `/api/anime/{id}`         | Update an anime record (partial update)  |
| DELETE | `/api/anime/{id}`         | Delete a specific anime record           |
| DELETE | `/api/deleteallanime`     | Delete all anime records                 |

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

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

Made with ❤️ by Flack 
