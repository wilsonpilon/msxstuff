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

/* MSX-BIOS function library */
/* By Tatsuhiko Syoji(Tatsu) 1999-2000 */

#ifndef	__MSXBIOS_H
#define	__MSXBIOS_H

/* MSX-C Library compatible function */
void disscr(void);
void enascr(void);
void screen(char);
void gicini(void);
void sound(char,unsigned char);
unsigned char rdpsg(char);
int chsns(void);
void chput(unsigned char);
unsigned char chget(void);
unsigned char lptout(char);
unsigned char lptstt(void);
unsigned char *pinlin(void);
unsigned char *inlin(void);
char breakx(void);
void beep(void);
void cls(void);
void locate(char,char);
void erafnk(void);
void dspfnk(void);
char gtstck(char);
char gttrig(char);
char gtpad(char);
unsigned char gtpdl(char);
char chgsnd(char);
unsigned char snsmat(char);
void kilbuf(void);
void ei(void);
void di(void);
/* MSX Turbo R BIOS function */
void chgcpu(char);
char getcpu(void);

/* Tatsu original */
/* Special I/O function for FM-TOWNS gamepad */
char gstex();

#endif
