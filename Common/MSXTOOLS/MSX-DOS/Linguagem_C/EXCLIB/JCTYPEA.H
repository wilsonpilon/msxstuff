#ifndef __JCTYPE_H
#define __JCTYPE_H

/* defines for poor linker */
#define iskanji iskan
#define iskanji2 iskan2
#define isalkana iak
#define ispnkana ipnk
#define isalnmkana isank
#define isprkana iprk
#define isgrkana igk
#define jmstojis mstojis
#define jistojms jistoms
#define isalhira iah
#define iskanahira ikh
#define isalkanahira iakh

/* ’è” */
#define	CT_ANK	0
#define	CT_KJ1	1
#define	CT_KJ2	2
#define	CT_ILGL	-1

extern unsigned _CHMAP[];


/* For stand alone */
#define iscsymf iscsmf
#define isalnum(c)		(isalpha(c) || isdigit(c))
#define isascii(c)		((c) <= '\177')
#define isxdigit(c)		(isdigit(c) || ('A' <= (c) && (c) <= 'F') || ('a' <= (c) && (c) <= 'f'))
BOOL ispunct(),isprint(),isgraph(),iscsymf(),iscsym();

int chkctype(),nthctype();
/* Kana character functions. */
BOOL iskanji(),iskanji2(),iskana(),iskpun(),iskmoji(),isalkana();
BOOL ispnkana(),isalnmkana(),isprkana(),isgrkana();
/* 2Byte code type */
BOOL jiszen(),jisl0(),jisl1(),jisl2(),jisupper(),jislower(),jisalpha();
BOOL jisdigit(),jiskata(),jishira(),jiskigou(),jisspace();
/* 2Byte code transfar */
unsigned jtoupper(),jtolower(),jtohira(),jtokata(),hantozen(),zentohan();
unsigned mstojis(),jistoms();
/* MSX 1byte hiragana functions */
BOOL ishira(),isalhira(),iskatahira(),isalkanahira();
/* 2Byte code transfar for MSX 1Byte code */
unsigned h2zen2(),z2han2();

#endif
