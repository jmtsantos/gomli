package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func NormalizePackageName(name string) string {
	className := strings.ReplaceAll(name, "/", ".")
	className = strings.TrimSuffix(className, ";")
	className = strings.TrimPrefix(className, "L")
	return className
}

// public getPackageName()Ljava/lang/String;
func NormalizeFunctionName(name string) string {
	// methodName := strings.ReplaceAll(name, "(", "")
	// methodName = strings.ReplaceAll(methodName, ")", "")

	methodName := strings.Split(name, "(")[0]

	return methodName
}

func ReadDirectory(directory string) {
	var (
		parsedSmali Smali
		err         error
	)

	log.Println("Reading directory", directory)

	app = make(map[string]Smali)

	if err = filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if strings.HasSuffix(path, ".smali") {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				parsedSmali = parseSmali(path)
				app[parsedSmali.ClassName] = parsedSmali
			}
			return err

		}); err != nil {
		log.Fatal(err)
	}
	log.Println("Done reading")

}

func parseSmali(path string) (smali Smali) {
	var (
		fileContentBty []byte
		err            error
	)

	if fileContentBty, err = ioutil.ReadFile(path); err != nil {
		log.WithField("path", path).Fatal("error reading file", err)
	}

	smali.Path = path
	smali.Raw = string(fileContentBty)
	smali.RawSlc = strings.Split(string(fileContentBty), "\n")

	smali.ParseClassName()
	smali.ParseProperties()
	smali.ParseMethods()

	fillGraph(smali)

	return smali
}

func fillGraph(smli Smali) {
	appGraph.AddNode(&Node{smli.ClassName})

	for _, instruction := range smli.Instructions {
		if instruction.OpCode == 0x71 {
			appGraph.AddEdge(&Node{smli.ClassName}, &Node{instruction.Verbs[3]})
		}
	}
}
