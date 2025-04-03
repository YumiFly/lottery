#!/bin/bash



# 检查是否存在 .env 文件，用于加载环境变量

if [ -f ".env" ]; then

set -a

source .env

set +a

else

echo "Error: .env file not found."

exit 1

fi



# 从环境变量中获取 PostgreSQL 连接信息

PG_HOST=${DB_HOST}

PG_PORT=${DB_PORT}

PG_USER=${DB_USER}

PG_PASSWORD=${DB_PASSWORD}

PG_DB=${DB_NAME}



# 验证是否缺少必要的环境变量

if [ -z "$PG_HOST" ] || [ -z "$PG_PORT" ] || [ -z "$PG_USER" ] || [ -z "$PG_PASSWORD" ] || [ -z "$PG_DB" ]; then

echo "Error: Missing required environment variables (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)."

exit 1

fi



# 设置 PostgreSQL 密码环境变量，以便 psql 命令使用

export PGPASSWORD=$PG_PASSWORD



# 测试与 PostgreSQL 数据库的连接

echo "Testing connection to remote PostgreSQL at $PG_HOST:$PG_PORT..."

psql -h $PG_HOST -p $PG_PORT -U $PG_USER -d postgres -c "SELECT 1;" > /dev/null 2>&1

if [ $? -ne 0 ]; then

echo "Error: Cannot connect to $PG_HOST:$PG_PORT. Check host, port, username, password, or network settings."

exit 1

fi

echo "Connection successful."

