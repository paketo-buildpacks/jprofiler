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

package helper_test

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/jprofiler/v4/helper"
)

func testProperties(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		p = helper.Properties{}
	)

	it("returns if $BPL_JPROFILER_ENABLED is not set", func() {
		Expect(p.Execute()).To(BeNil())
	})

	context("$BPL_JPROFILER_ENABLED", func() {
		it.Before(func() {
			Expect(os.Setenv("BPL_JPROFILER_ENABLED", "")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BPL_JPROFILER_ENABLED")).To(Succeed())
		})

		it("returns error if $BPI_JPROFILER_AGENT_PATH is not set", func() {
			_, err := p.Execute()
			Expect(err).To(MatchError("$BPI_JPROFILER_AGENT_PATH must be set"))
		})

		context("$BPI_JPROFILER_AGENT_PATH", func() {
			it.Before(func() {
				Expect(os.Setenv("BPI_JPROFILER_AGENT_PATH", "test-path")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("BPI_JPROFILER_AGENT_PATH")).To(Succeed())
			})

			it("contributes configuration", func() {
				Expect(p.Execute()).To(Equal(map[string]string{
					"JAVA_TOOL_OPTIONS": "-agentpath:test-path=port=8849,nowait",
				}))
			})

			context("$BPL_JPROFILER_PORT", func() {
				it.Before(func() {
					Expect(os.Setenv("BPL_JPROFILER_PORT", "8850")).To(Succeed())
				})

				it.After(func() {
					Expect(os.Unsetenv("BPL_JPROFILER_PORT")).To(Succeed())
				})

				it("contributes port configuration from $BPL_JPROFILER_PORT", func() {
					Expect(p.Execute()).To(Equal(map[string]string{
						"JAVA_TOOL_OPTIONS": "-agentpath:test-path=port=8850,nowait",
					}))
				})
			})

			context("$BPL_JPROFILER_NOWAIT", func() {
				it.Before(func() {
					Expect(os.Setenv("BPL_JPROFILER_NOWAIT", "false")).To(Succeed())
				})

				it.After(func() {
					Expect(os.Unsetenv("BPL_JPROFILER_NOWAIT")).To(Succeed())
				})

				it("contributes suspend configuration from $BPL_JPROFILER_NOWAIT", func() {
					Expect(p.Execute()).To(Equal(map[string]string{
						"JAVA_TOOL_OPTIONS": "-agentpath:test-path=port=8849",
					}))
				})
			})

			context("$JAVA_TOOL_OPTIONS", func() {
				it.Before(func() {
					Expect(os.Setenv("JAVA_TOOL_OPTIONS", "test-java-tool-options")).To(Succeed())
				})

				it.After(func() {
					Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
				})

				it("contributes configuration appended to existing $JAVA_TOOL_OPTIONS", func() {
					Expect(p.Execute()).To(Equal(map[string]string{
						"JAVA_TOOL_OPTIONS": "test-java-tool-options -agentpath:test-path=port=8849,nowait",
					}))
				})
			})

		})
	})

}
