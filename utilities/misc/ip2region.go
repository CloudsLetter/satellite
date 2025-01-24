package misc

import (
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"satellite/resources"
)

// var ip2regionBuf
var ip2regionSearcher *xdb.Searcher

func InitIp2Region() {
	var err error
	ip2regionSearcher, err = xdb.NewWithBuffer(resources.GetIp2Region())
	if err != nil {
		panic(err)
	}
}
func GetIp2RegionSearcher() *xdb.Searcher {
	return ip2regionSearcher
}
