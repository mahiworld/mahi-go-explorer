# Mahi Go Explorer  

**Exploring Go with Gin by building a scalable backend ecosystem featuring API microservices, multi-database support, caching, testing, comprehensive documentation, and more.**  

## Features  
- Built using [Gin Web Framework](https://gin-gonic.com/).  
- Multi-database support (currently MongoDB).  
- Modular architecture for scalability.  
- Includes Makefile for streamlined development, testing, and building.  

## Getting Started  

### Prerequisites  
- Go (version 1.20 or later recommended)  
- MongoDB (running locally or accessible via the provided `DB_URL`)  

### Environment Variables  
Create a `.env` file in the root directory and add the following:  
```env
DB_URL="mongodb://localhost:27017/"
DB_NAME="go-admin"
