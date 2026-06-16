#ifndef __CTYPE_H
#define __CTYPE_H

#define iscsymf iscsmf

/* ASCII character type functions. (<ctype.h>) */

/* この３つの関数についてはMSX-C Ver.1.2でもこのマクロと同様の */
/* 処理が行われていると推測する。 */
#define isalnum(c)		(isalpha(c) || isdigit(c))
#define isascii(c)		((c) <= '\177')
#define isxdigit(c)		(isdigit(c) || ('A' <= (c) && (c) <= 'F') || ('a' <= (c) && (c) <= 'f'))

BOOL ispunct(),isprint(),isgraph(),iscsymf(),iscsym();
char toascii(),_tolower(),_toupper();

#endif
