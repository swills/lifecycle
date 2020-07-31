package acceptance

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/lifecycle/acceptance/variables"
	h "github.com/buildpacks/lifecycle/testhelpers"
)

var (
	rootBuilderBinaryDir     = filepath.Join("testdata", "root-builder", "image", "container", "cnb", "lifecycle")
	rootBuilderDockerContext = filepath.Join("testdata", "root-builder", "image")
	rootBuilderImage         = "lifecycle/acceptance/root-builder"
	rootBuilderPath          = "/cnb/lifecycle/root-builder"
	stackpackFixtureDir      = filepath.Join("testdata", "root-builder", "stackpack")
)

func TestRootBuilder(t *testing.T) {
	h.SkipIf(t, runtime.GOOS == "windows", "These tests need to be adapted to work on Windows")
	rand.Seed(time.Now().UTC().UnixNano())

	info, err := h.DockerCli(t).Info(context.TODO())
	h.AssertNil(t, err)
	daemonOS = info.OSType

	// Setup test container

	h.MakeAndCopyLifecycle(t, daemonOS, rootBuilderBinaryDir)
	h.DockerBuild(t,
		rootBuilderImage,
		rootBuilderDockerContext,
		h.WithFlags("-f", filepath.Join(rootBuilderDockerContext, variables.DockerfileName)),
	)
	defer h.DockerImageRemove(t, rootBuilderImage)

	spec.Run(t, "acceptance-root-builder", testRootBuilder, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testRootBuilder(t *testing.T, when spec.G, it spec.S) {
	when("called", func() {
		it("does something", func() {
			h.SkipIf(t, runtime.GOOS == "windows", "Not relevant on Windows")

			output := h.DockerRun(t,
				rootBuilderImage,
				h.WithBash(fmt.Sprintf("%s -group /cnb/group.toml -plan /cnb/plan.toml; tar tvf /layers/example_stack.tgz", rootBuilderPath)),
			)

			h.AssertMatch(t, output, ".wh.sbin")
			h.AssertMatch(t, output, "bin/exe-to-snapshot")
		})
	})
}