package validator

type language struct {
	Required string // 必填
}

func init() {
	lang = language{
		Required: "${notes}为必填",
	}
}
