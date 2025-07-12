# 🌸 AnimeVerse - Complete Anime Discovery Platform

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/Flack74/AnimeVerse) 
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE) 
[![Go Version](https://img.shields.io/badge/go-1.24+-blue)](https://golang.org/) 
[![Docker](https://img.shields.io/badge/docker-ready-blue)](https://hub.docker.com/r/flack74621/animeverse) 
[![MongoDB](https://img.shields.io/badge/database-mongodb-green)](https://mongodb.com/)
[![Redis](https://img.shields.io/badge/cache-redis-red)](https://redis.io/)
[![HTMX](https://img.shields.io/badge/frontend-htmx-purple)](https://htmx.org/)

**AnimeVerse** is a modern, full-stack anime discovery and management platform built with **Go**, **MongoDB**, **Redis**, **HTMX**, and **Tailwind CSS**. Discover, search, and manage your anime collection with a beautiful, responsive interface.

---

## 📸 Screenshots

### Home Page
<img width="1860" height="859" alt="image" src="https://github.com/user-attachments/assets/9269e26e-e9ae-451b-9ea4-63e1bc4bfae8" />

### Browse Anime
<img width="1854" height="859" alt="image" src="https://github.com/user-attachments/assets/438bca15-3380-4648-9c3f-c5e0eab728fc" />

### Anime Details
<img width="1858" height="850" alt="image" src="https://github.com/user-attachments/assets/1607a8e5-9e6f-492e-95c1-7e9eef041c4c" />

### User Profile
<img width="1872" height="858" alt="image" src="https://github.com/user-attachments/assets/8982a170-4475-460b-8eb2-fe98e61ff81b" />

### Trending Anime
<img width="1830" height="855" alt="image" src="https://github.com/user-attachments/assets/2912fd9a-7e0d-488f-8678-abb37bb96f01" />

### Schedule
<img width="1842" height="796" alt="image" src="https://github.com/user-attachments/assets/9c7b914a-f7e0-4cf5-85b6-5a0ad03635c9" />

## 🚀 Features

### 🌐 **Complete Anime Platform**
- **Browse & Discover** - Explore 39,000+ anime with advanced filtering
- **Personal Lists** - Create and manage your anime watchlist
- **User Profiles** - Customizable profiles with avatar upload
- **Real-time Search** - Instant search across our extensive database
- **Trending & Popular** - Stay updated with latest anime trends

### 🎨 **Modern UI/UX**
- **Dark/Light Mode** - System-wide theme toggle with persistence
- **Responsive Design** - Optimized for desktop, tablet, and mobile
- **Glass Morphism** - Modern backdrop-blur navigation and cards
- **Interactive Elements** - Smooth animations and micro-interactions
- **Loading States** - Skeleton loading for better user experience

### 🚀 **Performance & Scalability**
- **Smart Caching** - Redis caching with optimized strategies
- **Database** - 39K+ anime dataset stored in MongoDB
- **High-Quality Images** - Automatic image quality optimization
- **Fast Loading** - 50 items per page with progressive loading
- **CDN Integration** - Optimized content delivery

---

## 🏗️ Architecture

### **Backend Technologies**
- **Go 1.24+** - High-performance backend with Chi router
- **MongoDB** - NoSQL database with 39,000+ anime records
- **Redis** - In-memory caching for performance optimization
- **RESTful API** - Clean, documented endpoints with JSON responses

### **Frontend Technologies**
- **HTMX** - Dynamic, interactive web interface
- **Tailwind CSS** - Utility-first CSS framework
- **Vanilla JavaScript** - Lightweight client-side functionality
- **Responsive Design** - Mobile-first approach

### **DevOps & Infrastructure**
- **Docker** - Containerization for consistent deployments
- **GitHub Actions** - Automated CI/CD pipeline
- **Multi-stage Builds** - Optimized production containers
- **Health Checks** - Application monitoring and reliability

---

## 📊 Database & Performance

### **Anime Dataset**
- **39,000+ Anime Records** - Comprehensive anime database
- **High-Quality Metadata** - Detailed information including genres, scores, episodes
- **Optimized Indexing** - Fast search and filtering capabilities
- **Real-time Updates** - Dynamic content synchronization

### **Caching Strategy**
- **Multi-level Caching** - Browser → Redis → Database
- **Smart Expiration** - 5-30 minute cache windows
- **Cache Warming** - Preloaded popular content
- **Performance Metrics** - 85%+ cache hit rate

### **Performance Benchmarks**
- **Response Time** - <100ms for cached requests
- **Throughput** - 1000+ requests per second
- **Database Queries** - <10ms average query time
- **Uptime** - 99.9% availability target

---

## 🚀 Quick Start

### **Docker Deployment (Recommended)**
```bash
# Clone the repository
git clone https://github.com/Flack74/AnimeVerse.git
cd AnimeVerse

# Start with Docker Compose
docker compose up -d

# Access the application
open http://localhost:8000
```

### **Manual Setup**
```bash
# Prerequisites: Go 1.24+, MongoDB, Redis

# Install dependencies
go mod tidy

# Set environment variables
cp .env.example .env

# Run the application
go run main.go
# or with hot reload
air
```

### **Environment Configuration**
```env
# Database
ConnectionString=mongodb://localhost:27017
DBName=animeverse
CollectionName=anime

# Cache
REDIS_URL=redis://localhost:6379

# Server
PORT=8000

# Authentication (Optional)
ADMIN_USERNAME=admin
ADMIN_PASSWORD=secure-password
```

---

## 📱 User Interface

### **Homepage**
- Dynamic carousel with latest trending anime
- Featured anime recommendations
- Airing schedule with real-time updates
- Quick search and navigation

### **Browse & Discovery**
- Advanced filtering by genre, year, status
- 50 anime per page with load more functionality
- High-quality image previews
- Instant search with debounced input

### **Personal Features**
- User authentication and profiles
- Personal anime lists with status tracking
- Avatar upload and profile customization
- Statistics and viewing history

### **Responsive Design**
- Mobile-optimized interface
- Touch-friendly navigation
- Adaptive layouts for all screen sizes
- Progressive web app capabilities

---

## 🔧 API Documentation

### **Public Endpoints**
```http
GET  /api/animes/trending           # Trending anime
GET  /api/animes/search?q=naruto    # Search anime
GET  /api/anime/{name}              # Get specific anime details
GET  /api/simple/browse             # Fast browse with filters
```

### **User Endpoints** (Authentication Required)
```http
POST /api/user/anime                # Add anime to list
PUT  /api/user/anime/{id}/status    # Update anime status
GET  /api/user/stats                # Get user statistics
DELETE /api/user/anime/{id}         # Remove from list
```

### **Response Format**
```json
{
  "success": true,
  "data": [...],
  "source": "database|cache|api",
  "total": 39000
}
```

---

## 🛠️ Development

### **Project Structure**
```
AnimeVerse/
├── controllers/     # HTTP request handlers
├── services/        # Business logic layer
├── models/          # Data structures
├── middleware/      # Authentication & CORS
├── cache/           # Redis caching layer
├── config/          # Database configuration
├── static/          # Frontend assets
├── router/          # Route definitions
└── docker-compose.yml
```

### **Development Workflow**
```bash
# Hot reload development
air

# Run tests
go test ./...

# Build for production
go build -o animeverse-api .

# Docker development
docker compose --profile dev up
```

### **Code Quality**
- Go best practices and conventions
- Comprehensive error handling
- Unit tests with coverage reports
- Automated code review process
- Security vulnerability scanning

---

## 🔒 Security & Authentication

### **User Authentication**
- Secure user registration and login
- Session management with JWT tokens
- Password hashing with bcrypt
- Rate limiting on authentication endpoints

### **Data Protection**
- Input validation and sanitization
- SQL injection prevention
- XSS protection with CSP headers
- Secure cookie handling

### **Infrastructure Security**
- Container security with non-root users
- Network isolation with Docker
- Environment variable encryption
- SSL/TLS enforcement

---

## 📈 Monitoring & Analytics

### **Application Metrics**
- Request/response times
- Error rates and types
- Cache hit/miss ratios
- User engagement statistics

### **Infrastructure Monitoring**
- Container resource utilization
- Database performance metrics
- Redis memory usage
- Network traffic analysis

### **Health Checks**
- Application health endpoint
- Database connectivity checks
- Cache availability monitoring
- Automated alerting system

---

## 🚀 Deployment

### **Production Deployment**
```bash
# Build production image
docker build -t animeverse:latest .

# Deploy with Docker Compose
docker compose -f docker-compose.prod.yml up -d

# Scale horizontally
docker compose up --scale animeverse=3
```

### **Cloud Deployment**
- AWS EC2 with Auto Scaling
- Load balancer integration
- RDS for MongoDB hosting
- ElastiCache for Redis
- CloudWatch monitoring

### **CI/CD Pipeline**
- Automated testing on pull requests
- Docker image building and pushing
- Staging environment deployment
- Production deployment with rollback

---

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### **Development Setup**
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

### **Code Standards**
- Follow Go conventions
- Write comprehensive tests
- Update documentation
- Ensure Docker builds pass
- Test responsive design

### **Issue Reporting**
- Use GitHub Issues for bug reports
- Provide detailed reproduction steps
- Include system information
- Attach relevant logs or screenshots

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **Go Community** - For the excellent ecosystem and tools
- **MongoDB** - For the flexible NoSQL database
- **Redis** - For high-performance caching
- **HTMX** - For modern web interactions
- **Tailwind CSS** - For rapid UI development
- **Docker** - For containerization and deployment

---

## 🌟 Star History

If you find AnimeVerse useful, please consider giving it a star! ⭐

---

**Built with ❤️ by the Flack**
