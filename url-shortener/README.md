# Go URL Shortener 🔗

A secure, user-based URL shortener built in Go with MySQL backend, user authentication, customizable slugs, click tracking, and a modern frontend interface.

Built as a full-stack learning project combining Go’s `net/http`, session management, dynamic templates, and clean UI — ready for local or containerized deployment.

---

## ✅ Features

- 🔐 **User Authentication** (signup, login, logout)
- 🔗 **Create short URLs** with optional custom slugs
- 🚫 **Slug conflict prevention** (global uniqueness)
- 📊 **Click tracking** (view clicks per short link)
- 🧾 **Dashboard UI** (see and delete your own links)
- 🖼️ **Modern responsive frontend** using pure HTML/CSS
- 🧠 **Session middleware** for route protection
- 🚫 **Cache-control** to block back-button access after logout
- 🧼 **Error & success messages** (inline and styled)
- 📦 **Modular Go structure**
- 💾 **MySQL persistent storage** (no more `data.json`)

---

## 🖼️ Screenshots

> Login / Signup / Dashboard view  
![screenshot](preview.png)

> Custom slug and short link result  
![demo](demo.gif)

---

## 🧑‍💻 Tech Stack

- Language: **Go 1.21+**
- DB: **MySQL**
- Session Management: `gorilla/sessions`
- Password Hashing: `bcrypt`
- Env Loader: `joho/godotenv`
- Templating: `html/template`

---

## 🚀 How to Run Locally

### 🧰 Prerequisites

- Go 1.18+
- MySQL (local or Docker)
- Git

### ⚙️ 1. Clone the Project

```bash
git clone https://github.com/YOUR_USERNAME/go-url-shortener.git
cd go-url-shortener

⚙️ 2. Create the MySQL Database

CREATE DATABASE url_shortener;

-- Run this to create tables
USE url_shortener;

CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE short_links (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  original_url TEXT NOT NULL,
  short_slug VARCHAR(255) UNIQUE NOT NULL,
  click_count INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

⚙️ 3. Setup Environment

Create a config.env file:
DB_URL=root:your_mysql_password@tcp(127.0.0.1:3306)/url_shortener
SECRET=your_session_secret

⚙️ 4. Run the Project

go run main.go routes.go

Visit:
http://localhost:8080

📁 Project Structure
url-shortener/
├── main.go                # Server setup
├── routes.go              # Route registration
├── db/                    # DB connection
│   └── db.go
├── handler/               # Handlers for login, signup, dashboard, links
│   ├── auth.go
│   ├── link.go
│   └── dashboard.go
├── middleware/            # Session middleware
├── model/                 # DB models
├── static/                # CSS / JS (optional)
├── templates/             # HTML templates (login, signup, dashboard)
├── utils/                 # Slug generator, password utils
├── go.mod
├── config.env


🌍 Access Behavior
/signup → Register new user

/login → Login to your dashboard

/dashboard → View + manage your links

/logout → Ends session and clears access

/{slug} → Redirect to original URL (and increment click count)

🔒 Security Notes
Sessions are cookie-based (secure + SameSite=Lax)

Back button after logout is blocked with JS + cache headers

Passwords are hashed with bcrypt

SQL uses prepared statements (? placeholders)

💡 Future Improvements
📋 Copy to clipboard buttons

🌓 Light/dark theme toggle

📈 Charts for click activity

🧪 Unit tests

🐳 Docker + Docker Compose

🌍 Deployment on Fly.io or Railway

