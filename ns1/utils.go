package ns1

// Map to String Map
func expandStringMap(v map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for key, val := range v {
		m[key] = val.(string)
	}

	return m
}

// StringList to StringSlice
func expandStringList(v []interface{}) []string {
	var vs []string

	for _, v := range v {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}

	return vs
}
