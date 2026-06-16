cpppp < %1.c > $.C
cf %2 $
ren $.tco %1.tco
del $.C
fpc %1 mlib lib
cg -k %3 %1
m80 =%1/z
del %1.mac
l80 ck,%1,mlib/s,clib/s,crun/s,cend,%1/n/e:xmain
