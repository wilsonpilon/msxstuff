#include	<stdio.h>

/* #define LSIC 1 */

#ifdef LSIC
#include	<string.h>
typedef int BOOL;
#endif

#include	<jstring.h>

#pragma nonrec

/* defines for poor linker */
#define iskanji iskan
#define iskanji2 iskan2

BOOL iskanji(),iskanji2();

/* function in jctype
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
*/

int btom(buf,n)
char *buf;
int n;
{
	int i,len;
	char zf;
	
	len = 0;
	zf = 0;
	
	for (i = 0;i < n;i++){
		if (zf){
			if (iskanji2(*(buf+i))){
				len++;
			}else{
				len += 2;
			}
			zf = 0;
		}else{
			if (iskanji(*(buf+i))){
				zf = 1;
			}else{
				len++;
			}
		}
	}
	return(len);
}

int mtob(buf,n)
char *buf;
int n;
{
	int i,len;
	char zf;
	
	len = 0;
	zf = 0;
	for (i = 0;i < n;i++){
		if (iskanji(*buf)){
			buf++;
			len++;
			if (iskanji2(*buf)){
				buf++;
				len++;
			}
		}else{
			buf++;
			len++;
		}
	}
	return(len);
}

char *jstradv(buf,n)
char *buf;
int n;
{
	return(buf + mtob(buf,n));
}

char *jstrchr(buf,c)
char *buf;
unsigned c;
{
	unsigned xc;
	char kf;
	
	while(*buf){
		kf = 0;
		xc = (unsigned)(*buf);
		if (iskanji(*buf)){
			if (iskanji2(*(buf+1))){
				xc = (xc << 8) + (unsigned)(*(buf+1));
				kf = 1;
			}
		}
		if (xc == c){
			return(buf);
		}
		buf++;
		if (kf){
			buf++;
		}
	}
	return(NULL);
}

char *jstrrchr(buf,c)
char *buf;
unsigned c;
{
	unsigned xc;
	int i,max;

	max = strlen(buf);
	while(*(buf+1)){
		buf++;
	}
	for (i = 0;i < max;i++){
		xc = (unsigned)(*buf);
		if (iskanji(*buf)){
			if (iskanji2(*(buf+1))){
				xc = (xc << 8) + (unsigned)(*(buf+1));
			}
		}
		if (xc == c){
			return(buf);
		}
		buf--;
	}
	return(NULL);
}

int jstrcmp(s,t)
char *s;
char *t;
{
	unsigned sx,tx;

	while(*s){
		sx = (unsigned)(*s);
		if (iskanji(*s)){
			if (iskanji2(*(s+1))){
				sx = (sx << 8) + (unsigned)(*(s+1));
				s++;
			}
		}
		tx = (unsigned)(*t);
		if (iskanji(*t)){
			if (iskanji2(*(t+1))){
				tx = (tx << 8) + (unsigned)(*(t+1));
				t++;
			}
		}
		if (sx > tx){
			return(sx - tx);
		}
		if (tx > sx){
			return(-(int)(tx - sx));
		}
		s++;
		t++;
	}
	return(0);
}

int jstrncmp(s,t,n)
char *s;
char *t;
int n;
{
	unsigned sx,tx;

	for(;*s && n;n--){
		sx = (unsigned)(*s);
		if (iskanji(*s)){
			if (iskanji2(*(s+1))){
				sx = (sx << 8) + (unsigned)(*(s+1));
				s++;
			}
		}
		tx = (unsigned)(*t);
		if (iskanji(*t)){
			if (iskanji2(*(t+1))){
				tx = (tx << 8) + (unsigned)(*(t+1));
				t++;
			}
		}
		if (sx > tx){
			return(sx - tx);
		}
		if (tx > sx){
			return(-(int)(tx - sx));
		}
		s++;
		t++;
	}
	return(0);
}

int jstrlen(buf)
char *buf;
{
	int len;
	
	len = 0;
	
	while (*buf){
		if (iskanji(*buf)){
			if (iskanji2(*(buf+1))){
				buf++;
				len++;
			}else{
				len++;
			}
		}else{
			len++;
		}
		buf++;
	}
	return(len);
}

char *jstrmatch(s,t)
char *s;
char *t;
{
	unsigned sx,tx;
	char bf,*t2;

	while(*s){
		bf = 0;
		sx = (unsigned)(*s);
		if (iskanji(*s)){
			if (iskanji2(*(s+1))){
				sx = (sx << 8) + (unsigned)(*(s+1));
				bf = 1;
			}
		}
		t2 = t;
		while(*t2){
			tx = (unsigned)(*t2);
			if (iskanji(*t2)){
				if (iskanji2(*(t2+1))){
					tx = (tx << 8) + (unsigned)(*(t2+1));
					t2++;
				}
			}
			if (sx == tx){
				return(s);
			}
			t2++;
		}
		s++;
		if (bf){
			s++;
		}
	}
	return(NULL);
}

