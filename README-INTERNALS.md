This document describes some general patterns in the implementation.

# Resource reclamation.

The library is written such that all resources will be released when the
garbage collector claims an object. However, we also export Close()
methods on each type that needs explicit cleanup. This is because the Go
garbage collector [isn't always able to make the right decisions about
objects with pointers to C memory][1].

## Tracking dependencies

Each notmuch object has a corresponding wrapper object:
notmuch_database_t is wrapped by DB, notmuch_query_t is wrapped by Query
and so on. Each of these wrappers is an alias for the type cStruct,
which holds a pointer to the underlying C object, and also to the
wrappers for any objects referenced by the underlying C object (via the
`parent` field).  This keeps the GC from collecting parent objects if
the children are still in use.

## Creating objects

When creating an object, the caller should set the `parent` field to a
pointer to the object's immediate parent object, and set cptr to the
underlying c pointer. Finally, calling setGcClose on the object will
cause it to be released properly by the garbage collector.

## Cleaning up

Calling the `Close()` method on an object a second time is a no-op, and
thus always safe. The primary reason for this is that it makes dealing
with mixing GC and manual reclamation simple.

The Close methods also clear all of their references to parent objects
explicitly. While this isn't strictly necessary, it means the GC
will know that they are unreachable sooner, if that becomes the case.
Per the documentation for `runtime.SetFinalizer`, once the finalizer is
called, the object will stick around until the next time the GC looks at
it. Because of this, it won't otherwise even consider parent objects
until the third pass at least.

In some cases, invoking the `*_destroy` function on an object also
cleans up its children, in which case it becomes unsafe to then invoke
their `*_destroy` function. cStruct handles all of the bookkeeping and
synchronization necessary to deal with this; derived types just need to
make proper use of doClose() in their Close() implementation (see the
comments in cstruct.go)

[1]: https://gist.github.com/dwbuiten/c9865c4afb38f482702e#cleaning
