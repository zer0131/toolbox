module github.com/zer0131/toolbox

go 1.14

require (
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gomodule/redigo v1.8.3
	github.com/gorilla/handlers v1.5.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/jmoiron/sqlx v1.2.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/sirupsen/logrus v1.7.0
	github.com/tebeka/strftime v0.1.5 // indirect
	github.com/zer0131/logfox v1.2.1
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/grpc v1.30.0
	gopkg.in/olivere/elastic.v5 v5.0.86
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
