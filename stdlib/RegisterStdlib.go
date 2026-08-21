package stdlib

import (
	"github.com/Moritisimor/gomad/stdlib/arithmetics"
	"github.com/Moritisimor/gomad/stdlib/functions"
	"github.com/Moritisimor/gomad/value"
)

func RegisterStdlib(env *value.Env) {
	arithmetics.RegisterFuns(env)
	functions.RegisterFuns(env)
}
