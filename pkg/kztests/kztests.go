package kztests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type TestKit struct {
	Arg  any
	Arg1 any
	Arg2 any
	Arg3 any
	Arg4 any
	Arg5 any
	Arg6 any
	Arg7 any
	Args []any

	Result  any
	Result1 any
	Result2 any
	Result3 any
	Result4 any
	Result5 any
	Result6 any
	Result7 any
	Results []any
}

func RunTests(t *testing.T, functions any, tkits []TestKit) {
	fffs := []reflect.Value{}
	functionsRV := reflect.ValueOf(functions)
	switch functionsRV.Kind() {
	case reflect.Func:
		fffs = append(fffs, functionsRV)
	case reflect.Slice:
		for i := 0; i < functionsRV.Len(); i++ {
			functionRV := functionsRV.Index(i)
			if functionRV.Kind() == reflect.Interface {
				//functionRV = reflect.Indirect(functionRV)
				functionRV = reflect.ValueOf(functionRV.Interface())
			}
			fffs = append(fffs, functionRV)
		}
	}

	for _, tkit := range tkits {
		tkitArgs := make([]reflect.Value, 0, 7+len(tkit.Args))
		tkitArgsStrings := make([]string, 0, 7+len(tkit.Args))
		if tkit.Arg != nil {
			tkitArgs = append(tkitArgs, reflect.ValueOf(tkit.Arg))
			tkitArgsStrings = append(tkitArgsStrings, convertValue2String(tkit.Arg))
		}
		if tkit.Arg1 != nil {
			tkitArgs = append(tkitArgs, reflect.ValueOf(tkit.Arg1))
			tkitArgsStrings = append(tkitArgsStrings, convertValue2String(tkit.Arg1))
		}
		if tkit.Arg2 != nil {
			tkitArgs = append(tkitArgs, reflect.ValueOf(tkit.Arg2))
			tkitArgsStrings = append(tkitArgsStrings, convertValue2String(tkit.Arg2))
		}
		if tkit.Arg3 != nil {
			tkitArgs = append(tkitArgs, reflect.ValueOf(tkit.Arg3))
			tkitArgsStrings = append(tkitArgsStrings, convertValue2String(tkit.Arg3))
		}
		if tkit.Arg4 != nil {
			tkitArgs = append(tkitArgs, reflect.ValueOf(tkit.Arg4))
			tkitArgsStrings = append(tkitArgsStrings, convertValue2String(tkit.Arg4))
		}
		if tkit.Arg5 != nil {
			tkitArgs = append(tkitArgs, reflect.ValueOf(tkit.Arg5))
			tkitArgsStrings = append(tkitArgsStrings, convertValue2String(tkit.Arg5))
		}
		if tkit.Arg6 != nil {
			tkitArgs = append(tkitArgs, reflect.ValueOf(tkit.Arg6))
			tkitArgsStrings = append(tkitArgsStrings, convertValue2String(tkit.Arg6))
		}
		if tkit.Arg7 != nil {
			tkitArgs = append(tkitArgs, reflect.ValueOf(tkit.Arg7))
			tkitArgsStrings = append(tkitArgsStrings, convertValue2String(tkit.Arg7))
		}
		for _, tkitArg := range tkit.Args {
			tkitArgs = append(tkitArgs, reflect.ValueOf(tkitArg))
			tkitArgsStrings = append(tkitArgsStrings, convertValue2String(tkitArg))
		}

		for _, fff := range fffs {
			results := fff.Call(tkitArgs)

			tkitResults := make([]reflect.Value, 0, 7+len(tkit.Results))
			tkitResultsStrings := make([]string, 0, 7+len(tkit.Results))
			if tkit.Result != nil {
				tkitResults = append(tkitResults, reflect.ValueOf(tkit.Result))
				tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkit.Result))
			}
			if tkit.Result1 != nil {
				tkitResults = append(tkitResults, reflect.ValueOf(tkit.Result1))
				tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkit.Result1))
			}
			if tkit.Result2 != nil {
				tkitResults = append(tkitResults, reflect.ValueOf(tkit.Result2))
				tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkit.Result2))
			}
			if tkit.Result3 != nil {
				tkitResults = append(tkitResults, reflect.ValueOf(tkit.Result3))
				tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkit.Result3))
			}
			if tkit.Result3 != nil {
				tkitResults = append(tkitResults, reflect.ValueOf(tkit.Result3))
				tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkit.Result3))
			}
			if tkit.Result4 != nil {
				tkitResults = append(tkitResults, reflect.ValueOf(tkit.Result4))
				tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkit.Result4))
			}
			if tkit.Result5 != nil {
				tkitResults = append(tkitResults, reflect.ValueOf(tkit.Result5))
				tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkit.Result5))
			}
			if tkit.Result6 != nil {
				tkitResults = append(tkitResults, reflect.ValueOf(tkit.Result6))
				tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkit.Result6))
			}
			if tkit.Result7 != nil {
				tkitResults = append(tkitResults, reflect.ValueOf(tkit.Result7))
				tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkit.Result7))
			}
			if len(tkit.Results) > 0 {
				for _, tkitResult := range tkit.Results {
					tkitResults = append(tkitResults, reflect.ValueOf(tkitResult))
					tkitResultsStrings = append(tkitResultsStrings, convertValue2String(tkitResult))
				}
			}

			isSucc := true
			isSucc = assert.Exactly(t, len(tkitResults), len(results))
			for i, tkitResultRV := range tkitResults {
				resultRV := results[i]
				if tkitResultRV.Kind() == reflect.Func {
					tmp := tkitResultRV.Call([]reflect.Value{reflect.ValueOf(t), resultRV})
					isSucc = isSucc && assert.Exactly(t, true, tmp[0].Interface())
				} else {
					isSucc = isSucc && assert.Exactly(t, tkitResultRV.Interface(), resultRV.Interface())
				}
			}

			if !isSucc {
				fffName := runtime.FuncForPC(fff.Pointer()).Name()
				fffNameParts := strings.Split(fffName, ".")
				fffName = fffNameParts[len(fffNameParts)-1]

				resultsStrings := make([]string, 0, len(results))
				for _, result := range results {
					resultsStrings = append(resultsStrings, convertValue2String(result.Interface()))
				}

				fmt.Printf("!!! func %s%s => %s <> %s\n",
					fffName,
					"("+strings.Join(tkitArgsStrings, ",")+")",
					"["+strings.Join(resultsStrings, ",")+"]",
					"["+strings.Join(tkitResultsStrings, ",")+"]",
				)
			}
		}
	}
}

func convertValue2String(iv any) (ret string) {
	switch v := iv.(type) {
	case string:
		ret = `"` + v + `"`
	case bool:
		ret = fmt.Sprintf("%v", v)
	case int:
		ret = fmt.Sprintf("%d", v)
	default:
		ret = fmt.Sprintf("%q", v)
	}
	return
}
