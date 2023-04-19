package strings

func DefaultWhenEmpty(value *string, defaultValue string) string {
	if value == nil || *value == "" {
		return defaultValue
	}

	return *value
}
