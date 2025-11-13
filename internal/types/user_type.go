package types

type UpdateProfileReq struct {
	Name string `json:"name" binding:"required,min=2,max=30"`
}
