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
}

type CreateEC2InstanceReq struct {
	InstanceName      string
	CredentialsID     uuid.UUID `json:"credentials_id" binding:"required,uuid"`
	EC2InstanceTypeID uuid.UUID `json:"ec2_instance_type_id" binding:"required,uuid"`
}

type DeployAWSGatewayPayload struct {
	GatewayID        uuid.UUID
	CredentialsID    uuid.UUID
	InstanceID       string
	UserID           uuid.UUID
	UnhashedPassword string
	Region           string
}
