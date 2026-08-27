package eval

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/lexer"
	"github.com/Moritisimor/gomad/parser"
	"github.com/Moritisimor/gomad/value"
)

func RegisterOSNatives(env *value.Env) {
	for _, r := range osNatives() {
		env.Define(r.name, r.fn)
	}
}

func osNatives() []nativeReg {
	return []nativeReg{
		reg("exec", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("exec", 1, len(args))
			}
			cmd, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			c := exec.Command("sh", "-c", cmd)
			err = c.Run()
			raw := 127.0
			if err == nil {
				raw = 0
			} else if ee, ok := err.(*exec.ExitError); ok {
				raw = float64(ee.ExitCode() * 256)
			}
			return value.Number{Val: raw}, nil
		}),
		reg("exit", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("exit", 1, len(args))
			}
			code, err := GetNumber(args[0], env)
			if err != nil {
				return nil, err
			}
			return nil, &value.Error{Kind: value.ErrExit, Code: floatToInt32(code)}
		}),
		reg("bye", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			return nil, &value.Error{Kind: value.ErrExit, Code: 0}
		}),
		reg("print_env", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 0 {
				return nil, arityErr("print_env", 0, len(args))
			}
			var b strings.Builder
			idx := 0
			for cur := env; cur != nil; cur = cur.Parent {
				fmt.Fprintf(&b, "Scope %d:\n", idx)
				keys := make([]string, 0, len(cur.Bindings))
				for k := range cur.Bindings {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Fprintf(&b, "\t%s: %s\n", k, cur.Bindings[k])
				}
				idx++
			}
			fmt.Print(b.String())
			return value.Unit{}, nil
		}),
		reg("include", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("include", 1, len(args))
			}
			pathSym, ok := args[0].(expr.Symbol)
			if !ok {
				return nil, value.EvalErr("include expects a symbol path")
			}
			content, err := readFileString(pathSym.Val)
			if err != nil {
				return nil, value.EvalErrf("Error while including '%s': %v", pathSym.Val, err)
			}
			tokens, err := lexer.Tokenize(content)
			if err != nil {
				return nil, &value.Error{Kind: value.ErrTokenize, Msg: err.Error()}
			}
			forms, err := parser.ParseProgram(tokens)
			if err != nil {
				return nil, &value.Error{Kind: value.ErrParse, Msg: err.Error()}
			}
			_, err = EvalSeq(forms, env)
			if err != nil {
				return nil, err
			}
			return value.Unit{}, nil
		}),
		reg("read_file", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("read_file", 1, len(args))
			}
			path, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			content, err := readFileString(path)
			if err != nil {
				return nil, value.EvalErrf("Error while reading '%s': %v", path, err)
			}
			return value.String{Val: content}, nil
		}),
		reg("write_file", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 2 {
				return nil, arityErr("write_file", 2, len(args))
			}
			path, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			content, err := GetString(args[1], env)
			if err != nil {
				return nil, err
			}
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return nil, value.EvalErrf("Couldn't write to '%s': %v", path, err)
			}
			return value.Unit{}, nil
		}),
		reg("remove_file", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("remove_file", 1, len(args))
			}
			path, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			if err := os.Remove(path); err != nil {
				return nil, value.EvalErrf("Couldn't remove file '%s': %v", path, err)
			}
			return value.Unit{}, nil
		}),
		reg("read_dir", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("read_dir", 1, len(args))
			}
			path, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil, value.EvalErrf("Couldn't read directory '%s': %v", path, err)
			}
			names := make([]string, 0, len(entries))
			for _, e := range entries {
				names = append(names, e.Name())
			}
			sort.Strings(names)
			vals := make([]value.Value, len(names))
			for i, n := range names {
				vals[i] = value.String{Val: n}
			}
			return value.FromVec(vals), nil
		}),
		reg("mkdir", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("mkdir", 1, len(args))
			}
			path, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			if err := os.MkdirAll(path, 0o755); err != nil && !os.IsExist(err) {
				return nil, value.EvalErrf("Couldn't create directory '%s': %v", path, err)
			}
			return value.Unit{}, nil
		}),
		reg("remove_dir", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("remove_dir", 1, len(args))
			}
			path, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			if err := os.Remove(path); err != nil {
				return nil, value.EvalErrf("Couldn't remove directory: '%s': %v", path, err)
			}
			return value.Unit{}, nil
		}),
		reg("chdir", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("chdir", 1, len(args))
			}
			path, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			if err := os.Chdir(path); err != nil {
				return nil, value.EvalErrf("Error while changing working directory to '%s': %v", path, err)
			}
			return value.Unit{}, nil
		}),
		reg("cwd", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 0 {
				return nil, arityErr("cwd", 0, len(args))
			}
			dir, err := os.Getwd()
			if err != nil {
				return nil, value.EvalErrf("Could not determine working directory: %v", err)
			}
			return value.String{Val: dir}, nil
		}),
		reg("get_env", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("get_env", 1, len(args))
			}
			key, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			v, ok := os.LookupEnv(key)
			if !ok {
				return nil, value.EvalErrf("Environment variable '%s' not found", key)
			}
			return value.String{Val: v}, nil
		}),
		reg("get_env_unit", func(args []expr.Expr, env *value.Env) (value.Value, error) {
			if len(args) != 1 {
				return nil, arityErr("get_env", 1, len(args))
			}
			key, err := GetString(args[0], env)
			if err != nil {
				return nil, err
			}
			if v, ok := os.LookupEnv(key); ok {
				return value.String{Val: v}, nil
			}
			return value.Unit{}, nil
		}),
	}
}

func readFileString(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func floatToInt32(x float64) int {
	if x != x {
		return 0
	}
	if x >= 1<<31-1 {
		return 1<<31 - 1
	}
	if x <= -(1 << 31) {
		return -(1 << 31)
	}
	return int(x)
}