# 检查目标数据库是否存在
echo "Checking if database $PG_DB exists..."
DB_EXISTS=$(psql -h $PG_HOST -p $PG_PORT -U $PG_USER -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$PG_DB'")
if [ -z "$DB_EXISTS" ]; then
    echo "Database $PG_DB does not exist. Creating it..."
    psql -h $PG_HOST -p $PG_PORT -U $PG_USER -d postgres -c "CREATE DATABASE $PG_DB;"
    if [ $? -eq 0 ]; then
        echo "Database $PG_DB created successfully."
    else
        echo "Failed to create database $PG_DB."
        exit 1
    fi
else
    echo "Database $PG_DB already exists."
fi

# 定义 SQL 语句：删除现有表
SQL_DROP_TABLES="
DROP TABLE IF EXISTS kyc_verification_histories CASCADE;
DROP TABLE IF EXISTS kyc_data CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS role_menus CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS lotteries CASCADE;
DROP TABLE IF EXISTS lottery_types CASCADE;
DROP TABLE IF EXISTS lottery_issues CASCADE;
DROP TABLE IF EXISTS lottery_tickets CASCADE;
DROP TABLE IF EXISTS winners CASCADE;

"

# 定义 SQL 语句：创建表（不设置主键和外键约束）
SQL_CREATE_TABLES="
-- 角色表
CREATE TABLE roles (
    role_id INTEGER,
    role_name VARCHAR(50) NOT NULL,
    role_type VARCHAR(50),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 角色菜单表
CREATE TABLE role_menus (
    role_menu_id INTEGER,
    role_id INTEGER NOT NULL,
    menu_name VARCHAR(50) NOT NULL,
    menu_path VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 用户表
CREATE TABLE customers (
    customer_address VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    verifier_address VARCHAR(255),
    verification_time TIMESTAMP WITH TIME ZONE,
    registration_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    role_id INTEGER NOT NULL,
    assigned_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- KYC 数据表
CREATE TABLE kyc_data (
    customer_address VARCHAR(255),
    name VARCHAR(100),
    birth_date DATE,
    nationality VARCHAR(50),
    residential_address TEXT,
    phone_number VARCHAR(20),
    email VARCHAR(255),
    document_type VARCHAR(50),
    document_number VARCHAR(50),
    file_path TEXT,
    submission_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    risk_level VARCHAR(20),
    source_of_funds TEXT,
    occupation VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- KYC 验证历史表
CREATE TABLE kyc_verification_histories (
    history_id INTEGER,
    customer_address VARCHAR(255) NOT NULL,
    verify_status VARCHAR(50),
    verifier_address VARCHAR(255),
    verification_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    comments TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建 lottery_types 表
CREATE TABLE lottery_types (
    type_id VARCHAR(50),
    type_name VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建 lottery_issues 表
CREATE TABLE lottery_issues (
    issue_id VARCHAR(50),
    lottery_id VARCHAR(50) NOT NULL,
    issue_number VARCHAR(50) NOT NULL,
    sale_end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    draw_time TIMESTAMP WITH TIME ZONE NOT NULL,
    prize_pool NUMERIC NOT NULL,
    winning_numbers VARCHAR(100),
    random_seed VARCHAR(100),
    draw_tx_hash VARCHAR(66),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 创建 lottery 表
CREATE TABLE lotteries (
    lottery_id VARCHAR(50),
    type_id VARCHAR(50) NOT NULL,
    ticket_name VARCHAR(255) NOT NULL,
    ticket_price NUMERIC NOT NULL,
    ticket_supply NUMERIC NOT NULL,
    betting_rules VARCHAR(1000) NOT NULL,
    prize_structure VARCHAR(1000) NOT NULL,
    registered_addr VARCHAR(255) NOT NULL,
    rollout_contract_address VARCHAR(255) NOT NULL,
    contract_address VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 创建 lottery_tickets 表
CREATE TABLE lottery_tickets (
    ticket_id VARCHAR(50),
    issue_id VARCHAR(50) NOT NULL,
    buyer_address VARCHAR(66) NOT NULL,
    purchase_time TIMESTAMP WITH TIME ZONE NOT NULL,
    bet_content VARCHAR(100) NOT NULL,
    purchase_amount NUMERIC NOT NULL,
    transaction_hash VARCHAR(66),
   created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 创建 winners 表
CREATE TABLE winners (
    winner_id VARCHAR(50),
    issue_id VARCHAR(50) NOT NULL,
    ticket_id VARCHAR(50) NOT NULL,
    address VARCHAR(66) NOT NULL,
    prize_level VARCHAR(50) NOT NULL,
    prize_amount NUMERIC NOT NULL,
    claim_tx_hash VARCHAR(66), 
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

"

# 定义 SQL 语句：插入初始数据
SQL_SEED_DATA="
-- 初始化角色
INSERT INTO roles (role_id, role_name, role_type, description) VALUES
    (1, 'admin', 'admin', 'Administrator for all lottery management,can manage all lottery'),
    (2, 'normal_user', 'user', 'Normal user with limited access'),
    (3, 'lottery_admin', 'lottery_admin', 'Administrator for only one lottery management,can not manage other lottery');

-- 初始化菜单
INSERT INTO role_menus (role_menu_id, role_id, menu_name, menu_path) VALUES
    (1, 1, 'lottery_management', '/lottery/manage'),
    (2, 1, 'purchase_page', '/lottery/purchase'),
    (3, 1, 'account_management', '/account'),
    (4, 2, 'purchase_page', '/lottery/purchase'),
    (5, 2, 'account_management', '/account');

-- 初始化用户
INSERT INTO customers (customer_address, is_verified, role_id, registration_time, assigned_date) VALUES
    ('0xAdminAddress123', TRUE, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('0xUserAddress456', FALSE, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- 初始化 KYC 数据
INSERT INTO kyc_data (customer_address, name, birth_date, nationality, residential_address, phone_number, email, document_type, document_number, file_path, risk_level, source_of_funds, occupation) VALUES
    ('0xAdminAddress123', 'Jane Doe', '1985-05-15', 'CN', '456 Elm St', '9876543210', 'jane@example.com', 'ID', 'ID123456789', '/path/to/id_image.jpg', 'Low', 'Salary', 'Engineer'),
    ('0xUserAddress456', 'John Doe', '1990-01-01', 'US', '123 Main St', '1234567890', 'john@example.com', 'Passport', 'PP987654321', '/path/to/passport_image.jpg', 'Medium', 'Investment', 'Trader');

-- 初始化 KYC 验证历史
INSERT INTO kyc_verification_histories (history_id, customer_address, verify_status, verifier_address, verification_date, comments) VALUES
    (1, '0xAdminAddress123', 'Approved', '0xVerifierAddress789', CURRENT_TIMESTAMP, 'KYC verification completed successfully');

-- 初始化彩票类型
INSERT INTO lottery_types (type_id, type_name, description, created_at, updated_at) VALUES
    (1, '数字型', '数字型', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, '乐透型', '乐透型', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, '基诺型', '基诺型', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (4, '福彩', '福彩', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (5, '体彩', '体彩', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
"

# 删除现有表
echo "Dropping existing tables..."
psql -h $PG_HOST -p $PG_PORT -U $PG_USER -d $PG_DB -c "$SQL_DROP_TABLES"
if [ $? -eq 0 ]; then
    echo "Tables dropped successfully."
else
    echo "Failed to drop tables."
    exit 1
fi

# 创建新表
echo "Creating tables..."
psql -h $PG_HOST -p $PG_PORT -U $PG_USER -d $PG_DB -c "$SQL_CREATE_TABLES"
if [ $? -eq 0 ]; then
    echo "Tables created successfully."
else
    echo "Failed to create tables."
    exit 1
fi

# 插入初始数据
echo "Seeding initial data..."
psql -h $PG_HOST -p $PG_PORT -U $PG_USER -d $PG_DB -c "$SQL_SEED_DATA"
if [ $? -eq 0 ]; then
    echo "Initial data seeded successfully."
else
    echo "Failed to seed initial data."
    exit 1
fi

# 完成提示
echo "Database initialization on remote Docker completed."