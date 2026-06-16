#include	<stdio.h>

/* #define LSIC 1 */

/* For LSI C-86 3.3 試食版 */
/* ちなみに、MSX-Cの原形はLSI C-80（セルフ版）の古いヴァージョンであることは */
/* CF.COMかCG.COMを覗くと分かる。 */
/* LSI C-80（セルフ版）はANSIに対応したのに、MSX-Cは対応していない。 */
/* ASCII社がMSX-Cをいかにサボっていたか、いかに軽視していたかが良く分かる。 */
#ifdef LSIC
#include	<stdlib.h>
typedef char BOOL;
typedef char VOID;
#endif

#pragma nonrec

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

char getcx()
{
	int c2;
	if ( (c2 = getchar()) == EOF){
		exit(0);
	}else{
		return((char)c2);
	}
}

VOID ccomen() /* skip C comment */
{
	char c,d;
	int n;
	
	n = 1;
	putchar('/');
	putchar('*');
	c = ' ';
	for(;;){
		d = c;
		c = getcx();
		putchar(c);
		if ((d == '*') && (c == '/')){
			n--;
			if (n == 0){
				break;
			}
			c = ' ';
		}
		if ((d == '/') && (c == '*')){
			n++;
			c = ' ';
		}
	}
	c = ' ';
}

VOID binhex() /* convert binary to hex. */
{
	int bit;
	unsigned int l;
	char c;
	
	l = 0;
	bit = 0;
	
	for (c = getcx();(c == '0') || (c == '1');c = getcx()){
		bit++;
		l = (l << 1);
		if (c == '1'){
			l++;
		}
	}

	if (bit < 9){ /* 8bit output mode */
		printf("x%02x",l);
	}else{ /* 16bit output mode */
		printf("x%04x",l);
	}
	ungetc(c,stdin);
}

VOID string(d,sc) /* skip & output string (or character) */
char d;
char sc;
{
	char c,cb;
	
	putchar(sc);
	for(c = getcx();c != d;c = getcx()){
		putchar(c);
		if ((c == '\\') || iskan(c)){ /* \の次は無条件出力、文字がシフトＪＩＳの時の処理 */
			cb = c;
			c = getcx();
			putchar(c);
			if (iskan(cb) && (c == '\\')){
				putchar('\\');
			}
		}
	}
	putchar(c);
}

int main()
{

	int c2;
	char c,d;

	fprintf(stderr,"MSX C ver 1.12PP   (pre-processer)\nCopyright (C) 1995 by TASCII Corporation (^^)\n");

	while ((c2 = getchar()) != EOF){
		c = (char)c2;
		if (c == '/'){
			c = getcx();
			switch (c){
				case '*':
					ccomen();
					c = ' ';
				break;
				case '/':
					putchar('/');
					putchar('*');
					for(c = getcx();c != '\n';c = getcx()){
						putchar(c);
					}
					putchar('*');
					putchar('/');
					/* '\n' ﾊ ｱﾄ ﾉ ﾙｰﾁﾝ ﾃﾞ ｼｭﾂﾘｮｸ */
				break;
				default:
					ungetc(c,stdin);
					c = '/';
				break;
			}
		}

		d = c;
		
		switch (c){
			case '"':
			case '\'':
				string(d,c);
			break;
			case '0':
				putchar(c);
				c = getcx();
				if ((c == 'b') || (c == 'B')){
					binhex();
				}else{
					ungetc(c,stdin);
				}
			break;
			default:
				putchar(c);
				if (iskan(c)){
					c = getcx();
					putchar(c);
					if (c == '\\'){
						putchar(c);
					}
				}
			break;
		}
		
	}

	return(0);
}
