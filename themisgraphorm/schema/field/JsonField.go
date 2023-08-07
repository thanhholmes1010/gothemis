package field

type jsonField struct {
	castTypeRuntime any
	name            string
	unique          bool
	ai              bool
	nullable        bool
	defaultVal      any
}

func JSONField(castTypeRuntime any) IField {
	return &jsonField{
		castTypeRuntime: castTypeRuntime,
	}
}

func (j *jsonField) Name(name string) IField {
	j.name = name
	return j
}

func (j *jsonField) Unique() IField {
	j.unique = true
	return j
}

func (j *jsonField) AI() IField {
	j.ai = true
	return j
}

func (j *jsonField) Unsigned() IField {
	return j
}

func (j *jsonField) Null(v bool) IField {
	j.nullable = v
	return j
}

func (j *jsonField) Default(v any) IField {
	j.defaultVal = v
	return j
}

func (j *jsonField) GetName() string {
	return j.name
}

func (j *jsonField) CanNull() bool {
	return j.nullable
}
