package request

type MotdRequest struct {
	Host string `uri:"host,required"`
	Port int    `uri:"port,required"`
	Type int    `uri:"type,required"`
}
