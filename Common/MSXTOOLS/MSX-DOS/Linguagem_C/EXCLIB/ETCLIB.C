#include	<stdio.h>
#include	<bdosfunc.h>
#include	<msxbios.h>
#include	<etclib.h>

#pragma nonrec

/* Dummy file for FPC.COM */
/* Library,for MSX-C ,by assembler */
/* ÇbÇ…íºÇ∑Ç∆Ç±Ç±Ç…ÇÕÇ¢ÇÈÇ◊Ç´ */
unsigned rotl(val,cnt)
unsigned val;
char cnt;
{
}

unsigned rotr(val,cnt)
unsigned val;
char cnt;
{
}

char *getdpb(drive)
char drive;
{
}
/* Dummy end */

/* Library,for MSX-C ,by MSX-C */
VOID getdate(date)
struct dosdate_t *date;
{
	REGS r;

	r.bc = (NAT)0x2a;
	callx(0x5,&r);

	date->day = (char)(r.de & 0xff);
	date->month = (char)(r.de >> 8);
	date->year = (unsigned int)r.hl;
	date->dayofweek = r.a;

}

STATUS setdate(date)
struct dosdate_t *date;
{
	REGS r;

	r.bc = (NAT)0x2b;
	callx(0x5,&r);

	r.hl = date->year;
	r.de = (unsigned)((unsigned)date->month >> 8) + (unsigned)date->year;
	return(r.a);
}

VOID gettime(time)
struct dostime_t *time;
{
	REGS r;
	
	r.bc = (NAT)0x2c;
	callx(0x5,&r);
	
	time->hour = (char)(r.hl >> 8);
	time->minute = (char)(r.hl & 0xff);
	time->second = (char)(r.de >> 8);
	time->hsecond = (char)(r.de & 0xff);
}

STATUS settime(time)
struct dostime_t *time;
{
	REGS r;

	r.bc = (NAT)0x2d;
	callx(0x5,&r);

	r.hl = (unsigned)((unsigned)time->hour >> 8) + (unsigned)time->minute;
	r.de = (unsigned)((unsigned)time->second >> 8) + (unsigned)time->hsecond;
	return(r.a);
}

VOID _setmk(buf,mask)
char *buf;
char *mask;
{
	char *w,*m2;
	int i,j;

	w = mask;
	m2 = mask;

	/* Drive name */
	while(*w){
		if (*w == ':'){
			*buf = toupper(*(w-1)) - 0x40; /* MSX-C's function in stdio.h */
			m2 = w + 1;
			break;
		}
		w++;
	}

	/* Space clear */
	buf++;
	w = buf;
	for (i = 0;i < 11;i++){
		*w = ' ';
		w++;
	}

	/* Set file name */
	for (i = 0;i < 8;i++){
		switch(*m2){
			case '*':
				for (j = i;j < 8;j++){
					buf[j] = '?';
				}
				i = 7;
				m2++;
				while(*m2){
					if (*m2 == '.'){
						m2++;
						break;
					}
					m2++;
				}
			break;
			case '.':
				i = 7;
				m2++;
			break;
			case '\0':
				return;
			break;
			default:
				buf[i] = *m2;
				m2++;
			break;
		}
	}
	for (i = 8;i < 11;i++){
		switch(*m2){
			case '*':
				for (j = i;j < 11;j++){
					buf[j] = '?';
				}
				return;
			break;
			case '.':
				i--;
				m2++;
			break;
			case '\0':
				return;
			break;
			default:
				buf[i] = *m2;
				m2++;
			break;
		}
	}
}

char	findfirst(mask,dir)
char *mask;
DIR *dir;
{
	char mask2[20],fcb[37],*tmp;
	int i;
	for (i = 0;i < 37;i++){
		fcb[i] = 0;
	}

	_setmk(fcb,mask);

	bdos(SETDMA,(int)dir,0);
	return(bdos(SEARF,(int)fcb,0));
}

char	findnext(dir)
DIR *dir;
{

	bdos(SETDMA,(int)dir,0);
	return(bdos(SEARN,0,0));
}

/* void *Çchar *Ç≈ë„ópÇµÇΩÇ∆Ç±ÇÎÇ™åãç\Ç†ÇÈÅBãÍÇµÇ¢ÅB */
char *calloc(n,s)
unsigned n;
unsigned s;
{
	unsigned i,j;
	char *ret;
	if (!n || !s)
		return(NULL);

	i = n * s;
	ret = (char *)malloc(i);
	for (j = 0;j < i;j++){
		*(ret + j) = 0;
	}
	return(ret);
}

/* ç\ë¢ëÃÇªÇÃÇ‡ÇÃÇï‘ÇπÇ»Ç¢îÒÇ`ÇmÇrÇhÇbÇÃêßå¿è„é¿ëïÇ≈Ç´Ç»Ç¢ä÷êîÇ… */
/* div()Ç™Ç†ÇÈÅB */

