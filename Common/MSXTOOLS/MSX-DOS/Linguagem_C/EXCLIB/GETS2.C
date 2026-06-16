/* ‚±‚ê‚ª‘¼‚Ì‚bŒ¾Œê‚Æ“¯‚¶“®ì‚ğ‚·‚égets() & puts() */

#include	<stdio.h>

#pragma nonrec

STATUS  puts(s)
char    *s;
{
    char    c;

    while (c = *s++) {
	if (putc(c, stdout) == ERROR)
	    return (ERROR);
    }
    putc('\n',stdout);
    return (OK);
}

char    *gets(s)
char    *s;
{
    int     c;
    char    *cptr;

    cptr = s;
    while ((c = getc(stdin)) != EOF) {
	if ((*cptr++ = c) == '\n')
		cptr--;
	    break;
    }
    *cptr = '\0';
    return (c == EOF && cptr == s)? NULL: s;
}

