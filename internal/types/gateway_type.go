// Package types
package types

import "github.com/google/uuid"

type CreateGatewayWithAWSReq struct {
	CredentialsID      uuid.UUID `json:"credentials_id" binding:"required,uuid"`
	EC2InstanceTypeID  uuid.UUID `json:"ec2_instance_type_id" binding:"required,uuid"`
	Region             string    `json:"region" binding:"required"`
	RPCURL             string    `json:"rpc_url" binding:"required,url"`
	Password           string    `json:"password" binding:"required,min=8"`
	GatewayType        string    `json:"gateway_type" binding:"required,oneof=ai transcoding"`
	GatewayName        string    `json:"gateway_name" binding:"required,min=3"`
	TranscodingProfile string    `json:"transcoding_profile" binding:"required,oneof=480p 720p 1080p"`
	Provider           string    `json:"provider" binding:"required,oneof=aws"`
}

type CreateEC2InstanceReq struct {
	CredentialsID     uuid.UUID `json:"credentials_id" binding:"required,uuid"`
	EC2InstanceTypeID uuid.UUID `json:"ec2_instance_type_id" binding:"required,uuid"`
}

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

type DeployAWSGatewayPayload struct {
	GatewayID        uuid.UUID
	UnhashedPassword string
}
