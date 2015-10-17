package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

import (
	"io"
	"runtime"
	"sync"
	"unsafe"
)

// cStruct is the common representation of almost all of our wrapper types.
//
// It does the heavy lifting of interacting with the garbage
// collector/*_destroy functions.
type cStruct struct {
	// A pointer to the underlying c object
	cptr unsafe.Pointer

	// Parent object. Holding a pointer to this in Go-land makes the reference
	// visible to the garbage collector, and thus prevents it from being
	// reclaimed prematurely.
	parent *cStruct

	// readers-writer lock for dealing with mixing manual calls to Close() with
	// GC.
	lock sync.RWMutex
}

// Recursively acquire read locks on this object and all parent objects.
func (c *cStruct) rLock() {
	c.lock.RLock()
	if c.parent != nil {
		c.parent.rLock()
	}
}

// Recursively release read locks this object and all parent objects.
func (c *cStruct) rUnlock() {
	if c.parent != nil {
		c.parent.rUnlock()
	}
	c.lock.RUnlock()
}

// Call f in a context in which it is safe to destroy the underlying C object.
// `f` will only be invoked if the underlying object is still live. When `f`
// is invoked,  The calling goroutine will hold the necessary locks to make
// destroying the underlying object safe.
//
// Typically, wrapper types will use this to implement their Close() methods;
// it handles all of the synchronization bits.
func (c *cStruct) doClose(f func() error) error {
	// Briefly:
	// 1. Acquire a write lock on ourselves.
	// 2. Acquire read locks for all of our ancestors in pre-order (the ordering
	//    is important to avoid deadlocks).
	// 3. Check if we're live, and call f if so.
	// 4. Clear all of our references to other objects, and release the locks
	var err error
	c.lock.Lock()
	if c.parent != nil {
		c.parent.rLock()
	}
	defer func() {
		if c.parent != nil {
			c.parent.rUnlock()
		}
		c.cptr = nil
		c.parent = nil
		c.lock.Unlock()
	}()
	if c.live() {
		err = f()
	}
	return err
}

// Returns true if and only if c's underlying object is live, i.e. neither it
// nor any of its ancestor objects have been finalized.
//
// Note that this method does no synchronization; the caller must separately
// acquire the necessary locks.
func (c *cStruct) live() bool {
	if c.cptr == nil {
		return false
	}
	if c.parent == nil {
		return true
	}
	return c.parent.live()
}

// Set a finalizer to invoke c.Close() when c is garbage collected.
func setGcClose(c io.Closer) {
	runtime.SetFinalizer(c, func(c io.Closer) {
		c.Close()
	})
}
