package markup

import (
	"github.com/yossoy/exciton/internal/markup"
)

func Style(name string, value string) MarkupOrChild {
	return markup.StyleApplyer{
		Name:  name,
		Value: value,
	}
}

func Attribute(name string, value interface{}) MarkupOrChild {
	//if name is start "data-", change name and return dataApplyer.
	if ds, ok := markup.AttrToDataSet(name); ok {
		return Data(ds, value.(string))
	}

	return markup.AttrApplyer{
		Name:  name,
		Value: value,
	}
}

func Classes(classes ...string) MarkupOrChild {
	return markup.ClassApplyer(classes)
}

func AttributeNS(name string, nameSpace string, value interface{}) MarkupOrChild {
	return markup.AttrApplyer{
		Name:      name,
		NameSpace: nameSpace,
		Value:     value,
	}
}

func Data(name string, value string) MarkupOrChild {
	return markup.DataApplyer{
		Name:  name,
		Value: value,
	}
}

func Property(name string, value interface{}) MarkupOrChild {
	return markup.PropApplyer{
		Name:  name,
		Value: value,
	}
}
