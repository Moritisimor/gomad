package stringlib

import (
	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) error {
	if _, err := eval.DoString(`
		(letfun has_prefix (s1 s2) (begins_with (chars s1) (chars s2)))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun has_suffix (s1 s2) (ends_with (chars s1) (chars s2)))
	`, env); err != nil {
		return err
	}

	return nil
}
