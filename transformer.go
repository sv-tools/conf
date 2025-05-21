package conf

// Transform is a function to transform the data
type Transform func(key string, value any, c Conf) any
