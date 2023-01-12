module github.com/quanxiang-cloud/organizations

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/alicebob/miniredis/v2 v2.14.1
	github.com/elliotchance/redismock/v8 v8.11.0
	github.com/gin-gonic/gin v1.7.7
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/zapr v1.2.2
	github.com/go-playground/validator/v10 v10.9.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang/mock v1.1.1
	github.com/olivere/elastic/v7 v7.0.30
	github.com/quanxiang-cloud/cabin v0.0.6
	github.com/quanxiang-cloud/search v0.0.0-20220324022408-21413b3d50fd
	github.com/stretchr/testify v1.7.0
	github.com/tealeg/xlsx v1.0.5
	go.uber.org/zap v1.19.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.2.2
	gorm.io/gorm v1.22.4
)

replace github.com/elliotchance/redismock/v8 v8.11.0 => github.com/vvlgo/redismock/v8 v8.11.2
replace (
	github.com/quanxiang-cloud/cabin v0.0.6 => ../cabin
)
