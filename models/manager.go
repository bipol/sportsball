package models

//Manager is also a user
type Manager struct {
	FullName string `json:"name"`
	ID       int64  `json:"id"`
	TeamID   int64  `json:"team_id"`
}

//AccessKeyResponse returns the access key of the manager to the user
type AccessKeyResponse struct {
	AccessKey string `json:"access_key"`
}
