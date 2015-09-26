// Package notmuch provides a Go language binding to the notmuch mail library.
//
// The design is similar enough to the underlying C library that familiarity with
// one will inform the other. There are some differences, however:
//
// * The Go binding arranges for the garbage collector to deal with objects
//   allocated by notmuch correctly. You should close the database manually,
//   but everything else will be garbage collected when it becomes unreachable,
//   and not before. Objects hold references to their parent objects to make this
//   go smoothly.
// * If notmuch returns NULL because of an out of memory error, Go will panic, as
//   it does with other out of memory errors.
// * Some of the names have been shortened or made more idiomatic. The documentation
//   indends to make it obvious when this is the case.
// * Functions which create a child object from a parent object are methods on the
//   parent object, rather than stand-alone functions.
// * Functions which in C return a status code and pass back a value via a pointer
//   argument now return a (value, error) pair.
package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.
