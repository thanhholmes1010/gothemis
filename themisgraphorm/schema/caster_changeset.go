package schema

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type ChangeCaster interface {
	ValidateRequired(fieldNames ...string) ChangeCaster
	ValidPattern(fieldName string, patternRegex []byte) ChangeCaster
	ValidUnique(fieldName string) ChangeCaster
	IsValidAll() error
}

type ValidateError uint8

var errTypeToString = [...]string{"RequiredError", "PatternError", "UniqueError"}

func (v ValidateError) toErrString(fieldName string, need interface{}, have interface{}) string {
	return fmt.Sprintf("[%v]: Validate on Field [%v] failed, "+
		"need [%v] but have [%v]\n",
		errTypeToString[v], fieldName, need, have)
}

const (
	ValidateRequiredError ValidateError = iota
	ValidatePatternError
	ValidateUniqueError
)

type boxFieldCaster struct {
	col          string
	val          any
	indexNotNull int
}

type changeCasterImpl struct {
	boxes        map[string]*boxFieldCaster
	validFns     []func() string
	indexNotNull uint32
}

func (c *changeCasterImpl) ValidateRequired(fieldNames ...string) ChangeCaster {
	c.validFns = append(c.validFns, func() string {
		if c.indexNotNull == 0 {
			return ""
		}
		var errString string = ""
		if len(fieldNames) == 0 {
			for fieldName, box := range c.boxes {
				if (c.indexNotNull & (1 << c.boxes[fieldName].indexNotNull)) != 0 {
					errString += ValidateRequiredError.toErrString(fieldName, "Not Null", box.val)
				}
			}
		} else {
			for _, fieldName := range fieldNames {
				if _, ok := c.boxes[fieldName]; ok {
					val := c.boxes[fieldName].val
					if (c.indexNotNull & (1 << c.boxes[fieldName].indexNotNull)) != 0 {
						errString += ValidateRequiredError.toErrString(fieldName, "Not Null", val)
					}
				}
			}
		}
		return errString
	})
	return c
}

// ValidPattern only used for field have value is varchar or char implement string type
func (c *changeCasterImpl) ValidPattern(fieldName string, patternRegex []byte) ChangeCaster {
	c.validFns = append(c.validFns, func() string {
		if _, ok := c.boxes[fieldName]; ok {
			fmt.Println(c.boxes[fieldName])
			if valString, ok := c.boxes[fieldName].val.(string); ok {
				regex, err := regexp.Compile(string(patternRegex))
				if err != nil {
					return ValidatePatternError.toErrString(fieldName,
						"pattern regex valid", string(patternRegex)+" not valid")
				}
				if regex.MatchString(valString) {
					return ""
				} else {
					return ValidatePatternError.toErrString(fieldName, string(patternRegex), valString)
				}
			} else {
				return ValidatePatternError.toErrString(fieldName, string(patternRegex), "Not String Type")
			}
		}
		return ValidatePatternError.toErrString(fieldName, "must exist on schema", fieldName+"not exist")
	})
	return c
}

func (c *changeCasterImpl) ValidUnique(fieldName string) ChangeCaster {
	return c
}

func (c *changeCasterImpl) IsValidAll() error {
	allErr := ""
	for _, validFn := range c.validFns {
		allErr += validFn()
	}
	if allErr == "" {
		return nil
	}
	return fmt.Errorf(allErr)
}

func Cast(fromMessage any, schemaMigrator Migrator) ChangeCaster {
	rvSchema := reflect.Indirect(reflect.ValueOf(schemaMigrator))
	rvMessage := reflect.Indirect(reflect.ValueOf(fromMessage))
	prefixSchemaName := rvSchema.Type().Name()
	indexBox := 0
	cs := &changeCasterImpl{
		boxes:        make(map[string]*boxFieldCaster),
		indexNotNull: 0,
	}
	for indexBox < len(schemaMigrator.DefineFields()) {
		fieldSchema := schemaMigrator.DefineFields()[indexBox]
		fieldName := fieldSchema.GetName()
		if _, ok := cs.boxes[fieldSchema.GetName()]; !ok {
			newBox := &boxFieldCaster{
				col:          fieldName,
				val:          nil,
				indexNotNull: indexBox,
			}
			cs.boxes[fieldSchema.GetName()] = newBox
			fmt.Printf("field name [%v] can null: %v\n", fieldName, fieldSchema.CanNull())
			if !fieldSchema.CanNull() /* this mean cannot nullable, must rebuild for required */ {
				cs.indexNotNull |= 1 << indexBox
			}
			indexBox++
		}
	}
	switch rvMessage.Kind() {
	case reflect.Struct:
		for i := 0; i < rvMessage.NumField(); i++ {
			fieldMessageNameHavePrefixSchemaName := rvMessage.Type().Field(i).Name
			rvMessageField := rvMessage.Field(i)
			if strings.HasPrefix(fieldMessageNameHavePrefixSchemaName, prefixSchemaName) {
				var fieldNameTrimedPrefixSchemaName = fieldMessageNameHavePrefixSchemaName[len(prefixSchemaName):]
				if _, ok := rvSchema.Type().FieldByName(fieldNameTrimedPrefixSchemaName); ok {
					// exist in schema with field name
					if rvMessageField.Type().Name() ==
						rvSchema.FieldByName(fieldNameTrimedPrefixSchemaName).Type().Name() {
						// check same type success
						if !rvMessageField.IsZero() {
							// only cast val not nil, not zero, not empty string
							if _, ok := cs.boxes[fieldNameTrimedPrefixSchemaName]; ok {
								// get box with fieldName
								rvSchema.FieldByName(fieldNameTrimedPrefixSchemaName).
									Set(reflect.ValueOf(rvMessageField.Interface()))
								cs.boxes[fieldNameTrimedPrefixSchemaName].val = rvMessageField.Interface()
								// clear index not null exist from it
								cs.indexNotNull &= ^(1 << cs.boxes[fieldNameTrimedPrefixSchemaName].indexNotNull)
							}
						}
					}
				}
			}
		}
	}
	return cs
}
