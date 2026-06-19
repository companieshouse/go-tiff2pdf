/*
 * libtiff source files compiled directly into the cgo build.
 *
 * Unlike the original 2016 build, tif_unix.c, tif_error.c and tif_warning.c
 * are now compiled in unchanged. The Go-specific bridge in tif_golang.c no
 * longer reimplements libtiff internals (allocators, error handlers); it only
 * adds GoTIFFFdOpen and registers Go error/warning handlers. This keeps us in
 * step with the modern libtiff API (re-entrant *ExtR handlers, *Ext
 * allocators) instead of shadowing it.
 *
 * tif_win32.c is intentionally excluded (POSIX build). The codec files for
 * features disabled at ./configure time (jpeg, lzma, jbig, zstd, lerc, webp,
 * pixarlog, zlib) compile to empty translation units, so they are harmless to
 * include and the list stays a stable superset across configure options.
 */
#include "tif_aux.c"
#include "tif_close.c"
#include "tif_codec.c"
#include "tif_color.c"
#include "tif_compress.c"
#include "tif_dir.c"
#include "tif_dirinfo.c"
#include "tif_dirread.c"
#include "tif_dirwrite.c"
#include "tif_dumpmode.c"
#include "tif_error.c"
#include "tif_extension.c"
#include "tif_fax3.c"
#include "tif_fax3sm.c"
#include "tif_flush.c"
#include "tif_getimage.c"
#include "tif_hash_set.c"
#include "tif_jbig.c"
#include "tif_jpeg.c"
#include "tif_jpeg_12.c"
#include "tif_lerc.c"
#include "tif_luv.c"
#include "tif_lzma.c"
#include "tif_lzw.c"
#include "tif_next.c"
#include "tif_ojpeg.c"
#include "tif_open.c"
#include "tif_packbits.c"
#include "tif_pixarlog.c"
#include "tif_predict.c"
#include "tif_print.c"
#include "tif_read.c"
#include "tif_strip.c"
#include "tif_swab.c"
#include "tif_thunder.c"
#include "tif_tile.c"
#include "tif_unix.c"
#include "tif_version.c"
#include "tif_warning.c"
#include "tif_webp.c"
#include "tif_write.c"
#include "tif_zip.c"
#include "tif_zstd.c"
