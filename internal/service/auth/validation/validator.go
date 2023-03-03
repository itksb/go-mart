package validation

type ValidateFunc func(v interface{}) (ok bool, err error)
