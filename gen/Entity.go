package gen

import "bytes"

type property struct {
	isPrimaryKey      bool
	upperProperty     string
	propertyType      string
	tableField        string
	tableFieldComment string
}

type Entity struct {
	stubContext         string
	table               Table
	primaryKey          string
	filename            string
	outfile             string
	imports             []string
	properties          []string
	upperProperties     []string
	newFileBuf          *bytes.Buffer
	fieldMaxLen         int
	propertyTypeMaxLen  int
	upperPropertyMaxLen int
}

func (this *Entity) setFieldMaxLen(fieldMaxLen int) *Entity {
	this.fieldMaxLen = fieldMaxLen
	return this
}

func (this *Entity) setPropertyTypeMaxLen(propertyTypeMaxLen int) *Entity {
	this.propertyTypeMaxLen = propertyTypeMaxLen
	return this
}

func (this *Entity) setUpperPropertyMaxLen(upperPropertyMaxLen int) *Entity {
	this.upperPropertyMaxLen = upperPropertyMaxLen
	return this
}

func (this *Entity) setPrimaryKey(primaryKey string) *Entity {
	this.primaryKey = primaryKey
	return this
}

func (this *Entity) setTable(table Table) *Entity {
	this.table = table
	return this
}

func newEntity() *Entity {
	return &Entity{imports: make([]string, 0), properties: make([]string, 0), upperProperties: make([]string, 0), newFileBuf: bytes.NewBufferString("")}
}

func (this *Entity) Write(p []byte) (n int, err error) {
	return this.newFileBuf.Write(p)
}

func (this *Entity) gen() {

}
