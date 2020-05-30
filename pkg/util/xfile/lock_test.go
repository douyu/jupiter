// Copyright 2020 Douyu
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

package xfile

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

var lockFilename = flag.String("lockfile", "", "File to lock.  Non-empty value pimples child process.")

func spawn(prog, filename string) ([]byte, error) {
	return exec.Command(prog, "-lockfile", filename, "-test.v",
		"-test.run=TestLock$").CombinedOutput()
}

// TestLock locks a file, spawns a second process that attempts to grab the
// lock to verify it fails.
// Then it closes the lock, and spawns a third copy to verify it can be
// relocked.
func TestLock(t *testing.T) {
	child := *lockFilename != ""
	var filename string
	if child {
		filename = *lockFilename
	} else {
		f, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatal(err)
		}
		filename = f.Name()
		defer os.Remove(filename)
	}

	// Avoid truncating an existing, non-empty file.
	fi, err := os.Stat(filename)
	if err == nil && fi.Size() != 0 {
		t.Fatalf("The file %s is not empty", filename)
	}

	t.Logf("Locking %s\n", filename)
	lock, err := Lock(filename)
	if err != nil {
		t.Fatalf("Could not lock %s: %v", filename, err)
	}

	if !child {
		t.Logf("Spawning child, should fail to grab lock.")
		out, err := spawn(os.Args[0], filename)
		if err == nil {
			t.Fatalf("Attempt to grab open lock should have failed.\n%s", out)
		}
		if !bytes.Contains(out, []byte("Could not lock")) {
			t.Fatalf("Child failed with unexpected output: %s\n", out)
		}
		t.Logf("Child failed to grab lock as expected.")
	}

	t.Logf("Unlocking %s", filename)
	if err := lock.Close(); err != nil {
		t.Fatalf("Could not unlock %s: %v", filename, err)
	}

	if !child {
		t.Logf("Spawning child, should successfully grab lock.")
		if out, err := spawn(os.Args[0], filename); err != nil {
			t.Fatalf("Attempt to re-open lock should have succeeded: %v\n%s",
				err, out)
		}
		t.Logf("Child grabbed lock.")
	}
}
