package util

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/leopardslab/dunner/internal/logger"
)

func init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	HomeDir = user.HomeDir
}

// HomeDir is the environment variable HOME
var HomeDir string
var log = logger.Log

var userDir = os.Getenv("user")

// progressReader is for indicating the download / upload progress on the console
type progressReader struct {
	io.Reader
	bytesTransfered   int64
	totalBytes        int64
	progress          float64
	progressDisplayed bool
}

// Read overrides the underlying io.Reader's Read method.
// io.Copy() will be calling this method.
func (w *progressReader) Read(p []byte) (int, error) {
	n, err := w.Reader.Read(p)
	if n > 0 {
		w.bytesTransfered += int64(n)
		percent := float64(w.bytesTransfered) * float64(100) / float64(w.totalBytes)
		if percent-w.progress > 4 {
			fmt.Print(".")
			w.progress = percent
			w.progressDisplayed = true
		}
	}
	return n, err
}

// DirExists returns true if the given param is a valid existing directory
func DirExists(dir string) bool {
	if strings.HasPrefix(dir, "~") {
		dir = path.Join(HomeDir, strings.Trim(dir, "~"))
	}
	src, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return src.IsDir()
}

// FileExists checks if the given file exists
// The argument can be full path of file or relative to working directory
func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// Download makes the http call to given URL and saves to filename given in current directory
func Download(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error downloading file: %s.\n%s", url, resp.Status)
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	progressReader := &progressReader{Reader: resp.Body, totalBytes: resp.ContentLength}
	_, err = io.Copy(out, progressReader)
	if progressReader.progressDisplayed {
		fmt.Println()
	}
	return nil
}

// GetURLContents gets the contents in given URL
func GetURLContents(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error downloading file %s: %s", url, resp.Status)
	}
	defer resp.Body.Close()
	contents, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("Failed to read file contents: %s", readErr.Error())
	}
	return contents, nil
}

// ExecuteSystemCommand executes the given system command in the working directory.
func ExecuteSystemCommand(command []string, outputStreamWriter io.Writer, errorStreamWriter io.Writer) (*exec.Cmd, error) {
	cmd := prepareCommand(command, outputStreamWriter, errorStreamWriter)
	err := cmd.Start()
	return cmd, err
}

func prepareCommand(command []string, outputStreamWriter io.Writer, errorStreamWriter io.Writer) *exec.Cmd {
	cmd := getExecutableCommand(command...)
	wd, _ := os.Getwd()
	cmd.Dir = wd
	cmd.Stdout = outputStreamWriter
	cmd.Stderr = errorStreamWriter
	cmd.Stdin = os.Stdin
	return cmd
}

// getExecutableCommand returns the path of the executable file
func getExecutableCommand(command ...string) *exec.Cmd {
	if len(command) == 0 {
		panic(fmt.Errorf("Invalid executable command"))
	}
	cmd := &exec.Cmd{Path: command[0]}
	if len(command) > 1 {
		cmd = exec.Command(command[0], command[1:]...)
		cmd.Args = append([]string{command[0]}, command[1:]...)
	} else {
		cmd = exec.Command(command[0])
		cmd.Args = append([]string{command[0]})
	}
	return cmd
}

// ShowLoadingMessage is qn util function to show an inline loading message while the process is being carried out.
// This MUST be run in a separate goroutine than the process.
func ShowLoadingMessage(loadingMsg string, finalLog string, done *chan bool, show *chan bool) {
	ticker := time.Tick(time.Second / 2)
	busyChars := []string{`-`, `\`, `|`, `/`}
	x := 0
loop:
	for true {
		select {
		case stop := <-*done:
			if stop {
				break loop
			}
		default:
			if flag.Lookup("test.v") == nil {
				x %= 4
				<-ticker
				fmt.Printf("\r%s... %s",
					loadingMsg,
					busyChars[x],
				)
				x++
			}
		}
	}
	fmt.Print("\r")
	log.Info(finalLog)
	if show != nil {
		*show <- true
	}
	return
}
