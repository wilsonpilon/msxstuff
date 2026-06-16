#include	<stdio.h>

/* defines for poor linker */
#define iskanji iskan
#define iskanji2 iskan2

#define	CT_ANK	0
#define	CT_KJ1	1
#define	CT_KJ2	2
#define	CT_ILGL	-1

/* #define LSIC 1 */
#ifdef LSIC /* LSI-C86 3.3試食版 */

#include	<ctype.h>
typedef int BOOL;

#else /* MSX-C Ver.1.1 */
/* ASCII character type functions. */
/* この３つの関数についてはMSX-C Ver.1.2でもこのマクロと同様の */
/* 処理が行われていると推測する。 */
#define isalnum(c)		(isalpha(c) || isdigit(c))
#define isascii(c)		((c) <= '\177')
#define isxdigit(c)		(isdigit(c) || ('A' <= (c) && (c) <= 'F') || ('a' <= (c) && (c) <= 'f'))

/* For poor linker */
#define iscsymf iscsmf

BOOL ispunct(c)
char c;
{
  return (!isalnum(c) && !iscntrl(c));
}

BOOL isprint(c)
char c;
{
  return (c>=32 && c<=126);
}

BOOL isgraph(c)
int c;
{
	return (c>=33 && c<=126);
}

BOOL iscsymf(c) /* C ﾉ ｼｷﾍﾞﾂｼ ﾉ 1 ﾓｼﾞ ﾒ ｶ ? */
char c;
{
	return(isalpha(c) || (c == '_'));
}

BOOL iscsym(c) /* C ﾉ ｼｷﾍﾞﾂｼ ｶ ? */
char c;
{
	return(isalnum(c) || (c == '_'));
}

#endif

/* Japanese character type functions. */
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

int chkctype(c,type)
char c;
int type;
{
	int btype;
	
	switch (type){
		case CT_ANK:
		case CT_KJ2:
		case CT_ILGL:
			if (iskanji(c)){
				btype = CT_KJ1;
			}else{
				btype = CT_ANK;
			}
		break;
		case CT_KJ1:
			if (iskanji2(c)){
				btype = CT_KJ2;
			}else{
				btype = CT_ILGL;
			}
		break;
	}
	return(btype);
}

int nthctype(s,n)
char *s;
int n;
{
	int type,i;
	
	type = CT_ANK;
	
	for (i = 0;i < (n-1);i++){
		type = chkctype(*s,type);
		if (type == CT_KJ1){
			s++;
			type = chkctype(*s,type);
			if (type = CT_ILGL){
				s--;
			}
		}
		s++;
	}
	type = chkctype(*s,type);
	return(type);
}

BOOL iskana(c)
char c;
{
	return((c >= 0xa1) && (c <= 0xdf));
}

BOOL iskpun(c)
char c;
{
	return((c >= 0xa1) && (c <= 0xa5));
}

BOOL iskmoji(c)
char c;
{
	return((c >= 0xa6) && (c <= 0xdf));
}

#define isalkana iak
BOOL isalkana(c)
char c;
{
	return(isalpha(c) || iskana(c));
}

#define ispnkana ipnk
BOOL ispnkana(c)
char c;
{
	return(ispunct(c) || iskpun(c));
}

#define isalnmkana isank
BOOL isalnmkana(c)
char c;
{
	return(isalnum(c) || iskana(c));
}

#define isprkana iprk
BOOL isprkana(c)
char c;
{
	return(isprint(c) || iskana(c));
}

#define isgrkana igk
BOOL isgrkana(c)
char c;
{
	return(isgraph(c) || iskpun(c));
}

BOOL jiszen(c)
unsigned c;
{
	return ((c >= 0x8140) && (c <= 0xeffc));
}

BOOL jisl0(c)
unsigned c;
{
	return (jiszen(c) && (c >= 0x8140) && (c <= 0x889e));
}

BOOL jisl1(c)
unsigned c;
{
	return (jiszen(c) && (c >= 0x889f) && (c <= 0x9872));
}

BOOL jisl2(c)
unsigned c;
{
	return (jiszen(c) && (c >= 0x989f) && (c <= 0xeffc));
}

BOOL jisupper(c)
unsigned c;
{
	return ((c >= 0x8260) && (c <= 0x8279));
}

BOOL jislower(c)
unsigned c;
{
	return ((c >= 0x8281) && (c <= 0x829a));
}

