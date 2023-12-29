package gno

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type command struct {
	name string
	opts []string
}

type gnoDets struct {
	buildDir string
	binName  string
	src      string
	commands []command
}

const separator = string(os.PathSeparator)

func logMsg(msg string, level string) {
	switch {
	case level == "error":
		log.Fatalf("[ERROR] %s\n", msg)
	case level == "info":
		log.Printf("[INFO]  %s\n", msg)
	case level == "warn":
		log.Printf("[WARN]  %s\n", msg)
	default:
		log.Fatalf("[ERROR] %s\n", msg)
	}
}

func New() *gnoDets {
	return &gnoDets{}
}

// Sets up the build location
// Provide the location to put the build artefacts if any
func (g *gnoDets) BootstrapBuild(buildDirLocation string, bin string, source string) {
	g.buildDir = buildDirLocation
	g.binName = bin
	g.src = source
	if len(g.buildDir) == 0 {
		logMsg("Build directory not provided", "error")
	} else {
		logMsg("Creating build directory", "info")
		err := os.Mkdir(g.buildDir, 0o770)
		if err != nil {
			logMsg(err.Error(), "warn")
			logMsg("Skipping build dir creation", "info")
		}
	}
}

// Copy resources to the final build dir
func (g gnoDets) CopyResources(src string) {
	copyDir(src, g.buildDir)
}

// Add commands to be executed
func (g *gnoDets) AddCommand(name string, bg bool, opts ...string) {

	if bg {
		opts = append(opts, "&")
	}

	c := &command{
		name: name,
		opts: opts,
	}
	g.commands = append(g.commands, *c)

}

func runCommands(g gnoDets) {
	if len(g.commands) >= 1 {
		for _, v := range g.commands {
			logMsg("Running command "+v.name, "info")
			cmd := exec.Command(v.name, v.opts...)
			res, err := cmd.Output()
			if err != nil {
				logMsg(err.Error(), "error")
			}
			fmt.Println("\n" + string(res))
		}
	} else {
		logMsg("No commands passed, skipping", "info")
	}
}

func (g gnoDets) Run() {
	buildBinary(g, false)
	runCommands(g)
}

func (g gnoDets) Build() {
	buildBinary(g, true)
	runCommands(g)
}

func buildBinary(g gnoDets, build bool) {
	if build {
		logMsg("Building Binary", "info")
		binLoc := g.buildDir + separator + g.binName
		err := exec.Command("go", "build", "-o", binLoc, g.src).Run()
		if err != nil {
			logMsg(err.Error(), "error")
		}
		m := fmt.Sprintf("Built Binary [%s]", binLoc)
		logMsg(m, "info")
	} else {
		exec.Command("go", "run", g.src).Run()
	}
}

func copyDir(src string, dest string) {
	ms := fmt.Sprintf("Copying [%s] to [%s]", src, dest)
	logMsg(ms, "info")

	if dest == src {
		logMsg("Cannot copy a folder into itself!", "error")
	}

	f, err := os.Open(src)
	if err != nil {
		logMsg(err.Error(), "error")
	}

	file, err := f.Stat()
	if err != nil {
		logMsg(err.Error(), "error")
	}

	if !file.IsDir() {
		msg := file.Name() + " is not a directory!"
		logMsg(msg, "error")
	}

	err = os.Mkdir(dest+separator+src, 0o770)
	if err != nil {
		logMsg(err.Error(), "warn")
		logMsg("Skipping dir creation", "info")
	} else {
		logMsg("Created "+dest+separator+src, "info")
	}

	files, err := os.ReadDir(src)
	if err != nil {
		logMsg(err.Error(), "error")
	}

	for _, f := range files {

		if f.IsDir() {
			copyDir(src+separator+f.Name(), dest+separator+src+f.Name())
		}

		if !f.IsDir() {

			content, err := os.ReadFile(src + separator + f.Name())
			if err != nil {
				logMsg(err.Error(), "error")

			}

			err = os.WriteFile(dest+separator+src+f.Name(), content, 0755)
			if err != nil {
				logMsg(err.Error(), "error")
			}
		}
	}
}
