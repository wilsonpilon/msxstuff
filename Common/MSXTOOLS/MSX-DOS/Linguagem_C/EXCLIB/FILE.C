/* ｺﾉ ﾌﾟﾛｸﾞﾗﾑ ﾊ MSX-C Ver.1.1 ﾖｳ ﾃﾞｽ */
/* ﾌｧｲﾙｱｸｾｽ ﾉ ｼｶﾀ ｶﾞ ﾁｶﾞｳ ﾉﾃﾞ */
/* MSX-C Ver.1.2 ﾃﾞﾊ ｾﾞｯﾀｲ ﾆ ｼﾞｯｺｳ ｼﾅｲ ﾃﾞ ｸﾀﾞｻｲ ｡ */

#include	<stdio.h>
#include	<bdosfunc.h>
#include	<math.h>
#include	<file.h>

#pragma nonrec

int fread(buf,size,len,fp)
char *buf;
size_t size;
size_t len;
FILE *fp;
{
	unsigned i,j;
	int c;
	char f;

	for (i = 0;i < len;i++){
		for (j = 0;j < size;j++){
			f = 0;
			if ((c = getc(fp)) == EOF){
				f = 1;
				break;
			}
			*buf = (char)c;
			buf++;
		}
		if (f){
			break;
		}
	}
	return(i);
}

int fwrite(buf,size,len,fp)
char *buf;
size_t size;
size_t len;
FILE *fp;
{
	unsigned i,j;
	char f;

	for (i = 0;i < len;i++){
		for (j = 0;j < size;j++){
			f = 0;
			if (putc(*buf,fp) == ERROR){
				f = 1;
				break;
			}
			buf++;
		}
		if (f){
			break;
		}
	}
	return(i);
}

int getw(fp)
FILE *fp;
{
	char bf;
	int c,r;

	bf = 1;
	if (fp->mode & _BINARY){
		bf = 0;
	}
	fp->mode = fp->mode | _BINARY;
	
	if ((c = getc(fp)) == EOF){
		if (bf){
		    fp->mode &= ~_BINARY;
		}
		return(EOF);
	}
	r = c;
	if ((c = getc(fp)) == EOF){
		if (bf){
		    fp->mode &= ~_BINARY;
		}
		return(EOF);
	}
	r += (c << 8);
	if (bf){
	    fp->mode &= ~_BINARY;
	}

	return(r);
}

int putw(w,fp)
int w;
FILE *fp;
{
	char bf;
	int c,r;

	bf = 1;
	if (fp->mode & _BINARY){
		bf = 0;
	}
	fp->mode = fp->mode | _BINARY;
	
	if (putc((char)(w & 0x00ff),fp) == ERROR){
		if (bf){
		    fp->mode &= ~_BINARY;
		}
		return(EOF);
	}
	if (putc((char)(w >> 8),fp) == ERROR){
		if (bf){
		    fp->mode &= ~_BINARY;
		}
		return(EOF);
	}
	if (bf){
	    fp->mode &= ~_BINARY;
	}

	return(w);
}

BOOL feof(fp)
FILE *fp;
{

	return(fp->mode & _EOF);
}

BOOL ferror(fp)
FILE *fp;
{

	return(fp->mode & _OVF);
}

VOID clearerr(fp)
FILE *fp;
{
	fp->mode = fp->mode & 0xf3; /* fp->mode & 11110011B */
}

int fcloseall()
{
	int i,files;
	
	files = 0;
	for (i = 3;i < _NFILES;i++){
		if (_iob[i].mode){
			if (fclose(&_iob[i]) == 0){
				files++;
			}
		}
	}
	if (files == 0){
		return(EOF);
	}else{
		return(files);
	}

}

FD fileno(fp)
FILE *fp;
{
	return(fp->fd);
}

