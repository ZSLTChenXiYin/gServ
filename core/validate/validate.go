package validate

func Validate(a any) error {
	return data_validator.Struct(a)
}
