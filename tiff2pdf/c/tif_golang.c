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
 *
 * This file is the bridge between libtiff and the Go I/O layer. libtiff's own
 * tif_unix.c / tif_error.c / tif_warning.c are compiled in unchanged (see
 * libtiff.h), so this file no longer reimplements allocators or error
 * handlers. It provides:
 *
 *   - GoTIFFFdOpen: opens a TIFF whose I/O is serviced by Go callbacks,
 *     identified by a small integer "fd" carried as the libtiff client data.
 *   - Go error/warning handlers, installed via libtiff's public handler API.
 *     libtiff routes every message (TIFFError, TIFFErrorExt and the
 *     re-entrant TIFFErrorExtR used internally) through the global *Ext
 *     handler when no per-TIFF handler is set, so registering these captures
 *     all diagnostics. The default stderr handlers are cleared so messages go
 *     only to Go.
 */

#include "tif_config.h"

#ifdef HAVE_SYS_TYPES_H
# include <sys/types.h>
#endif

#include <errno.h>
#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>

#ifdef HAVE_UNISTD_H
# include <unistd.h>
#endif

#ifdef HAVE_FCNTL_H
# include <fcntl.h>
#endif

#include "tiffiop.h"

/* The Go callbacks (//export in hooks.go) take the fd as a Go int. libtiff
 * carries it as thandle_t (a pointer-sized opaque handle), so convert via
 * intptr_t at each boundary. Names are prefixed _go to avoid colliding with
 * tif_unix.c's own static _tiff*Proc helpers (same translation unit). */

static tmsize_t
_goReadProc(thandle_t fd, void* buf, tmsize_t size)
{
	return (tmsize_t)GoTiffReadProc((int)(intptr_t)fd, buf, size);
}

static tmsize_t
_goNoopReadProc(thandle_t fd, void* buf, tmsize_t size)
{
	(void) fd; (void) buf; (void) size;
	return -1;
}

static tmsize_t
_goWriteProc(thandle_t fd, void* buf, tmsize_t size)
{
	return (tmsize_t)GoTiffWriteProc((int)(intptr_t)fd, buf, size);
}

static uint64_t
_goSeekProc(thandle_t fd, uint64_t off, int whence)
{
	return (uint64_t)GoTiffSeekProc((int)(intptr_t)fd, (int64_t)off, whence);
}

static int
_goCloseProc(thandle_t fd)
{
	return GoTiffCloseProc((int)(intptr_t)fd);
}

static uint64_t
_goSizeProc(thandle_t fd)
{
	return (uint64_t)GoTiffSizeProc((int)(intptr_t)fd);
}

static int
_goMapProc(thandle_t fd, void** pbase, toff_t* psize)
{
	(void) fd; (void) pbase; (void) psize;
	return (0);
}

static void
_goUnmapProc(thandle_t fd, void* base, toff_t size)
{
	(void) fd; (void) base; (void) size;
}

static void
_goErrorHandlerExt(thandle_t fd, const char* module, const char* fmt, va_list ap)
{
	char s[4096];
	(void) module;
	vsnprintf(s, sizeof(s), fmt, ap);
	GoTiffErrorExt((int)(intptr_t)fd, s);
}

static void
_goWarningHandlerExt(thandle_t fd, const char* module, const char* fmt, va_list ap)
{
	char s[4096];
	(void) module;
	vsnprintf(s, sizeof(s), fmt, ap);
	GoTiffWarningExt((int)(intptr_t)fd, s);
}

static void
_goInstallHandlers(void)
{
	static int installed = 0;
	if (installed)
		return;
	installed = 1;
	/* Clear the default handlers (which print to stderr) so libtiff
	 * diagnostics are delivered only to Go. */
	TIFFSetErrorHandler(NULL);
	TIFFSetWarningHandler(NULL);
	TIFFSetErrorHandlerExt(_goErrorHandlerExt);
	TIFFSetWarningHandlerExt(_goWarningHandlerExt);
}

/*
 * Open a TIFF for read/writing, backed by the Go I/O callbacks identified by fd.
 */
TIFF*
GoTIFFFdOpen(int fd, const char* name, const char* mode)
{
	TIFF* tif;
	TIFFReadWriteProc readproc = _goReadProc;

	_goInstallHandlers();

	if (strlen(name) >= 4 && 0 == strncmp(name + strlen(name) - 4, ".pdf", 4)) {
		readproc = _goNoopReadProc;
	}

	tif = TIFFClientOpen(name, mode,
	    (thandle_t)(intptr_t) fd,
	    readproc, _goWriteProc,
	    _goSeekProc, _goCloseProc, _goSizeProc,
	    _goMapProc, _goUnmapProc);
	if (tif)
		tif->tif_fd = fd;
	return (tif);
}

/* vim: set ts=8 sts=8 sw=8 noet: */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 8
 * fill-column: 78
 * End:
 */
