package types

type DeployInstancePayloadReq struct {
	RpcUrl             string `json:"rpc_url"`
	Password           string `json:"password"`
	GatewayType        string `json:"gateway_type"`
	GatewayName        string `json:"gateway_name"`
	TranscodingProfile string `json:"transcoding_profile"`
}