BOOL jisalpha(c)
unsigned c;
{
	return (jisupper(c) || jislower(c));
}

BOOL jisdigit(c)
unsigned c;
{
	return ((c >= 0x824f) && (c <= 0x8258));
}

BOOL jiskata(c)
unsigned c;
{
	return ((c >= 0x8340) && (c <= 0x8396) && (c != 0x837f));
}

BOOL jishira(c)
unsigned c;
{
	return ((c >= 0x829f) && (c <= 0x82f1));
}

BOOL jiskigou(c)
unsigned c;
{
	return ((c >= 0x8141) && (c <= 0x81ac) && (c != 0x817f));
}

BOOL jisspace(c)
unsigned c;
{
	return (c == 0x8140);
}

unsigned jtoupper(c)
unsigned c;
{
	if (jislower(c)){
		c -= 0x21;
	}
	return(c);
}

unsigned jtolower(c)
unsigned c;
{
	if (jisupper(c)){
		c += 0x21;
	}
	return(c);
}

unsigned jtohira(c)
unsigned c;
{
	if (jiskata(c)){
		if ((c >= 0x8380) && (c <= 0x8393)){ /* 謎のコードの隙間１文字 */
			c -= 0xa2;
		}else{
			c -= 0xa1;
		}
	}
	return(c);
}

unsigned jtokata(c)
unsigned c;
{
	if (jishira(c)){
		if (c >= 0x82de){
			c += 0xa2;
		}else{
			c += 0xa1;
		}
	}
	return(c);
}

unsigned _CHMAP[] = {
	0x8140,0x8149,0x8168,0x8194,0x8190,0x8193,0x8195,0x8166,0x8169,
	0x816a,0x8196,0x817b,0x8143,0x817c,0x8144,0x815e,0x824f,0x8250,
	0x8251,0x8252,0x8253,0x8254,0x8255,0x8256,0x8257,0x8258,0x8146,
	0x8147,0x8183,0x8181,0x8184,0x8148,0x8197,0x8260,0x8261,0x8262,
	0x8263,0x8264,0x8265,0x8266,0x8267,0x8268,0x8269,0x826a,0x826b,
	0x826c,0x826d,0x826e,0x826f,0x8270,0x8271,0x8272,0x8273,0x8274,
	0x8275,0x8276,0x8277,0x8278,0x8279,0x816d,0x818f,0x816e,0x814f,
	0x8151,0x8165,0x8281,0x8282,0x8283,0x8284,0x8285,0x8286,0x8287,
	0x8288,0x8289,0x828a,0x828b,0x828c,0x828d,0x828e,0x828f,0x8290,
	0x8291,0x8292,0x8293,0x8294,0x8295,0x8296,0x8297,0x8298,0x8299,
	0x829a,0x816f,0x8162,0x8170,0x8150
};

unsigned hantozen(c)
unsigned c;
{
	if ((c < 32) || (c >= 127)){
		return(c);
	}else{
		return(_CHMAP[c-32]);
	}
}

unsigned zentohan(c)
unsigned c;
{
	int i;
	for (i = 0;i < 95;i++){
		if (_CHMAP[i] == c){
			return((unsigned)(i+32));
		}
	}
	return(c);
}

#define jmstojis mstojis
unsigned mstojis(c)
unsigned c;
{
	unsigned jh,jl;

	jh = (c >> 8);
	jl = c & 0xff;

	if (jh < 0xa0){
		jh = jh - 0x71;
	}else{
		jh = jh - 0xb1;
	}
	jh = (jh << 1) + 1;
	
	if (jl > 0x7f){
		jl--;
	}
	if (jl < 0x9e){
		jl = jl - 0x1f;
	}else{
		jl = jl - 0x7d;
		jh++;
	}
	return((jh << 8) | jl);
}

#define jistojms jistoms
unsigned jistoms(c)
unsigned c;
{
	unsigned jh,jl;

	jh = (c >> 8);
	jl = c & 0xff;

	if (jh <= 0x5e){
		jh = ((jh - 1) >> 1) + 0x71;
	}else{
		jh = ((jh - 1) >> 1) + 0xb1;
	}
	
	if (jl % 2){
		jl = jl + 0x1f;
		if (jl >= 0x7f){
			jl++;
		}
	}else{
		jl = jl + 0x7e;
	}

	return((jh << 8) | jl);
}

