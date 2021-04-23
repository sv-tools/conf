package conf

// Transform is a function to transform the data
type Transform func(key, value interface{}, c Conf) interface{}
