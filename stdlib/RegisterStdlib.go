package stdlib

import (
	"github.com/Moritisimor/gomad/stdlib/arithmetics"
	"github.com/Moritisimor/gomad/stdlib/functions"
	"github.com/Moritisimor/gomad/stdlib/iolib"
	"github.com/Moritisimor/gomad/stdlib/lists"
	"github.com/Moritisimor/gomad/value"
)

func RegisterStdlib(env *value.Env) {
	arithmetics.RegisterFuns(env)
	functions.RegisterFuns(env)
	iolib.RegisterFuns(env)
	lists.RegisterFuns(env)
}
