package gen

type fieldsListCollect struct {
	upperProperty, propertyType, propertyTag string
	collect                                  []fieldRaw
}

func newFieldsListCollect() *fieldsListCollect {
	return &fieldsListCollect{collect: make([]fieldRaw, 0)}
}

func (this *fieldsListCollect) append(upperProperty, propertyType, propertyTag string) {
	has := false
	var index int
	for i, s := range this.collect {
		if s.upperProperty == upperProperty {
			has = true
			index = i
			break
		}
	}
	if !has {
		this.collect = append(this.collect, fieldRaw{upperProperty: upperProperty, propertyType: propertyType, propertyTag: propertyTag})
	} else {
		this.collect[index] = fieldRaw{upperProperty: upperProperty, propertyType: propertyType, propertyTag: propertyTag}
	}
}

type fieldRaw struct {
	upperProperty, propertyType, propertyTag string
}
