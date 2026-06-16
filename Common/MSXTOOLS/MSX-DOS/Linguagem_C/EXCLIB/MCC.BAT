cf %2 %1
fpc %1 exlib mlib lib
cg -k %3 %1
rem Xela Soft's Optimizer ¦ Â¶³ ËÄ Ê REM ¦ ÄÙ
rem optd1 %1.mac temp.mac
rem del %1.mac
rem ren temp.mac %1.mac
m80 =%1/z
del %1.mac
l80 ck,%1,exlib/s,mlib/s,clib/s,crun/s,cend,%1/n/e:xmain
