package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/debug"
	"github.com/puppetlabs/relay/pkg/integration/container/def"
	"github.com/puppetlabs/relay/pkg/integration/container/generator"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		// In all cases, any error here (most likely) means that the file doesn't
		// exist. There are some failure modes that are getting ignored but
		// they're, like, really not possible and catastrophic in nature.
		return false
	}

	return true
}

func findContainerFile(path string) string {
	if fileExists(filepath.Join(path, "container.yml")) {
		return filepath.Join(path, "container.yml")
	}

	if fileExists(filepath.Join(path, "container.yaml")) {
		return filepath.Join(path, "container.yaml")
	}

	return ""
}

type stepBuildFunc func(context, containerFile string) error

func forEachContainer(path string, cb stepBuildFunc) error {
	entries, err := ioutil.ReadDir(path)

	if err != nil {
		debug.Logf("failed to list directories when looking for containers: %v", err)
		return err
	}

	for _, info := range entries {
		// if we find a directory we need to look for a container.yml file inside
		// of it. if we find that, then it's a candidate for being built.
		if info.IsDir() {
			dir := filepath.Join(path, info.Name())
			file := findContainerFile(dir)

			if file != "" {
				cb(dir, file)
			}
		}
	}

	return nil
}

func buildContainer(dir, path string) error {
	// TODO: We should pull much of this out of Spindle/Nebula SDK.
	containerDef, err := def.NewFromFilePath(path)

	if err != nil {
		return err
	}

	gen := generator.New(
		containerDef.Container,
		generator.WithFilesRelativeTo(def.NewFileRef(dir)),

		// TODO: This should come from the configuration and default to
		// "puppetlabs" as far as I can tell.
		generator.WithRepoNameBase("puppetlabs"),
	)

	files, err := gen.Files()

	if err != nil {
		debug.Logf("failed to generate files: %v", err)
		return err
	}

	tmpdir, err := ioutil.TempDir("", "relay-integration-build")

	if err != nil {
		debug.Logf("failed to create a tempdir: %v", err)
		return err
	}

	dockerfile := filepath.Join(tmpdir, "Dockerfile")
	if err := ioutil.WriteFile(dockerfile, []byte(files[1].Content), 0644); err != nil {
		debug.Logf("failed to write Dockerfile to %s: %v", path, err)
		return err
	}

	cmd := exec.Command("docker", "build", "--file", dockerfile, dir)
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		debug.Logf("failed to get start docker build command: %v", err)
		return err
	}

	return nil
}

func Build(path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := ReadConfig(f); err != nil {
		return err
	} else {
		// This is the root directory of the integration. There will be standard
		// directories underneath this.
		dir := filepath.Dir(path)
		actionsdir := filepath.Join(dir, "actions")
		stepsdir := filepath.Join(actionsdir, "steps")

		// TODO: add support for building triggers?
		err = forEachContainer(stepsdir, func(context, containerYaml string) error {
			if err := buildContainer(context, containerYaml); err != nil {
				debug.Logf("error: %v", err)
			}
			return nil
		})

		if err != nil {
			debug.Logf("failed to build some container somewhere: %v", err)
			return err
		}
	}

	return nil
}
