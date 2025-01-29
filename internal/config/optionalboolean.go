package config

import (
	"fmt"
	"strconv"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/spf13/pflag"
)

// OptionalBoolean represents a boolean that can also be 'not set', eg. having three possible states.
// OptionalBoolean implements pflag.Value and flag.Value and as such, can be used with flag.Var() and friends.
// However, it does not allow to be set more than once, in order to support multiple flags touching the
// same boolean, eg. when the use case is --enable-something'/--disable-something
type OptionalBoolean int

// TODO: figure out if the default value handling could be done in a nicer way - maybe a separate type instead of trying
// to force one type to rule them all?
// Downside with this approach is that bool == True might be false which is kinda unexpected?
// This is required for 'boolean flags with default values', in order to be able to set the default value and not throw
// an error when trying to set a non-Unset OptionalBoolean

const (
	// FalseFlagValue false Bool String value
	FalseFlagValue = "false"
	// TrueFlagValue true Bool String value
	TrueFlagValue = "true"
)

const (
	// Unset is the OptionalBoolean value representing not set
	Unset OptionalBoolean = iota // 0
	// True is the OptionalBoolean value representing true
	True
	// False is the OptionalBoolean value representing false
	False

	// DefaultTrue is the OptionalBoolean value representing not set, but with a default value of true.
	// It returns false from IsSet, to allow for overriding the default (once)
	DefaultTrue
	// DefaultFalse is the OptionalBoolean value representing not set, but with a default value of false.
	// It returns false from IsSet, to allow for overriding the default (once)
	DefaultFalse
)

// AddEnableDisableFlags is a convenience method to generate --enable-something and --disable-something
// flags with the correct settings.
//
// *name* specifies the name of the flags (eg. '--enable-[name]')
// *subject* is used to create the usage description for the flags, in the form of 'Enable [subject]'.
func AddEnableDisableFlags(flags *pflag.FlagSet, target *OptionalBoolean, name, subject string) {
	flags.Var(target, fmt.Sprintf("enable-%s", name), fmt.Sprintf("Enable %s.", subject))
	flags.Var(target, fmt.Sprintf("disable-%s", name), fmt.Sprintf("Disable %s.", subject))
	flags.Lookup(fmt.Sprintf("enable-%s", name)).NoOptDefVal = TrueFlagValue
	flags.Lookup(fmt.Sprintf("enable-%s", name)).DefValue = ""
	flags.Lookup(fmt.Sprintf("disable-%s", name)).NoOptDefVal = FalseFlagValue
	flags.Lookup(fmt.Sprintf("disable-%s", name)).DefValue = ""
}

// AddEnableOrDisableFlag is a convenience method to generate --enable-something *or* --disable-something
// flag with the correct settings, to overrider the default value. eg. if default is true, flag --disable-something
// will be generated.
//
// *name* specifies the name of the flag (eg. '--enable-[name]')
// *subject* is used to create the usage description for the flag, in the form of 'Enable [subject]'.
func AddEnableOrDisableFlag(flags *pflag.FlagSet, target *OptionalBoolean, defaultValue bool, name, subject string) {
	if defaultValue {
		target.SetDefault(true)
		flags.Var(target, fmt.Sprintf("disable-%s", name), fmt.Sprintf("Disable %s.", subject))
		flags.Lookup(fmt.Sprintf("disable-%s", name)).NoOptDefVal = FalseFlagValue
		flags.Lookup(fmt.Sprintf("disable-%s", name)).DefValue = ""
	} else {
		target.SetDefault(false)
		flags.Var(target, fmt.Sprintf("enable-%s", name), fmt.Sprintf("Enable %s.", subject))
		flags.Lookup(fmt.Sprintf("enable-%s", name)).NoOptDefVal = TrueFlagValue
		flags.Lookup(fmt.Sprintf("enable-%s", name)).DefValue = ""
	}
}

// AddToggleFlag is a convenience method to generate --toggle type of a flag with the correct settings and
// a default value.
func AddToggleFlag(flags *pflag.FlagSet, target *OptionalBoolean, name string, defaultValue bool, usage string) {
	if defaultValue {
		target.SetDefault(true)
		flags.Var(target, name, usage)
		flags.Lookup(name).NoOptDefVal = FalseFlagValue
		flags.Lookup(name).DefValue = ""
	} else {
		target.SetDefault(false)
		flags.Var(target, name, usage)
		flags.Lookup(name).NoOptDefVal = TrueFlagValue
		flags.Lookup(name).DefValue = ""
	}
}

// String implements flag.Value
func (s OptionalBoolean) String() string {
	return fmt.Sprintf("%t", s.Value())
}

// Value returns the underlying bool for OptionalBoolean.
// nb. returns false if flag is not set
func (s OptionalBoolean) Value() bool {
	b := boolFromSetBool(s)
	return b
}

// Type implements pflag.Value
func (s OptionalBoolean) Type() string {
	// return empty type in order to not display the parameter in help texts
	// seems like this has no other meaning(?)
	return ""
}

// Set implements flag.Value
// nb. OptionalBoolean will not allow itself to be set twice, if you want to have an underlying
// default value, use SetDefault()
func (s *OptionalBoolean) Set(value string) error {
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

// IsSet returns true if OptionalBoolean has been set
func (s OptionalBoolean) IsSet() bool {
	return s == True || s == False
}

// OverrideNotSet returns an OptionalBoolean set to b if the original OptionalBoolean was not set
func (s OptionalBoolean) OverrideNotSet(b bool) OptionalBoolean {
	if s.IsSet() {
		return s
	}
	return setBoolFromBool(b)
}

// AsUpcloudBoolean return OptionalBoolean as upcloud.Boolean
// nb. DefaultTrue and DefaultEmpty return upcloud.Empty, as upcloud.Boolean has no concept of default values
func (s OptionalBoolean) AsUpcloudBoolean() upcloud.Boolean {
	switch s {
	case Unset:
		return upcloud.Empty
	case True:
		return upcloud.True
	case False:
		return upcloud.False
	case DefaultTrue:
		return upcloud.Empty
	case DefaultFalse:
		return upcloud.Empty
	}
	panic("unknown OptionalBoolean")
}

// SetDefault sets the default value of OptionalBoolean to b
// Default value is returned from Value() if the OptionalBoolean has not been set
func (s *OptionalBoolean) SetDefault(b bool) {
	if b {
		*s = DefaultTrue
	} else {
		*s = DefaultFalse
	}
}

func boolFromSetBool(s OptionalBoolean) bool {
	return s == True || s == DefaultTrue
}

func setBoolFromBool(b bool) OptionalBoolean {
	if b {
		return True
	}
	return False
}
