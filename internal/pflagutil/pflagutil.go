package pflagutil

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

var _ pflag.Value = &KeyValueFlag{}

// KeyValueFlag can be used for users to define variables with quotes.
// Similar to pflag.StringArrayVar and StringToStringVar, but with different quoting.
type KeyValueFlag struct {
	Map *map[string]string
}

func (v *KeyValueFlag) String() string {
	return strings.Join(v.GetSlice(), "\n")
}

// GetSlice returns the flag value list as an array of strings.
func (v *KeyValueFlag) GetSlice() []string {
	m := make([]string, 0, len(*v.Map))
	for k, v := range *v.Map {
		m = append(m, fmt.Sprintf("%s=%s", k, v))
	}
	return m
}

// Type returns a descriptive name for the underlying data type.
// Used in the help output of pflag/cobra.
func (v *KeyValueFlag) Type() string {
	return "key=value"
}

func (v *KeyValueFlag) Set(kv string) error {
	v.init()

	equals := strings.IndexByte(kv, '=')
	if equals == -1 {
		return fmt.Errorf("Missing \"=\": %q", kv)
	}
	(*v.Map)[kv[:equals]] = kv[equals+1:]

	return nil
}

func (v *KeyValueFlag) init() {
	if *v.Map == nil {
		*v.Map = make(map[string]string)
	}
}
