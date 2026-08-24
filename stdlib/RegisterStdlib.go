package stdlib

import (
	"github.com/Moritisimor/gomad/stdlib/arithmetics"
	"github.com/Moritisimor/gomad/stdlib/conditionals"
	"github.com/Moritisimor/gomad/stdlib/exceptions"
	"github.com/Moritisimor/gomad/stdlib/functions"
	"github.com/Moritisimor/gomad/stdlib/iolib"
	"github.com/Moritisimor/gomad/stdlib/lists"
	"github.com/Moritisimor/gomad/stdlib/macros"
	"github.com/Moritisimor/gomad/stdlib/records"
	"github.com/Moritisimor/gomad/stdlib/stringlib"
	"github.com/Moritisimor/gomad/stdlib/variables"
	"github.com/Moritisimor/gomad/value"
)

func RegisterStdlib(env *value.Env) {
	arithmetics.RegisterFuns(env)
	functions.RegisterFuns(env)
	iolib.RegisterFuns(env)
	lists.RegisterFuns(env)
	variables.RegisterFuns(env)
	conditionals.RegisterFuns(env)
	records.RegisterFuns(env)
	macros.RegisterFuns(env)
	exceptions.RegisterFuns(env)
	stringlib.RegisterFuns(env)
}
