#ifndef __FILE_H
#define __FILE_H

#ifndef __BDOSFUNC_H
/* FCB型を使うのに必要 */
#include	<bdosfunc.h>
#endif

/* fgetc(),fputc()が無いので同じ働きをする関数に名前を置き換える */
/* 本当は逆なのだが... */
#define fputc putc
#define fgetc getc

/* 内部関数ではあるが他の処理系では公式な関数になっている */
#define lseek _seek

/* for poor linker */
#define fcloseall fcloal

/* fseek()及びlseek()が使う */
#define SEEK_SET (char)0
/* ファイルの先頭からの絶対位置 */
#define SEEK_CUR (char)1
/* 現在の位置からの相対位置 */
#define SEEK_END (char)2
/* ファイルの最後からの相対位置 */

/* 関数ご一行様 */
int fread(),fwrite(),getw(),putw(),fcloseall(),fflush(),flushall(),eof();
int fseek();
BOOL feof(),ferror();
VOID clearerr(),rewind(),rdsec(),wtsec();
FD fileno();
char *fgets2();
char phread(),phwrit();
STATUS	seekl();
/* SLONG *telll(); math.hを使うときだけ定義する。 */

/* 関数内部で使う関数と定義 */
STATUS	_flushbuf();
STATUS	_seek();
FCB *_getfcb();

/* 似非関数 */
#define setverify(c) bdos((char)0x2e,c,0)

#define _NFCBS	20

#define _READ	 01
#define _WRITE	 02
#define _EOF	 04
#define _OVF	010
#define _BINARY 020

#endif
