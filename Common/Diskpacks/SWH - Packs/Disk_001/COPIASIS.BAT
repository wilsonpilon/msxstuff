rem Este arquivo BATCH cria um
rem disco de sistema. Voce precisara
rem de um disco em BRANCO.
rem Para abortar, tecle CTRL-STOP
rem
rem 1) Este disco devera estar no
rem    drive A: e insira um disco em
rem    branco no drive B: (se voce
rem    tiver apenas 1 drive, siga as 
rem    instrucoes do sistema opera-
rem    cional que simulara a exis-
rem    tencia de dois drives, isto
rem    e, voce devera trocar os 
rem    discos a medida que o siste-
rem    ma for pedindo).
rem
rem 2) O comando a seguir formata
rem    o disco. Responda com "B" e
rem    digite o numero correspon-
rem    dente ao formato do disco a
rem    ser gravado (geralmente "2").
rem
FORMAT
rem
rem 2) Copia os arquivos do sistema.
rem    (Se voce tiver apenas um drive
rem    lembre-se de ir trocando os
rem    discos a medida que o sistema
rem    operacional ir pedindo.
rem
copy ddxdos.sys b:
copy comando.com b:
