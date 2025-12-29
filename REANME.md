```
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/dgrijalva/jwt-go
go get golang.org/x/crypto/bcrypt

go get github.com/mattn/go-isatty  # 日志颜色适配
go get github.com/sirupsen/logrus  # 可选：更强大的日志库（推荐）
```

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzY3MjM3NTkzLCJpc3MiOiJnby1ibG9nLXN5c3RlbSJ9.606hvVlfZ__6UUKGyst3kvKPc7bD5hOKa-QNYSD8NRc
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
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzY3MjM3NTkzLCJpc3MiOiJnby1ibG9nLXN5c3RlbSJ9.606hvVlfZ__6UUKGyst3kvKPc7bD5hOKa-QNYSD8NRc",
    "user_id": 1,
    "username": "testuser"
  }
}
```
3. 认证接口测试（需带 Token）
```bash
curl -X GET http://localhost:8080/api/profile \
-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzY3MjM3NTkzLCJpc3MiOiJnby1ibG9nLXN5c3RlbSJ9.606hvVlfZ__6UUKGyst3kvKPc7bD5hOKa-QNYSD8NRc"
返回示例：
json
{
  "message": "获取个人信息成功",
  "data": {
    "user_id": 1,
    "username": "testuser"
  }
}
1. 创建文章（需 Token）
```bash
curl -X POST http://localhost:8080/api/posts \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzY3MjM3NTkzLCJpc3MiOiJnby1ibG9nLXN5c3RlbSJ9.606hvVlfZ__6UUKGyst3kvKPc7bD5hOKa-QNYSD8NRc" \
-H "Content-Type: application/json" \
-d '{
	"title": "我的第一篇博客",
	"content": "使用Gin+GORM开发个人博客系统"
}'
返回示例：
json
{
  "message": "文章创建成功",
  "data": {
    "ID": 1,
    "CreatedAt": "2025-01-01T12:00:00+08:00",
    "UpdatedAt": "2025-01-01T12:00:00+08:00",
    "DeletedAt": null,
    "title": "我的第一篇博客",
    "content": "使用Gin+GORM开发个人博客系统",
    "user_id": 1,
    "user": {
      "ID": 1,
      "CreatedAt": "2025-01-01T11:00:00+08:00",
      "UpdatedAt": "2025-01-01T11:00:00+08:00",
      "DeletedAt": null,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```
2. 获取所有文章
```bash
curl -X GET http://localhost:8080/api/posts
```
3. 更新文章（仅作者）
```bash
curl -X PUT http://localhost:8080/api/posts/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzY3MjM3NTkzLCJpc3MiOiJnby1ibG9nLXN5c3RlbSJ9.606hvVlfZ__6UUKGyst3kvKPc7bD5hOKa-QNYSD8NRc" \
-H "Content-Type: application/json" \
-d '{
	"title": "我的第一篇博客（更新）"
}'
```
4. 删除文章（仅作者）
```bash
curl -X DELETE http://localhost:8080/api/posts/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzY3MjM3NTkzLCJpc3MiOiJnby1ibG9nLXN5c3RlbSJ9.606hvVlfZ__6UUKGyst3kvKPc7bD5hOKa-QNYSD8NRc"
```
1. 发表评论（查询参数方式）
```bash
curl -X POST http://localhost:8080/api/comments?post_id=1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzY3MjM3NTkzLCJpc3MiOiJnby1ibG9nLXN5c3RlbSJ9.606hvVlfZ__6UUKGyst3kvKPc7bD5hOKa-QNYSD8NRc" \
-H "Content-Type: application/json" \
-d '{
	"content": "这篇博客写得很棒！"
}'
```
2. 获取文章评论列表
```bash
curl -X GET http://localhost:8080/api/comments?post_id=1
```