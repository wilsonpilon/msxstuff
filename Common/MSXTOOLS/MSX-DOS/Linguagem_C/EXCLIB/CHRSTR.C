#include	<stdio.h>
#include	<string.h>
#include	<ctype.h>

#pragma nonrec

char *strncpy(s1,s2,n)
char *s1;
char *s2;
int n;
{
	int i;
	for (i = 0;*(s2+i) && (i < n);i++){
		*(s1+i) = *(s2+i);
	}
	if (i < n){
		for (;i < n;i++){
			*(s1+i) = '\0';
		}
	}
	return (s1);
}

char *strncat(s1,s2,n)
char *s1;
char *s2;
int n;
{
	int i,s;
	s = strlen(s1);
	s1 += s;
	
	for (;*s2 && n;n--){
		*s1 = *s2;
		s1++;
		s2++;
	}
	*s1 = '\0';
	return (s1);
}

int     strncmp(s,t,n)
char    *s, *t;
int n;
{
	int i;
    for (; (*s == *t) && n; s++, t++ ,n-- ){
    	if (*s == '\0'){
    		return(0);
    	}
    }
	
	if (n){
		return((int)*s - (int)*t);
	}else{
		return(0);
	}
	
}

int     strcmpi(s, t)
char    *s, *t;
{
    for (; tolower(*s) == tolower(*t); s++, t++){
		if (*s == '\0'){
		   return (0);
		}
	}
    return ((int) tolower(*s) - (int) tolower(*t));
}

char *strchr(str,chr)
char *str;
char chr;
{

	while(*str){
		if (*str == chr){
			return(str);
		}
		str++;
	}
	if (*str == chr){
		return(str);
	}
	return(NULL);
}

char *strstr(s,s2)
char *s;
char *s2;
{
	while(*s){
		if (*s == *s2){
			if (strncmp(s,s2,strlen(s2)) == 0){
				return(s);
			}
		}
		s++;
	}
	return(NULL);

}

int     strncmpi(s,t,n)
char    *s, *t;
int n;
{
    for (; tolower(*s) == tolower(*t) && n; s++, t++ ,n-- ){
    	if (*s == 0){
    		return(0);
    	}
    }
	
	if (n){
		return((int)tolower(*s) - (int)tolower(*t));
	}else{
		return(0);
	}
	
}

char *strrev(str)
char *str;
{
	int i,len;
	char tmp;
	len = strlen(str);
	for (i = 0;i < (len / 2);i++){ /* 0123456 */
		tmp = str[i];
		str[i] = str[len-1-i];
		str[len-1-i] = tmp;
	}
	return (str);
}

char *strrchr(str,ch)
char *str,ch;
{
	int i;
	for (i = strlen(str);(*(str+i) != ch) && (i > -1);i--){}

	if (i > -1){
		return(str+i);
	}else{
		return(NULL);
	}
}

char *strupr(s)
char *s;
{
	char *s2;
	s2 = s;
	
	while(*s){
		*s = toupper(*s);
		s++;
	}
	return (s2);
	
}

char *strlwr(s)
char *s;
{
	char *s2;
	s2 = s;
	
	while(*s){
		*s = tolower(*s);
		s++;
	}
	return (s2);
	
}

char *strtok(buf,mask)
char *buf;
char *mask;
{
	static char *bp;
	char *m2,*rp;
	
	if (buf != NULL){
		bp = buf;
	}
	
	for(;*bp != '\0';bp++){
		for (m2 = mask;*m2 != '\0';m2++){
			if (*bp == *m2){
				break;
			}
		}
		if (*m2 == '\0'){
			rp = bp;
			break;
		}
	}
	if (*bp == '\0'){
		return(NULL);
	}

	for(;*bp != '\0';bp++){
		for (m2 = mask;*m2 != '\0';m2++){
			if (*bp == *m2){
				*bp = '\0';
				break;
			}
		}
		if (*m2 != '\0'){
			bp++;
			break;
		}
	}
	return(rp);
}

char *strpbrk(s1,s2)
char *s1;
char *s2;
{
	char *sw;
	for (;*s1 != '\0';s1++){
		for (sw = s2;*sw != '\0';sw++){
			if (*s1 == *sw){
				return(s1);
			}
		}
	}
	return(NULL);
}

char *strnset(str,ch,n)
char *str;
char ch;
int  n;
{
	char *ptr;
	ptr = str;
	while(*ptr && n--){
		*ptr = ch;
		ptr++;
	}
	return(str);
}

char *strset(str,ch)
char *str;
char ch;
{
	char *ptr;
	ptr = str;
	while(*ptr){
		*ptr = ch;
		ptr++;
	}
	return(str);
}

char *strdup(s)
char *s;
{
	char *buf;
	buf = malloc(strlen(s)+1);
	if (buf != NULL){
		strcpy(buf,s);
	}
	return(buf);
}

char *stpcpy(s,t)
char *s;
char *t;
{
	while(*t){
		*s = *t;
		s++;
		t++;
	}
	*s = *t;
	return(s);
}

char *memccpy(dest,src,c,n)
char *dest;
char *src;
char c;
unsigned n;
{
	while(n > 0){
		*dest = *src;
		if (*dest == c){
			return(dest+1);
		}
		dest++;
		src++;
		n--;
	}
	return(NULL);
}

char *memchr(addr,byte,count)
char *addr;
char byte;
size_t count;
{
	while(count > 0){
		if (*addr == byte){
			return(addr);
		}
		count--;
		addr++;
	}
	return(NULL);
}

int     memcmp(s,t,n)
char *s;
char *t;
int n;
{
	int i;
    for (;n > 0; s++, t++ ,n-- ){
    	if (*s != *t){
			return((int)*s - (int)*t);
    	}
	}
	return(0);
	
}

int     memicmp(s,t,n)
char *s;
char *t;
int n;
{
	int i;
    for (;n > 0; s++, t++ ,n-- ){
    	if (tolower(*s) != tolower(*t)){
			return((int)tolower(*s) - (int)tolower(*t));
    	}
	}
	return(0);
	
}

unsigned strspn(s,t)
char *s;
char *t;
{
	char *t2;
	unsigned len;
	len = 0;
	while(*s){
		t2 = t;
		while(*t2){
			if (*s == *t2){
				break;
			}
		}
		if (*t2 == '\0'){ /* ﾌｸﾏﾚﾅｲ */
			return(len);
		}
		s++;
		len++;
	}
	return(len);
}

unsigned strcspn(s,t)
char *s;
char *t;
{
	char *t2;
	unsigned len;
	len = 0;
	while(*s){
		t2 = t;
		while(*t2){
			if (*s == *t2){
				return(len);
			}
		}
		s++;
		len++;
	}
	return(len);
}

/* strerror()はMS-DOSに、movedata()は8086CPU のアーキテクチャーに */
/* 依存するため、本ライブラリパッケージでは実装を断念する。 */

/* ASCII character type functions. */
/* この３つの関数についてはMSX-C Ver.1.2でもこのマクロと同様の */
/* 処理が行われていると推測する。 */

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

char toascii(c)
char c;
{
	return(c & 0x7f);
}

char _tolower(c)
char c;
{
	return(c + 0x20);
}

char _toupper(c)
char c;
{
	return(c - 0x20);
}

