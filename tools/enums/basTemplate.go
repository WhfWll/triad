package enums

var BasTemplateEnum basTemplateEnum

type basTemplateEnum struct {
}

const (
	BasTemplateIsDefaultN = 0
	BasTemplateIsDefaultY = 1
)

func (b *basTemplateEnum) AllIsDefault() map[int]string {
	return map[int]string{
		BasTemplateIsDefaultN: "否",
		BasTemplateIsDefaultY: "是",
	}
}
func (b *basTemplateEnum) GetIsDefault(isDefault int) string {
	return b.AllIsDefault()[isDefault]
}
