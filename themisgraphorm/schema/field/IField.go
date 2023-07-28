package field

type IField interface {
	Name(name string) IField
	Unique() IField   // set unique constraint, same create unique indexing
	AI() IField       // set auto_increment
	Unsigned() IField // set unsigned
	Null(v bool) IField
	Default(v any) IField
	GetName() string
	CanNull() bool
}

func Integer(bitSize int) IField {
	return newIntField(bitSize)
}

func Varchar(size int) IField {
	return newVarcharField(size)
}

func newVarcharField(size int) *varcharField {
	return &varcharField{
		size: size,
	}
}
