#ifndef __STRING_H
#define __STRING_H

#define strcmpi scmpi /* for poor linker with 6 chars symbol */
#define stricmp scmpi /* for poor linker with 6 chars symbol */
#define strncmpi sncmpi /* for poor linker with 6 chars symbol */
#define strnicmp sncmpi /* for poor linker with 6 chars symbol */
#define strpbrk strpbr /* for poor linker with 6 chars symbol */
#define strnset strnse /* easy operation for poor linker with 6 chars symbol */

char *strncpy(),*strncat(),*strchr(),*strstr(),*strrev(),*strrchr(),*strupr();
char *strlwr(),*strtok(),*strpbrk(),*strnset(),*strset(),*strdup(),*stpcpy();
char *memccpy(),*memchr();
int strncmp(),strcmpi(),strncmpi(),memcmp(),memicmp();
unsigned strspn(),strcspn();

VOID memmove(); /* Function in mlib.rel (or curses.rel) */

#endif
