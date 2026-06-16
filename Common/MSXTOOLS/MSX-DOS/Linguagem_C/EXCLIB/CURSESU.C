#include	<stdio.h>
#include	<msxbios.h>
#include	<curses.h>

#pragma nonrec

VOID offcur(),oncur();

VOID _boxtit(win,title)
WINDOW *win;
char *title;
{
	box(win,'|','-');
	mvwaddstr(win,0,1,title);
}

VOID _winclose(win)
WINDOW *win;
{
	wclear(win);
	wrefresh(win);
	delwin(win);
}

int menuwin(title,item,argc,wx,wy,sx,sy)
char *title;
char *item[];
int argc;
int wx;
int wy;
int sx; /* 5 ﾓｼﾞ ｲｼﾞｮｳ */
int sy; /* 4 ﾓｼﾞ ｲｼﾞｮｳ */
{
	int i,sel,rcx;
	char key;
	WINDOW *win;
	
	if ((win = newwin(sy,sx,wy,wx)) == NULL){
		return(-2); /* Error */
	}
	_boxtit(win,title);
	for (i = 0;i < argc;i++){
		mvwaddstr(win,i+1,2,item[i]);
	}
	sel = 0;
	rcx = sx - 2;
	mvwaddch(win,sel+1,1,'[');
	mvwaddch(win,sel+1,rcx,']');
	wrefresh(win);
	while(1){
		offcur();
		while(chsns() == 0){}
		key = chget();
		mvwaddch(win,sel+1,1,' ');
		mvwaddch(win,sel+1,rcx,' ');
		switch (key){
			case 13:
				_winclose(win);
				return(sel);
			break;
			case 27:
				_winclose(win);
				return(-1); /* Cancel */
			break;
			case 30: /* Up */
				if (sel == 0){
					sel = argc - 1;
				}else{
					sel--;
				}
			break;
			case 31: /* Down */
				sel++;
				if (sel == argc){
					sel = 0;
				}
			break;
			default:;
			break;
		}
		mvwaddch(win,sel+1,1,'[');
		mvwaddch(win,sel+1,rcx,']');
		wrefresh(win);
	}
}

int ynwin(title,wx,wy)
char *title;
int wx;
int wy;
{
	static int cx[4] = {1,5,6,10};
	int i,sel;
	char key;
	WINDOW *win;
	
	i = strlen(title);
	if (i < 10){
		i = 12;
	}else{
		i += 2;
	}

	if ((win = newwin(3,i,wy,wx)) == NULL){
		return(-2); /* Error */
	}
	_boxtit(win,title);
	mvwaddstr(win,1,2,"Yes");
	mvwaddstr(win,1,7,"No ");
	sel = 0;
	mvwaddch(win,1,cx[sel*2],'[');
	mvwaddch(win,1,cx[sel*2+1],']');
	wrefresh(win);

	while(1){
		offcur();
		while(chsns() == 0){}
		key = chget();
		mvwaddch(win,1,cx[sel*2],' ');
		mvwaddch(win,1,cx[sel*2+1],' ');
		switch (key){
			case 13:
				_winclose(win);
				return(sel);
			break;
			case 27:
				_winclose(win);
				return(-1); /* Cancel */
			break;
			case 29: /* Left */
				if (sel == 0){
					sel = 1;
				}else{
					sel--;
				}
			break;
			case 28: /* Right */
				if (sel == 1){
					sel = 0;
				}else{
					sel++;
				}
			break;
			default:;
			break;
		}
		mvwaddch(win,1,cx[sel*2],'[');
		mvwaddch(win,1,cx[sel*2+1],']');
		wrefresh(win);
	}
}

VOID _bkdel(str,x,max)
char *str;
int x;
int max;
{
	int i;
	for (i = x;i < max;i++){
		*(str+i) = *(str+i+1);
	}
}

/* function in jcurses.rel */
char _kcount();
#define iskanji iskan
#define iskanji2 iskan2
BOOL iskanji(),iskanji2();