VOID getdrive(drive)
unsigned *drive;
{
	char a;
	a = bdos((char)0x19,0,0);
	*drive = (unsigned)a;
}

VOID setdrive(drive,drives)
unsigned drive;
unsigned *drives;
{
	char a;
	int hl,i;
	a = bdos((char)0x0e,(int)drive,0);
	
	hl = bdosh((char)0x18,0,0);
	*drives = 0;
	for (i = 0;i < 16;i++){
		if (hl & 0x0001){
			(*drives)++;
		}
		hl = hl >> 1;
	}
}

char putch(c)
char c;
{
	char a;
	if (c != 0xff){
		a = bdos((char)0x06,(int)c,0);
	}
	return(c);
}

/* cprintf(),cscanf()ÇÃstdin,stdoutÇÕí†êKçáÇÌÇπÇÃÉ_É~Å[Ç¡ÉXÅB */
STATUS  _spr();

STATUS _cputx(c,fp)
char c;
FILE *fp;
{
	if (c == '\n'){
		chput(0x0d);
		chput(0x0a);
	}else{
		chput(c);
	}
	return(OK);
}

STATUS  cprintf(nargs, format)
int     nargs;
char    *format;
{
    return (_spr(&format, _cputx, stdout));
}

int     _scn();

int _cgetx(fp)
FILE *fp;
{
	return((int)chget());
}

int _nungc(c,fp)
char c;
FILE *fp;
{
	return(1);
}

int     cscanf(nargs, format)
int     nargs;
char    *format;
{
    return (_scn(&format, _cgetx, stdin, _nungc));
}

char *cgets(s)
char *s;
{
	char a;
	int i;
	
	a = bdos((char)0x0a,(int)s,0);
	for (i = 0;i < s[1];i++){
		if (s[i+2] == 0x0d){
			s[i+2] = '\0';
			break;
		}
	}
	
	return(&s[2]);
}

VOID cputs(s)
char *s;
{
	while(*s){
		chput(*s); /* Ç±ÇÍÇ»ÇÁê‚ëŒÉRÉìÉ\Å[ÉãÇ…èoóÕÇ∑ÇÈÇµÅA'\n'ÇÃïœä∑Ç‡ÇµÇ»Ç¢ */
		s++;
	}
}

VOID swab(s,t,n)
char *s;
char *t;
int n;
{
	while(n > 0){
		*t = *(s+1);
		*(t+1) = *s;
		s += 2;
		t += 2;
		n -= 2;
	}
}

/*  chmod()ÇÕÇlÇrÇwÅ|ÇcÇnÇrÇPÇ…ÉtÉ@ÉCÉãëÆê´Ç™ñ≥Ç¢Ç©ÇÁå©ëóÇË */

unsigned getdiskfree(device,diskspace)
unsigned device;
struct diskfree_t *diskspace;
{
	DPB *dpb;

	dpb = (DPB *)getdpb((char)(device));
	if (dpb == NULL){
		return(1);
	}
	diskspace->total_clusters = (unsigned)(dpb->la_cls - 1);
	diskspace->avail_clusters = (unsigned)bdosh(0x1b,device,0);
	diskspace->sectors_per_cluster = (unsigned)(dpb->c_mask + 1);
	diskspace->bytes_per_sector = (unsigned)(dpb->si_fat);
	return(0);
}

/* ÇQï™íTçıÅBä»åâÇ≥ÇèdéãÇµÇƒçƒãAåƒèoÇégÇ¡ÇƒâÇ¢ÇøÇ·Ç¡ÇΩÅB */
/* Ç≈Ç‡ÅAâÇ¢ÇΩêlÇ™èüÇøÇ»ÇÃÇ≥ÅBÇ÷Ç÷Å[ÇÒÅB */
#pragma recursive
char *bsearch(key,base,membs,size,compar)
char *key;
char *base;
unsigned membs;
unsigned size;
int     (*compar)();
{
	unsigned md,ls;
	int cp;
	
	md = membs / 2;
	ls = membs;
	switch(membs){
		case 1:
			cp = (*compar)(base,key);
			if (cp == 0){
				return(base);
			}
		break;
		case 2:
			cp = (*compar)(base,key);
			if (cp == 0){
				return(base);
			}else{
				cp = (*compar)(base+size,key);
				if (cp == 0){
					return(base+size);
				}
			}
		break;
		default:
			cp = (*compar)(base+md*size,key);
			if (cp == 0){
				return(base+size*md);
			}
			if (cp > 0){ /* key < */
				ls = ls - md - (ls % 2);
				return(bsearch(key,base,ls,size,compar));
			}else{ /* key > */
				ls = ls - md - 1;
				base = base + (md+1) * size;
				return(bsearch(key,base,ls,size,compar));
			}
		break;
	}

	return(NULL);

}

