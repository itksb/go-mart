package auth

import "errors"

var ErrOptsLoggerIsNil = errors.New("logger is nil")
var ErrOptsIdentityProviderIsNil = errors.New("identityProvider is nil")
var ErrOptsHashAlgoIsNil = errors.New("hashAlgo is nil")
var ErrOptsTokenCreateIsNil = errors.New("tokenCreate is nil")
var ErrOptsTokenParseIsNil = errors.New("tokenParse is nil")
var ErrOptsSecretIsNil = errors.New("secretReader is nil")
var ErrOptsNowTimeIsNil = errors.New("nowTime is nil")
var ErrInvalidCredentials = errors.New("invalid user credentials")
