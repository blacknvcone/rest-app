package validations

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// StructValidation const
const (
	StructValidationTimeAfterNow                      = "time_after_now"
	StructValidationTimeAfterField                    = "time_after_field"
	StructValidationMinimumIfFieldEqual               = "min_if_field_eq"
	StructValidationMaximumIfFieldEqual               = "max_if_field_eq"
	StructValidationLessThanEqualFieldIfFieldEqual    = "lte_field_if_field_eq"
	StructValidationGreaterThanEqualFieldIfFieldEqual = "gte_field_if_field_eq"
	StructValidationMinimumFieldIfFieldEqual          = "min_field_if_field_eq"
	StructValidationMaximumFieldIfFieldEqual          = "max_field_if_field_eq"
)

// InitStructValidation init struct validation
func InitStructValidation() {
	structValidation := map[string]func(fl validator.FieldLevel) bool{
		StructValidationTimeAfterNow:                      TimeAfterNow,
		StructValidationTimeAfterField:                    TimeAfterField,
		StructValidationMinimumIfFieldEqual:               MinIfFieldEqual,
		StructValidationMaximumIfFieldEqual:               MaxIfFieldEqual,
		StructValidationLessThanEqualFieldIfFieldEqual:    LTEFieldIfFieldEqual,
		StructValidationGreaterThanEqualFieldIfFieldEqual: GTEFieldIfFieldEqual,
		StructValidationMinimumFieldIfFieldEqual:          MinFieldIfFieldEqual,
		StructValidationMaximumFieldIfFieldEqual:          MaxFieldIfFieldEqual,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for tag, validationFunc := range structValidation {
			err := v.RegisterValidation(tag, validationFunc)
			if err != nil {
				panic(fmt.Errorf("can not register validation function: %s", tag))
			}
		}
	}
}
