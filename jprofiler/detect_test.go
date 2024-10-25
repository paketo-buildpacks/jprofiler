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
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/jprofiler/v4/jprofiler"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect jprofiler.Detect
	)

	it("fails without BP_JPROFILER_ENABLED", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{Pass: false}))
	})

	it("fails with BP_JPROFILER_ENABLED set to false", func() {
		Expect(os.Setenv("BP_JPROFILER_ENABLED", "false")).To(Succeed())
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{Pass: false}))
	})

	context("$BP_JPROFILER_ENABLED", func() {
		it.Before(func() {
			Expect(os.Setenv("BP_JPROFILER_ENABLED", "true")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_JPROFILER_ENABLED")).To(Succeed())
		})

		it("passes with BP_DEBUG_ENABLED", func() {
			Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
				Pass: true,
				Plans: []libcnb.BuildPlan{
					{
						Provides: []libcnb.BuildPlanProvide{
							{Name: "jprofiler"},
						},
						Requires: []libcnb.BuildPlanRequire{
							{Name: "jprofiler"},
							{Name: "jvm-application"},
						},
					},
				},
			}))
		})
	})
}
