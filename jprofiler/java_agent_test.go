/*
 * Copyright 2018-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jprofiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/jprofiler/v4/jprofiler"
)

func testJavaAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Layers.Path = t.TempDir()
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	context("BP_ARCH=amd64", func() {
		it.Before(func() {
			t.Setenv("BP_ARCH", "amd64")
		})

		it("contributes Java agent", func() {
			dep := libpak.BuildpackDependency{
				URI:    "https://localhost/stub-jprofiler-agent.tar.gz",
				SHA256: "9ec6fd679560481ff82d59397ffa289028e2c68df41802d172b35884b84b304d",
			}
			dc := libpak.DependencyCache{CachePath: "testdata"}

			j, _ := jprofiler.NewJavaAgent(dep, dc)
			layer, err := ctx.Layers.Layer("test-layer")
			Expect(err).NotTo(HaveOccurred())

			layer, err = j.Contribute(layer)
			Expect(err).NotTo(HaveOccurred())

			Expect(layer.Launch).To(BeTrue())
			Expect(filepath.Join(layer.Path, "fixture-marker")).To(BeARegularFile())
			Expect(layer.LaunchEnvironment["BPI_JPROFILER_AGENT_PATH.default"]).To(
				Equal(filepath.Join(layer.Path, "bin", "linux-x64", "libjprofilerti.so")))
		})
	})

	context("BP_ARCH=arm64", func() {
		it.Before(func() {
			t.Setenv("BP_ARCH", "arm64")
		})

		it("contributes Java agent", func() {
			dep := libpak.BuildpackDependency{
				URI:    "https://localhost/stub-jprofiler-agent.tar.gz",
				SHA256: "9ec6fd679560481ff82d59397ffa289028e2c68df41802d172b35884b84b304d",
			}
			dc := libpak.DependencyCache{CachePath: "testdata"}

			j, _ := jprofiler.NewJavaAgent(dep, dc)
			layer, err := ctx.Layers.Layer("test-layer")
			Expect(err).NotTo(HaveOccurred())

			layer, err = j.Contribute(layer)
			Expect(err).NotTo(HaveOccurred())

			Expect(layer.Launch).To(BeTrue())
			Expect(filepath.Join(layer.Path, "fixture-marker")).To(BeARegularFile())
			Expect(layer.LaunchEnvironment["BPI_JPROFILER_AGENT_PATH.default"]).To(
				Equal(filepath.Join(layer.Path, "bin", "linux-aarch64", "libjprofilerti.so")))
		})
	})
}
