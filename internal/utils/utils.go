package utils

import "yema.dev/pkg/cfgstruct"

func IsDev() bool {
	env := cfgstruct.DefaultsType()
	return env == "" || env == "dev"
}
