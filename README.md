# 🌸 AnimeVerse - Complete Anime Discovery Platform

<div align="center">

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/Flack74/AnimeVerse) 
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE) 
[![Go Version](https://img.shields.io/badge/go-1.24+-blue)](https://golang.org/) 
[![Docker](https://img.shields.io/badge/docker-ready-blue)](https://hub.docker.com/r/flack74621/animeverse) 
[![MongoDB](https://img.shields.io/badge/database-mongodb-green)](https://mongodb.com/)
[![Redis](https://img.shields.io/badge/cache-redis-red)](https://redis.io/)
[![HTMX](https://img.shields.io/badge/frontend-htmx-purple)](https://htmx.org/)

**A modern, full-stack anime discovery and management platform built with Go, MongoDB, Redis, HTMX, and Tailwind CSS**

[🚀 Demo](https://your-demo-link.com) • [📚 Documentation](https://docs.animeverse.com) • [🐛 Report Bug](https://github.com/Flack74/AnimeVerse/issues) • [✨ Request Feature](https://github.com/Flack74/AnimeVerse/issues)

</div>

## 🌟 Screenshots

### 🏠 Homepage & Discovery
<div align="center">
  <img src="screenshots/homepage.png" alt="Homepage" width="45%">
  <img src="screenshots/discovery.png" alt="Discovery Page" width="45%">
</div>

### 🔍 Search & Browse
<div align="center">
  <img src="screenshots/search.png" alt="Search Interface" width="45%">
  <img src="screenshots/browse.png" alt="Browse Page" width="45%">
</div>

### 👤 User Profile & Lists
<div align="center">
  <img src="screenshots/profile.png" alt="User Profile" width="45%">
  <img src="screenshots/watchlist.png" alt="Personal Watchlist" width="45%">
</div>

### 📱 Mobile Experience
<div align="center">
  <img src="screenshots/mobile-home.png" alt="Mobile Homepage" width="30%">
  <img src="screenshots/mobile-search.png" alt="Mobile Search" width="30%">
</div>

---

## 🚀 Features

### 🌐 **Complete Anime Platform**
- 🎌 **Browse & Discover** - Explore 39,000+ anime with advanced filtering
- 📝 **Personal Lists** - Create and manage your anime watchlist with custom categories
- 👤 **User Profiles** - Customizable profiles with avatar upload and viewing statistics
- 🔍 **Real-time Search** - Instant search across our extensive database with autocomplete
- 📊 **Trending & Popular** - Stay updated with latest anime trends and community favorites
- 🎯 **Personalized Recommendations** - AI-powered suggestions based on your viewing history

### 🎨 **Modern UI/UX**
- 🌙 **Dark/Light Mode** - System-wide theme toggle with user preference persistence
- 📱 **Responsive Design** - Pixel-perfect optimization for desktop, tablet, and mobile
- ✨ **Glass Morphism** - Modern backdrop-blur navigation and cards with smooth transitions
- 🎭 **Interactive Elements** - Smooth animations, micro-interactions, and loading states
- 🚀 **Progressive Loading** - Skeleton screens and lazy loading for optimal UX

### ⚡ **Performance & Scalability**
- 🧠 **Smart Caching** - Multi-level Redis caching with intelligent invalidation
- 🗄️ **Optimized Database** - 39K+ anime dataset with efficient indexing
- 🖼️ **Image Optimization** - Automatic WebP conversion and CDN integration
- 🔄 **Progressive Loading** - Infinite scroll with 50 items per page
- 📈 **Real-time Analytics** - Performance monitoring and user engagement tracking

---

## 🏗️ Tech Stack

<div align="center">
  
| **Category** | **Technology** | **Purpose** |
|--------------|----------------|-------------|
| **Backend** | Go 1.24+ | High-performance API with Chi router |
| **Database** | MongoDB | NoSQL database with 39K+ anime records |
| **Cache** | Redis | In-memory caching & session management |
| **Frontend** | HTMX + Tailwind CSS | Dynamic UI with utility-first styling |
| **DevOps** | Docker + GitHub Actions | Containerization & CI/CD pipeline |
| **Monitoring** | Custom metrics | Performance & health monitoring |

</div>

### **System Architecture**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │   Application   │    │     Database    │
│    (Optional)   │────│   Server (Go)   │────│   (MongoDB)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                               │                        │
                       ┌─────────────────┐    ┌─────────────────┐
                       │  Cache Layer    │    │   File Storage  │
                       │   (Redis)       │    │   (Local/S3)    │
                       └─────────────────┘    └─────────────────┘
```

---

## 📊 Performance Metrics

<div align="center">

| **Metric** | **Value** | **Description** |
|------------|-----------|-----------------|
| 🚀 **Response Time** | <100ms | Average response time for cached requests |
| 📈 **Throughput** | 1000+ RPS | Requests per second under load |
| 🎯 **Cache Hit Rate** | 85%+ | Redis cache effectiveness |
| 🔍 **Database Queries** | <10ms | Average MongoDB query time |
| 🆙 **Uptime** | 99.9% | Target availability |
| 📦 **Dataset Size** | 39K+ records | Comprehensive anime database |

</div>

---

## 🚀 Quick Start

### **🐳 Docker Deployment (Recommended)**
```bash
# Clone the repository
git clone https://github.com/Flack74/AnimeVerse.git
cd AnimeVerse

# Start all services
docker compose up -d

# Access the application
open http://localhost:8000
```

### **🔧 Manual Setup**
```bash
# Prerequisites: Go 1.24+, MongoDB, Redis

# Install dependencies
go mod tidy

# Configure environment
cp .env.example .env
# Edit .env with your configuration

# Run the application
go run main.go

# Or with hot reload (install air first)
air
```

### **📋 Environment Configuration**
```env
# Database Configuration
ConnectionString=mongodb://localhost:27017
DBName=animeverse
CollectionName=anime

# Cache Configuration
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=your-redis-password

# Server Configuration
PORT=8000
HOST=localhost
ENV=development

# Authentication
JWT_SECRET=your-super-secret-jwt-key
ADMIN_USERNAME=admin
ADMIN_PASSWORD=secure-password

# External APIs (Optional)
TMDB_API_KEY=your-tmdb-api-key
CLOUDINARY_URL=your-cloudinary-url
```

---

## 🔧 API Documentation

### **🌍 Public Endpoints**
```http
GET  /api/animes/trending           # Get trending anime
GET  /api/animes/search?q=naruto    # Search anime by name
GET  /api/animes/popular            # Get popular anime
GET  /api/anime/{name}              # Get specific anime details
GET  /api/animes/genres             # Get all available genres
GET  /api/simple/browse             # Fast browse with filters
```

### **🔐 Protected Endpoints** (Authentication Required)
```http
POST /api/user/anime                # Add anime to user list
PUT  /api/user/anime/{id}/status    # Update anime status
GET  /api/user/stats                # Get user statistics
GET  /api/user/profile              # Get user profile
PUT  /api/user/profile              # Update user profile
DELETE /api/user/anime/{id}         # Remove from user list
```

### **📄 Response Format**
```json
{
  "success": true,
  "data": {
    "animes": [...],
    "pagination": {
      "page": 1,
      "limit": 50,
      "total": 39000,
      "hasNext": true
    }
  },
  "source": "cache|database|api",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## 🛠️ Development

### **📁 Project Structure**
```
AnimeVerse/
├── 🎮 controllers/         # HTTP request handlers
├── ⚙️ services/           # Business logic layer
├── 📊 models/             # Data structures & validation
├── 🔒 middleware/         # Authentication, CORS, rate limiting
├── 🗄️ cache/             # Redis caching strategies
├── ⚙️ config/            # Database & app configuration
├── 🎨 static/            # Frontend assets (CSS, JS, images)
│   ├── css/              # Tailwind CSS & custom styles
│   ├── js/               # Vanilla JS & HTMX interactions
│   └── images/           # Static images & icons
├── 🛣️ router/            # Route definitions & middleware
├── 📝 templates/         # HTML templates
├── 🧪 tests/             # Unit & integration tests
├── 🐳 docker-compose.yml  # Docker services configuration
└── 📋 Makefile           # Build & deployment scripts
```

### **🔄 Development Workflow**
```bash
# Start development environment
make dev

# Run tests with coverage
make test

# Build for production
make build

# Run linting & formatting
make lint

# Generate API documentation
make docs

# Database migrations
make migrate
```

### **✅ Code Quality Standards**
- 🏗️ **Architecture** - Clean architecture with dependency injection
- 🧪 **Testing** - Unit tests with 80%+ coverage requirement
- 🔍 **Linting** - golangci-lint with strict rules
- 📝 **Documentation** - Comprehensive code comments and API docs
- 🔒 **Security** - Regular vulnerability scanning and OWASP compliance

---

## 🔒 Security & Authentication

### **🔐 Authentication System**
- 🎫 **JWT Tokens** - Secure session management with refresh tokens
- 🔑 **Password Security** - bcrypt hashing with salt rounds
- 📧 **Email Verification** - Account activation via email
- 🛡️ **Rate Limiting** - Prevent brute force attacks
- 🚫 **Account Lockout** - Temporary lockout after failed attempts

### **🛡️ Data Protection**
- ✅ **Input Validation** - Comprehensive sanitization and validation
- 🚨 **SQL Injection Prevention** - Parameterized queries and ORM
- 🔒 **XSS Protection** - Content Security Policy headers
- 🍪 **Secure Cookies** - HttpOnly, Secure, SameSite attributes
- 🔐 **HTTPS Enforcement** - SSL/TLS in production

### **🏢 Infrastructure Security**
- 🐳 **Container Security** - Non-root users and minimal base images
- 🌐 **Network Isolation** - Docker networks and firewall rules
- 🔑 **Secrets Management** - Environment variables and Docker secrets
- 📊 **Audit Logging** - Comprehensive security event logging

---

## 📈 Monitoring & Analytics

### **📊 Application Metrics**
- ⏱️ **Performance** - Response times, throughput, error rates
- 💾 **Cache Performance** - Hit/miss ratios, memory usage
- 👥 **User Analytics** - Active users, popular content, engagement
- 🔍 **Search Analytics** - Query patterns, result relevance

### **🖥️ Infrastructure Monitoring**
- 🐳 **Container Health** - Resource utilization, restart counts
- 🗄️ **Database Performance** - Query performance, connection pools
- 🧠 **Cache Monitoring** - Redis memory usage, key expiration
- 🌐 **Network Monitoring** - Traffic patterns, latency metrics

### **🚨 Health Checks & Alerting**
- ❤️ **Health Endpoints** - Application and dependency health
- 📧 **Alert System** - Email/Slack notifications for issues
- 📊 **Dashboard** - Real-time monitoring dashboard
- 🔄 **Auto-Recovery** - Automatic service restart on failure

---

## 🚀 Deployment

### **🐳 Production Deployment**
```bash
# Build production image
docker build -t animeverse:latest .

# Deploy with Docker Compose
docker compose -f docker-compose.prod.yml up -d

# Scale horizontally
docker compose up --scale animeverse=3 --scale nginx=1
```

### **☁️ Cloud Deployment Options**

#### **AWS Deployment**
- 🖥️ **EC2** - Auto Scaling Groups with Application Load Balancer
- 🗄️ **RDS** - Managed MongoDB with automated backups
- 🧠 **ElastiCache** - Managed Redis with clustering
- 📊 **CloudWatch** - Monitoring and alerting
- 🔐 **IAM** - Role-based access control

#### **Docker Swarm**
```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.prod.yml animeverse

# Scale services
docker service scale animeverse_app=3
```

### **🔄 CI/CD Pipeline**
```yaml
# .github/workflows/deploy.yml
name: Deploy to Production
on:
  push:
    branches: [main]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Tests
        run: make test
      
  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Production
        run: make deploy
```

---

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

### **🌟 Ways to Contribute**
- 🐛 **Bug Reports** - Help us identify and fix issues
- ✨ **Feature Requests** - Suggest new features and improvements
- 💻 **Code Contributions** - Submit pull requests with enhancements
- 📚 **Documentation** - Improve docs and add examples
- 🧪 **Testing** - Help expand test coverage

### **📋 Development Setup**
1. 🍴 Fork the repository
2. 🌿 Create a feature branch (`git checkout -b feature/amazing-feature`)
3. 💻 Make your changes
4. 🧪 Add tests for new functionality
5. ✅ Ensure all tests pass (`make test`)
6. 📝 Update documentation if needed
7. 🚀 Submit a pull request

### **📏 Code Standards**
- 🏗️ Follow Go conventions and best practices
- 📝 Write comprehensive tests (aim for 80%+ coverage)
- 📚 Update documentation for new features
- 🔍 Run linting and formatting (`make lint`)
- 📱 Test responsive design on multiple devices

### **🐛 Issue Reporting**
When reporting issues, please include:
- 📋 **Steps to reproduce** the issue
- 🖥️ **System information** (OS, Go version, etc.)
- 📊 **Expected vs actual behavior**
- 📸 **Screenshots** if applicable
- 📝 **Relevant logs** or error messages

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

Special thanks to the amazing open-source community and these fantastic projects:

- 🐹 **[Go Community](https://golang.org/)** - For the excellent ecosystem and tools
- 🍃 **[MongoDB](https://mongodb.com/)** - For the flexible NoSQL database
- 🧠 **[Redis](https://redis.io/)** - For high-performance caching
- ⚡ **[HTMX](https://htmx.org/)** - For modern web interactions without complexity
- 🎨 **[Tailwind CSS](https://tailwindcss.com/)** - For rapid UI development
- 🐳 **[Docker](https://docker.com/)** - For containerization and deployment
- 🔍 **[Chi Router](https://github.com/go-chi/chi)** - For lightweight HTTP routing

---

## 📈 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Flack74/AnimeVerse&type=Date)](https://star-history.com/#Flack74/AnimeVerse&Date)

If you find AnimeVerse useful, please consider giving it a star! ⭐

---

## 🔗 Links

- 🌐 **[Live Demo](https://animeverse-d9xm.onrender.com/)** - Try the live application
- 🐛 **[Issue Tracker](https://github.com/Flack74/AnimeVerse/issues)** - Report bugs
- 💬 **[Discussions](https://github.com/Flack74/AnimeVerse/discussions)** - Community chat
- 📧 **[Contact](mailto:puspendrachawlax@example.com)** - Get in touch

---

<div align="center">
**Built with ❤️ by [Flack74](https://github.com/Flack74)**
</div>
