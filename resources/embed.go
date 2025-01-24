package resources

import _ "embed"

//go:embed ip2region.xdb
var ip2region []byte

func GetIp2Region() []byte {
	return ip2region
}
