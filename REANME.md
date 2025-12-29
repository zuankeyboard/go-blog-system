1. 注册接口测试
```bash
curl -X POST http://localhost:8080/api/register \
-H "Content-Type: application/json" \
-d '{
	"username": "testuser",
	"password": "123456",
	"email": "test@example.com"
}'
返回示例：
json
{
  "message": "注册成功",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```
2. 登录接口测试
```bash
curl -X POST http://localhost:8080/api/login \
-H "Content-Type: application/json" \
-d '{
	"username": "testuser",
	"password": "123456"
}'
返回示例：
json
{
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": 1,
    "username": "testuser"
  }
}
{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzY3MjI2MjE5LCJpYXQiOjE3NjY5NjcwMTksImlzcyI6ImJsb2ctYmFja2VuZCJ9.CF4_3w46oadxff4IBsktEci0isOiNjt05gD60fRM-eU","user_id":1,"username":"testuser"},"message":"登录成功"}
```
3. 认证接口测试（需带 Token）
```bash
curl -X GET http://localhost:8080/api/profile \
-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
curl -X GET http://localhost:8080/api/profile \
-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzY3MjI2MjE5LCJpYXQiOjE3NjY5NjcwMTksImlzcyI6ImJsb2ctYmFja2VuZCJ9.CF4_3w46oadxff4IBsktEci0isOiNjt05gD60fRM-eU"
返回示例：
json
{
  "message": "获取个人信息成功",
  "data": {
    "user_id": 1,
    "username": "testuser"
  }
}