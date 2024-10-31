# Try it out

```
# Run databaase + cache
docker compose up

# Run backend
cd backend
go mod tidy
go run main.go

# Run frontend
cd frontend
yarn install
yarn start
```

# TO do
- [ ] Send "campaign" notification on the fly
- [ ] SMS abandoned cart automation
- [ ] WhatsApp abandoned cart automation
- [ ] Email abandoned cart automation
- [ ] Web push notification abandoned cart automation

# Thoughts

- How about making a "shop" the tenant of the app?
