# 🌸 AnimeVerse API v3.0

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/Flack74/AnimeVerse) [![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE) [![Go Version](https://img.shields.io/badge/go-1.24+-blue)](https://golang.org/) [![Chi Router](https://img.shields.io/badge/router-chi-orange)](https://github.com/go-chi/chi)

Welcome to **AnimeVerse** – your ultimate anime management platform with a modern web interface and powerful RESTful API! Built with **Go**, **MongoDB**, **HTMX**, and **Tailwind CSS**, this full-stack application provides everything you need to manage and explore your anime collection. 🎉

---

## 🚀 Features

### 🎨 **Modern Frontend**
- **HTMX + Tailwind CSS:** AniList-inspired dark theme design
- **Real-time Search:** 500ms debounced search with instant results
- **Dynamic Filters:** Genre, Year, Season, Format, Status dropdowns
- **Episode Tracking:** One-click increment/decrement episode progress
- **Status Management:** Toggle between all watch statuses
- **Responsive Design:** Works perfectly on desktop, tablet, and mobile
- **Loading States:** Visual feedback during all operations
- **Modal Details:** Click any anime for detailed information

### 🔧 **Backend API**
- **RESTful Design:** Clean, intuitive, and standardized endpoints
- **CRUD Operations:** Create, Read, Update, Delete anime records
- **Bulk Operations:** Import multiple anime in single request
- **Advanced Filtering:** Search by name, genre, year, season, format, status
- **Data Import:** Fetch fresh anime from MyAnimeList API
- **Duplicate Prevention:** Smart handling of existing entries

### 🛡️ **Security & Performance**
- **Basic Authentication:** Protected admin endpoints
- **CORS Support:** Cross-origin requests enabled
- **Request Compression:** Automatic gzip compression
- **Timeout Protection:** Request timeout handling
- **Graceful Shutdown:** Proper server shutdown handling
- **Error Handling:** Comprehensive error responses

### 🚀 **DevOps & Deployment**
- **CI/CD Pipeline:** Automated testing and Docker Hub deployment
- **AWS Integration:** One-click EC2 deployment with Terraform
- **Docker Support:** Multi-stage builds for dev and production
- **Hot Reload:** Development mode with Air
- **Environment Config:** Flexible configuration management

---

## 🎨 Modern Frontend

AnimeVerse now features a beautiful, modern web interface built with HTMX + Tailwind CSS:

### 🌐 **Access the Frontend**
- **Main UI:** `http://localhost:8000/` - Modern anime listing interface
- **API Docs:** `http://localhost:8000/api-home` - Original API documentation

### 🔍 **Search & Filter Features**
- **Real-time Search:** Type to search anime instantly
- **Dynamic Filters:** Genre, Year, Season, Format, Status dropdowns
- **Live Updates:** Results update without page refresh
- **Loading Indicators:** Visual feedback during searches

### 📺 **Episode Tracking**
- **Click any anime card** to open detailed modal
- **+ Episode:** Increment watched episodes
- **- Episode:** Decrement watched episodes  
- **Toggle Status:** Cycle through all watch statuses
- **Real-time Updates:** Changes reflect immediately

### 🎨 **UI Sections**
1. **Header:** AnimeVerse branding with import buttons
2. **Search Bar:** Instant search with 500ms delay
3. **Filter Dropdowns:** Genre, Year, Season, Format, Status
4. **Trending Now:** Horizontal scrollable top-rated anime
5. **Popular This Season:** Horizontal scrollable completed anime
6. **All Anime Grid:** Responsive card layout with hover effects
7. **Modal Details:** Click cards for detailed anime information

### 🔐 **Authentication**
- **Public Access:** Browse, search, and view anime details
- **Admin Features:** Episode updates and status changes require login
- **Credentials:** Set in `.env` file (`ADMIN_USERNAME`, `ADMIN_PASSWORD`)

---

## 📥 Quick Start

### **Option 1: Docker Compose (Recommended)**

```bash
# Clone repository
git clone https://github.com/Flack74/AnimeVerse.git
cd AnimeVerse

# Configure environment
cp .env.example .env
# Edit .env with your MongoDB connection and admin credentials

# Run production mode
docker-compose up

# OR run development mode with hot reload
docker-compose --profile dev up animeverse-dev
```

### **Option 2: Local Development**

```bash
# Install dependencies
go mod tidy

# Install Air for hot reloading
go install github.com/air-verse/air@latest

# Run with hot reload
air

# OR run directly
go run main.go
```

**Access Points:**
- **Frontend:** http://localhost:8000 (Modern UI)
- **API:** http://localhost:8000/api-home (Documentation)
- **Health:** http://localhost:8000/health

---

## 🔄 CI/CD Pipeline

AnimeVerse includes automated CI/CD using GitHub Actions with two workflows:

### 📦 **Continuous Integration** (`ci-cd.yml`)
Automatically triggered on every push to main:
- **🧪 Runs Tests:** Executes `go test ./...`
- **🐳 Builds Docker Image:** Multi-stage production build
- **📦 Pushes to Docker Hub:** Deploys as `flack74621/animeverse:latest`
- **⚡ Zero Downtime:** Automated deployment pipeline

### 🚀 **AWS Deployment** (`deploy.yml`)
Manual deployment workflow with environment selection:
- **☁️ Terraform Infrastructure:** Provisions EC2, Security Groups
- **🖥️ EC2 Deployment:** Automated Docker container deployment
- **🌐 Public Access:** Provides deployment URL
- **🎯 Environment Support:** Production/Staging environments

### Setup Requirements

Add these secrets to your GitHub repository:
- `DOCKER_USERNAME` - Your Docker Hub username
- `DOCKER_PASSWORD` - Your Docker Hub password/token
- `AWS_ACCESS_KEY_ID` - Your AWS access key
- `AWS_SECRET_ACCESS_KEY` - Your AWS secret key

### Pull the Latest Image

```bash
docker pull flack74621/animeverse:latest
docker run -p 8000:8000 flack74621/animeverse:latest
```

### 🐳 **Docker Compose Usage**

```bash
# Production mode (port 8000)
docker-compose up

# Development mode with hot reload (port 8001)
docker-compose --profile dev up animeverse-dev

# Stop all services
docker-compose down
```

---

## ☁️ AWS Deployment

Deploy AnimeVerse API to AWS EC2 with one click using Terraform infrastructure as code:

### 🚀 **Quick Deployment**

1. Go to **Actions** tab in your GitHub repository
2. Select **Deploy to AWS** workflow
3. Click **Run workflow**
4. Choose environment (production/staging)
5. Click **Run workflow**

### 🏗️ **Infrastructure Components**

**Terraform Configuration:**
- **Provider:** AWS (us-east-1 region)
- **EC2 Instance:** t2.micro (Free Tier eligible)
- **Security Group:** HTTP (8000) and SSH (22) access
- **Auto-deployment:** Docker container with latest image
- **User Data Script:** Automated Docker installation and app startup

**Deployment Process:**
1. **Terraform Init:** Initialize backend and providers
2. **Terraform Plan:** Preview infrastructure changes
3. **Terraform Apply:** Create AWS resources
4. **Docker Deployment:** Pull and run latest container
5. **Output URLs:** Provide public access endpoint

### 🌐 **Access Your Deployment**

After deployment completes, check the workflow logs for:
```
🚀 Application deployed at: http://YOUR-EC2-IP:8000
```

### 🏛️ **Architecture Overview**

```
🌐 Internet
    │
    ↓
🔒 AWS Security Group
    ├── Port 8000 (HTTP)
    └── Port 22 (SSH)
    │
    ↓
💻 EC2 t2.micro Instance
    ├── Amazon Linux 2023
    ├── Docker Engine
    └── AnimeVerse Container
```

### 🧹 **Resource Cleanup**

To avoid AWS charges, destroy resources when done:

```bash
# Manual cleanup
cd terraform
terraform destroy -auto-approve

# Or use AWS Console to terminate EC2 instance
```

### 🔧 **Troubleshooting**

**Deployment Issues:**
- ✅ Verify AWS credentials in GitHub secrets
- ✅ Check AWS account has EC2/VPC permissions
- ✅ Ensure Docker image exists on Docker Hub
- ✅ Review workflow logs for specific errors

**Application Access Issues:**
- ⏱️ Wait 2-3 minutes for EC2 boot and Docker startup
- 🔒 Verify Security Group allows port 8000
- 🐳 SSH to instance and check: `docker ps`
- 📊 Check application logs: `docker logs animeverse-app`

---

## 🌟 What's New in v3.0

### 🆕 Major Features Added

- **🎨 Modern Frontend:** HTMX + Tailwind CSS with AniList-inspired design
- **🔍 Real-time Search:** Dynamic filtering with live HTMX updates
- **📺 Episode Tracking:** One-click episode increment/decrement
- **🎯 Status Management:** Toggle between all watch statuses
- **📊 Data Import:** Fetch anime from MyAnimeList API
- **📦 Bulk Operations:** Create multiple anime records in one request
- **🖼️ Enhanced Data Model:** Added `bannerUrl`, `imageUrl`, `year`, `season` fields
- **🔄 CI/CD Pipeline:** Automated testing and Docker Hub deployment
- **☁️ AWS Deployment:** One-click EC2 deployment with Terraform
- **🌐 Production Ready:** Public demo deployment capability

## 🌟 What's New in v2.0

### 🔄 **Migration to Chi Router**

- **Upgraded from Gorilla Mux to Chi Router** for better performance and middleware support
- **Improved Request Handling** with Chi's lightweight and fast routing
- **Enhanced Middleware Stack** including CORS, compression, and logging
- **Better Error Handling** with standardized JSON responses
- **Timeout Protection** for all API endpoints

### 🚀 **Performance Improvements**

- **Response Compression:** Automatic gzip compression reduces bandwidth usage
- **Request Logging:** Comprehensive logging for better debugging and monitoring
- **Graceful Shutdown:** Proper server shutdown handling for production environments

---

## 📊 Anime Data Structure

Each anime record follows this enhanced structure:

```json
{
  "_id": "6858f43b802fc0a3285a680e",
  "name": "Attack on Titan",
  "type": "TV",
  "score": 9,
  "progress": {
    "watched": 87,
    "total": 87
  },
  "status": "completed",
  "genre": ["Action", "Drama", "Fantasy"],
  "notes": "Epic story about humanity's fight against titans",
  "year": 2013,
  "season": "Spring",
  "imageUrl": "https://cdn.myanimelist.net/images/anime/10/47347.jpg",
  "bannerUrl": "https://s4.anilist.co/file/anilistcdn/media/anime/banner/16498.jpg"
}
```

### Field Descriptions

| Field | Type | Description |
|-------|------|-------------|
| `_id` | String | Unique MongoDB ObjectId |
| `name` | String | Anime title (required) |
| `type` | String | Format: TV, Movie, OVA, ONA |
| `score` | Number | Personal rating (0-10) |
| `progress.watched` | Number | Episodes watched |
| `progress.total` | Number | Total episodes |
| `status` | String | watching, completed, on-hold, dropped, plan-to-watch |
| `genre` | Array | List of genres |
| `notes` | String | Synopsis or personal notes |
| `year` | Number | Release year |
| `season` | String | Winter, Spring, Summer, Fall |
| `imageUrl` | String | Cover image URL |
| `bannerUrl` | String | Banner image URL |

---

## 🖼️ API Endpoints

### 🌐 **Frontend Routes**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Modern frontend interface |
| `GET` | `/api-home` | API documentation page |
| `GET` | `/health` | Health check endpoint |

### 📖 **Public API Endpoints**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/animes` | Retrieve all anime records |
| `GET` | `/api/animes/filter` | Filter anime by search, genre, year, etc |
| `GET` | `/api/animes/trending` | Get top 5 trending anime (by score) |
| `GET` | `/api/animes/popular` | Get top 5 popular completed anime |
| `GET` | `/api/anime/{name}` | Retrieve specific anime by name |

### 🔒 **Protected Admin Endpoints**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/admin/anime` | Create anime (requires auth) |
| `POST` | `/api/admin/addmultipleanimes` | Bulk create anime (requires auth) |
| `PUT` | `/api/admin/anime/{id}` | Update anime (requires auth) |
| `DELETE` | `/api/admin/anime/{id}` | Delete anime (requires auth) |
| `DELETE` | `/api/admin/deleteallanime` | Delete all anime (requires auth) |
| `POST` | `/api/admin/anime/{id}/episode/increment` | Increment episode count |
| `POST` | `/api/admin/anime/{id}/episode/decrement` | Decrement episode count |
| `POST` | `/api/admin/anime/{id}/status/toggle` | Toggle watch status |
| `POST` | `/api/admin/import/trending` | Import trending anime from MyAnimeList |
| `POST` | `/api/admin/import/seasonal` | Import seasonal anime |

---

## 💡 Usage Examples

### 🎨 **Frontend Usage**

**Getting Fresh Anime Data:**
1. Go to `http://localhost:8000`
2. Click **"Import Trending"** - Gets top 25 anime from MyAnimeList
3. Click **"Import Seasonal"** - Gets current season anime

**Search & Filter:**
- Type in search bar: "attack" → Real-time results
- Select filters: Genre="Action", Year="2023" → Dynamic filtering

**Episode Tracking:**
1. Click any anime card to open modal
2. Click "+ Episode" to increment watched count
3. Click "Toggle Status" to change watch status
4. Changes update instantly via HTMX

### 🔧 **API Usage**

**Create Anime:**
```bash
curl -X POST http://localhost:8000/api/admin/anime \
  -u admin:password \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Hero Academia",
    "type": "TV",
    "score": 9,
    "progress": {"watched": 25, "total": 88},
    "status": "watching",
    "genre": ["Action", "School", "Superhero"],
    "year": 2016,
    "season": "Spring"
  }'
```

**Filter Anime:**
```bash
# Search and filter
GET /api/animes/filter?search=demon&genre=Action&year=2023

# Get trending
GET /api/animes/trending

# Get popular
GET /api/animes/popular
```

**Import Fresh Data:**
```bash
# Import trending anime
curl -X POST http://localhost:8000/api/admin/import/trending \
  -u admin:password

# Import seasonal anime
curl -X POST "http://localhost:8000/api/admin/import/seasonal?year=2024&season=winter" \
  -u admin:password
```

**Bulk Insert:**
```bash
curl -X POST http://localhost:8000/api/admin/addmultipleanimes \
  -u admin:password \
  -H "Content-Type: application/json" \
  -d '[
    {
      "name": "Attack on Titan",
      "type": "TV",
      "score": 9,
      "genre": ["Action", "Drama"],
      "year": 2013
    },
    {
      "name": "Demon Slayer", 
      "type": "TV",
      "score": 8,
      "genre": ["Action", "Supernatural"],
      "year": 2019
    }
  ]'
```

---

## 🛠️ Development

### **Project Structure**
```
AnimeVerse/
├── .github/workflows/  # CI/CD automation
│   ├── ci-cd.yml      # Docker build & push
│   └── deploy.yml     # AWS deployment
├── terraform/         # Infrastructure as Code
│   ├── main.tf       # AWS resources
│   └── variables.tf  # Configuration
├── controllers/      # HTTP handlers
├── models/          # Data structures
├── services/        # Business logic
├── config/          # Database configuration
├── router/          # Route definitions
├── middleware/      # Authentication middleware
├── docker-compose.yml # Multi-environment setup
├── Dockerfile       # Multi-stage build
└── main.go         # Application entry point
```

### **Environment Variables**
```env
# MongoDB Configuration
ConnectionString=mongodb+srv://user:pass@cluster.mongodb.net/
DBName=anime
CollectionName=watchlist

# Server Configuration
PORT=8000

# Admin Authentication
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-secure-password
```

### **Development Commands**
```bash
# Hot reload development
air

# Run tests
go test ./...

# Build for production
go build -o animeverse .

# Docker development
docker-compose --profile dev up

# Docker production
docker-compose up
```

---

## 🤝 Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [Chi Router](https://github.com/go-chi/chi) for the excellent HTTP router
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) for database connectivity
- [HTMX](https://htmx.org/) for seamless frontend interactivity
- [Tailwind CSS](https://tailwindcss.com/) for beautiful styling
- [Jikan API](https://jikan.moe/) for anime data
- [AniList](https://anilist.co/) for design inspiration

---

**Made with ❤️ by Flack. 🚀 Enjoy managing your anime collection with AnimeVerse!**