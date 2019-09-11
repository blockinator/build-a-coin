package source

import (
    "buildacoin/data"
    "buildacoin/source/types"
    "buildacoin/template"
    "errors"
    "io"
)

var (
    // Error when no new dependencies can be met
    ErrDepDeadlock error = errors.New("dependency deadlock")
)

// Return a reader of the source archive for a new coin.
// `meta` describes the substitutions to be made in the template (see
// data.Meta).
// `values` are the user-provided template string values which will be
// converted from strings if necessary according to their types in `meta`.
// `base` is the archive template itself
func Generate(meta *data.Meta, base io.Reader,
        values map[string]string) (io.Reader, error) {

    filterMap, err := BuildFilterMap(meta, values)
    if err != nil {
        return nil, err
    }

    return template.NewFilter(base, filterMap), nil
}

// Build a mapping acceptable to the template filter for given metadata and
// user values 
func BuildFilterMap(meta *data.Meta,
        values map[string]string) (template.FilterMap, error) {

    output := make(template.FilterMap)

    allSubs := meta.Subs()
    subs := make([]data.Sub, len(allSubs))
    unmet := make([]data.Sub, 0, len(allSubs))
    depArgs := make([]string, 10)

    // for each pass over the substitutions, subs with unmet dependencies are
    // left until either the list of subs doesn't shrink or shrinks to zero
    copy(subs, allSubs)
    deploop:
    for {
        // for each substitution field described by the base coin metadata. . .
        passloop:
        for _, sub := range subs {
            // does this substitution field have dependencies?
            depArgs = depArgs[:0]
            deps := sub.Deps
            if len(deps) > 0 {
                for _, depIdx := range deps {
                    // if any dependency isn't available, add this sub to the
                    // unmet list and move on
                    depBytes, ok := output[depIdx]
                    if !ok {
                        unmet = append(unmet, sub)
                        continue passloop
                    }
                    depArgs = append(depArgs, string(depBytes))
                }
            }

            // find the type of the field (specifically the string -> type
            // conversion func)
            valueType, ok := types.Map[sub.Type]
            if !ok {
                return nil, ErrUnknownType { sub.Input, sub.Type }
            }

            // look for user input for the field in the supplied values map,
            // and use the default value if it wasn't given
            var input string
            if inputValue, ok := values[sub.Input]; sub.Input != "" && ok {
                input = inputValue
            } else {
                input = sub.Default
            }

            // convert the string input to a string representing the typed
            // value
            inputs := append([]string { input }, depArgs...)
            valueString, err := valueType.Produce(inputs...)
            if err != nil {
                return nil, ErrBadFieldValue { sub.Input, input, sub.Type }
            }

            // place the vetted value in the filter map at the appropriate
            // substitution index
            output[sub.Idx] = []byte(valueString)
        }
        if len(unmet) < 1 {
            break deploop
        }
        if len(unmet) >= len(subs) {
            return nil, ErrDepDeadlock
        }
        // the subs whose dependencies were unmet this pass are the source of
        // the next pass
        subs = subs[:len(unmet)]
        copy(subs, unmet)
        unmet = unmet[:0]
    }

    return output, nil
}

// Error when a provided substitution field is of an unknown type.
type ErrUnknownType struct { field string; typeName string }
func (tt ErrUnknownType) Error() string {
    return "field '" + tt.field + "' has unknown type '" + tt.typeName + "'"
}
// Error when the conversion for a substitution field from a string fails.
type ErrBadFieldValue struct { field string; value string; typeName string }
func (tt ErrBadFieldValue) Error() string {
    return "bad value '" + tt.value + "' for field '" + tt.field +
            "' of type '" + tt.typeName + "'"
}
