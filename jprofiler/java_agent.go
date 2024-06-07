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

package jprofiler

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
)

type JavaAgent struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJavaAgent(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (JavaAgent, libcnb.BOMEntry) {
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return JavaAgent{LayerContributor: contributor}, entry
}

func (j JavaAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Expanding to %s", layer.Path)

		if err := crush.Extract(artifact, layer.Path, 1); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand JProfiler\n%w", err)
		}

		arch := archForAgent()
		layer.LaunchEnvironment.Default("BPI_JPROFILER_AGENT_PATH",
			filepath.Join(layer.Path, "bin", fmt.Sprintf("linux-%s", arch), "libjprofilerti.so"))

		return layer, nil
	})
}

func (j JavaAgent) Name() string {
	return j.LayerContributor.LayerName()
}

func archForAgent() string {
	archFromEnv, ok := os.LookupEnv("BP_ARCH")
	if !ok {
		archFromEnv = runtime.GOARCH
	}

	if archFromEnv == "amd64" {
		return "x64"
	} else if archFromEnv == "arm64" {
		return "aarch64"
	} else {
		return archFromEnv
	}
}
