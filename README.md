Sassy
======

Go language interface for the libsass library.

I was disappointed when the various sass libraries depended on 
libsass being built as a shared library, and being not really
idiomatic Go. This library allows you to collect the scss files,
compile them to css, then serve the files over http.

Missing Functionality
---------------------

Sass formatted files may not be read at this point.
