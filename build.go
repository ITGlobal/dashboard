package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/daviddengcn/go-colortext"
	"github.com/mholt/archiver"
)

func main() {
	if len(os.Args) < 2 {
		error("Not enough arguments")
		os.Exit(1)
	}

	dir := os.Args[1]
	wd, err := os.Getwd()
	if err != nil {
		error(fmt.Sprintf("os.Getwd() failed: %s", err))
		os.Exit(1)
	}

	err = os.Chdir(dir)
	if err != nil {
		error(fmt.Sprintf("os.Chdir('%s') failed: %s", dir, err))
		os.Exit(1)
	}

	defer func() {
		os.Chdir(wd)
	}()

	build()
}

// ============================================================================
// Build script
// ============================================================================

func build() {
	usage()

	version := get_version()
	var goOS, goArch string

	if len(os.Args) >= 3 {
		goOS = os.Args[2]
	} else {
		goOS = runtime.GOOS
		message(fmt.Sprintf("auto detected OS type - %s", goOS))
	}

	if len(os.Args) >= 4 {
		goArch = os.Args[3]
	} else {
		goArch = runtime.GOARCH
		message(fmt.Sprintf("auto detected OS arch - %s", goArch))
	}

	run_build(version, goOS, goArch)

	ok()
}

func usage() {
	script := "./build.sh"
	if runtime.GOOS == "windows" {
		script = "build"
	}

	header("*********************")
	header("* dash build script *")
	header("*********************")
	message("Usage:")
	message(fmt.Sprintf("  %s               - compile for current system (best option for debugging)", script))
	message(fmt.Sprintf("  %s windows 386   - compile for Windows x86", script))
	message(fmt.Sprintf("  %s windows amd64 - compile for Windows x64", script))
	message(fmt.Sprintf("  %s linux 386     - compile for Linux x86", script))
	message(fmt.Sprintf("  %s linux amd64   - compile for Linux x64", script))
	message(fmt.Sprintf("  %s linux arm     - compile for Linux ARM6", script))
	message("")
}

func get_version() string {
	header("fetching version number")
	version := executeStr("git", "describe", "--tags")
	message(fmt.Sprintf("version = '%s'", version))
	return version
}

func run_build(version, goOS, goArch string) {

	ext := ""
	if goOS == "windows" {
		ext = ".exe"
	}

	ldflags := fmt.Sprintf("-X main.Version=%s -X main.BuildConfiguration=%s-%s", version, goOS, goArch)
	dir := fmt.Sprintf("./out/%s-%s", goOS, goArch)
	dashd := fmt.Sprintf("%s/dashd%s", dir, ext)
	dasht := fmt.Sprintf("%s/dasht%s", dir, ext)
	pkg := fmt.Sprintf("./out/dash-%s-%s.zip", goOS, goArch)

	os.Setenv("GOOS", goOS)
	os.Setenv("GOARCH", goArch)

	header("")
	header("fetching packages")
	execute("go", "get", "./daemon")
	execute("go", "get", "./terminal")

	header("")
	header(fmt.Sprintf("compiling for %s/%s", goOS, goArch))
	execute("go", "build", "-o", dashd, "-ldflags", ldflags, "./daemon")
	execute("go", "build", "-o", dasht, "-ldflags", ldflags, "./terminal")

	header("")
	header("packaging artifacts")
	zip(pkg, []string{dashd, dasht, "./dashctl"})

	message(fmt.Sprintf("created package %s", pkg))
}

// ============================================================================
// Helper functions
// ============================================================================

func header(msg string) {

	ct.Foreground(ct.Cyan, false)
	fmt.Print(">>> ")
	ct.Foreground(ct.Cyan, true)
	fmt.Print(msg)
	ct.ResetColor()
	fmt.Println()
}

func message(msg string) {
	ct.Foreground(ct.Yellow, false)
	fmt.Print(">>> ")
	ct.Foreground(ct.Yellow, true)
	fmt.Print(msg)
	ct.ResetColor()
	fmt.Println()
}

func printcmd(cmd string, arg ...string) {
	ct.Foreground(ct.White, false)
	fmt.Print(">>> ")
	ct.Foreground(ct.White, true)
	fmt.Printf("%s", cmd)
	for _, c := range arg {
		fmt.Printf(" %s", c)
	}
	ct.ResetColor()
	fmt.Println()
}

func error(msg string) {
	ct.Foreground(ct.Red, false)
	fmt.Print(">>> ")
	ct.Foreground(ct.Red, true)
	fmt.Print(msg)
	ct.ResetColor()
	fmt.Println()
}

func ok() {
	ct.Foreground(ct.Green, false)
	fmt.Print(">>> ")
	ct.Foreground(ct.Green, true)
	fmt.Print("OK")
	ct.ResetColor()
	fmt.Println()
}

func execute(cmd string, arg ...string) {
	printcmd(cmd, arg...)

	command := exec.Command(cmd, arg...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		error(fmt.Sprintf("command '%s' failed: %s", cmd, err))
		os.Exit(-1)
	}
}

func executeStr(cmd string, arg ...string) string {
	printcmd(cmd, arg...)

	command := exec.Command(cmd, arg...)

	out, err := command.CombinedOutput()
	if err != nil {
		error(fmt.Sprintf("command '%s' failed: %s", cmd, err))
		os.Exit(-1)
	}

	str := string(out)
	str = strings.TrimSpace(str)

	return str
}

func zip(zipfile string, files []string) {
	err := archiver.Zip(zipfile, files)
	if err != nil {
		error(fmt.Sprintf("unable to zip '%s' failed: %s", zipfile, err))
		os.Exit(-1)
	}

}
