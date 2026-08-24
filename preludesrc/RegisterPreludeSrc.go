// preludesrc is the package containing functions that are written in gomad itself, not gomad callbacks
// The reason this is a seperate package is because these functions depend on native callbacks that may not be registered.
package preludesrc

import (
	"fmt"

	"github.com/Moritisimor/gomad/preludesrc/lists"
	"github.com/Moritisimor/gomad/preludesrc/stringlib"
	"github.com/Moritisimor/gomad/value"
)

func RegisterPreludeSrc(env *value.Env) error {
	if err := lists.RegisterFuns(env); err != nil {
		return fmt.Errorf("Error while registering list functions: %s", err.Error())
	}

	if err := stringlib.RegisterFuns(env); err != nil {
		return fmt.Errorf("Error while registering string functions: %s", err.Error())
	}

	return nil
}