int fflush(fp)
FILE *fp;
{
	if (fp->serial == 0){
		if (fp->mode & _READ){
			if (_seek(fp->fd,-(fp->count),1) == OK){
				fp->count = 0;
				return(0);
			}else{
				return(EOF);
			}
		}else{
			if (fp->mode & _WRITE){
				if (_flushbuf(fp) == OK){
					return(0);
				}
			}else{
				return(EOF);
			}
		}
	}else{
		return(0); /* バッファリングが無いのでとっくにflushしてるよーん */
	}
	return(EOF);
}

int flushall()
{
	int i,files;
	
	files = 0;
	for (i = 3;i < _NFILES;i++){
		if (_iob[i].mode){
			if (fflush(&_iob[i]) == 0){
				files++;
			}
		}
	}
	if (files == 0){
		return(EOF);
	}else{
		return(files);
	}

}

int fseek(fp,offset,ptrname)
FILE    *fp;
int offset;
char ptrname;
{
	size_t  sz;
	
	if (fp->serial == 0){
		fflush(fp);
		if (_seek(fp->fd,offset,ptrname) == OK){
			fp->ptr = fp->base;
			if (fp->mode & _WRITE){
				fp->count = fp->bufsiz;
			}else{
				fp->count = 0;
			}
			return (0);
		}
	}
	return(1); /* Failed */

}

/* ftell()はlong型が無いので実現不可能です。MSX-C の馬鹿野郎！！ */
/* fgetpos(),fsetpos()も同様です。 */

VOID rewind(fp)
FILE *fp;
{
	fseek(fp, 0, SEEK_SET);
	clearerr(fp);
}

int eof(fd)
FD fd;
{
    FCB     *fcb;

    if ((fcb = _getfcb(fd)) == NULL ){
		return (ERROR);
	}
	if ((fcb->recpos[0] == fcb->filesize[0]) && (fcb->recpos[1] == fcb->filesize[1])){
		return(1);
	}
	return(0);
}

/* isatty()関数について。MSX-C Ver.1.1では低水準入出力は表向き */
/* ディスクのファイルに対してのみ行うことになっているので */
/* 意味が無いことになるので実装しない。 */

VOID rdsec(buf,drv,sec,len)
char *buf;
int drv;
unsigned sec;
int len;
{
	int secx;
	bdos((char)0x1a,buf,0);
	secx = (len << 8) + drv;
	bdos((char)0x2f,sec,secx);
}

VOID wtsec(buf,drv,sec,len)
char *buf;
int drv;
unsigned sec;
int len;
{
	int secx;
	bdos((char)0x1a,buf,0);
	secx = (len << 8) + drv;
	bdos((char)0x30,sec,secx);
}

/* Added  1995/06/28 */
char    *fgets2(s, n, fp)
char    *s;
int     n;
FILE    *fp;
{
    int     c;
    char    *cptr;

    cptr = s;
    while (--n != 0 && (c = getc(fp)) != EOF) {
		if ((*cptr++ = c) == '\n'){
			cptr--;
		    break;
		}
    }
    *cptr = '\0';

    return( (c == EOF && cptr == s)? NULL: s );
}

/* DOS1 Only SLONG version _seek() */
STATUS	seekl(fd, offset, mode)
FD		fd;
SLONG	*offset;
TINY	mode;
{
	FCB		*fcb;
	SLONG *rpos2;

	if (2 < mode || (fcb = _getfcb(fd)) == NULL){
		return (ERROR);
	}

	rpos2 = (SLONG *)(fcb->recpos);

	if (mode == 0) {
		slcpy(rpos2,offset);
	}else{
		if (mode == 2) {
			fcb->recpos[0] = fcb->filesize[0];
			fcb->recpos[1] = fcb->filesize[1];
		}
		sladd(rpos2,rpos2,offset);
	}

	return (OK);
}

/* DOS1 Only SLONG version tell() */
SLONG *telll(ret,fd)
SLONG *ret;
FD fd;
{
	FCB		*fcb;
	SLONG *rpos2,t;

	if ((fcb = _getfcb(fd)) == NULL){
		slcpy(ret,atosl(&t,"-1"));
	}else{
		rpos2 = (SLONG *)(fcb->recpos);
		slcpy(ret,rpos2);
	}
	return(ret);
}

