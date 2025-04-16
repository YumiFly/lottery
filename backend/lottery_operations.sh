#!/bin/bash

# 设置合约地址和 IP 地址
CONTRACT_ADDRESS="0x70e0bA845a1A0F2DA3359C97E0285013525FFC49"
REMOTE_IP="159.13.40.88"
REMOTE_PORT="58008"

# 设置 JWT 令牌
JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9hZGRyZXNzIjoiMHhBZG1pbkFkZHJlc3MxMjMiLCJleHAiOjE3NDQwNzU4NjMsInJvbGUiOiJhZG1pbiJ9.5RG_Ia5bvlDJvNH6cG2UXWbZmuKLYHp8ziFTx7QHKqo"

# 创建彩票期号 (CreateIssue)
issue_number="20250422"
issue_response=""
while [ -z "$issue_response" ]; do
  issue_response=$(curl -s -X POST http://localhost:8080/lottery/issues \
                   -H "Content-Type: application/json" \
                   -d "{\"lottery_id\":\"6ef1ecde-a58d-4377-933f-34a93760257e\",\"issue_number\":\"$issue_number\",\"sale_end_time\":\"2025-04-16T12:00:00Z\"}")

  # 检查 curl 命令是否成功执行
  if [ $? -ne 0 ]; then
    echo "Error: curl command failed"
    exit 1
  fi

  # 检查响应是否为有效的 JSON
  if [[ "$issue_response" == *"\"issue_id\""* ]]; then
    break # 响应包含 "issue_id"，退出循环
  else
    echo "Waiting for valid response..."
    sleep 2 # 等待 2 秒后重试
  fi
done

# 打印原始响应
echo "Raw response: $issue_response"

# 提取 issue_id
issue_id=$(echo "$issue_response" | jq -r '.data.issue_id')

# 检查 issue_id 是否成功提取
if [ -z "$issue_id" ]; then
  echo "Error: Failed to extract issue_id"
  exit 1
fi

echo "Extracted issue_id: $issue_id"

# 购买彩票 (BuyTicket) - 使用提取的 issue_id
curl -X POST http://localhost:8080/lottery/tickets \
     -H "Content-Type: application/json" \
     -d "{\"issue_id\":\"$issue_id\",\"buyer_address\":\"0x70997970C51812dc3A010C7d01b50e0d17dc79C8\",\"bet_content\":\"6,12,16\",\"purchase_amount\":4.0}"

curl -X POST http://localhost:8080/lottery/tickets \
     -H "Content-Type: application/json" \
     -d "{\"issue_id\":\"$issue_id\",\"buyer_address\":\"0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC\",\"bet_content\":\"6,10,19\",\"purchase_amount\":12.0}"

curl -X POST http://localhost:8080/lottery/tickets \
     -H "Content-Type: application/json" \
     -d "{\"issue_id\":\"$issue_id\",\"buyer_address\":\"0x90F79bf6EB2c4f870365E785982E1f101E93b906\",\"bet_content\":\"2,5,17\",\"purchase_amount\":4.0}"

# 开奖 (DrawLottery) - 使用提取的 issue_id
# 设置合约地址 (SetContractAddress)
# curl -X POST "http://$REMOTE_IP:$REMOTE_PORT/setContractAddress" \
#      -H "Content-Type: application/json" \
#      -d "{\"address\": \"$CONTRACT_ADDRESS\",\"timeout\":-1}"
curl -X POST "http://localhost:8080/lottery/draw?issue_id=$issue_id" \
     -H "Authorization: Bearer $JWT_TOKEN"