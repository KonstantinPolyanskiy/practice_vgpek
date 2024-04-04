package permissions

type GetObjectResp struct {
	Id   int    `json:"object_id"`
	Name string `json:"object_name"`
}

type GetObjectsResp struct {
	Objects []GetObjectResp
}

type ObjectEntity struct {
	Id   int    `db:"internal_object_id"`
	Name string `db:"internal_object_name"`
}

type ObjectDTO struct {
	Name string
}

type AddObjectReq struct {
	Name string `json:"name"`
}

type AddObjectResp struct {
	Name string `json:"name"`
}
