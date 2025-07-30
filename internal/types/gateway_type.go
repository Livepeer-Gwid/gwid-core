// Package types
package types

type DeployGatewayPayload struct {
	RPCURL             string `json:"rpc_url" binding:"required,url"`
	Password           string `json:"password" binding:"required,min=8"`
	GatewayType        string `json:"gateway_type" binding:"required,oneof=ai transcoding"`
	GatewayName        string `json:"gateway_name" binding:"required,min=3"`
	TranscodingProfile string `json:"transcoding_profile" binding:"required,oneof=480p 720p 1080p"`
	Provider           string `json:"provider" binding:"required,oneof=aws"`
}

type DeployGatewayPayloadReq struct {
	DeployGatewayPayload
	Region string `json:"region" binding:"required"`
}
