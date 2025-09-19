# Chirpy
A simple Twitter-like microblogging platform built as part of the Boot.dev curriculum.  
Users can create accounts, post short messages ("chirps"), and view others' posts.

---

## Table of Contents
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Setup](#setup)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Tests](#tests)
- [License](#license)

---

## Features
- User registration and authentication  
- Create, read, and delete chirps  
- Simple REST API with JSON responses  
- PostgreSQL support  

---

## Tech Stack
- **Go** (main programming language)  
- **PostgreSQL**  
- **JSON REST API**  

---

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/chirpy.git
   cd chirpy
2. Install dependencies 
    go mod tidy
3. Set Up your database 
    Createdb chirpy
4. Run the server
    go run main.go

## Usage

Start the server and send requerts using curl or Postman
Example: Create a new chirp

curl -X POST http://localhost:8080/chirps -d '{"body":"Hello world!"}'

## API Endpoints 

* POST /api/users – Create a new user
* POST /api/login – Authenticate a user
* POST /api/chirps – Create a chirp
* GET /api/chirps – List all chirps
* GET /api/chirps/{id} – Get a single chirp
* DELETE /api/chirps/{id} – Delete a chirp

## Tests
go test ./..

## Licence

This project is part of the Boot.dev curriculum. Please refer to Boot.dev for licensing information