#include	<stdio.h>
#include	<glib.h>
#include	<uty.h>

#pragma nonrec

char inkey()
{
	if (chsns() == NO){
		return('\0');
	}else{
		return(chget());
	}
}

int maxpos(x,lim)
int x;
int lim;
{

	if (x < 0){
		x = lim;
	}else{
		if (x > lim){
			x = 0;
		}
	}
	return(x);
}

VOID kbcom(buf) /* 39 char */
char *buf;
{
	int i;
	for (i = 0;(i < 39) && (buf[i] != '\0');i++){
		*(char *)(0xfbf0 + i) = buf[i];
	}
	*(unsigned *)(0xf3fa) = 0xfbf0;
	*(unsigned *)(0xf3f8) = 0xfbf0 + i;
}

VOID kbcomr(buf) /* 38 char */
char *buf;
{
	int i;
	for (i = 0;(i < 38) && (buf[i] != '\0');i++){
		*(char *)(0xfbf0 + i) = buf[i];
	}
	*(char *)(0xfbf0 + i) = 13;
	*(unsigned *)(0xf3fa) = 0xfbf0;
	*(unsigned *)(0xf3f8) = 0xfbf0 + i + 1;
}

STATUS bsaves(str,s_adr,e_adr)
char *str;
unsigned s_adr;
unsigned e_adr;
{
	FD fp;
	char *buf;
	unsigned i,cnt;

	if ((buf = malloc(2048)) == NULL){
		return(ERROR);
	}
	if ((fp = creat(str)) == ERROR){
		free(buf);
		return(ERROR);
	}
	*buf = 0xfe;
	/* Header - check code 0xfe,start addr.,end addr.,run addr. - */
	if (write(fp,buf,1) < 1){ 
		close(fp);
		free(buf);
		return(ERROR);
	}
	if (write(fp,(char *)&s_adr,2) < 2){
		close(fp);
		free(buf);
		return(ERROR);
	}
	if (write(fp,(char *)&e_adr,2) < 2){
		close(fp);
		free(buf);
		return(ERROR);
	}
	*buf = 0;*(buf+1) = 0;
	if (write(fp,buf,2) < 2){
		close(fp);
		free(buf);
		return(ERROR);
	}
	/* Data */
	cnt = 0;
	setrd(s_adr);
	for (i = s_adr;i <= e_adr;i++){
		*(buf + cnt) = invdp();
		cnt++;
		if (i == e_adr || cnt == 2048){
			if (write(fp,buf,cnt) < (int)cnt){
				close(fp);
				free(buf);
				return(ERROR);
			}
			cnt = 0;
		}
	}
	close(fp);
	free(buf);
	return(OK);
}

STATUS bloads(str)
char *str;
{
	FD fp;
	char *buf;
	unsigned cnt2,cnt3,s_adr,e_adr;
	if ((buf = malloc(2048)) == NULL){
		return(ERROR);
	}
	if ((fp = open(str,0)) == ERROR){
		free(buf);
		return(ERROR);
	}
	*buf = 0xfe;
	/* Header - check code 0xfe,start addr.,end addr.,run addr. - */
	if (read(fp,buf,1) < 1){
		close(fp);
		free(buf);
		return(ERROR);
	}
	if (buf[0] != 0xfe){
		close(fp);
		free(buf);
		return(ERROR);
	}
	if (read(fp,(char *)&s_adr,2) < 2){
		close(fp);
		free(buf);
		return(ERROR);
	}
	if (read(fp,(char *)&e_adr,2) < 2){
		close(fp);
		free(buf);
		return(ERROR);
	}
	*buf = 0;
	*(buf+1) = 0;
	if (read(fp,buf,2) < 2){
		close(fp);
		free(buf);
		return(ERROR);
	}
	/* Data */
	setwrt(s_adr);
	while ((cnt2 = read(fp,buf,2048)) == 2048){
		for (cnt3 = 0;cnt3 < cnt2;cnt3++){
			outvdp(*(buf + cnt3));
		}
	}
	for (cnt3 = 0;cnt3 < cnt2;cnt3++){
		outvdp(*(buf + cnt3));
	}
	close(fp);
	free(buf);
	return(OK);
}
