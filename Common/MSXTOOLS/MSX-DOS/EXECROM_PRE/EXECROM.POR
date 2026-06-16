-------------------------------------------------------------------------
ExecROM vr. 0.92b
para o sistema MSX

(c) 1996,1997 A&L Software
-------------------------------------------------------------------------

        00. Indice
            ------

            01. O que e' o ExecROM
            02. O que voce precisa para executar arquivos ROM
            03. O que e' o cartucho MegaRAM
            04. Como usar o ExecROM
            05. Modificando cartuchos ROM com o ExecROM
            06. Solucionando problemas no ExecROM
            07. Creditos, uso e copia
        

        01. O que e' o ExecROM
            ------------------

        O ExecROM e' um programa que lhe permite executar arquivos ROM no
seu MSX1, MSX2, MSX2+ ou MSX Turbo R. Ele pode executar arquivos  ROM  de
16, 32, 48, 64 (?), 128, 256, 512 e 768Kb (?!).
        Um arquivo ROM e' apenas uma copia de um cartucho MSX  e pode ser
encontrado em varios sites na Internet. Um arquivo de imagem ROM pode ser
apenas um cartucho de 16k ou 32k ou ate' mesmo  um  cartucho  MegaROM  de
512k.
        Os arquivos ROM eram originalmente para se  usar  nos  emuladores
MSX,mas com o ExecROM voce pode executa-los agora em seu MSX REAL (que e'
bem melhor que um emulador :)!


        02. O que voce precisa para executar arquivos ROM
            -------------------------------------------

        Para executar arquivos ROM voce vai precisar:
        
        * do programa ExecROM;
        * do sistema operacional MSX-DOS;
        * de um computador MSX;
        * de um cartucho MegaRAM;

        Se voce quer executar cartuchos MegaROM or  usar  o  Game  Master
voce REALMENTE precisara' de uma MegaRAM.Mas se voce quer apenas executar
cartuchos de 16k ou 32k sem o Game Master, apenas um MSX, o MSX-DOS  e  o
ExecROM sao necessarios.
        O ExecROM e' um arquivo COM do MSX-DOS, que e' necessario para  a
execucao do programa. Voce nao pode executar o ExecROM sob o BASIC.
        Originalmente escrito para ser executado em um computador MSX,  o
ExecROM pode ser executado com relativo sucesso nos emuladores (e' claro,
voce nao pode executar jogos MegaRAM ou usar o Game Master por que os emu
ladores nao emulam a MegaRAM...),e foi testado com sucesso com o emulador
MSX CJS versoes 0.99b e 1.0b e com o emulador fMSX versoes 1.4 e 1.6 para
MS-DOS.


        03. O que e' o cartucho MegaRAM
            ---------------------------
        
        A MegaRAM e' um cartucho de memoria  parecido  com  os  cartuchos
MegaROM. A diferenca entre a MegaRAMs e os MegaROMs e' obvia: enquanto os
cartuchos MegaROM so' podem ser lidos, as MegaRAMs  podem  ser  lidas  ou
gravadas. Voce pode encontrar MegaRAMs com o tamanho de 256, 512 e 768Kb.
        Voce tambem pode encontrar a MegaRAM Disk, que e' uma MegaRAM com
a ROM de disco interna, usando a MegaRAM como uma RAM  Disk. Voce pode u-
sar a MegaRAM Disk como disckette ou como uma MegaRAM normal (tambem cha-
mada de MegaRAM Game).
        A MegaRAM e' uma invencao brasileira  e e' quase desconhecida  na
Europa, Japao e outros paises. Ao contrario da Memory Mapper,  a  MegaRAM 
nao e' reconhecida como a verdadeira expansao de memoria do MSX.


        04. Como usar o ExecROM
            -------------------

        E' muito facil usar o ExecROM. Para apenas  executar  um  arquivo 
no seu MSX, digite, no prompt do DOS,

        ExecROM arquivo

        onde 'arquivo' e' o nome do arquivo ROM que voce quer executar. A
extensao '.ROM' pode ser omitida.
        Apenas espere alguns segundos e o cartucho sera' executado. Todas
as acoes (deteccao, analise, carregamento, etc), incluindo a deteccao  da
MegaRAM, sao automaticas.
        O ExecROM possui algumas opcoes, que podem ser selecionadas antes
do nome do arquivo. Elas estao listadas abaixo.

