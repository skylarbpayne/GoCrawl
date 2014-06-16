GoCrawl
=======

A concurrent web crawler in Go.

You *should* be able to clone this repository into your Go workspace and use it. However, note that it is a work in progress.

Once released as a package, this concurrent crawler will support definitions of your own page parser/fetcher and multithreading enabling this package to suit most if not all web crawling needs.

Left to do:
	-Debug some no host errors (has to do with allowing too many connections at once)
	-Implement multithreading

If you have any questions or suggestions, feel free to contact me.