char *jstrncat(s,t,n)
char *s;
char *t;
int n;
{
	char *s2;
	s2 = s;

	while(*s2){
		s2++;
	}
	for(;n && *t;n--){
		*s2 = *t;
		if (iskanji(*t)){
			if (iskanji2(*(t+1))){
				s2++;
				t++;
				*s2 = *t;
			}
		}
		s2++;
		t++;
	}
	*s2 = '\0';
	return(s);

}

char *jstrncpy(s,t,n)
char *s;
char *t;
int n;
{
	char *s2;
	s2 = s;

	for(;n && *t;n--){
		*s2 = *t;
		if (iskanji(*t)){
			if (iskanji2(*(t+1))){
				s2++;
				t++;
				*s2 = *t;
			}
		}
		s2++;
		t++;
	}
	if (n > 0){
		*s2 = '\0';
	}
	return(s);

}

char *jstrrev(s) /* ½ºÞ¸ ¸Û³ ¼À */
char *s;
{

	unsigned bx;
	int i,j,len,en,bt;
	char *b,*b2,*s2,bf;

	b2 = s;
	s2 = s;
	while (*(b2+1)){
		b2++;
	}
	
	len = jstrlen(s);
	bt = strlen(s) - 1;
	en = 0;

	for (i = 0;i < (len-1);i++){
		bf = 0;
		b = jstradv(s,len-1); /* ‘O‚©‚çŒ©‚È‚¢‚Æ®‡«‚ðŽ¸‚¤ */
		bx = (unsigned)(*b);
		if (iskanji(*b)){
			if (iskanji2(*(b+1))){
				bx = (bx << 8) + (unsigned)(*(b+1));
				bf = 1;
			}
		}
		for (j = bt;j > en;j--){
			s[j] = s[j-1];
		}
		en++;
		if (bf){
			for (j = bt;j > en;j--){
				s[j] = s[j-1];
			}
			en++;
			*s2 = (char)(bx >> 8);
			s2++;
			*s2 = (char)bx;
			s2++;
		}else{
			*s2 = (char)bx;
			s2++;
		}
		/* printf("bf %d:%s\n",(int)bf,s); for debug */
	}
	
	return(s);
}

char *jstrskip(s,t)
char *s;
char *t;
{
	unsigned sx,tx;
	char bf,*t2;

	while(*s){
		bf = 0;
		sx = (unsigned)(*s);
		if (iskanji(*s)){
			if (iskanji2(*(s+1))){
				sx = (sx << 8) + (unsigned)(*(s+1));
				bf = 1;
			}
		}
		t2 = t;
		while(*t2){
			tx = (unsigned)(*t2);
			if (iskanji(*t2)){
				if (iskanji2(*(t2+1))){
					tx = (tx << 8) + (unsigned)(*(t2+1));
					t2++;
				}
			}
			if (sx == tx){
				break;
			}
			t2++;
		}
		if (*t2 == '\0'){
			return(s);
		}
		s++;
		if (bf){
			s++;
		}
	}
	return(s);
}

char *jstrstr(s,t)
char *s;
char *t;
{

	while(*s){
		if (jstrncmp(s,t,jstrlen(t)) == 0){
			return(s);
		}
		s++;
	}
	return(NULL);
}

char *jstrtok(buf,mask)
char *buf;
char *mask;
{
	static char *bp;
	char *m2,*rp,kf;
	unsigned bx,mx;
	
	if (buf != NULL){
		bp = buf;
	}
	
	for(;*bp != '\0';bp++){
		bx = *bp;
		if (iskanji(*bp)){
			if (iskanji2(*(bp+1))){
				bx = (bx << 8) + (unsigned)(bp+1);
				kf = 1;
			}
		}
		for (m2 = mask;*m2;m2++){
			mx = *m2;
			if (iskanji(*m2)){
				if (iskanji2(*(m2+1))){
					mx = (mx << 8) + (unsigned)(m2+1);
					mx++;
				}
			}
			if (bx == mx){
				break;
			}
		}
		if (*m2 == '\0'){
			rp = bp;
			break;
		}
		if (kf){
			bp++;
		}
	}
	if (*bp == '\0'){
		return(NULL);
	}

	bp = jstrmatch(bp,mask);
	if (bp != NULL){
		if (iskanji(*bp)){
			*bp = '\0';
			if (iskanji2(*(bp+1))){
				bp++;
			}
		}else{
			*bp = '\0';
		}
		bp++;
	}
	
	return(rp);
}

/*
int main()
{
	char buf[100];
	
	while(gets(buf) != NULL){
		printf("%d:%d:%s\n",jstrlen(buf),strlen(buf),jstrrev(buf));
	}

	return(0);
}
*/
