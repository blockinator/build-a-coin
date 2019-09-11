package data

import (
    "encoding/json"
    "errors"
    "io"
    "os"
)

const (
    // Name of the metadata file within a base coin directory
    MetaFileName = "metadata.json"
)

var (
    // Error when a Meta tries to refer to it's associated Conf, but none
    // exists
    ErrNoConf error = errors.New("configuration missing")
    // Error when the ID of a Meta does not match the name of the base coin it
    // is a part of
    ErrMetaMismatch error = errors.New("retrieved metadata is not requested metadata")
)

// two-layer struct necessary to have unexported fields for immutabilty but
// still use relatively easy un-JSONing.
// see http://stackoverflow.com/questions/11126793/golang-json-and-dealing-with-unexported-fields
type meta_ struct {
    Id string
    Label string
    Version string
    InGroups []string `json:"input groups"`
    Inputs []Input `json:"user inputs"`
    Subs []Sub `json:"substitutions"`
}

// Metadata for a base coin, describing template inputs and outputs
type Meta struct {
    meta_
    conf *Conf
}

// A single user input for a base coin
type Input struct {
    // The input group this input will be displayed with ("basic"/"advanced"
    // etc.)
    Group string
    // Unique terse name for the input
    Id string
    // Longer descriptive name to be shown to users
    Label string
    // Default value for this input if an explicit value is not supplied
    Default string
}
// A single template substitution field in a base coin (an output)
type Sub struct {
    // Numeric identifier for this substitution as it will appear in the
    // template stream
    Idx uint `json:"substitution index"`
    // Input to pull value from or "" to always use default
    Input string
    // Description of field for human maintainers
    Comment string
    // Value to use if no input is associated
    Default string
    // Type that the value must conform to (see buildacoin/source/types)
    Type string
    Deps []uint `json:"dependencies"`
}

// Load a metadata file for the named base coin.
func LoadMeta(conf *Conf, name string) (*Meta, error) {
    output := new(Meta)

    filename := conf.BasesDir() + "/" + name + "/" + MetaFileName
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }

    dec := json.NewDecoder(file)
    err = dec.Decode(output)
    if err != nil && err != io.EOF {
        return nil, err
    }

    // Finding someone else's metadata in the base coin's directory is an error
    if output.Id() != name {
        return nil, ErrMetaMismatch
    }

    output.conf = conf

    return output, nil
}

// Get the template stream associated with this Meta (the base coin's template
// stream), along with the stream type
func (tt *Meta) Template() (io.ReadCloser, string, error) {
    if tt.conf == nil {
        return nil, "", ErrNoConf
    }
    return LoadTemplate(tt.conf, tt.Id())
}

//
// Accessors
//

// Get the base coin ID.
func (tt *Meta) Id() string {
    return tt.meta_.Id
}
// Get the base coin version.
func (tt *Meta) Version() string {
    return tt.meta_.Version
}
// Get the base coin's input groups.  Input groups group user options logically
// ("basic"/"advanced" etc).
func (tt *Meta) InGroups() []string {
    output := make([]string, len(tt.meta_.InGroups))
    copy(output, tt.meta_.InGroups)
    return output
}
// Get the number of input groups for this base coin.
func (tt *Meta) InGroupCount() int {
    return len(tt.meta_.InGroups)
}
// Get the input group at idx for this base coin.
func (tt *Meta) InGroup(idx int) string {
    return tt.meta_.InGroups[idx]
}
// Get the user input fields for this base coin.
func (tt *Meta) Inputs() []Input {
    output := make([]Input, len(tt.meta_.Inputs))
    copy(output, tt.meta_.Inputs)
    return output
}
// Get the number of user input fields for this base coin.
func (tt *Meta) InputCount() int {
    return len(tt.meta_.Inputs)
}
// Get the user input field at idx for this base coin.
func (tt *Meta) Input(idx int) Input {
    return tt.meta_.Inputs[idx]
}
// Get the substitutions (template outputs) for this base coin.
func (tt *Meta) Subs() []Sub {
    output := make([]Sub, len(tt.meta_.Subs))
    copy(output, tt.meta_.Subs)
    return output
}
// Get the number of substitutions (template outputs) for this base coin.
func (tt *Meta) SubCount() int {
    return len(tt.meta_.Subs)
}
// Get the substitution (template output) at idx for this base coin.
func (tt *Meta) Sub(idx int) Sub {
    return tt.meta_.Subs[idx]
}

//
// JSON serialization
//

// Populate this Meta from the JSON representation b.
func (tt *Meta) UnmarshalJSON(b []byte) error {
    return json.Unmarshal(b, &tt.meta_)
}
// Produce a JSON representation for this Meta.
func (tt *Meta) MarshalJSON() ([]byte, error) {
    return json.Marshal(&tt.meta_)
}

//
// Testing
//

// Produce a hand-crafted Meta instance.  Mostly useful for testing.
func NewMeta(id, label, version string, inGroups []string, inputs []Input,
        subs []Sub) *Meta {
    // This is a good place to have the less readable field:value struct
    // literal format.  Having it is nice because it will break the build if
    // another field is added to the struct, which makes me think twice about
    // accessors and other new-field necessities.
    return &Meta { meta_ { id, label, version, inGroups, inputs, subs },
            nil }
}
