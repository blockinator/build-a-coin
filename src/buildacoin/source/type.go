package source

import "buildacoin/source/types"

// A type constrains template inputs to legal values for a field
type Type interface {
    types.Type
}
