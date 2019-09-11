package data

import (
    "encoding/json"
    "errors"
    "io"
    "os"
)

const (
    // Default filesystem location for all build-a-coin data
    DefaultRootDir = "/usr/local/build-a-coin"
)
// default conf values
const (
    // Default filesystem location of configuration file
    DefaultConfPath = DefaultRootDir + "/build-a-coin.json"
    // Default filesystem location of base coins
    DefaultBasesDir = DefaultRootDir + "/bases"
    // Default IPv4 address to listen on
    DefaultListenAddr4 = "0.0.0.0"
    // Default TCP port to listen on
    DefaultListenPort = 39369
    // Default filesystem location of built-in assets
    DefaultAssetsDir = DefaultRootDir + "/assets"
    // Default database name
    DefaultDBName = "buildacoin"
    // Default database connection IP host (address or domain name)
    DefaultDBHost = "localhost"
    // Default database connection TCP port
    DefaultDBPort = 5432
    // Default database user
    DefaultDBUser = "buildacoin"
    // Default proxy status
    DefaultProxied = false
)

// two-layer struct necessary to have unexported fields for immutabilty but
// still use relatively easy un-JSONing.
// see http://stackoverflow.com/questions/11126793/golang-json-and-dealing-with-unexported-fields
type conf_ struct {
    ConfPath string `json:"-"`
    BasesDir string `json:"code templates dir"`
    ListenAddr4 string `json:"ipv4 listen address"`
    ListenPort uint16 `json:"listen port"`
    BaseCoins []string `json:"base coins"`
    Debug bool `json:"debug mode"`
    AssetsDir string `json:"assets dir"`
    DBName string `json:"database name"`
    DBHost string `json:"database host"`
    DBPort uint16 `json:"database port"`
    DBUser string `json:"database user"`
    DBPass string `json:"database password"`
    Proxied bool `json:"behind proxy"`
}

// Immutable configuration type containing config data for all packages
type Conf struct {
    conf_
}

// Load default config from file or fall back to built-in defaults.
func LoadConf() (*Conf) {
    conf := new(Conf)
    *conf = *DefaultConf
    conf.fromFile(DefaultConfPath)
    return conf
}

// Load a specific config file.  It is an error not to load from the file.
func LoadConfExplicit(path string) (*Conf, error) {
    conf := new(Conf)
    *conf = *DefaultConf
    err := conf.fromFile(path)
    if err != nil {
        return nil, err
    }
    return conf, nil
}

// If path is "", try to load the default config path and fall back to built-in
// defaults if necessary.  Otherwise, load the config at path successfully or
// return an error.
func LoadConfFromArg(path string) (*Conf, error) {
    var conf *Conf
    var err error
    if path == "" {
        conf = LoadConf()
    } else {
        conf, err = LoadConfExplicit(path)
        if err != nil {
            return nil, errors.New(
                "failed to load explicit config: " + err.Error())
        }
    }
    return conf, nil
}

func (tt *Conf) fromFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }

    backup := *tt

    dec := json.NewDecoder(file)
    err = dec.Decode(tt)
    if err != nil && err != io.EOF {
        *tt = backup
        return err
    }

    return nil
}

//
// Accessors
//

// Get the location of the general configuration file.
func (tt *Conf) ConfPath() string {
    return tt.conf_.ConfPath
}
// Get the location of the base coins
func (tt *Conf) BasesDir() string {
    return tt.conf_.BasesDir
}
// Get the IPv4 listen address
func (tt *Conf) ListenAddr4() string {
    return tt.conf_.ListenAddr4
}
// Get the TCP listen port
func (tt *Conf) ListenPort() uint16 {
    return tt.conf_.ListenPort
}
// Get the number of available base coins
func (tt *Conf) BaseCoinCount() int {
    return len(tt.conf_.BaseCoins)
}
// Get the base coin at index
func (tt *Conf) BaseCoin(idx int) string {
    return tt.conf_.BaseCoins[idx]
}
// Get whether debug mode is active
func (tt *Conf) Debug() bool {
    return tt.conf_.Debug
}
// Get the location of built-in assets
func (tt *Conf) AssetsDir() string {
    return tt.conf_.AssetsDir
}
// Get the name of the database
func (tt *Conf) DBName() string {
    return tt.conf_.DBName
}
// Get the database connection host
func (tt *Conf) DBHost() string {
    return tt.conf_.DBHost
}
// Get the database connection TCP port
func (tt *Conf) DBPort() uint16 {
    return tt.conf_.DBPort
}
// Get the database user
func (tt *Conf) DBUser() string {
    return tt.conf_.DBUser
}
// Get the database password
func (tt *Conf) DBPass() string {
    return tt.conf_.DBPass
}

// Get whether running behind a proxy
func (tt *Conf) Proxied() bool {
    return tt.conf_.Proxied
}

//
// Mutators: these functions return copies of the Conf they operate on rather
// than actually mutating.  Mostly for testing.
//

func (tt *Conf) dup() *Conf {
    output := new(Conf)
    *output = *tt
    return output
}

// Return a new Conf with ConfPath set to path
func (tt *Conf) WithConfPath(path string) *Conf {
    output := tt.dup()
    output.conf_.ConfPath = path
    return output
}
// Return a new Conf with BasesDir set to path
func (tt *Conf) WithBasesDir(path string) *Conf {
    output := tt.dup()
    output.conf_.BasesDir = path
    return output
}
// Return a new Conf with DBPass set to path
func (tt *Conf) WithDBPass(pass string) *Conf {
    output := tt.dup()
    output.conf_.DBPass = pass
    return output
}

//
// JSON serialization
//

// Populate the Conf from its JSON representation b.
func (tt *Conf) UnmarshalJSON(b []byte) error {
    return json.Unmarshal(b, &tt.conf_)
}
// Produce a JSON representation of the Conf.
func (tt *Conf) MarshalJSON() ([]byte, error) {
    return json.Marshal(&tt.conf_)
}

//
// Default config: not particularly useful to read since the struct layout is
// above and all values are consts which are also above
//

var DefaultConf = &Conf {
    conf_ {
        DefaultConfPath,
        DefaultBasesDir,
        DefaultListenAddr4,
        DefaultListenPort,
        nil, // BaseCoins
        false, // Debug
        DefaultAssetsDir,
        DefaultDBName,
        DefaultDBHost,
        DefaultDBPort,
        DefaultDBUser,
        "", // DBPass
        DefaultProxied,
    },
}