opcao /P: muda um byte no cartucho. E' util em jogos para se conseguir vi
          das infinitas, imortalidade, etc. A sintaxe e':
          /Ppp,xxxx,yy
          onde:  pp - e' a pagina a ser mudada. E' apenas necessaria para
                      jogos MegaROM, que sao divididos em  paginas  (0  a
                      0Fh para cartuchos de 128k, 0 a 1Fh para  cartuchos
                      256k e de 0 a 3Fh para cartuchos de 512k).
                 xxxx - e' o endereco a ser modificado. O endereco  e'  o
                        deslocamento relativo ao comeco do cartucho.  Nas 
                        imagens MegaROM o endereco e' o deslocamento rela
                        tivo ao comeco da pagina. O comeco do cartucho  e
                        da pagina e' assumido como endereco 0.
                 yy - e' o byte a colocar no endereco.

opcao /F: preenche uma area do cartucho. E' util para jogos tambem. A sin
          taxe e':
          /Fpp,qq,xxxx,yy,zz
          onde:  pp - e' a pagina inicial a ser preenchida. E' necessaria
                      apenas para jogos MegaROM.
                 qq - e' a pagina final   a ser preenchida. E' necessaria
                      apenas para jogos MegaROM.
                 xxxx - e' o endereco inicial a ser preenchido. O endere-
                        co e' o deslocamento relativo ao inicio do cartu-
                        cho (nos cartuchos MegaROMs e' relativo ao comeco
                        da  pagina). O comeco do cartucho e da pagina  e'
                        assumido como endereco 0.
                 yyyy - e' o endereco  final  a ser preenchido. O endere-
                        co e' o deslocamento relativo ao inicio do cartu-
                        cho (nos cartuchos MegaROMs e' relativo ao comeco
                        da  pagina). O comeco do cartucho e da pagina  e'
                        assumido como endereco 0.
                 zz - e' o byte a preencher a area.
           
          Nota importante: na opcao /F a imagem MegaROM e'  tratada, ape-
                           sar da divisao em paginas, como linear. Assim,
                           se voce requisitar /F1,2,7,10,8 a pagina 1 se-
                           ra' preenchida do endereco 7 ao final  com   o
                           byte 8 e a pagina 2 sera' preenchida do  ende-
                           reco 0 ao 10 com o byte 8. Se voce  colocar  5
                           no lugar do 2 as paginas 2,3 e 4 serao comple-
                           tamente preenchidas com o byte 8.

                           Representacao hipotetica da opcao /F1,2,7,10,8
                           ... pag. 1                       pag. 2    ...
                           |--------------------|------------------------
                           |??888888888888888888|8888????????????????????
                           ||-|-----------------||--|--------------------
                            0 7                  0  10        (enderecos)

                           Representacao hipotetica da opcao /F1,5,7,10,8
                           ... pag. 1 ... pags. 2,3  e  4 ... pag. 5  ...
                           |-----------|-----|------|------|-------------
                           |??888888888|88888|888888|888888|8888?????????
                           ||-|--------|-----|------|------||--|---------
                            0 7              (enderecos)    0  10        

opcao /R: coloca assinatura de cartucho no arquivo ROM carregado.
          Se voce esta' carregando cartuchos de expansao do BASIC, progra
          mas residentes (Eddy2, HotAsm, MegaAssembler, Logo, etc) ou mes
          mo alguns jogos, voce deve usar a opcao /R. Isto e'  devido  ao 
          fato de que o ExecROM desabilita a assinatura de cartucho do ar
          quivo ao carrega-lo. A nao existencia da assinatura de cartucho 
          lhe permite restar seu MSX (e nao desliga-lo) sem voltar ao jo-
          go.

opcao /B: nao executa o arquivo diretamente, mas reseta o MSX para  isso.
          Alguns jogos e programas residentes precisam desta opcao.

opcao /H: mostra uma pequena ajuda do programa (em ingles).

opcao /G: carrega o Game Master I da Konami com o arquivo ROM. O  arquivo
          ROM do Game Master DEVE possuir o nome de GMASTER.ROM.
          ATENCAO: o Game Master da Konami so' funciona com os  jogos  da
                   Konami (que podem ser de qualquer  tamanho,  inclusive
                   de 16 e 32k).

OBSERVACOES: as paginas, enderecos e bytes DEVEM ser dados em formato he-
             xadecimal. Eles nao necessitam ter exatamente dois ou quatro
             bytes (por exemplo: 0Dh pode ser 0D ou D; 12Fh pode ser  12F
             ou 012F).

             apenas forneca indicacao de pagina em jogos MegaROM.  Cartu-
             chos com 16k or 32k nao requerem indicacao de pagina por se-
             rem tratados de forma linear.

