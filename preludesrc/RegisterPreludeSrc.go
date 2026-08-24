// preludesrc is the package containing functions that are written in gomad itself, not gomad callbacks
// The reason this is a seperate package is because these functions depend on native callbacks that may not be registered.
package preludesrc

import (
	"github.com/Moritisimor/gomad/preludesrc/lists"
	"github.com/Moritisimor/gomad/preludesrc/stringlib"
	"github.com/Moritisimor/gomad/value"
)

func RegisterPreludeSrc(env *value.Env) {
	lists.RegisterFuns(env)
	stringlib.RegisterFuns(env)
}
