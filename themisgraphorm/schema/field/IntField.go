package field

// base implement IField, base type is uint32
type intField struct {
	unsigned   bool
	name       string
	ai         bool
	unique     bool
	bitSize    int
	nullable   bool
	defaultVal any
}

func (i *intField) Null(v bool) IField {
	i.nullable = v // default true is null, false // is notnullable
	return i
}

func (i *intField) Default(v any) IField {
	i.defaultVal = v
	return i
}

func (i *intField) Unsigned() IField {
	i.unsigned = true
	return i
}

func (i *intField) Name(name string) IField {
	i.name = name
	return i
}

func (i *intField) Unique() IField {
	i.unique = true
	return i
}

func (i *intField) AI() IField {
	i.ai = true
	return i
}

func newIntField(bitSize int) *intField {
	return &intField{
		bitSize: bitSize,
	}
}
