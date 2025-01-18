package request

type MotdRequest struct {
	Host      string `uri:"host,required"`
	Port      int    `uri:"port,required"`
	Type      int    `uri:"type,required"`
	Raw       bool   `uri:"raw,required"`
	DesRaw    bool   `uri:"desraw,required"`
	Anonymous bool   `uri:"anonymous,required"`
}
