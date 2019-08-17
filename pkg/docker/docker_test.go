package docker

import (
	"fmt"
	"testing"

	"context"

	"github.com/docker/docker/client"
	"github.com/leopardslab/dunner/internal/settings"
	"github.com/spf13/viper"
)

func TestExecWithInvalidImageName(t *testing.T) {
	imageName := "^&^(^(*_invalid"
	step := Step{Image: imageName}

	err := step.Exec()

	expectedErr := fmt.Sprintf("Failed to pull image %s: invalid reference format", imageName)
	if err == nil || err.Error() != expectedErr {
		t.Fatalf("expected error: %s, got: %s", expectedErr, err)
	}
}

func ExampleStep_Exec() {
	settings.Init()
	var testNodeVersion = "10.15.0"
	step := &Step{
		Task:     "test",
		Name:     "node",
		Image:    "node:" + testNodeVersion,
		Commands: [][]string{{"node", "--version"}},
		Env:      nil,
		Volumes:  nil,
	}

	err := step.Exec()
	if err != nil {
		panic(err)
	}
	// Output: OUT: v10.15.0
}

func ExampleStep_workingDirAbs() {
	var testNodeVersion = "10.15.0"
	var absPath = "/go"
	err := runCommand([]string{"pwd"}, absPath, testNodeVersion)

	if err != nil {
		panic(err)
	}
	// Output: OUT: /go
}

func Example_workingDirRel() {
	var testNodeVersion = "10.15.0"
	var relPath = "./"
	err := runCommand([]string{"pwd"}, relPath, testNodeVersion)
	if err != nil {
		panic(err)
	}
	// Output: OUT: /dunner
}

func runCommand(command []string, dir string, nodeVer string) error {
	settings.Init()
	step := &Step{
		Task:    "test",
		Name:    "node",
		Image:   "node:" + nodeVer,
		Command: command,
		Env:     nil,
		Volumes: nil,
		WorkDir: dir,
	}

	return step.Exec()
}

func TestStep_execWithErr(t *testing.T) {
	var testNodeVersion = "10.15.0"
	var relPath = "./"
	err := runCommand([]string{"ls", "/invalid_dir" +
		""}, relPath, testNodeVersion)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
	expectedErr := "Command execution failed with exit code 2"
	if err.Error() != expectedErr {
		t.Errorf("expected error: %s, got: %s", expectedErr, err.Error())
	}
}

func TestStepExecSuccess(t *testing.T) {
	var testNodeVersion = "10.15.0"

	err := runCommand([]string{"node", "--version"}, "./", testNodeVersion)

	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
}

func ExampleStep_execDryRun() {
	dryRun := viper.GetBool("Dry-run")
	viper.Set("Dry-run", true)

	defer viper.Set("Async", dryRun)
	var testNodeVersion = "10.15.0"
	var relPath = "./"
	err := runCommand([]string{"ls", "/invalid_dir" +
		""}, relPath, testNodeVersion)
	if err != nil {
		panic(err)
	}
	// Output:
}

func TestCheckImageExist_local(t *testing.T) {
	testImg := "dunner/test-image"
	check, err := checkImage(testImg, true)
	if err != nil {
		t.Fatal(err)
	}
	if !check {
		t.Fatal("Prebuilt image could not be identified")
	}
}

func TestCheckImageExist_notPresent(t *testing.T) {
	testImg := "random-image"
	check, err := checkImage(testImg, true)
	if err != nil {
		t.Fatal(err)
	}
	if check {
		t.Fatal("Wrong identification, result is false positive")
	}
}

func TestCheckImageExist_invalid(t *testing.T) {
	testImg := "random-image:tag:invalid:format"
	_, err := checkImage(testImg, true)
	expectedErr := fmt.Errorf(`docker: incorrect format for image name`)
	if err == nil || err.Error() != expectedErr.Error() {
		t.Fatal("Wrong image name format did not return appropriate error")
	}
}

func checkImage(img string, notag bool) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	cli.NegotiateAPIVersion(ctx)
	return CheckImageExist(ctx, cli, img, notag)
}
