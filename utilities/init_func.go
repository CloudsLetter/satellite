package utilities

import (
	"github.com/getsentry/sentry-go"
	"runtime"
	"satellite/middlewares"
	"satellite/utilities/database"
	"satellite/utilities/misc"
	"sync"
)

var once sync.Once

func InitFunc() {
	once.Do(func() {
		initDatabase()
		initCheckCurrentLimiting()
		misc.InitIp2Region()
		initSentry()
	})
}
func initDatabase() {
	if runtime.GOOS == "windows" {
		database.InitDB("host=127.0.0.1 user=postgres password=88323095.zxc dbname=satellitetmp port=5432 sslmode=disable TimeZone=Asia/Shanghai")
	} else {
		database.InitDB("host=postgres user=postgres password=88323095.zxc dbname=satellite port=5432 sslmode=disable TimeZone=Asia/Shanghai")
	}
	database.AutoMigrate(database.GetDBContext())
}
func initCheckCurrentLimiting() {
	for i := 0; i < 3; i++ {
		go middlewares.WriteDataToDB()
	}
	go middlewares.CheckCurrentLimiting()
}
func initSentry() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           "https://55268fcd7867979562654ef2f4660340@o926347.ingest.us.sentry.io/4508697946423296",
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		panic(err)
	}
}
