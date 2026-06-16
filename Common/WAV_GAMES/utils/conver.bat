@echo off
echo          Conversor de imagenes en Screen 2 a WAV
echo        -------------------------------------------
echo  .
echo  Autor: Angel Manuel G¢mez Garcia
echo  Para m s informacion visita:
echo         www.msxtape.tk  ¢  www.msxtape.cjb.net
echo  .
echo  Convirtiendo fichero %1.bin
util\msxf2w %1.bin %1.wav -2
util\wavglue %1_final.wav util\loader.wav %1.wav util\mostrar.wav -p
del %1.wav