char *jinputw(title,wx,wy,len,str)
char *title;
int wx;
int wy;
int len;
char *str;
{
	int i,sel,x,max;
	char key;
	WINDOW *win;
	

	if ((win = newwin(3,len+2,wy,wx)) == NULL){
		return(NULL); /* Error */
	}
	_boxtit(win,title);
	noecho();
	
	x = 0;
	max = 0;
	*str = '\0';
	wmove(win,1,x+1);
	wrefresh(win); /* 勝手にカーソルをオンにしてくれる。 */
	
	while(1){
		key = wgetch(win);
		switch (key){
			case 28:
				if (x < max){
					if (x < (max-1)){
						if (iskanji(*(str+x))){
							if (iskanji2(*(str+x+1))){
								x++;
							}
						}
					}
					x++;
				}
			break;
			case 29:
				if (x > 0){
					x--;
					if (x > 0){
						if (_kcount(str,x)){
							x--;
						}
					}
				}
			break;
			case 8:
				if ((x > 0) && (max > 0)){
					x--;
					_bkdel(str,x,max);
					max--;
					mvwaddch(win,1,max+1,' ');
					if ((x > 0) && (max > 0)){
						if (_kcount(str,x)){
							x--;
							_bkdel(str,x,max);
							max--;
							mvwaddch(win,1,max+1,' ');
						}
					}
				}
			break;
			case 127:
				if ((x < max) && (max > 0)){
					if (x < (max-1)){
						if (iskanji(*(str+x))){
							if (iskanji2(*(str+x+1))){
								_bkdel(str,x,max);
								max--;
								mvwaddch(win,1,max+1,' ');
							}
						}
					}
					_bkdel(str,x,max);
					max--;
					mvwaddch(win,1,max+1,' ');
				}
			break;
			case '\n':
				_winclose(win);
				return(str);
			break;
			case 27:
				_winclose(win);
				return(NULL);
			break;
			default:
				if ((max < (len-1)) && (key > 31)){
					for (i = max+1;i > x;i--){
						*(str+i) = *(str+i-1);
					}
					*(str+x) = key;
					x++;
					max++;
				}
			break;
		}
		mvwaddstr(win,1,1,str);
		wmove(win,1,x+1);
		wrefresh(win);
	}
}

char *inputw(title,wx,wy,len,str)
char *title;
int wx;
int wy;
int len;
char *str;
{
	int i,sel,x,max;
	char key;
	WINDOW *win;
	

	if ((win = newwin(3,len+2,wy,wx)) == NULL){
		return(NULL); /* Error */
	}
	_boxtit(win,title);
	noecho();
	
	x = 0;
	max = 0;
	*str = '\0';
	wmove(win,1,x+1);
	wrefresh(win); /* 勝手にカーソルをオンにしてくれる。 */
	
	while(1){
		key = wgetch(win);
		switch (key){
			case 28:
				if (x < max){
					x++;
				}
			break;
			case 29:
				if (x > 0){
					x--;
				}
			break;
			case 8:
				if ((x > 0) && (max > 0)){
					x--;
					_bkdel(str,x,max);
					max--;
					mvwaddch(win,1,max+1,' ');
				}
			break;
			case 127:
				if ((x < max) && (max > 0)){
					_bkdel(str,x,max);
					max--;
					mvwaddch(win,1,max+1,' ');
				}
			break;
			case '\n':
				_winclose(win);
				return(str);
			break;
			case 27:
				_winclose(win);
				return(NULL);
			break;
			default:
				if ((max < (len-1)) && (key > 31)){
					for (i = max+1;i > x;i--){
						*(str+i) = *(str+i-1);
					}
					*(str+x) = key;
					x++;
					max++;
				}
			break;
		}
		mvwaddstr(win,1,1,str);
		wmove(win,1,x+1);
		wrefresh(win);
	}
}

/*
char *mes[] = {
	"inputw",
	"jinputw",
	"Carp",
	"Swallows",
	"Tigers",
	"Baystars",
	"End"
};

int main()
{
	WINDOW *stw;
	int i;
	char s[50];
	
	if ((stw = initscr()) == NULL){
		fprintf(stderr,"Error!!\n");
		exit(1);
	}
	for (i = 0;i != 6;){
		i = menuwin("Your team?",mes,7,0,0,12,9);
		switch(i){
			case 0:
				if (inputw("Your opinion?",20,0,20,s) != NULL){
					mvaddstr(10,30,s);
					refresh();
				}
			break;
			case 1:
				if (jinputw("ご意見は？",20,0,20,s) != NULL){
					mvaddstr(10,30,s);
					refresh();
				}
			break;
			default:
				if ((i < 6) && (i > -1)){
					mvaddstr(0,20,mes[i]);
					refresh();
					switch(ynwin("Ok ?",20,5)){
						case 0:
							mvaddstr(10,20,"Yes !!");
						break;
						case 1:
							mvaddstr(10,20,"No ...");
						break;
						case -1:
							mvaddstr(10,20,"Cancel ");
						break;
						default:
							mvaddstr(10,20,"Error !!");
						break;
					}
					mvaddstr(10,30,mes[i]);
					refresh();
				}
			break;
		}
		
	}
	endwin();

}
*/
