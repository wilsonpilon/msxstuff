/*
 * Copyright (c) 1999-2000 Tatsuhiko Syoji, Japan 982-0221. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 * 1 Redistributions of source code must retain the above copyright notice,
 *  this list of conditions and the following disclaimer as the first lines
 *  of this file unmodified.
 *
 * 2 Redistributions in binary form must reproduce the above copyright
 * notice, this list of conditions and the following disclaimer in the
 * documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY Tatsuhiko Syoji ``AS IS'' AND ANY EXPRESS
 * ORIMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 *  WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 *  DISCLAIMED. IN NO EVENT SHALL Tatsuhiko Syoji BE LIABLE FOR ANY DIRECT,
 *  INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN
 *  ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

/* Library for MSX Games and Graphics */
/* By Tatsuuhiko Syoji(Tatsu) 1999,2000 */

#ifndef	__GLIB_H
#define	__GLIB_H

/* Graphics */
/* MSX-C Library compatible */
void ginit(void);
void interlace(char);

void setrd(unsigned);
char invdp(void);
void setwrt(unsigned);
void outvdp(unsigned char);
unsigned char vpeek(unsigned int);
void vpoke(unsigned int,unsigned char);

void wrtvdp(char,unsigned char);
void vdp(char,unsigned char);
unsigned char rdvdp(char);
unsigned char rdvsts(char);
unsigned char vdpstat(char);

void color(char,char,char);
void iniplt(void);
void setplt(char,unsigned int);

void inispr(char);
unsigned int calpat(char);
unsigned int calatr(char);
void sprite(char,void *);
void colspr(char,void *);
void putspr(char,int,int,char,char);

void totext(void);
void grpprt(char,char);
void glocate(int,int);

void setpg(char,char);

/* Tatsu original */
void fsetrd(char,unsigned int);
void fsetwt(char,unsigned int);
void lsetrd(unsigned long);
void lsetwt(unsigned long);
void vdpcom(void *);

#define	PSET	0x00
#define	AND		0x01
#define	OR		0x02
#define	XOR		0x03
#define	PRESET	0x04
#define	TPSET	0x08
#define	TAND	0x09
#define	TOR		0x0a
#define	TXOR	0x0b
#define	TPRESET	0x0c

#define c_screen (*((unsigned char *)(0xfcaf)))
#define c_fore	(*((unsigned char *)(0xf3e9)))
#define	c_back	(*((unsigned char *)(0xf3ea)))
#define	c_bord	(*((unsigned char *)(0xf3eb)))

#include	<msxbios.h>

#endif	/* __GLIB_H */
