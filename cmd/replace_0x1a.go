package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	v8 "rogchap.com/v8go"
)

func ReplaceString(jsScript string) {
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
	gomliSearch.Calls = DuplicateClasses(gomliSearch.Calls)
	log.WithFields(log.Fields{
		"calls": len(gomliSearch.Calls),
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
							"sourceClass":  nodes.String(),
							"sourceMethod": instruction.Method,
							"decrypted":    transResult.String(),
						}).Println("Transformed")

						// Saving
						appCopy := app[nodes.String()]
						appCopy.Instructions[i].Verbs[2] = instruction.Verbs[2] + ","

						// hack to fix escape sequences
						var tempStr string
						if tempStr, err = strconv.Unquote("\"" + transResult.String() + "\""); err != nil {
							log.WithError(err).Errorln("error unqouting transform")
							tempStr = transResult.String()
						}
						appCopy.Instructions[i].Verbs[3] = fmt.Sprintf("\"%s\"", tempStr)
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
