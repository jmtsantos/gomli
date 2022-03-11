package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

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

func Replace(jsScript string) {
	var (
		v8Target    *v8.Value
		gomliSearch GomliSearch

		compareResult, transResult *v8.Value
		obj                        *v8.Object

		err error
	)

	ctx := initJavascript(jsScript)
	defer ctx.Close()

	// fetch config from JS
	if v8Target, err = ctx.RunScript("GetGomli()", jsScript); err != nil {
		log.WithError(err).Fatal("error getting GomliSearch object from JS")
	}
	if err = json.Unmarshal([]byte(v8Target.String()), &gomliSearch); err != nil {
		log.WithError(err).Fatal("error unmarshalling GomliSearch")
	}
	log.WithFields(log.Fields{
		"rawcalls": len(gomliSearch.RawTargetCalls),
		"calls":    len(gomliSearch.Calls),
	}).Println("Got search params from JS")

	// start search
	var subGraph []*Node
	for _, c := range gomliSearch.Calls {
		subGraph = append(subGraph, appGraph.edges[Node{c.Class}]...)
	}

	log.Printf("Starting call search on %v nodes", len(subGraph))

	for _, nodes := range subGraph {
		for k1, v1 := range app[nodes.String()].Instructions {

			// this is super slow, but versitle
			obj = ctx.Global()
			jsmessage, _ := json.Marshal(v1)
			obj.Set("Message", base64.StdEncoding.EncodeToString(jsmessage))
			if compareResult, err = ctx.RunScript("Compare()", jsScript); err != nil {
				log.WithError(err).Fatal("error executing Compare()")
			}

			if compareResult.Boolean() {

				// get arguments and search back for the declaration of that single variable
				for i := k1; i > 0; i-- {

					instruction := app[nodes.String()].Instructions[i]
					if instruction.OpCode == 0x1A && instruction.Verbs[2] == v1.Verbs[2] {

						instruction.Verbs[3], _ = strconv.Unquote("\"" + instruction.Verbs[3] + "\"")

						// Sending for generic transform
						obj = ctx.Global()
						jsmessage, _ = json.Marshal([]Instruction{v1, instruction})
						obj.Set("Message", base64.StdEncoding.EncodeToString(jsmessage))
						if transResult, err = ctx.RunScript("Transform()", jsScript); err != nil {
							log.WithError(err).Fatal("error executing Transform()")
						}

						log.WithFields(log.Fields{
							"sourceClass": nodes.String(),
							"decrypted":   transResult.String(),
						}).Println("Transformed")

						// Saving
						appCopy := app[nodes.String()]
						appCopy.Instructions[i].Verbs[2] = instruction.Verbs[2] + ","
						appCopy.Instructions[i].Verbs[3] = fmt.Sprintf("\"%s\"", strconv.Quote(transResult.String()))
						appCopy.Instructions[i].Raw = strings.Join(app[nodes.String()].Instructions[i].Verbs[1:], " ")
						appCopy.GenerateRAWSmali()

						app[nodes.String()] = appCopy
						break
					}
				}
			}

			appCopy := app[nodes.String()]
			if err = appCopy.SaveSmali(); err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"class": appCopy.ClassName,
					"path":  appCopy.Path,
				}).Errorln("error saving smali to file")
			}
		}

	}
}
