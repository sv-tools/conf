package conf

// Transform is a function to transform the data
type Transform func(key string, value interface{}, c Conf) interface{}
