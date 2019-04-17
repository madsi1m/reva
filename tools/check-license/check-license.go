// Copyright 2018-2019 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var fix = flag.Bool("fix", false, "add header if not present")

var licenseText = `// Copyright 2018-2019 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

`

var license = regexp.MustCompile(licenseText)

const prefix = "// Copyright "

var skip = map[string]bool{}

func main() {
	flag.Parse()
	err := filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
		if skip[path] {
			return nil
		}

		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		src, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		// Check if license is at the top of the file.
		if !bytes.HasPrefix(src, []byte(prefix)) {
			err := fmt.Errorf("%v: license header not present or not at the top, to fix run: go run tools/check-license/check-license.go -fix", path)
			if *fix == true {
				newSrc := licenseText + string(src)
				ioutil.WriteFile(path, []byte(newSrc), 644)
				src = []byte(newSrc)
			} else {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
