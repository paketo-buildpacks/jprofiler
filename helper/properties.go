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

package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/paketo-buildpacks/libpak/bard"
)

type Properties struct {
	Logger bard.Logger
}

func (p Properties) Execute() (map[string]string, error) {
	if _, ok := os.LookupEnv("BPL_JPROFILER_ENABLED"); !ok {
		return nil, nil
	}

	var err error

	agentPath, ok := os.LookupEnv("BPI_JPROFILER_AGENT_PATH")
	if !ok {
		return nil, fmt.Errorf("$BPI_JPROFILER_AGENT_PATH must be set")
	}

	port := "8849"
	if s, ok := os.LookupEnv("BPL_JPROFILER_PORT"); ok {
		port = s
	}

	nowait := true
	if s, ok := os.LookupEnv("BPL_JPROFILER_NOWAIT"); ok {
		nowait, err = strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("unable to parse $BPL_DEBUG_SUSPEND\n%w", err)
		}
	}

	s := fmt.Sprintf("JProfiler enabled on port %s", port)
	if !nowait {
		s = fmt.Sprintf("%s, suspended on start", s)
	}
	p.Logger.Info(s)

	var values []string
	if s, ok := os.LookupEnv("JAVA_TOOL_OPTIONS"); ok {
		values = append(values, s)
	}

	s = fmt.Sprintf("-agentpath:%s=port=%s", agentPath, port)
	if nowait {
		s = fmt.Sprintf("%s,nowait", s)
	}

	values = append(values, s)

	return map[string]string{"JAVA_TOOL_OPTIONS": strings.Join(values, " ")}, nil
}
