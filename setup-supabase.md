# Supabase Setup Guide

## 1. Create Supabase Project

1. Go to [https://supabase.com](https://supabase.com)
2. Click "Start your project"
3. Create new project
4. Note down:
   - Project URL: `https://your-project.supabase.co`
   - Anon Key: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
   - JWT Secret: `your-jwt-secret`

## 2. Enable OAuth Providers

1. Go to Authentication > Providers
2. Enable:
   - Google OAuth
   - GitHub OAuth
   - Twitter OAuth (optional)
   - Facebook OAuth (optional)

3. Add redirect URLs:
   - `http://localhost:8000`
   - `https://your-domain.com` (for production)

## 3. Update .env File

```env
# Supabase Authentication
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_JWT_SECRET=your-jwt-secret
ADMIN_EMAIL=your@email.com
```

## 4. Update Frontend Config

In `controllers/controller.go`, update:

```javascript
const supabaseUrl = 'https://your-project.supabase.co';
const supabaseKey = 'your-anon-key';
```

## 5. Test Authentication

1. Start app: `docker-compose up`
2. Go to `http://localhost:8000`
3. Click "Sign In" → Should redirect to Google/GitHub
4. After login → Should return to app with user logged in

## 6. Make User Admin

Connect to MongoDB and run:

```javascript
db.users.updateOne(
  {"email": "your@email.com"}, 
  {"$set": {"role": "admin"}}
)
```

## Features Working:

✅ Google OAuth login
✅ GitHub OAuth login  
✅ JWT token validation
✅ User creation in MongoDB
✅ Admin role management
✅ Personal anime lists
✅ Episode tracking
✅ Smart search with external APIs