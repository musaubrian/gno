package gno

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		err := os.Mkdir(g.buildDir, 0o770)
		if err != nil {
			logMsg(err.Error(), "warn")
			logMsg("Skipping build dir creation", "info")
		} else {
			logMsg("Created build directory", "info")
		}
	}
}

// Copy resources to the final build dir
func (g gnoDets) CopyResources(src string) {
	copyDir(src, g.buildDir)
}

// Add commands to be executed
func (g *gnoDets) AddCommand(name string, opts ...string) {
	c := &command{
		name: name,
		opts: opts,
	}
	g.commands = append(g.commands, *c)

}

func runCommands(g gnoDets) {
	if len(g.commands) >= 1 {
		for _, v := range g.commands {
			ms := fmt.Sprintf("Running command: `%s`", v.name)
			logMsg(ms, "info")
			res, err := exec.Command(v.name, v.opts...).Output()
			if err != nil {
				logMsg(err.Error(), "error")
			} else {
				fmt.Println(string(res))
			}
		}
	} else {
		logMsg("No commands, skipping", "info")
	}
}

func (g gnoDets) Run() {
	buildBinary(g, false)
	runCommands(g)
}

// Builds the binary and runs the commands Synchronously
// So they need to be ordered correctly
func (g gnoDets) Build() {
	buildBinary(g, true)
	runCommands(g)
}

func buildBinary(g gnoDets, build bool) {
	if build {
		binLoc := filepath.Join(g.buildDir, g.binName)
		res, err := exec.Command("go", "build", "-o", binLoc, g.src).CombinedOutput()
		if err != nil {
			logMsg(string(res), "warn")
			logMsg(err.Error(), "error")
		}
		m := fmt.Sprintf("Built Binary -> %s", binLoc)
		logMsg(m, "info")
	} else {
		exec.Command("go", "run", g.src).Run()
	}
}

func copyDir(src string, dest string) {
	if dest == src {
		logMsg("Cannot copy a folder into itself!", "error")
	}
	files, err := os.ReadDir(src)
	if err != nil {
		logMsg(err.Error(), "error")
	}
	for _, f := range files {
		var destPath string
		srcPath := filepath.Join(src, f.Name())
		if !strings.Contains(dest, strings.Split(srcPath, "/")[0]) {
			destPath = filepath.Join(dest, strings.Split(srcPath, "/")[0], f.Name())
		} else {
			destPath = filepath.Join(dest, f.Name())
		}
		if f.IsDir() {
			copyDir(srcPath, destPath)
		} else {
			content, err := os.ReadFile(srcPath)
			if err != nil {
				logMsg(err.Error(), "error")
			}
			err = os.MkdirAll(filepath.Dir(destPath), 0770)
			if err != nil {
				logMsg(err.Error(), "error")
				logMsg("Skipping dir creation", "info")
			}
			err = os.WriteFile(destPath, content, 0644)
			if err != nil {
				logMsg(err.Error(), "error")
			}
			logMsg(fmt.Sprintf("Copied %s -> %s", srcPath, destPath), "info")
		}
	}
}
