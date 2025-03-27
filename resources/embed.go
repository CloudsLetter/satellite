package resources

import _ "embed"

//go:embed ip2region.xdb
var ip2region []byte

//go:embed al.txt
var sensitive_word string

func GetSensitiveWord() string {
	return sensitive_word
}

func GetIp2Region() []byte {
	return ip2region
}
