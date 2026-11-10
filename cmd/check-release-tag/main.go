/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Command check-release-tag checks that a release tag is for the DCGM
// version pinned in hack/versions.env. Release tags are
// v1.<DCGM version digits><package revision>.<patch>: DCGM 4.7.1 with package
// revision 2 gives v1.4712.0, v1.4712.1 and so on.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// number matches a decimal number without leading zeros, as Go module
// versions and the encoded DCGM pin need.
const number = `(0|[1-9]\d*)`

var (
	dcgmVersionPattern     = regexp.MustCompile(`^` + number + `\.` + number + `\.` + number + `$`)
	packageRevisionPattern = regexp.MustCompile(`^` + number + `$`)
)

// main reads the tag and the DCGM pin from flags and exits with status 1 when
// the tag does not match the pin, or 2 when the arguments are wrong.
func main() {
	tag := flag.String("tag", "", "release tag to check, such as v1.4712.3")
	dcgmVersion := flag.String("dcgm-version", "", "pinned DCGM version, such as 4.7.1")
	packageRevision := flag.String("dcgm-package-revision", "", "pinned DCGM package revision, such as 2")
	flag.Parse()

	if flag.NArg() != 0 || *tag == "" || *dcgmVersion == "" || *packageRevision == "" {
		fmt.Fprintln(os.Stderr, "usage: check-release-tag -tag <tag> -dcgm-version <X.Y.Z> -dcgm-package-revision <N>")
		os.Exit(2)
	}

	if err := checkReleaseTag(*tag, *dcgmVersion, *packageRevision); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	fmt.Printf("Release tag %s matches DCGM %s-%s\n", *tag, *dcgmVersion, *packageRevision)
}

// releaseTagPattern returns the pattern that release tags for DCGM
// dcgmVersion (X.Y.Z) with package revision packageRevision must match. It
// fails when either value has the wrong form.
func releaseTagPattern(dcgmVersion, packageRevision string) (*regexp.Regexp, error) {
	if !dcgmVersionPattern.MatchString(dcgmVersion) {
		return nil, fmt.Errorf("DCGM version %q is not X.Y.Z without leading zeros", dcgmVersion)
	}
	if !packageRevisionPattern.MatchString(packageRevision) {
		return nil, fmt.Errorf("DCGM package revision %q is not a number without leading zeros", packageRevision)
	}

	minor := strings.ReplaceAll(dcgmVersion, ".", "") + packageRevision

	return regexp.MustCompile(`^v1\.` + minor + `\.` + number + `$`), nil
}

// checkReleaseTag returns an error that names the expected form when tag is
// not a release tag for DCGM dcgmVersion with package revision
// packageRevision.
func checkReleaseTag(tag, dcgmVersion, packageRevision string) error {
	pattern, err := releaseTagPattern(dcgmVersion, packageRevision)
	if err != nil {
		return err
	}

	if !pattern.MatchString(tag) {
		return fmt.Errorf("release tag %s does not match %s for DCGM %s-%s",
			tag, pattern, dcgmVersion, packageRevision)
	}

	return nil
}
