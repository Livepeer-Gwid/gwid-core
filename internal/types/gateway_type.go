package types

type DeployGatewayPayload struct {
	RPCURL             string `json:"rpc_url" binding:"required"`
	Password           string `json:"password" binding:"required,min=8"`
	GatewayType        string `json:"gateway_type" binding:"required"`
	GatewayName        string `json:"gateway_name" binding:"required"`
	TranscodingProfile string `json:"transcoding_profile" binding:"required"`
	Provider           string `json:"provider" binding:"required"`
}

type DeployGatewayPayloadReq struct {
	DeployGatewayPayload
	Region string `json:"region" binding:"required"`
}
