// prelude is the package containing common functions a programming language might need
// It includes functions for stuff such as arithmetics, conditionals, variables, defining functions variables etc.
package prelude

import (
	"github.com/Moritisimor/gomad/prelude/arithmetics"
	"github.com/Moritisimor/gomad/prelude/conditionals"
	"github.com/Moritisimor/gomad/prelude/exceptions"
	"github.com/Moritisimor/gomad/prelude/functions"
	"github.com/Moritisimor/gomad/prelude/iolib"
	"github.com/Moritisimor/gomad/prelude/lists"
	"github.com/Moritisimor/gomad/prelude/macros"
	"github.com/Moritisimor/gomad/prelude/records"
	"github.com/Moritisimor/gomad/prelude/stringlib"
	"github.com/Moritisimor/gomad/prelude/variables"
	"github.com/Moritisimor/gomad/value"
)

func RegisterCommonPrelude(env *value.Env) {
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
