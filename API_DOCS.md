# 📚 AnimeVerse API Documentation

## Base URL
```
http://localhost:8000
```

## Response Format
All API responses follow a consistent JSON structure:

```json
{
  "success": boolean,
  "message": "string",
  "data": object|array|null,
  "error": "string" // only present when success is false
}
```

## Endpoints

### 🏠 Home & Health

#### GET `/`
Returns the homepage with API information.

**Response:** HTML page

#### GET `/health`
Health check endpoint for monitoring.

**Response:**
```json
{
  "success": true,
  "message": "API is healthy",
  "data": {
    "status": "healthy",
    "version": "2.0",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### 🎌 Anime Endpoints

#### GET `/api/animes`
Retrieve all anime records.

**Response:**
```json
{
  "success": true,
  "message": "Animes retrieved successfully",
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "name": "Attack on Titan",
      "type": "TV",
      "score": 10,
      "progress": {
        "watched": 87,
        "total": 87
      },
      "status": "completed",
      "genre": ["action", "drama", "fantasy"],
      "notes": "Masterpiece anime with incredible storytelling"
    }
  ]
}
```

#### GET `/api/anime/{animeName}`
Retrieve a specific anime by name. Supports URL-friendly formats:
- Spaces can be replaced with hyphens: `attack-on-titan`
- Spaces can be replaced with underscores: `attack_on_titan`

**Parameters:**
- `animeName` (string): The name of the anime (case-insensitive)

**Response:**
```json
{
  "success": true,
  "message": "Anime retrieved successfully",
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "name": "Attack on Titan",
    "type": "TV",
    "score": 10,
    "progress": {
      "watched": 87,
      "total": 87
    },
    "status": "completed",
    "genre": ["action", "drama", "fantasy"],
    "notes": "Masterpiece anime"
  }
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "Anime not found"
}
```

#### POST `/api/anime`
Create a new anime record.

**Request Body:**
```json
{
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
```

**Required Fields:**
- `name` (string): Anime name (1-200 characters)

**Optional Fields:**
- `type` (string): "TV", "Movie", or "ONA"
- `score` (integer): Rating from 0-10
- `progress.watched` (integer): Episodes watched
- `progress.total` (integer): Total episodes
- `status` (string): "watching", "completed", "on-hold", "dropped", "plan-to-watch"
- `genre` (array): Array of genre strings
- `notes` (string): Personal notes (max 500 characters)

**Response (201):**
```json
{
  "success": true,
  "message": "Anime created successfully",
  "data": {
    "_id": "507f1f77bcf86cd799439012",
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

**Error Response (409 - Conflict):**
```json
{
  "success": false,
  "error": "Anime with this name already exists"
}
```

#### POST `/api/addmultipleanimes`
Create multiple anime records in a single request (bulk insertion).

**Request Body:**
```json
[
  {
    "name": "Attack on Titan",
    "type": "TV",
    "score": 9,
    "genre": ["Action", "Drama", "Fantasy"],
    "status": "completed",
    "bannerUrl": "https://example.com/banner.jpg",
    "imageUrl": "https://example.com/cover.jpg"
  },
  {
    "name": "Demon Slayer",
    "type": "TV",
    "score": 8,
    "genre": ["Action", "Supernatural"],
    "status": "watching"
  }
]
```

**Features:**
- Automatically skips duplicate anime (by name)
- Uses MongoDB's InsertMany for optimal performance
- Returns detailed response with inserted IDs and duplicates

**Response (201 - All inserted):**
```json
{
  "success": true,
  "message": "All animes created successfully",
  "data": {
    "inserted_count": 2,
    "inserted_ids": ["507f1f77bcf86cd799439012", "507f1f77bcf86cd799439013"]
  }
}
```

**Response (206 - Partial success with duplicates):**
```json
{
  "success": true,
  "message": "Animes inserted with some duplicates skipped",
  "data": {
    "inserted_count": 1,
    "inserted_ids": ["507f1f77bcf86cd799439012"],
    "duplicates": ["Attack on Titan"]
  }
}
```

**Error Response (400):**
```json
{
  "success": false,
  "error": "No animes provided"
}
```

#### PUT `/api/anime/{id}`
Update an existing anime record (partial update supported).

**Parameters:**
- `id` (string): MongoDB ObjectID of the anime

**Request Body (partial update example):**
```json
{
  "score": 10,
  "status": "completed",
  "progress": {
    "watched": 88,
    "total": 88
  },
  "notes": "Finished watching - absolutely incredible!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Anime updated successfully",
  "data": {
    "matched": 1,
    "modified": 1,
    "id": "507f1f77bcf86cd799439011"
  }
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "Anime not found"
}
```

#### DELETE `/api/anime/{id}`
Delete a specific anime record.

**Parameters:**
- `id` (string): MongoDB ObjectID of the anime

**Response:**
```json
{
  "success": true,
  "message": "Anime deleted successfully",
  "data": {
    "deleted_id": "507f1f77bcf86cd799439011"
  }
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "Anime not found"
}
```

#### DELETE `/api/deleteallanime`
Delete all anime records. **Use with caution!**

**Response:**
```json
{
  "success": true,
  "message": "All animes deleted successfully",
  "data": {
    "deleted_count": 15
  }
}
```

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request data |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 500 | Internal Server Error - Server error |

## Data Types

### AnimeType
- `"TV"` - Television series
- `"Movie"` - Anime movie
- `"ONA"` - Original Net Animation

### WatchStatus
- `"watching"` - Currently watching
- `"completed"` - Finished watching
- `"on-hold"` - Temporarily paused
- `"dropped"` - Stopped watching
- `"plan-to-watch"` - Planning to watch

## Examples

### cURL Examples

**Get all anime:**
```bash
curl -X GET http://localhost:8000/api/animes
```

**Create new anime:**
```bash
curl -X POST http://localhost:8000/api/anime \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Demon Slayer",
    "type": "TV",
    "score": 9,
    "status": "watching",
    "genre": ["action", "supernatural"]
  }'
```

**Create multiple anime (bulk):**
```bash
curl -X POST http://localhost:8000/api/addmultipleanimes \
  -H "Content-Type: application/json" \
  -d '[
    {
      "name": "One Piece",
      "type": "TV",
      "score": 9,
      "genre": ["adventure", "comedy"]
    },
    {
      "name": "Naruto",
      "type": "TV",
      "score": 8,
      "genre": ["action", "ninja"]
    }
  ]'
```

**Update anime:**
```bash
curl -X PUT http://localhost:8000/api/anime/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{
    "score": 10,
    "status": "completed"
  }'
```

**Delete anime:**
```bash
curl -X DELETE http://localhost:8000/api/anime/507f1f77bcf86cd799439011
```

### JavaScript Fetch Examples

**Get all anime:**
```javascript
fetch('http://localhost:8000/api/animes')
  .then(response => response.json())
  .then(data => {
    if (data.success) {
      console.log('Animes:', data.data);
    } else {
      console.error('Error:', data.error);
    }
  });
```

**Create new anime:**
```javascript
fetch('http://localhost:8000/api/anime', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: 'One Piece',
    type: 'TV',
    score: 9,
    status: 'watching',
    genre: ['adventure', 'comedy', 'shounen']
  })
})
.then(response => response.json())
.then(data => {
  if (data.success) {
    console.log('Created:', data.data);
  } else {
    console.error('Error:', data.error);
  }
});
```

**Create multiple anime (bulk):**
```javascript
fetch('http://localhost:8000/api/addmultipleanimes', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify([
    {
      name: 'Jujutsu Kaisen',
      type: 'TV',
      score: 9,
      genre: ['action', 'supernatural']
    },
    {
      name: 'Chainsaw Man',
      type: 'TV',
      score: 8,
      genre: ['action', 'horror']
    }
  ])
})
.then(response => response.json())
.then(data => {
  if (data.success) {
    console.log(`Inserted ${data.data.inserted_count} animes`);
    if (data.data.duplicates) {
      console.log('Duplicates skipped:', data.data.duplicates);
    }
  } else {
    console.error('Error:', data.error);
  }
});
```

## Rate Limiting & CORS

- **CORS:** Enabled for all origins (`*`)
- **Timeout:** 60 seconds per request
- **Compression:** Gzip compression enabled
- **Rate Limiting:** Not implemented (consider adding for production)

## Development

The API includes comprehensive logging and error handling. All requests are logged with Chi's built-in logger middleware.

For development with hot reload:
```bash
air
```

For production:
```bash
go run main.go
```