module go-authentication-boilerplate

go 1.15

replace go-authentication-boilerplate => ./

require (
	github.com/SherClockHolmes/webpush-go v1.3.0
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/bold-commerce/go-shopify/v4 v4.5.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-pg/pg/v10 v10.13.0
	github.com/gofiber/fiber/v2 v2.1.1
	github.com/google/uuid v1.1.2
	github.com/joho/godotenv v1.3.0
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.3.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/resend/resend-go/v2 v2.9.0
	github.com/satori/go.uuid v1.2.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gorm.io/driver/postgres v1.0.5
	gorm.io/gorm v1.20.5
)
