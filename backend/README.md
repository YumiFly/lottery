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

   solc --abi token/LOTToken.sol -o ../../backend/build
   solc --abi --base-path . --include-path ./node_modules/ ./token/LOTToken.sol -o ./token/build
   abigen --abi ./token/build/LOTToken.abi --pkg blockchain --type LOTToken  --out ../backend/blockchain/lottery/lottoken.go


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
      {"message":"Login successful","code":200,"data":{"customer_address":"0xAdminAddress123","role":"lottery_admin","menus":[{"role_menu_id":1,"role_id":1,"menu_name":"lottery_management","menu_path":"/lottery/manage"},{"role_menu_id":2,"role_id":1,"menu_name":"purchase_page","menu_path":"/lottery/purchase"},{"role_menu_id":3,"role_id":1,"menu_name":"account_management","menu_path":"/account"}],"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9hZGRyZXNzIjoiMHhBZG1pbkFkZHJlc3MxMjMiLCJleHAiOjE3NDM0NzMyODYsInJvbGUiOiJsb3R0ZXJ5X2FkbWluIn0.zYAhvvKXLUs_sWSPrNimDBiyCZfebFe0-LdwUtntzNE"}}%  
      ```
      用户注册：
      ```bash
      `curl -X POST http://localhost:8080/customers \
                   -H "Content-Type: application/json" \
                   -d '{
                  "customer_address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
                  "is_verified": false,
                  "verifier_address": "",
                  "verification_time": "0001-01-01T00:00:00Z",
                  "registration_time": "2025-03-24T12:00:00Z",
                  "role_id": 0,
                  "assigned_date": "2025-03-24T12:00:00Z",
                  "kyc_data": {
                     "customer_address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
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
      {"message":"Customer registered successfully","code":200,"data":{"customer_address":"0xNewUser789","is_verified":false,"verifier_address":"","verification_time":"0001-01-01T00:00:00Z","registration_time":"2025-03-24T12:00:00Z","role_id":0,"assigned_date":"2025-03-24T12:00:00Z","kyc_data":{"customer_address":"0xNewUser789","name":"Alice Smith","birth_date":"1995-08-20T00:00:00Z","nationality":"UK","residential_address":"789 Oak St","phone_number":"5551234567","email":"alice@example.com","document_type":"Passport","document_number":"PP123456789","file_path":"/path/to/new_passport_image.jpg","submission_date":"2025-03-24T12:00:00Z","risk_level":"Low","source_of_funds":"Savings","occupation":"Designer"},"kyc_verifications":[],"role":{"role_id":0,"role_name":"","role_type":"","description":"","menus":null}}}%  
      ```
   - 用户信息查询和更新
      用户KYC验证：
      ```bash
        ` curl -X POST http://localhost:8080/auth/verify \
                        -H "Content-Type: application/json" \
                        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9hZGRyZXNzIjoiMHhBZG1pbkFkZHJlc3MxMjMiLCJleHAiOjE3NDM3MzQwNTIsInJvbGUiOiJhZG1pbiJ9.z4UGie8yrJ5_-C-Y7E7ymybfOY04JRQ7IOzWrI3E5l8" \
                        -d '{
                           "history_id": 2,
                           "customer_address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
                           "verify_status": "Approved",
                           "verifier_address": "0xAdminAddress123",
                           "verification_date": "2025-03-24T12:00:00Z",
                           "comments": "KYC verification passed"
                        }'`
         {"message":"Verification successful","code":200,"data":null}%
      ```
      获取用户列表 
      ```bash
         `curl -X GET http://localhost:8080/customers
         {"message":"Customers retrieved successfully","code":200,"data":[{"customer_address":"0xAdminAddress123","is_verified":true,"verifier_address":"","verification_time":"0001-01-01T00:00:00Z","registration_time":"2025-03-31T02:06:04.37407Z","role_id":1,"assigned_date":"2025-03-31T02:06:04.37407Z","kyc_data":{"customer_address":"0xAdminAddress123","name":"Jane Doe","birth_date":"1985-05-15T00:00:00Z","nationality":"CN","residential_address":"456 Elm St","phone_number":"9876543210","email":"jane@example.com","document_type":"ID","document_number":"ID123456789","file_path":"/path/to/id_image.jpg","submission_date":"2025-03-31T02:06:04.37407Z","risk_level":"Low","source_of_funds":"Salary","occupation":"Engineer"},"kyc_verifications":[{"history_id":1,"customer_address":"0xAdminAddress123","verify_status":"Approved","verifier_address":"0xVerifierAddress789","verification_date":"2025-03-31T02:06:04.37407Z","comments":"KYC verification completed successfully"}],"role":{"role_id":1,"role_name":"lottery_admin","role_type":"admin","description":"Administrator for lottery management","menus":[{"role_menu_id":1,"role_id":1,"menu_name":"lottery_management","menu_path":"/lottery/manage"},{"role_menu_id":2,"role_id":1,"menu_name":"purchase_page","menu_path":"/lottery/purchase"},{"role_menu_id":3,"role_id":1,"menu_name":"account_management","menu_path":"/account"}]}},{"customer_address":"0xUserAddress456","is_verified":false,"verifier_address":"","verification_time":"0001-01-01T00:00:00Z","registration_time":"2025-03-31T02:06:04.37407Z","role_id":2,"assigned_date":"2025-03-31T02:06:04.37407Z","kyc_data":{"customer_address":"0xUserAddress456","name":"John Doe","birth_date":"1990-01-01T00:00:00Z","nationality":"US","residential_address":"123 Main St","phone_number":"1234567890","email":"john@example.com","document_type":"Passport","document_number":"PP987654321","file_path":"/path/to/passport_image.jpg","submission_date":"2025-03-31T02:06:04.37407Z","risk_level":"Medium","source_of_funds":"Investment","occupation":"Trader"},"kyc_verifications":[],"role":{"role_id":2,"role_name":"normal_user","role_type":"user","description":"Normal user with limited access","menus":[{"role_menu_id":4,"role_id":2,"menu_name":"purchase_page","menu_path":"/lottery/purchase"},{"role_menu_id":5,"role_id":2,"menu_name":"account_management","menu_path":"/account"}]}},{"customer_address":"0xNewUser789","is_verified":true,"verifier_address":"0xAdminAddress123","verification_time":"2025-03-24T12:00:00Z","registration_time":"2025-03-24T12:00:00Z","role_id":2,"assigned_date":"2025-03-31T10:09:57.359221Z","kyc_data":{"customer_address":"0xNewUser789","name":"Alice Smith","birth_date":"1995-08-20T00:00:00Z","nationality":"UK","residential_address":"789 Oak St","phone_number":"5551234567","email":"alice@example.com","document_type":"Passport","document_number":"PP123456789","file_path":"/path/to/new_passport_image.jpg","submission_date":"2025-03-24T12:00:00Z","risk_level":"Low","source_of_funds":"Savings","occupation":"Designer"},"kyc_verifications":[{"history_id":2,"customer_address":"0xNewUser789","verify_status":"Approved","verifier_address":"0xAdminAddress123","verification_date":"2025-03-24T12:00:00Z","comments":"KYC verification passed"}],"role":{"role_id":2,"role_name":"normal_user","role_type":"user","description":"Normal user with limited access","menus":[{"role_menu_id":4,"role_id":2,"menu_name":"purchase_page","menu_path":"/lottery/purchase"},{"role_menu_id":5,"role_id":2,"menu_name":"account_management","menu_path":"/account"}]}}]}%  
      ```
      根据ID获取用户信息 
      ```bash   
         `curl -X GET http://localhost:8080/customers/0xNewUser789
          {"message":"Customer retrieved successfully","code":200,"data":{"customer_address":"0xNewUser789","is_verified":true,"verifier_address":"0xAdminAddress123","verification_time":"2025-03-24T12:00:00Z","registration_time":"2025-03-24T12:00:00Z","role_id":2,"assigned_date":"2025-03-31T10:09:57.359221Z","kyc_data":{"customer_address":"0xNewUser789","name":"Alice Smith","birth_date":"1995-08-20T00:00:00Z","nationality":"UK","residential_address":"789 Oak St","phone_number":"5551234567","email":"alice@example.com","document_type":"Passport","document_number":"PP123456789","file_path":"/path/to/new_passport_image.jpg","submission_date":"2025-03-24T12:00:00Z","risk_level":"Low","source_of_funds":"Savings","occupation":"Designer"},"kyc_verifications":[{"history_id":2,"customer_address":"0xNewUser789","verify_status":"Approved","verifier_address":"0xAdminAddress123","verification_date":"2025-03-24T12:00:00Z","comments":"KYC verification passed"}],"role":{"role_id":2,"role_name":"normal_user","role_type":"user","description":"Normal user with limited access","menus":[{"role_menu_id":4,"role_id":2,"menu_name":"purchase_page","menu_path":"/lottery/purchase"},{"role_menu_id":5,"role_id":2,"menu_name":"account_management","menu_path":"/account"}]}}}%    

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
         {"message":"Lottery type created successfully","code":200,"data":{"type_id":"899c4312-89b6-421e-bf13-32f9a4804268","type_name":"乐透型","description":"A simple lottery type","created_at":"2025-03-31T10:14:35.825968+08:00","updated_at":"2025-03-31T10:14:35.825968+08:00"}}%  
      ```
      获取彩票类型列表
      ```bash
         `curl -X GET http://localhost:8080/lottery/types
         {"message":"Lottery types retrieved successfully","code":200,"data":[{"type_id":"3e3ff670-9201-4f17-9ff2-972a785cb40f","type_name":"简单型","description":"A simple lottery type","created_at":"2025-03-31T10:12:02.458191Z","updated_at":"2025-03-31T10:12:02.458192Z"},{"type_id":"dacb6f6f-e799-487a-ab02-8cbdf130d7b0","type_name":"即开型","description":"A simple lottery type","created_at":"2025-03-31T10:14:18.210076Z","updated_at":"2025-03-31T10:14:18.210076Z"},{"type_id":"899c4312-89b6-421e-bf13-32f9a4804268","type_name":"乐透型","description":"A simple lottery type","created_at":"2025-03-31T10:14:35.825968Z","updated_at":"2025-03-31T10:14:35.825968Z"},{"type_id":"90936ac0-24bc-4171-8680-09135a488b0f","type_name":"基诺型","description":"A simple lottery type","created_at":"2025-03-31T10:14:53.349044Z","updated_at":"2025-03-31T10:14:53.349044Z"}]}% 
      ```
      创建彩票
      ```bash
         `curl -X POST http://localhost:8080/lottery/lottery \
               -H "Content-Type: application/json" \
               -d '{"type_id":"1","ticket_name":"乐可可","ticket_price":2.0,"ticket_supply":10,
               "betting_rules":"Choose 3 numbers between 1 and 36","prize_structure":"1st Prize: 50% of pool",
               "registered_addr":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
               "rollout_contract_address":"0x0E801D84Fa97b50751Dbf25036d067dCf18858bF"}'`
         {"message":"Lottery created successfully","code":200,"data":{"lottery_id":"a0ccdbef-0f74-4096-b69e-012e882a7f65","type_id":"3e3ff670-9201-4f17-9ff2-972a785cb40f","ticket_name":"SimpleTicket","ticket_price":"0.1","betting_rules":"Choose 3 numbers between 1 and 36","prize_structure":"1st Prize: 50% of pool","registered_addr":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","contract_address":"0x1234567890abcdef1234567890abcdef12345678","created_at":"2025-03-31T10:21:06.091355+08:00","updated_at":"2025-03-31T10:21:06.091355+08:00"}}% 
      ```
      获取彩票列表
      ```bash
         `curl -X GET http://localhost:8080/lottery/lottery
         {"message":"Lotteries retrieved successfully","code":200,"data":[{"lottery_id":"a0ccdbef-0f74-4096-b69e-012e882a7f65","type_id":"3e3ff670-9201-4f17-9ff2-972a785cb40f","ticket_name":"SimpleTicket","ticket_price":"0.1","betting_rules":"Choose 3 numbers between 1 and 36","prize_structure":"1st Prize: 50% of pool","contract_address":"0x1234567890abcdef1234567890abcdef12345678","created_at":"2025-03-31T10:21:06.091355Z","updated_at":"2025-03-31T10:21:06.091355Z"}]}% 
      ```
      
      创建彩票期号 (CreateIssue)
      ```bash
         `curl -X POST http://localhost:8080/lottery/issues \
               -H "Content-Type: application/json" \
                -d '{"lottery_id":"3e4ea146-ca83-4dff-83ac-601f3548218e","issue_number":"20250409","sale_end_time":"2025-04-09T12:00:00Z"}'
         {"message":"Issue created successfully","code":200,"data":{"issue_id":"c5b8edda-6d38-4b85-ac3f-9e5ba84d5848","lottery_id":"a0ccdbef-0f74-4096-b69e-012e882a7f65","issue_number":"20250405","sale_end_time":"2025-04-05T12:00:00Z","draw_time":"0001-01-01T00:00:00Z","prize_pool":"","winning_numbers":"","random_seed":"","draw_tx_hash":"","created_at":"2025-03-31T10:29:20.406934+08:00","updated_at":"2025-03-31T10:29:20.406935+08:00"}}%
      ```
     
      获取根据彩票ID获取最近的发行信息
      ```bash
         `curl -X GET http://localhost:8080/lottery/issues/latest/aa1c1fc9-3316-4d65-876b-5b83390097f1
         {"message":"Latest issue retrieved successfully","code":200,"data":{"issue_id":"86db3427-f58d-4047-a280-574852df374f","lottery_id":"a0ccdbef-0f74-4096-b69e-012e882a7f65","issue_number":"20250331","sale_end_time":"2025-03-28T12:00:00Z","draw_time":"2025-03-27T12:00:00Z","prize_pool":"100","winning_numbers":"","random_seed":"","draw_tx_hash":"","created_at":"2025-03-31T10:25:33.964414Z","updated_at":"2025-03-31T10:25:33.964415Z"}}%  
      ```

      购买彩票 (BuyTicket)
      ```bash
         `curl -X POST http://localhost:8080/lottery/tickets \
               -H "Content-Type: application/json" \
               -d '{"ticket_id":"TK001","issue_id":"24f8445c-d831-4447-a2b8-e37eff8437c3","buyer_address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","bet_content":"6,12,16","purchase_amount":4.0}'
         {"message":"Ticket purchased successfully","code":200,"data":{"ticket_id":"7caf92e5-b6d2-4a81-9753-06843d4113ed","issue_id":"c5b8edda-6d38-4b85-ac3f-9e5ba84d5848","buyer_address":"0xabcdef1234567890abcdef1234567890abcdef12","purchase_time":"2025-03-31T10:30:12.724715+08:00","bet_content":"6,11,16","purchase_amount":"0.1","transaction_hash":"","claim_tx_hash":"","created_at":"2025-03-31T10:30:12.724715+08:00","updated_at":"2025-03-31T10:30:12.724716+08:00"}}

      curl -X POST http://localhost:8080/lottery/tickets \
               -H "Content-Type: application/json" \
               -d '{"ticket_id":"TK001","issue_id":"24f8445c-d831-4447-a2b8-e37eff8437c3","buyer_address":"0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC","bet_content":"6,10,19","purchase_amount":12.0}'
      curl -X POST http://localhost:8080/lottery/tickets \
               -H "Content-Type: application/json" \
               -d '{"ticket_id":"TK001","issue_id":"24f8445c-d831-4447-a2b8-e37eff8437c3","buyer_address":"0x90F79bf6EB2c4f870365E785982E1f101E93b906","bet_content":"2,5,17","purchase_amount":6.0}'
      ```
      获取用户购买的彩票列表 (GetUserTickets)
      ```bash
         `curl -X GET http://localhost:8080/lottery/tickets/customer/0xabcdef1234567890abcdef1234567890abcdef12
         {"message":"Purchased tickets retrieved successfully","code":200,"data":[{"ticket_id":"7caf92e5-b6d2-4a81-9753-06843d4113ed","issue_id":"c5b8edda-6d38-4b85-ac3f-9e5ba84d5848","buyer_address":"0xabcdef1234567890abcdef1234567890abcdef12","purchase_time":"2025-03-31T10:30:12.724715Z","bet_content":"6,11,16","purchase_amount":"0.1","transaction_hash":"","claim_tx_hash":"","created_at":"2025-03-31T10:30:12.724715Z","updated_at":"2025-03-31T10:30:12.724716Z"}]}%
      ```
      
      获取总奖池
      ```bash
         curl -X GET "http://localhost:8080/lottery/pools"
         {"message":"Pools retrieved successfully","code":200,"data":0}%
      ```

      开奖
      ```bash
         `curl -X POST "http://localhost:8080/lottery/draw?issue_id=24f8445c-d831-4447-a2b8-e37eff8437c3" \
               -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9hZGRyZXNzIjoiMHhBZG1pbkFkZHJlc3MxMjMiLCJleHAiOjE3NDQwNzU4NjMsInJvbGUiOiJhZG1pbiJ9.5RG_Ia5bvlDJvNH6cG2UXWbZmuKLYHp8ziFTx7QHKqo"
      ```
      
      获取开奖结果
      ```bash
         `curl -X GET "http://localhost:8080/lottery/draw?issue_id=IS001"
      ```
      获取最近开奖结果以及彩票信息
      ```bash
         `curl -X GET "http://localhost:8080/lottery/draw/latest"
      ```

   - 角色权限检查
   - Token 黑名单功能
   - 日志持久化功能

   --需要提供一个获取彩票剩余量的合约接口
   curl -X POST http://127.0.0.1:8888/setContractAddress -H "Content-Type: application/json" -d '{"address": "0xbDA5747bFD65F08deb54cb465eB87D40e51B197E","timeout":-1}'

   curl -X POST http://159.13.40.88:58008/setContractAddress -H "Content-Type: application/json" -d '{"address": "0x70e0bA845a1A0F2DA3359C97E0285013525FFC49","timeout":-1}'
  


