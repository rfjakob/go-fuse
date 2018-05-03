// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/hanwen/go-fuse/fuse"
)

// Check that loopback Utimens() works as expected.
// Called by TestLoopbackFileUtimens and TestLoopbackFileSystemUtimens.
//
// Parameters:
//   path ........ path to the backing file
//   utimensFn ... Utimens() function that acts on the backing file
func TestLoopbackUtimens(t *testing.T, path string, utimensFn func(atime *time.Time, mtime *time.Time) fuse.Status) {
	// Arbitrary date: 05/02/2018 @ 7:57pm (UTC)
	t0sec := int64(1525291058)

	// Read original timestamp
	var st syscall.Stat_t
	err := syscall.Stat(path, &st)
	if err != nil {
		t.Fatal(err)
	}
	// FromStat handles the differently-named Stat_t fields on Linux and
	// Darwin
	var a1 fuse.Attr
	a1.FromStat(&st)

	// Change atime, keep mtime
	t0 := time.Unix(t0sec, 0)
	status := utimensFn(&t0, nil)
	if !status.Ok() {
		t.Fatal(status)
	}
	err = syscall.Stat(path, &st)
	if err != nil {
		t.Fatal(err)
	}
	var a2 fuse.Attr
	a2.FromStat(&st)
	if a1.Mtime != a2.Mtime {
		t.Errorf("mtime has changed: %v -> %v", a1.Mtime, a2.Mtime)
	}
	if a2.Atime != uint64(t0.Unix()) {
		t.Errorf("wrong atime")
	}

	// Change mtime, keep atime
	t1 := time.Unix(t0sec+123, 0)
	status = utimensFn(nil, &t1)
	if !status.Ok() {
		t.Fatal(status)
	}
	err = syscall.Stat(path, &st)
	if err != nil {
		t.Fatal(err)
	}
	var a3 fuse.Attr
	a3.FromStat(&st)
	if a2.Atime != a3.Atime {
		t.Errorf("atime has changed: %v -> %v", a2.Atime, a3.Atime)
	}
	if a3.Mtime != uint64(t1.Unix()) {
		t.Errorf("wrong mtime")
	}

	// Change both mtime and atime
	ta := time.Unix(t0sec+456, 0)
	tm := time.Unix(t0sec+789, 0)
	status = utimensFn(&ta, &tm)
	if !status.Ok() {
		t.Fatal(status)
	}
	err = syscall.Stat(path, &st)
	if err != nil {
		t.Fatal(err)
	}
	var a4 fuse.Attr
	a4.FromStat(&st)
	if a4.Atime != uint64(ta.Unix()) {
		t.Errorf("wrong atime")
	}
	if a4.Mtime != uint64(tm.Unix()) {
		t.Errorf("wrong mtime")
	}
}

// Check that loopbackFile.Utimens() works as expected
func TestLoopbackFileUtimens(t *testing.T) {
	f2, err := ioutil.TempFile("", "TestLoopbackFileUtimens")
	if err != nil {
		t.Fatal(err)
	}
	path := f2.Name()
	defer os.Remove(path)
	defer f2.Close()

}
