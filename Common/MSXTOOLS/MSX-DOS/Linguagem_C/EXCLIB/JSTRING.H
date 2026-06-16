#ifndef __JSTRING_H
#define __JSTRING_H

/* for poor linker */
#define jstradv jsadv
#define jstrchr jschr
#define jstrrchr jsrchr
#define jstrcmp jscmp
#define jstrncmp jsncmp
#define jstrlen jslen
#define jstrmatch jsmat
#define jstrncat jsncat
#define jstrncpy jsncpy
#define jstrrev jsrev
#define jstrskip jsskip
#define jstrstr jsstr
#define jstrtok jstok

int btom();
int mtob();
char *jstradv();
char *jstrchr();
char *jstrrchr();
int jstrcmp();
int jstrncmp();
int jstrlen();
char *jstrmatch();
char *jstrncat();
char *jstrncpy();
char *jstrrev(); /* ½ºŞ¸ ¸Û³ ¼À */
char *jstrskip();
char *jstrstr();
char *jstrtok();

#endif
