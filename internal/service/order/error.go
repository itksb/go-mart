package order

import "errors"

var ErrOrderAlreadyCreatedByCurUser = errors.New("order already created by current user")
var ErrOrderAlreadyCreatedByAnotherUser = errors.New("order already created by another user")
var ErrOrderIncorrectOrderNumber = errors.New("order number is incorrect (luhn algo)")
