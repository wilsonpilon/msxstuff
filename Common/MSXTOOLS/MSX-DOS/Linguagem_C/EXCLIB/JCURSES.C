#pragma nonrec
#define EXTERN

#include    <stdio.h>
#include    <curses.h>

VOID _putcc();


#define iskanji iskan
#define iskanji2 iskan2
BOOL iskan(n)
char n;
{
	return((0x81<=n && n<=0x9F) ||  (0xE0<=n && n<=0xFC));
}

BOOL iskan2(n)
char n;
{
	return((0x40<=n && n<=0x7E) ||  (0x80<=n && n<=0xFC));
}

char _kcount(s,n)
char *s;
int n;
{
	char kf;
	kf = 0;
	
	for (;n;n--){
		if (kf){ /* 1Â Ï´ Ê ¼ÌÄJIS É 1 ÊÞ²Ä Ò */
			kf = 0; /* ¼ÌÄJIS É 2 ÊÞ²Ä Ò ÃÞÓ Ì¾² Å Ó¼Þ ÃÞÓ ÂÅ¶Þ×Å² */
		}else{
			if (iskanji(*s)){ /* ¼ÌÄ JIS É 1 ÊÞ²Ä Ò */
				kf = 1;
			}
		}
		s++;
	}
	return(kf);
}

VOID _refine(idpc, idps, begy, begx, y)
INDEXLIN *idpc, *idps; /* index line of current window & window */
TINY begy, begx, y; /* Home position of window & line for refine */
{
    int firstch, lastch, x;
    char *stdy, *cury,kf;
    TINY count;

    count = 0;
    firstch = idps[0]._firstch;
    lastch = idps[0]._lastch;

	stdy = (idps[0]._y + firstch);
	
	kf = 0;
	if (firstch > 0){
		if (_kcount(idps[0]._y,firstch)){
			firstch--;
			stdy--;
			kf = 1;
			putch((char)29);
		}
	}

    for (cury = (idpc[0]._y + firstch + begx),x = firstch; x <= lastch;stdy++, cury++, x++) { /*Start of nested loop*/
	    putch(*stdy);
/*
		if(*stdy != *cury) {
		    *cury = *stdy;
		    if(count){
				_putcc(&count, y, x, begy, begx, &stdy);
			}
		    putch(*stdy);
		} else{
			if (kf){
			    putch(*stdy);
			}else{
				count++;        /* Logical cursor moved. */
			}
		}
*/
    } /*End of nested loop*/
/*
    if (kf){
		putch((char)28);
	}
*/
}

STATUS jwdelch(win)
WINDOW *win;
{
    TINY maxx, curx, cury;
    INDEXLIN *idp;
    char *cp, *sp, *ep;
    int *firstch;

    cury = win->_cury;
    curx = win->_curx;
    idp = win->_index + cury;

    firstch = &(idp[0]._firstch);

    maxx = win->_maxx;
    ep = idp[0]._y + (maxx - 1);
    sp = idp[0]._y + curx + 1;
    
	if ((curx+1) < maxx){
		if (_kcount(idp[0]._y,(int)(curx+1))){
			sp = idp[0]._y + curx + 2;
			cp = sp - 1;
			while(cp < ep){
				*cp++ = *sp++;
			}
		    *cp = SPACE; /* cp == ep */
		}
	}

    sp = idp[0]._y + curx + 1;
    cp = sp - 1;
    while(cp < ep){
		*cp++ = *sp++;
	}

    *cp = SPACE; /* cp == ep */

    idp[0]._lastch = maxx - 1;
    if(*firstch == _NOCHANGE || *firstch > curx){
		*firstch = curx;
	}
    return(OK);
}

STATUS jwmove(win, y, x)
WINDOW  *win;
int     y, x;
{
    if(x >= (int)win->_maxx || y >= (int)win->_maxy || x < 0 || y < 0)
	return(ERROR);

    win->_curx = (TINY)x;
    win->_cury = (TINY)y;
    
    if (x > 0){
	    if (_kcount(win->_index->_y,x)){
			x--;
			win->_curx = (TINY)x;
		}
	}

    return(OK);
}

unsigned jwinch(win)
WINDOW *win;
{
    char c1;
    unsigned c2;
    
    c2 = 0;
    c1 = win->_index[win->_cury]._y[win->_curx];
    if (iskanji(c1)){
    	c2 = (unsigned)(win->_index[win->_cury]._y[win->_curx+1]);
    }
    c2 = (c2 << 8) + (unsigned)c1;
    
    return(c2);
}
