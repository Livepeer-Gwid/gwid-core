package types

type AWSCredentialsReq struct {
	AccessKeyID     string `json:"access_key_id" binding:"required,min=16,max=128"`
	SecretAccessKey string `json:"secret_access_key" binding:"required,min=16"`
}

type AWSCredentailsProfile struct {
	ProfileName string
	ProfileARN  string
	RoleName    string
	RoleARN     string
}