Exemplo: para executar o jogo MegaROM GAME01.ROM, mudando o byte no ende-
         reco 3Eh para 64h e preenchendo a area de 1850h a 18A0h da pagi-
         na 5h com o byte 40h voce so' precisa digitar

         ExecROM /p3e,64 /f5,5,1850,18a0,40 game01

         para executar o jogo de 32Kb GAME02.ROM, mudando o byte no ende-
         reco 21h para 0Fh e preenchendo a area de 2320h a  30A0h  com  o 
         byte 10h voce so' precisa digitar
         
         ExecROM /p21,f /f2320,30a0,10 game02


        05. Modificando cartuchos ROM com o ExecROM
            ---------------------------------------

           Foi mostrado que o ExecROM pode mudar arquivos ROM, como jogos 
        para conseguir vidas infinitas, imunidade, etc. Mas a maioria des
        tas adulteracoes sao para jogos em formato binario, que e`  dife-
        rente do formato ROM.
           Para usar estas adulteracoes no ExecROM voce tem que  recalcu-
        lar o endereco a alterar. Segue um exemplo para  um  jogo  de 32k
        e que contem os procedimentos basicos para a adaptacao.
           Suponha que voce tenha um jogo de 32k em formato binario e seu
        carregador BASIC seja o  listado abaixo. Suponha que as linhas 30
        e 40 modificam o jogo para lhe dar imunidade. E suponha que  voce
        tem este mesmo jogo no formato ROM.
                
                10 cls:print"Carregando jogo..."
                20 bload"game-1.bin"
                25 ' * Inicio da alteracao para imunidade
                30 poke &ha320,&h30
                40 defusr=peek(&hfcbf)+256*peek(&hfcc0):a=usr(0)
                42 ' A linha 40 apenas executa o bloco 01
                45 ' * Fim da mudanca para imunidade
                50 bload"game-2.bin",r

           Se o jogo e' um cartucho de 16k ou 32k inicialmente voce  deve
        achar os bytes de identificacao de cartucho ("AB") na memoria. E'
        isso que faz a linha abaixo.

                bload"game-1.bin":for i=&h8000 to &hc000:if  peek(i)=&h41
                and peek(i+1)=&h42 then print"Inicio:&h";hex$(i):end else
                next

           Suponha que a rotina lhe de^ o endereco &h9000. Apenas subtra-
        ia esse endereco do endereco de imunidade para obter  o  endereco
        pro ExecROM. Neste exemplo: &hA320-&h9000=&h1320.  Apenas  digite
        no BASIC:

                print hex$(&hA320-&h9000)

           Agora digite no DOS:

                ExecROM /p1320,30 game.rom

           Mas se o endereco a mudar estiver depois do carregamento do se
        gundo bloco, como abaixo?

                10 cls:print"Carregando jogo..."
                20 bload"game-1.bin",r
                30 bload"game-2.bin"
                35 ' * Inicio da mudanca para imunidade
                40 poke &ha320,&h30
                45 ' * Fim da mudanca para imunidade
                50 defusr=peek(&hfcbf)+256*peek(&hfcc0):a=usr(0)
                52 ' A linha 50 apenas executa o bloco 02

           Neste caso o procedimento mais comum e' subtrair &h8000 do en-
        dereco a modificar. Neste exemplo obteriamos &h2320.
           Se nao funcionar, tente subtrair o endereco do comeco do bloco
        (que pode ser visto com qualquer leitor de cabecalho  de  arquivo
        binario) do endereco a modificar.
           Para jogos Memory Mapper adaptados para MegaRAM  nao  ha'  uma
        formula especifica por que os enderecos dos blocos mudam de  jogo
        para jogo e de bloco para bloco. Apenas um usario  experiente  de
        MSX pode faze-lo.
           Para jogos MegaRAM a regra geral e' subtrair do endereco a mo-
        dificar o comeco do bloco (use um leitor de cabecalho de  arquivo
        binario). Se os blocos sao de 8k de tamanho a pagina e' o  numero
        do bloco menos um (assumindo o primeiro bloco como bloco 1).Se os
        blocos sao de 16k, cada bloco contem 2 paginas. A pagina e' o do-
        bro do numero do bloco menos 2 (assumindo que o primeiro bloco e'
        o bloco 1). Se o endereco a modificar for maior que 2000h, incre-
        mente o numero da pagina e subtraia 2000h do endereco.


        06. Solucionando problemas no ExecROM
            ---------------------------------

        * ATENCAO:

            a) Nao execute o ExecROM com MegaRAM Disk no slot 1;
               se voce tem a MegaRAM Disk no slot 1, de boot pelo disco.

            b) Nao execute o ExecROM a partir da MegaRAM RAMDisk,  ou  um
               erro travara' seu computador.
               
            c) Geralmente, resete seu MSX com a tecla CONTROL apertada an
               tes de executar o ExecROM se voce tem mais de um drive.
               
        * SOLUCAO GERAL DE PROBLEMAS:
            
            Se voce esta tendo problemas com o ExecROM usando uma MegaRAM
            Disk e/ou tem mais que um drive conectado no seu  MSX,  tente 
            ligar seu MSX com a tecla DELETE apertada (para  desativar  a
            MegaRAM Disk e torna-la MegaRAM Game) e/ou a tecla CONTROL a-
            pertada (para desativar os drives logicos).
           
        * SOLUCAO DE PROBLEMAS ESPECIFICOS:
        
        a) Depois do carregamento o ExecROM trava o computador na  mensa-
           gem "Executing ROM cartridge"

           Alguns cartuchos  requerem algumas "condicoes especiais",  que
           podem ser providas peloExecROM, para executarem. Alguns preci
           sam do identificador de cartucho ou do reset do MSX ou de  am-
           bos. Estes casos sao resolvidos com as opcoes /R (que coloca o
           identificador de cartucho no arquivo ROM) e /B (que reseta seu
           MSX para executar o arquivo ROM).Tente usar uma das duas ou as
           duas opcoes com o arquivo ROM. Se ainda assim ele  nao  rodar,
           veja o "SOLUCAO GERAL DE PROBLEMAS".
           Aqui' esta' uma lista de cartuchos que requerem opcoes especi-
           ais para executar:
		                Cartucho                        Opcoes

 		                 Elevator Action                 /R
		                 Kung Fu Master                  /R /B
		                 Star Force                      /R /B

	b) Quando o arquivo ROM e' executado imagens  extranhas  aparecem
	   na tela e o computador trava
	   
	   Veja o topico anterior.
	   
	c) Cartuchos residentes nao sao habilitados quando carregados com
	   o ExecROM
 
 	   Para cartuchos residentes voce deve usar a opcao /R no ExecROM
 	   ou eles nao serao habilitados. Veja "Como usar o ExecROM" para
 	   detalhes.
 	   
 	d) As modificacoes (patchs) no arquivo ROM nao sao feitas
 
 	   i)  . A pagina, endereco e byte do patch estao corretas?
 	   ii) . Voce digitou os dados em formato hexadecimal?
 	   iii). Voce nao recebeu a mensagem "Error in patch request"?
 	   
 	e) Depois de um erro no ExecROM o computador trava
 	
 	   Voce provavelmente esta' usando uma MegaRAM Disk e executou  o
 	   ExecROM no drive da MegaRAM Disk. Execute o  ExecROM a  partir
 	   do diskette ou resete seu MSX com a tecla DELETE apertada.  Se
 	   o problema persistir, reinicie seu MSX a partir do   diskette, 
 	   nao a partir da MegaRAM Disk.
 	   
 	f) O MSX trava apos a identificacao do computador
 	
 	   Veja o topico anterior.
 	   
	g) O ExecROM e' reexecutado ou mensagens estranhas aparecem quan-
	   do  as  mensagens  "Stopping  the  drive"  ou  "Executing  ROM 
	   cartridge" aparecem

	   Veja "SOLUCAO GERAL DE PROBLEMAS".
	   		   
        h) O Game Master nao funciona com um jogo
        
 	   O Game Master so' funciona com jogos da Konami.
 	   
 	i) O ExecROM nao carrega o Game Master I
 	
 	   i)  . O arquivo do Game Master tem nome de GMASTER.ROM?       
	   ii) . Voce tem certeza que o arquivo e' o Game Master I?
          iii). Voce tem uma MegaRAM para usar o Game Master?
          

        07. Creditos, uso e copia
            -----------------------

        O programa ExecROM e' freeware, ou seja, e' gratis, voce nao pre-
cisa pagar por ele. Mas  voce so'pode usa-lo, copia-lo ou  distribui'-lo
se nao  remover  o  nome  do  autor  e  incluir sempre  os  3  arquivos:
ExecROM.COM, ExecROM.ENG e ExecROM.POR.
        Esta versao 0.92b e' totalmente funcional com todos os  tipos  de
ROM. Eu nao tive nenhum problema em todos os meus testes mas se voce des-
cobrir qualquer bug ou tiver qualquer ideia ou comentario, deixe-me sabe-
lo. Se um arquivo ROM nao funciona com voce (por favor, veja a solucao de
problemas primeiro!), conte-me e, se possivel, mande-me uma copia do  seu
ROM para que eu possaresolver o problema mais rapido.
        Eu gostaria de agradecer a todo o  pessoal  ativo  da  lista  MSX
Broadcast da Unesp, SP, Brasil pelos testes, ideias, bugs e comentarios.
	E eu gostaria de agradecer especialmente ao Marco  Antonio  Simon
Dal Poz pela sua ajuda sobre a MegaRAM e  sobre  o  bug  no  tamanho  das
MegaRAM Disk de 512 e 768k.
        O ExecROM e' copyright (c) 1996,1997 A&L Software. Contacte em:

        adrcunha@dcc.unicamp.br
        adrcunha@fee.unicamp.br
        http://www.dcc.unicamp.br/~adrcunha

        E lembre-se: MSX para sempre!