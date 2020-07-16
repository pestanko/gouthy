package web_utils

func NewDeepLink(path string, params interface{}) DeepLink {
	return DeepLink{
		Link:   path,
		Params: params,
	}
}

func DecodeDeepLink(data string) (DeepLink, error) {
	dl := DeepLink{}
	err := JsonDecode(data, &dl)
	return dl, err
}

type DeepLink struct {
	Link   string      `json:"l"`
	Params interface{} `json:"p"`
}

func (dl *DeepLink) Encode() (string, error) {
	return JsonEncode(dl)
}
