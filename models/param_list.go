package models

// ParamList :
type ParamList struct {
	Page       int    `json:"page" valid:"Required"`
	PerPage    int    `json:"per_page" valid:"Required"`
	Search     string `json:"search,omitempty"`
	InitSearch string `json:"init_search,omitempty"`
	SortField  string `json:"sort_field,omitempty"`
}

type ParamListGeo struct {
	GeoBarber
	ParamList
}

type ParamListOrder struct {
	BarberId   int    `json:"barber_id"`
	Page       int    `json:"page" valid:"Required"`
	PerPage    int    `json:"per_page" valid:"Required"`
	Search     string `json:"search,omitempty"`
	InitSearch string `json:"init_search,omitempty"`
	SortField  string `json:"sort_field,omitempty"`
}

type ParamDynamicList struct {
	ParamList
	ParamView string `json:"param_view,omitempty"`
}

type PostMulti struct {
	MenuUrl string      `json:"menu_url" valid:"Required"`
	LineNo  int         `json:"line_no,omitempty"`
	InData  interface{} `json:"in_data,omitempty"`
}
