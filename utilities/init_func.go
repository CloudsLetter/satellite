package utilities

import (
	"runtime"
	"satellite/middlewares"
	"satellite/utilities/database"
	"satellite/utilities/misc"
	"sync"
	"time"
)

var once sync.Once

func InitFunc() {
	once.Do(func() {
		initDatabase()
		misc.InitIp2Region()
		misc.NewGoroutinePool(10, 500, 1*time.Minute)
		misc.InitACAutomaton()
		initCheckCurrentLimiting()
	})
}
func initDatabase() {
	if runtime.GOOS == "windows" {
		database.InitDB("host=127.0.0.1 user=postgres password=88323095.zxc dbname=satellitetmp port=5432 sslmode=disable TimeZone=Asia/Shanghai")
		database.InitRedisDB0("127.0.0.1:6379")

	} else {
		database.InitDB("host=postgres user=postgres password=88323095.zxc dbname=satellite port=5432 sslmode=disable TimeZone=Asia/Shanghai")
		database.InitRedisDB0("redis:6379")
	}
}
func initCheckCurrentLimiting() {
	for i := 0; i < 3; i++ {
		err := misc.GetGoroutinePool().Submit(middlewares.WriteDataToDB)
		if err != nil {
			return
		}
	}
	err := misc.GetGoroutinePool().Submit(middlewares.CheckCurrentLimiting)
	if err != nil {
		return
	}
}
