package main

import (
	"flag"
	"fmt"
	"github.com/collinvandyck/gesture/core"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type PluginInfo struct {
	Name string
	Path string
}

// Generates a gesture.go script that you can use to build your very own loving,
// caring, Go-powered IRC bot.
//
// Bot generation is done through some simple string templating. While, in the
// future, it may make sense for this script to validate each specified plugin
// package, we'll leave that to the compiler for now.
func main() {
	templateFile := flag.String("template", "", "The path to the robot template to use")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Printf("usage: build_bot [-template template_file] config\n")
		return
	}
	configFile := args[0]

	config, err := core.ReadConfig(configFile)
	if err != nil {
		fmt.Printf("error: failed to read config file: %s\n", err)
		return
	}
	pluginInfo := parsePluginInfo(config)

	var t *template.Template
	if *templateFile != "" {
		t, err = loadTemplate(*templateFile)
		if err != nil {
			fmt.Printf("error loading template: %s\n", err)
			return
		}
	} else {
		t, err = template.New("gesture").Parse(defaultTemplate)
		if err != nil {
			fmt.Printf("error parsing default template: %s\n", err)
			return
		}
	}
	t.Execute(os.Stdout, pluginInfo)
}

func parsePluginInfo(config *core.Config) []PluginInfo {
	info := make([]PluginInfo, len(config.Plugins))

	i := 0
	for pluginPath, _ := range config.Plugins {
		info[i] = PluginInfo{ packageName(pluginPath), pluginPath }
		i++
	}
	
	return info
}

func packageName(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx == -1 {
		return path
	}
	return path[idx+1:]
}

func loadTemplate(filename string) (*template.Template, error) {
	// Otherwise, load the template from a file
	// TODO: This should be able to be done in a call to template.ReadFiles.
	// Apparently it returns an empty template right now though? WTF?
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	t, err := template.New("gesture").Parse(string(body))
	if err != nil {
		return nil, err
	}

	return t, nil	
}


const defaultTemplate = 
`package main

import (
	"flag"
	"fmt"
	"github.com/collinvandyck/gesture/core"
	"log"{{range .}}
	"{{.Path}}"{{end}}
	"math/rand"
	"time"
)

const banner = "\n\n" +
" / ___| ___  ___| |_ _   _ _ __ ___  \n" + 
" | |  _ / _ \\/ __| __| | | | '__/ _ \\ \n" +
" | |_| |  __/\\__ \\ |_| |_| | | |  __/\n" +
"  \\____|\\___||___/\\__|\\__,_|_|  \\___|\n"


func loadPlugins(bot *core.Gobot) {{"{"}}{{range .}}
	{{.Name}}.Create(bot){{end}}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Parse command-line arguments in logging package
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalln("usage: gesture conf_file")
	}

	config, err := core.ReadConfig(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	bot := core.CreateGobot(config)
	loadPlugins(bot)
	fmt.Println(banner)
	quit, err := bot.Connect(config.Hostname)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}

	<-quit
}
`
