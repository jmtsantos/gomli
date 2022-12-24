package cmd

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	v8 "rogchap.com/v8go"
)

type GomliSearch struct {
	RawTargetCalls []string
	Calls          []Call
}

func initJavascript(jsScript string) (ctx *v8.Context) {
	var (
		jsScriptBty []byte
		jqueryBty   []byte
		err         error
	)

	// init javascript
	if jqueryBty, err = ioutil.ReadFile("assets/utils.js"); err != nil {
		log.WithError(err).Fatal("error reading jquery file", jsScript)
	}
	if jsScriptBty, err = ioutil.ReadFile(jsScript); err != nil {
		log.WithError(err).Fatal("error reading javascript file", jsScript)
	}

	ctx = v8.NewContext()

	if _, err = ctx.RunScript(string(jqueryBty), "utils.js"); err != nil {
		log.WithError(err).Fatal("error loading jquery")
	}
	if _, err = ctx.RunScript(string(jsScriptBty), jsScript); err != nil {
		log.WithError(err).Fatal("error loading custom script")
	}

	return
}

func DuplicateClasses(intSlice []Call) []Call {
	keys := make(map[string]bool)
	list := []Call{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry.Class]; !value {
			keys[entry.Class] = true
			list = append(list, entry)
		}
	}
	return list
}
