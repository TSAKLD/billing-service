package entity

type Parameters struct {
	UserID *int `json:"user_id"`
	Offset *int `json:"page"`
	Limit  *int `json:"limit"`
}

func NewParameters() Parameters {
	offset := 0
	limit := 20

	return Parameters{
		UserID: nil,
		Offset: &offset,
		Limit:  &limit,
	}
}
