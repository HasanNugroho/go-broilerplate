package roles

type RoleModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Permissions string `json:"permission"`
}

type RoleUpdateModel struct {
	Name        string `json:"name"`
	Permissions string `json:"permission"`
}