/* MSX hiragana character functions */
BOOL ishira(c) /* is hiragana ? */
char c;
{
	return ( ((c >= 0x86) && (c <= 0x8f)) || ((c >= 0x91) && (c <= 0x9f)) || ((c >= 0xe1) && (c <= 0xfd)) );
}

#define isalhira iah
BOOL isalhira(c) /* is hiragana or alphabet ? */
char c;
{
	return(isalpha(c) || ishira(c));
}

#define iskanahira ikh
BOOL iskanahira(c) /* is hiragana or katakana ? */
char c;
{
	return(ishira(c) || iskana(c));
}

#define isalkanahira iakh
BOOL isalkanahira(c) /* is hiragana or alphabet or katakana ? */
char c;
{
	return(isalpha(c) || iskana(c) || ishira(c));
}

unsigned _CMAPM[] = {
	0x8140,0x8149,0x8168,0x8194,0x8190,0x8193,0x8195,0x8166, /* ASCII */
	0x8169,0x816a,0x8196,0x817b,0x8143,0x817c,0x8144,0x815e,
	0x824f,0x8250,0x8251,0x8252,0x8253,0x8254,0x8255,0x8256,
	0x8257,0x8258,0x8146,0x8147,0x8183,0x8181,0x8184,0x8148,
	0x8197,0x8260,0x8261,0x8262,0x8263,0x8264,0x8265,0x8266,
	0x8267,0x8268,0x8269,0x826a,0x826b,0x826c,0x826d,0x826e,
	0x826f,0x8270,0x8271,0x8272,0x8273,0x8274,0x8275,0x8276,
	0x8277,0x8278,0x8279,0x816d,0x818f,0x816e,0x814f,0x8151,
	0x8165,0x8281,0x8282,0x8283,0x8284,0x8285,0x8286,0x8287,
	0x8288,0x8289,0x828a,0x828b,0x828c,0x828d,0x828e,0x828f,
	0x8290,0x8291,0x8292,0x8293,0x8294,0x8295,0x8296,0x8297,
	0x8298,0x8299,0x829a,0x816f,0x8162,0x8170,0x8150,0x8140,
	0x8140,0x8140,0x8140,0x8140,0x819b,0x819c,0x82f0,0x829f, /* Japanese */
	0x82a1,0x82a3,0x82a5,0x82a7,0x82e1,0x82e3,0x82e5,0x82c1,
	0x8140,0x82a0,0x82a2,0x82a4,0x82a6,0x82a8,0x82a9,0x82ab,
	0x82ad,0x82af,0x82b1,0x82b3,0x82b5,0x82b7,0x82b9,0x82bb,
	0x8140,0x8142,0x8175,0x8176,0x8141,0x8145,0x8392,0x8340,
	0x8342,0x8344,0x8346,0x8348,0x8383,0x8385,0x8387,0x8362,
	0x815b,0x8341,0x8343,0x8345,0x8347,0x8349,0x834a,0x834c,
	0x834e,0x8350,0x8352,0x8354,0x8356,0x8358,0x835a,0x835c,
	0x835e,0x8360,0x8363,0x8365,0x8367,0x8369,0x836a,0x836b,
	0x836c,0x836d,0x836e,0x8371,0x8374,0x8377,0x837a,0x837d,
	0x837e,0x8380,0x8381,0x8382,0x8384,0x8386,0x8388,0x8389,
	0x838a,0x838b,0x838c,0x838d,0x838f,0x8393,0x814a,0x814b,
	0x82bd,0x82bf,0x82c2,0x82c4,0x82c6,0x82c8,0x82c9,0x82ca,
	0x82cb,0x82cc,0x82cd,0x82d0,0x82d3,0x82d6,0x82d9,0x82dc,
	0x82dd,0x82de,0x82df,0x82e0,0x82e2,0x82e4,0x82e6,0x82e7,
	0x82e8,0x82e9,0x82ea,0x82eb,0x82ed,0x82f1,0x8140,0x8140,
};

unsigned h2zen2(c)
unsigned c;
{
	if (c < 32){
		return(c);
	}else{
		return(_CMAPM[c-32]);
	}
}

unsigned z2han2(c)
unsigned c;
{
	int i;
	for (i = 0;i < 224;i++){
		if (_CMAPM[i] == c){
			return((unsigned)(i+32));
		}
	}
	return(c);
}

