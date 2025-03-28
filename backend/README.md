# Backend KYC API

A RESTful API for KYC user data management built with Go 1.23.2, Gin, and PostgreSQL.

## Prerequisites
- Go 1.23.2
- PostgreSQL (remote Docker)
- Git

## Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd backend

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure PostgreSQL in db/database.go and init_kyc_db.sh.
   Create a .env file in the root directory with the following content
   ```bash
      DB_HOST=<your-db-host>
      DB_PORT=<your-db-port>
      DB_USER=<your-db-user>
      DB_PASSWORD=<your-db-password>
      DB_NAME=<your-db-name>
      DB_SSLMODE=disable
      DB_TIMEZONE=Asia/Shanghai
   ```
4. Initiate the database:
   ```bash
   chmod +x init_kyc_db.sh
   ./init_kyc_db.sh
   ```

5. Run the API server:
   ```bash
   go run main.go
   ```
Project Structure
   main.go: Entry point.
   controllers/: API handlers.
   models/: Database models.
   middleware/: Middleware functions.
   routes/: API route definitions.
   db/: Database connection.
   utils/: Utility functions (e.g., response formatting).
   config/: Configuration loading logic.
   .env: Environment variables for configuration.
   init_kyc_db.sh: Database initialization script.

## API Endpoints
- `POST /customers`: Create a new user.
- `GET /customers`: Retrieve all users.


Response Format
   All responses are in JSON format:
   ```json
      {
      "message": "string",   // Operation result message
      "code": "int",         // HTTP status code (e.g., 200, 400, 401)
      "data": "any"          // Data or error details (null if no data)
      }
   ```

## Error Handling
   Errors are returned with a status code and an error message in the response body:
   ```json
   {
     "status": "error",
     "message": "Error message"
   }
   ```

# in Production
   you can use pathway to deploy the application:
   ```bash
      export DB_HOST=example.com
      export DB_PORT=5432
      export DB_USER=postgres
      export DB_PASSWORD=your_password
      export DB_NAME=kyc_db
      export DB_SSLMODE=disable
      export DB_TIMEZONE=Asia/Shanghai
      go run main.go
   ```
next, you can use the API endpoints to manage KYC user data.

TODO:
   Token 黑名单：实现注销功能，记录失效的 Token。
   日志持久化：将日志写入文件或外部服务（如 ELK）。
   性能优化：缓存角色权限检查结果。
   安全性：增加 HTTPS 支持，使用 JWT 进行身份验证和授权。
   文档：编写 API 文档，使用 Swagger 或其他工具生成在线文档。
   测试：编写单元测试和集成测试，确保代码质量和稳定性。

 区块链生成绑定文件：
   ```bash
   solc --bin KYC.sol -o build
   abigen --bin build/KYC.bin --abi build/KYC.abi --pkg blockchain --type KYC --out blockchain/kyc.go
   ```


 测试用例：
   - 用户注册和登录  
      用户登录：
       ```bash
      `curl -X POST http://localhost:8080/login \
                  -H "Content-Type: application/json" \
                  -d '{
                     "wallet_address": "0xAdminAddress123"
                  }'`
      ```
      用户注册：
      ```bash
      `curl -X POST http://localhost:8080/customers \
                   -H "Content-Type: application/json" \
                   -d '{
                  "customer_address": "0xNewUser789",
                  "is_verified": false,
                  "verifier_address": "",
                  "verification_time": "0001-01-01T00:00:00Z",
                  "registration_time": "2025-03-24T12:00:00Z",
                  "role_id": 0,
                  "assigned_date": "2025-03-24T12:00:00Z",
                  "kyc_data": {
                     "customer_address": "0xNewUser789",
                     "name": "Alice Smith",
                     "birth_date": "1995-08-20T00:00:00Z",
                     "nationality": "UK",
                     "residential_address": "789 Oak St",
                     "phone_number": "5551234567",
                     "email": "alice@example.com",
                     "document_type": "Passport",
                     "document_number": "PP123456789",
                     "file_path": "/path/to/new_passport_image.jpg",
                     "submission_date": "2025-03-24T12:00:00Z",
                     "risk_level": "Low",
                     "source_of_funds": "Savings",
                     "occupation": "Designer"
                  },
                  "kyc_verifications": []
               }'`
      ```
   - 用户信息查询和更新
      用户KYC验证：
      ```bash
        ` curl -X POST http://localhost:8080/auth/verify \
                        -H "Content-Type: application/json" \
                        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9hZGRyZXNzIjoiMHhBZG1pbkFkZHJlc3MxMjMiLCJleHAiOjE3NDMwNTMzMzcsInJvbGUiOiJsb3R0ZXJ5X2FkbWluIn0.xft0eIk1gLGmtwxkXOfmeShQ52wVCqkAolmWl1h9mvU" \
                        -d '{
                           "history_id": 2,
                           "customer_address": "0xNewUser789",
                           "verify_status": "Approved",
                           "verifier_address": "0xAdminAddress123",
                           "verification_date": "2025-03-24T12:00:00Z",
                           "comments": "KYC verification passed"
                        }'`
      ```
      获取用户列表 
      ```bash
         `curl -X GET http://localhost:8080/customers
      ```
      根据ID获取用户信息 
      ```bash   
         `curl -X GET http://localhost:8080/customers/0xAdminAddress123
      ```
      根据钱包地址获取用户信息：
      ```bash
         `curl http://localhost:8080/customers/address/456%20Elm%20St
      ```
   - 角色权限检查
   - Token 黑名单功能
   - 日志持久化功能

  


