/*
Copyright 2016 The Kubernetes Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"runtime"
	"testing"
)

func TestVerifyCmd(t *testing.T) {

	statExe := "stat"
	statPathMsg := "no such file or directory"
	statFileMsg := statPathMsg
	if runtime.GOOS == "windows" {
		statExe = "GetFileAttributesEx"
		statPathMsg = "The system cannot find the path specified."
		statFileMsg = "The system cannot find the file specified."
	}

	tests := []struct {
		name      string
		cmd       string
		expect    string
		wantError bool
	}{
		{
			name:      "verify requires a chart",
			cmd:       "verify",
			expect:    "\"helm verify\" requires 1 argument\n\nUsage:  helm verify PATH [flags]",
			wantError: true,
		},
		{
			name:      "verify requires that chart exists",
			cmd:       "verify no/such/file",
			expect:    fmt.Sprintf("%s no/such/file: %s", statExe, statPathMsg),
			wantError: true,
		},
		{
			name:      "verify requires that chart is not a directory",
			cmd:       "verify testdata/testcharts/signtest",
			expect:    "unpacked charts cannot be verified",
			wantError: true,
		},
		{
			name:      "verify requires that chart has prov file",
			cmd:       "verify testdata/testcharts/compressedchart-0.1.0.tgz",
			expect:    fmt.Sprintf("could not load provenance file testdata/testcharts/compressedchart-0.1.0.tgz.prov: %s testdata/testcharts/compressedchart-0.1.0.tgz.prov: %s", statExe, statFileMsg),
			wantError: true,
		},
		{
			name:      "verify validates a properly signed chart",
			cmd:       "verify testdata/testcharts/signtest-0.1.0.tgz --keyring testdata/helm-test-key.pub",
			expect:    "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := executeCommand(nil, tt.cmd)
			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, but got none: %q", out)
				}
				if err.Error() != tt.expect {
					t.Errorf("Expected error %q, got %q", tt.expect, err)
				}
				return
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
			if out != tt.expect {
				t.Errorf("Expected %q, got %q", tt.expect, out)
			}
		})
	}
}
