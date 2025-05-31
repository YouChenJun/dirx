package utils

type Result struct {
	Url      string `json:"url"`
	Body     string `json:"body"`
	Code     string `json:"code"`
	Location string `json:"location"`
	Ctype    string `json:"ctype"`
	Server   string `json:"server"`
	Status   string `json:"status"`
	Size     string `json:"size"`
	Time     string `json:"time"`
}
