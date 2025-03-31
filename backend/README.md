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

   solc --abi sample_rollout.sol -o ../../backend/build
   abigen --abi build/SimpleRollout.abi --pkg blockchain --type SimpleRollout --out blockchain/simple_rollout.go

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
                        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9hZGRyZXNzIjoiMHhBZG1pbkFkZHJlc3MxMjMiLCJleHAiOjE3NDM0NzMyODYsInJvbGUiOiJsb3R0ZXJ5X2FkbWluIn0.zYAhvvKXLUs_sWSPrNimDBiyCZfebFe0-LdwUtntzNE" \
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
         `curl -X GET http://localhost:8080/customers/0xNewUser789
      ```
      根据钱包地址获取用户信息：
      ```bash
         `curl http://localhost:8080/customers/address/456%20Elm%20St
      ```
      创建彩票类型
      ```bash
         `curl -X POST http://localhost:8080/lottery/types \
               -H "Content-Type: application/json" \
               -d '{"type_name":"简单型","description":"A simple lottery type"}'`
      ```
      获取彩票类型列表
      ```bash
         `curl -X GET http://localhost:8080/lottery/types
      ```
      创建彩票
      ```bash
         `curl -X POST http://localhost:8080/lottery/lottery \
               -H "Content-Type: application/json" \
               -d '{"type_id":"3e3ff670-9201-4f17-9ff2-972a785cb40f","ticket_name":"SimpleTicket","ticket_price":"0.1","betting_rules":"Choose 3 numbers between 1 and 36","prize_structure":"1st Prize: 50% of pool","contract_address":"0x1234567890abcdef1234567890abcdef12345678"}'`
      ```
      获取彩票列表
      ```bash
         `curl -X GET http://localhost:8080/lottery/lottery
      ```
      
      创建彩票期号 (CreateIssue)
      ```bash
         `curl -X POST http://localhost:8080/lottery/issues \
               -H "Content-Type: application/json" \
                -d '{"lottery_id":"a0ccdbef-0f74-4096-b69e-012e882a7f65","issue_number":"20250405","sale_end_time":"2025-04-05T12:00:00Z"}'
      ```
     
      获取根据彩票ID获取最近的发行信息
      ```bash
         `curl -X GET http://localhost:8080/lottery/issues/latest/a0ccdbef-0f74-4096-b69e-012e882a7f65
      ```

      购买彩票 (BuyTicket)
      ```bash
         `curl -X POST http://localhost:8080/lottery/tickets \
               -H "Content-Type: application/json" \
               -d '{"ticket_id":"TK001","issue_id":"c5b8edda-6d38-4b85-ac3f-9e5ba84d5848","buyer_address":"0xabcdef1234567890abcdef1234567890abcdef12","bet_content":"6,11,16","purchase_amount":"0.1"}'
      ```
      获取用户购买的彩票列表 (GetUserTickets)
      ```bash
         `curl -X GET http://localhost:8080/lottery/tickets/customer/0xabcdef1234567890abcdef1234567890abcdef12
      ```
      
      获取总奖池
      ```bash
         curl -X GET "http://localhost:8080/lottery/pools"
      ```

      开奖
      ```bash
         `curl -X POST "http://localhost:8080/lottery/draw?issue_id=IS001" \
               -H "Authorization: Bearer <admin_token>"
      ```
      
      获取开奖结果
      ```bash
         `curl -X GET "http://localhost:8080/lottery/draw?issue_id=IS001"
      ```
      

   - 角色权限检查
   - Token 黑名单功能
   - 日志持久化功能

  


