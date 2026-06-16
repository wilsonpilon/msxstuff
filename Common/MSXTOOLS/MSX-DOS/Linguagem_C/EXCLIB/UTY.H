#ifndef __UTY_H
#define __UTY_H

/* C-Lang */
char inkey();
int maxpos();
VOID kbcom(),kbcomr();
STATUS bsaves(),bloads();

/* ASM */
VOID fsetrd(),fsetwt();
unsigned rndx();
VOID srand();
char gstex();
VOID chgcpu();
char getcpu();
VOID vdpcom();

#define getkmode()	calla(0xffca,0,0,0,0x1100)
#define chkkmode()	calla(0xffca,0xff,0,0,0x1100)
#define getwidth()	*(char *)0xf3b0
#define TIME	*(unsigned *)(0xfc9e)

BOOL _isfn();
STATUS _chkdrv();
VOID _swp();
char *_skipsp();

typedef char void;
typedef char tiny;
typedef char bool;
typedef unsigned nat;

#endif
