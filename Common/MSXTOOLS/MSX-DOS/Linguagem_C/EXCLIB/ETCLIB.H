#ifndef __ETCLIB_H
#define __ETCLIB_H

typedef struct find_t{
	char drive;
	char name[8];
	char ext[3];
	char atr;
	char dummy[10];
	int time;
	int date;
	int clst;
	int size[2];
} DIR;

struct  dosdate_t       {
	char   day;       /* 日(1-31)        */
	char   month;     /* 月(1-12)        */
	unsigned int    year;      /* 年(1980-2099)   */
	char   dayofweek; /* 曜日(0:日曜...) */
};

struct  dostime_t       {
	char hour;    /* 時(0-23) */
	char minute;  /* 分(0-59) */
	char second;  /* 秒(0-59) */
	char hsecond; /* 1/100秒(0-99) */
};

struct  diskfree_t      {
	unsigned total_clusters;      /* 総クラスタ数             */
	unsigned avail_clusters;      /* 使用可能クラスタ数       */
	unsigned sectors_per_cluster; /* クラスタあたりのセクタ数 */
	unsigned bytes_per_sector;    /* セクタあたりのバイト数   */
};

typedef struct dpb{
	char drive;
	char id;
	int sec_s;
	char d_mask;
	char d_sft;
	char c_mask;
	char c_sft;
	int fi_fat;
	char fat_c;
	char ents;
	int fi_dat;
	int la_cls;
	char si_fat;
	int fi_dir;
	char *fatbuf;
} DPB;

VOID getdate(),gettime(),getdrive(),setdrive(),cputs(),swab();
STATUS setdate(),settime(),cprintf();
char findfirst(),findnext(),*calloc(),putch(),*cgets(),*getdpb(),*bsearch();
int cscanf();
unsigned rotl(),rotr(),getdiskfree();

#endif
