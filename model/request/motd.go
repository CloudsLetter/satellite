package request

type MotdRequest struct {
	Host string `json:"host,required"`
	Port int    `json:"port,required"`
	Type int    `json:"type,required"`
}
