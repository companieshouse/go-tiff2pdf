/* $Id$ */

/*
 * Copyright (c) 2014 Companies House
 *
 * Permission to use, copy, modify, distribute, and sell this software and
 * its documentation for any purpose is hereby granted without fee, provided
 * that (i) the above copyright notices and this permission notice appear in
 * all copies of the software and related documentation, and (ii) the names of
 * Companies House may not be used in any advertising or publicity relating
 * to the software without the specific, prior written permission of Companies House.
 *
 * THE SOFTWARE IS PROVIDED "AS-IS" AND WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS, IMPLIED OR OTHERWISE, INCLUDING WITHOUT LIMITATION, ANY
 * WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE.
 *
 * IN NO EVENT SHALL COMPANIES HOUSE BE LIABLE FOR ANY SPECIAL, INCIDENTAL,
 * INDIRECT OR CONSEQUENTIAL DAMAGES OF ANY KIND, OR ANY DAMAGES WHATSOEVER
 * RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER OR NOT ADVISED OF THE
 * POSSIBILITY OF DAMAGE, AND ON ANY THEORY OF LIABILITY, ARISING OUT OF OR IN
 * CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

/*
 * TIFF Library Golang-specific Routines.
 */

#include "tif_config.h"

#ifdef HAVE_SYS_TYPES_H
# include <sys/types.h>
#endif

#include <errno.h>

#include <stdarg.h>
#include <stdlib.h>
#include <sys/stat.h>

#ifdef HAVE_UNISTD_H
# include <unistd.h>
#endif

#ifdef HAVE_FCNTL_H
# include <fcntl.h>
#endif

#ifdef HAVE_IO_H
# include <io.h>
#endif

#include "tiffiop.h"

static tmsize_t
_tiffReadProc(thandle_t fd, void* buf, tmsize_t size)
{
	return (tmsize_t)GoTiffReadProc(fd, buf, size);
}

static tmsize_t
_tiffNoop(thandle_t fd, void* buf, tmsize_t size)
{
	return -1;
}

static tmsize_t
_tiffWriteProc(thandle_t fd, void* buf, tmsize_t size)
{
	return GoTiffWriteProc(fd, buf, size);
}

static uint64
_tiffSeekProc(thandle_t fd, uint64 off, int whence)
{
	return (uint64)GoTiffSeekProc(fd, off, whence);
}

static int
_tiffCloseProc(thandle_t fd)
{
	return GoTiffCloseProc(fd);
}

static uint64
_tiffSizeProc(thandle_t fd)
{
	return GoTiffSizeProc(fd);
}

#ifdef HAVE_MMAP
#include <sys/mman.h>

static int
_tiffMapProc(thandle_t fd, void** pbase, toff_t* psize)
{
	return GoTiffMapProc(fd, pbase, psize);
}

static void
_tiffUnmapProc(thandle_t fd, void* base, toff_t size)
{
	GoTiffUnmapProc(fd, base, size);
}
#else /* !HAVE_MMAP */
static int
_tiffMapProc(thandle_t fd, void** pbase, toff_t* psize)
{
	//GoTiffMapProc(fd, pbase, psize);
	(void) fd; (void) pbase; (void) psize;
	return (0);
}

static void
_tiffUnmapProc(thandle_t fd, void* base, toff_t size)
{
	//GoTiffUnmapProc(fd, base, size);
	(void) fd; (void) base; (void) size;
}
#endif /* !HAVE_MMAP */

/*
 * Open a TIFF file descriptor for read/writing.
 */
TIFF*
TIFFFdOpen(int fd, const char* name, const char* mode)
{
	TIFF* tif;
	TIFFReadWriteProc readproc = _tiffReadProc;

	if (0 == strncmp(name+strlen(name)-4, ".pdf", 4)) {
		readproc = _tiffNoop;
	}

	tif = TIFFClientOpen(name, mode,
	    (thandle_t) fd,
	    readproc, _tiffWriteProc,
	    _tiffSeekProc, _tiffCloseProc, _tiffSizeProc,
	    _tiffMapProc, _tiffUnmapProc);
	if (tif)
		tif->tif_fd = fd;
	return (tif);
}

void*
_TIFFmalloc(tmsize_t s)
{
        if (s == 0)
                return ((void *) NULL);

	return (malloc((size_t) s));
}

void
_TIFFfree(void* p)
{
	free(p);
}

void*
_TIFFrealloc(void* p, tmsize_t s)
{
	return (realloc(p, (size_t) s));
}

void
_TIFFmemset(void* p, int v, tmsize_t c)
{
	memset(p, v, (size_t) c);
}

void
_TIFFmemcpy(void* d, const void* s, tmsize_t c)
{
	memcpy(d, s, (size_t) c);
}

int
_TIFFmemcmp(const void* p1, const void* p2, tmsize_t c)
{
	return (memcmp(p1, p2, (size_t) c));
}

static void
unixWarningHandler(const char* module, const char* fmt, va_list ap)
{
	if (module != NULL)
		fprintf(stderr, "%s: ", module);
	fprintf(stderr, "Warning, ");
	vfprintf(stderr, fmt, ap);
	fprintf(stderr, ".\n");
}
TIFFErrorHandler _TIFFwarningHandler = unixWarningHandler;

static void
unixErrorHandler(const char* module, const char* fmt, va_list ap)
{
	if (module != NULL)
		fprintf(stderr, "%s: ", module);
	vfprintf(stderr, fmt, ap);
	fprintf(stderr, ".\n");
}
TIFFErrorHandler _TIFFerrorHandler = unixErrorHandler;

void
t2p_disable(TIFF *tif)
{
	T2P *t2p = (T2P*) TIFFClientdata(tif);
	GoOutputDisable(t2p);
}

void
t2p_enable(TIFF *tif)
{
	T2P *t2p = (T2P*) TIFFClientdata(tif);
	GoOutputEnable(t2p);
}

/* vim: set ts=8 sts=8 sw=8 noet: */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 8
 * fill-column: 78
 * End:
 */
