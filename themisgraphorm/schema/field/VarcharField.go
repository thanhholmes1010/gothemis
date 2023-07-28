package field

type varcharField struct {
	name       string
	unique     bool
	ai         bool
	unsigned   bool
	size       int
	nullable   bool
	defaultVal any
}

func (v *varcharField) Null(value bool) IField {
	v.nullable = value
	return v
}

func (v *varcharField) Default(value any) IField {
	v.defaultVal = value
	return v
}

func (v *varcharField) Name(name string) IField {
	v.name = name
	return v
}

func (v *varcharField) Unique() IField {
	v.unique = true
	return v
}

func (v *varcharField) AI() IField {
	v.ai = true
	return v
}

func (v *varcharField) Unsigned() IField {
	v.unsigned = true
	return v
}
