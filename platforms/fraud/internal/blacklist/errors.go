package blacklist

import "errors"

var ErrEntryNotFound = errors.New("fraud: blacklist entry not found")
var ErrEntryExists = errors.New("fraud: blacklist entry already exists")
