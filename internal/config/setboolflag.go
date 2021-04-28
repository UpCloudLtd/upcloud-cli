package config

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/pflag"
	"strconv"
)

// SetBool represents a boolean that can also be 'not set', eg. having three possible states.
// SetBool implements pflag.Value and flag.Value and as such, can be used with flag.Var() and friends.
// However, it does not allow to be set more than once, in order to support multiple flags touching the
// same boolean, eg. when the use case is --enable-something'/--disable-something
type SetBool int

const (
	// Unset is the SetBool value representing not set
	Unset SetBool = iota // 0
	// True is the SetBool value representing true
	True
	// False is the SetBool value representing false
	False
)

// AddEnableDisableFlags is a convenience method to generate --enable-something and --disable-something
// flags with the correct settings. *name* specifies the name of the flags (eg. '--enable-[name]') and
// *subject* is used to create the usage description for the flags, in the form of 'Enable [subject]'.
func AddEnableDisableFlags(flags *pflag.FlagSet, target *SetBool, name, subject string) {
	flags.Var(target, fmt.Sprintf("enable-%s", name), fmt.Sprintf("Enable %s.", subject))
	flags.Var(target, fmt.Sprintf("disable-%s", name), fmt.Sprintf("Disable %s.", subject))
	flags.Lookup(fmt.Sprintf("enable-%s", name)).NoOptDefVal = "true"
	flags.Lookup(fmt.Sprintf("enable-%s", name)).DefValue = ""
	flags.Lookup(fmt.Sprintf("disable-%s", name)).NoOptDefVal = "false"
	flags.Lookup(fmt.Sprintf("disable-%s", name)).DefValue = ""
}

// String implements flag.Value
func (s SetBool) String() string {
	return fmt.Sprintf("%t", s.Value())
}

// Value returns the underlying bool for SetBool.
// nb. returns false if flag is not set
func (s SetBool) Value() bool {
	b := boolFromSetBool(s)
	return b
}

// Type implements pflag.Value
func (s SetBool) Type() string {
	return ""
}

// Set implements flag.Value
func (s *SetBool) Set(value string) error {
	if s.IsSet() {
		return fmt.Errorf("cannot set twice")
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	*s = setBoolFromBool(boolVal)
	return nil
}

// IsSet returns true if SetBool is set
func (s SetBool) IsSet() bool {
	return s != Unset
}

// ApplyDefault returns a SetBool set to b if the original SetBool was not set
func (s SetBool) ApplyDefault(b bool) SetBool {
	if s.IsSet() {
		return s
	}
	return setBoolFromBool(b)
}

// AsUpcloudBoolean return SetBool as upcloud.Boolean
func (s SetBool) AsUpcloudBoolean() upcloud.Boolean {
	switch s {
	case Unset:
		return upcloud.Empty
	case True:
		return upcloud.True
	case False:
		return upcloud.False
	}
	panic("unknown SetBool")
}

func boolFromSetBool(s SetBool) bool {
	return s == True
}

func setBoolFromBool(b bool) SetBool {
	if b {
		return True
	}
	return False
}
