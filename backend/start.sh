export $(cat /Users/aditya/Documents/OSS/acidRain/backend/.env | xargs)

docker compose up db -d
go run main.go
