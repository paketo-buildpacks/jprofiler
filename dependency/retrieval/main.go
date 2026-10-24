// Copyright 2018-2026 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/paketo-buildpacks/libdependency/retrieve"
	"github.com/paketo-buildpacks/libdependency/upstream"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

const (
	id   = "jprofiler"
	name = "JProfiler Agent"
	purl = "pkg:generic/jprofiler-agent"

	changelog = "https://www.ej-technologies.com/jprofiler/changelog"
	download  = "https://download-gcdn.ej-technologies.com/jprofiler"
)

var versionPattern = regexp.MustCompile(`Release (\d+)\.(\d+)\.(\d+)`)

func main() {
	retrieve.NewMetadata(id, getAllVersions, generateMetadata)
}

func getAllVersions() (versionology.VersionFetcherArray, error) {
	resp, err := http.Get(changelog)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch %s\n%w", changelog, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to fetch %s: status code %d", changelog, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s\n%w", changelog, err)
	}

	var versions versionology.VersionFetcherArray
	seen := map[string]bool{}
	for _, m := range versionPattern.FindAllStringSubmatch(string(body), -1) {
		versionString := fmt.Sprintf("%s.%s.%s", m[1], m[2], m[3])
		if seen[versionString] {
			continue
		}
		seen[versionString] = true

		version, err := semver.NewVersion(versionString)
		if err != nil {
			fmt.Printf("Skipping %s: unable to parse version\n", versionString)
			continue
		}

		versions = append(versions, versionology.NewSimpleVersionFetcher(version))
	}

	return versions, nil
}

func generateMetadata(versionFetcher versionology.VersionFetcher) ([]versionology.Dependency, error) {
	versionString := versionFetcher.Version().String()

	uri := fmt.Sprintf("%s/jprofiler_linux_%s.tar.gz",
		download, regexp.MustCompile(`\.`).ReplaceAllString(versionString, "_"))
	checksum, err := upstream.GetSHA256OfRemoteFile(uri)
	if err != nil {
		return nil, fmt.Errorf("unable to checksum %s\n%w", uri, err)
	}

	dependency := cargo.ConfigMetadataDependency{
		Checksum: fmt.Sprintf("sha256:%s", checksum),
		CPE:      fmt.Sprintf("cpe:2.3:a:jprofiler:agent:%s:*:*:*:*:*:*:*", versionString),
		ID:       id,
		Licenses: []interface{}{
			map[string]string{
				"uri": "https://www.ej-technologies.com/buy/jprofiler/licensing.html",
			},
		},
		Name:    name,
		PURL:    fmt.Sprintf("%s@%s", purl, versionString),
		SHA256:  checksum,
		Stacks:  []string{"*"},
		URI:     uri,
		Version: versionString,
	}

	return versionology.NewDependencyArray(dependency, "")
}
