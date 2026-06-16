package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// InitDB lê os arquivos de 1 a 10 e cria tabelas correspondentes no SQLite
func InitDB() error {
	logDebug("Banco: Iniciando InitDB (limpando e reinicializando tabelas consolidada e discos 1 a 10)...")
	// Garante que o diretório data existe
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório data: %w", err)
	}

	// Abre/Cria o banco de dados sqlite
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir banco de dados: %w", err)
	}
	defer db.Close()

	// Recria a tabela consolidada msxstuffs
	_, err = db.Exec("DROP TABLE IF EXISTS msxstuffs")
	if err != nil {
		return fmt.Errorf("erro ao remover tabela consolidada msxstuffs: %w", err)
	}

	// Cria a tabela configuracoes se não existir
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS configuracoes (
		chave TEXT PRIMARY KEY,
		valor TEXT
	)`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela configuracoes: %w", err)
	}

	rootDir, err := filepath.Abs(".")
	if err != nil {
		rootDir = "."
	}

	// Insere configurações padrão se não existirem
	defaultConfigs := map[string]string{
		"theme":             "One Dark",
		"language":          "en",
		"emulador":          "openmsx",
		"options":           "-machine Gradiente_Expert_GPC-1 -ext DDX_3.0",
		"volume":            "50",
		"debug":             "false",
		"raiz":              rootDir,
		"pictures":          filepath.Join(rootDir, "pictures"),
		"dsks":              filepath.Join(rootDir, "DSK"),
		"roms":              filepath.Join(rootDir, "ROM"),
		"msxmania":          filepath.Join(rootDir, "Common", "MSX_MANIA"),
		"msxmania_pictures": filepath.Join(rootDir, "pictures", "msxmania", "MSX"),
		"database":          filepath.Join(rootDir, "data", "msxstuff.db"),
		"volume_inicial":    "1",
		"openmsx":           "",
		"fmsx":              "",
		"bluemsx":           "",
		"rumsx":             "",
		"openmsx_maquina":   "Gradiente_Expert_GPC-1",
		"openmsx_opcoes":    "",
		"openmsx_extensao1": "DDX_3.0",
		"openmsx_extensao2": "",
		"openmsx_extensao3": "",
		"openmsx_extensao4": "",
		"fmsx_opcoes":       "",
		"bluemsx_maquina":   "",
		"bluemsx_opcoes":    "",
		"rumsx_opcoes":      "",
	}

	for k, v := range defaultConfigs {
		_, err = db.Exec("INSERT OR IGNORE INTO configuracoes (chave, valor) VALUES (?, ?)", k, v)
		if err != nil {
			return fmt.Errorf("erro ao inicializar configuração padrão %s: %w", k, err)
		}
	}

	_, err = db.Exec(`CREATE TABLE msxstuffs (
		disco TEXT,
		descricao TEXT,
		cdnumero INTEGER,
		raiz TEXT,
		tipo TEXT,
		emulador TEXT,
		options TEXT
	)`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela consolidada msxstuffs: %w", err)
	}

	// Itera sobre os CDs de 1 a 10
	for i := 1; i <= 10; i++ {
		tableName := fmt.Sprintf("disco%d", i)

		// Nomes possíveis do arquivo devido a variações de maiúsculas/minúsculas
		fileName1 := fmt.Sprintf("MSXSTUFF%d.txt", i)
		fileName2 := fmt.Sprintf("Msxstuff%d.txt", i)

		filePath := filepath.Join("data", fileName1)
		file, err := os.Open(filePath)
		if err != nil {
			// Tenta com o outro nome se falhar
			filePath = filepath.Join("data", fileName2)
			file, err = os.Open(filePath)
			if err != nil {
				return fmt.Errorf("erro ao abrir arquivo para o CD %d (tentou %s e %s): %w", i, fileName1, fileName2, err)
			}
		}

		// Recria a tabela para ser idempotente
		_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
		if err != nil {
			file.Close()
			return fmt.Errorf("erro ao remover tabela %s: %w", tableName, err)
		}

		_, err = db.Exec(fmt.Sprintf("CREATE TABLE %s (disco TEXT, descricao TEXT, cdnumero INTEGER)", tableName))
		if err != nil {
			file.Close()
			return fmt.Errorf("erro ao criar tabela %s: %w", tableName, err)
		}

		// Transação para inserção rápida
		tx, err := db.Begin()
		if err != nil {
			file.Close()
			return fmt.Errorf("erro ao iniciar transação para %s: %w", tableName, err)
		}

		stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (disco, descricao, cdnumero) VALUES (?, ?, ?)", tableName))
		if err != nil {
			tx.Rollback()
			file.Close()
			return fmt.Errorf("erro ao preparar insert em %s: %w", tableName, err)
		}

		stmtConsolidated, err := tx.Prepare("INSERT INTO msxstuffs (disco, descricao, cdnumero, raiz, tipo, emulador, options) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			stmt.Close()
			tx.Rollback()
			file.Close()
			return fmt.Errorf("erro ao preparar insert consolidado para %s: %w", tableName, err)
		}

		scanner := bufio.NewScanner(file)
		isFirstLine := true
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			// Ignora a primeira linha se for disco=descricao
			if isFirstLine {
				isFirstLine = false
				if strings.ToLower(line) == "disco=descricao" {
					continue
				}
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue // Ignora linhas inválidas
			}

			discoVal := strings.TrimSpace(parts[0])
			descricaoVal := strings.TrimSpace(parts[1])

			// Calcula os campos adicionais
			ext := filepath.Ext(discoVal)
			raiz := strings.TrimSuffix(discoVal, ext)
			tipo := strings.TrimPrefix(ext, ".")

			// Calcula as opções do emulador baseado no tipo / extensão
			optionsVal := ""
			upperTipo := strings.ToUpper(tipo)
			if upperTipo == "ROM" {
				optionsVal = "-machine Gradiente_Expert_GPC-1"
			} else if upperTipo == "DSK" {
				optionsVal = "-machine Gradiente_Expert_GPC-1 -ext DDX_3.0"
			} else {
				optionsVal = ""
			}

			_, err = stmt.Exec(discoVal, descricaoVal, i)
			if err != nil {
				stmt.Close()
				stmtConsolidated.Close()
				tx.Rollback()
				file.Close()
				return fmt.Errorf("erro ao inserir dados em %s: %w", tableName, err)
			}

			_, err = stmtConsolidated.Exec(discoVal, descricaoVal, i, raiz, tipo, "openmsx", optionsVal)
			if err != nil {
				stmt.Close()
				stmtConsolidated.Close()
				tx.Rollback()
				file.Close()
				return fmt.Errorf("erro ao inserir dados consolidados em msxstuffs para %s: %w", tableName, err)
			}
		}

		stmt.Close()
		stmtConsolidated.Close()
		if err := tx.Commit(); err != nil {
			file.Close()
			return fmt.Errorf("erro ao efetuar commit em %s: %w", tableName, err)
		}

		file.Close()
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("erro na leitura do scanner para o CD %d: %w", i, err)
		}
	}

	return nil
}

// GetConfig consulta o valor de uma configuração na tabela configuracoes
func GetConfig(key string) (string, error) {
	logDebug("GetConfig: Buscando configuração para a chave '%s'", key)
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var val string
	err = db.QueryRow("SELECT valor FROM configuracoes WHERE chave = ?", key).Scan(&val)
	if err != nil {
		logDebug("GetConfig: Chave '%s' não encontrada ou erro: %v", key, err)
		return "", err
	}
	logDebug("GetConfig: Chave '%s' retornou valor '%s'", key, val)
	return val, nil
}

// SetConfig grava ou atualiza o valor de uma configuração
func SetConfig(key string, val string) error {
	logDebug("SetConfig: Gravando chave '%s' com valor '%s'", key, val)
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Garante que a tabela existe
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS configuracoes (
		chave TEXT PRIMARY KEY,
		valor TEXT
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT OR REPLACE INTO configuracoes (chave, valor) VALUES (?, ?)", key, val)
	if err != nil {
		logDebug("SetConfig: Erro ao gravar chave '%s': %v", key, err)
	} else {
		logDebug("SetConfig: Chave '%s' gravada com sucesso.", key)
	}
	return err
}

// GetVolumes retorna os volumes (cdnumero) distintos existentes na tabela msxstuffs
func GetVolumes() ([]string, error) {
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Verifica se a tabela existe antes de consultar
	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='msxstuffs'").Scan(&tableExists)
	if err != nil || tableExists == 0 {
		// Se a tabela não existe ainda, retorna de 1 a 10 como padrão
		return []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, nil
	}

	rows, err := db.Query("SELECT DISTINCT cdnumero FROM msxstuffs ORDER BY cdnumero ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumes []string
	for rows.Next() {
		var vol int
		if err := rows.Scan(&vol); err != nil {
			return nil, err
		}
		volumes = append(volumes, fmt.Sprintf("%d", vol))
	}

	// Se não achou nada, retorna fallback
	if len(volumes) == 0 {
		return []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, nil
	}

	return volumes, nil
}

// Game representa um registro consolidado da tabela msxstuffs
type Game struct {
	Disco     string
	Descricao string
	CdNumero  int
	Raiz      string
	Tipo      string
	Emulador  string
	Options   string
}

// GetGamesByVolume busca os jogos do respectivo volume (cdnumero) filtrados por uma busca textual opcional
func GetGamesByVolume(category string, volume int, filter string) ([]Game, error) {
	logDebug("SQL query: Iniciando busca de jogos no SQLite. Categoria='%s', Volume=%d, Busca='%s'", category, volume, filter)
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var query string
	var args []any

	if category == "MSX Mania" {
		query = "SELECT disco, descricao, cdnumero, emulador, options FROM msxmania WHERE 1=1"
		if filter != "" {
			query += " AND (descricao LIKE ? OR disco LIKE ?)"
			likeFilter := "%" + filter + "%"
			args = append(args, likeFilter, likeFilter)
		}
		query += " ORDER BY cdnumero ASC, descricao ASC"
	} else if category == "Good MSX 1" {
		query = "SELECT disco, descricao, cdnumero, emulador, options FROM goodmsx1 WHERE 1=1"
		if filter != "" {
			query += " AND (descricao LIKE ? OR disco LIKE ?)"
			likeFilter := "%" + filter + "%"
			args = append(args, likeFilter, likeFilter)
		}
		query += " ORDER BY descricao ASC"
	} else {
		query = "SELECT disco, descricao, cdnumero, raiz, tipo, emulador, options FROM msxstuffs WHERE cdnumero = ?"
		args = append(args, volume)
		if filter != "" {
			query += " AND (descricao LIKE ? OR disco LIKE ?)"
			likeFilter := "%" + filter + "%"
			args = append(args, likeFilter, likeFilter)
		}
		query += " ORDER BY descricao ASC"
	}

	logDebug("SQL comando: %s com parâmetros: %v", query, args)
	rows, err := db.Query(query, args...)
	if err != nil {
		logDebug("SQL erro: %v", err)
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var g Game
		if category == "MSX Mania" || category == "Good MSX 1" {
			err = rows.Scan(&g.Disco, &g.Descricao, &g.CdNumero, &g.Emulador, &g.Options)
			if err != nil {
				return nil, err
			}
			ext := filepath.Ext(g.Disco)
			g.Raiz = strings.TrimSuffix(g.Disco, ext)
			g.Tipo = strings.TrimPrefix(ext, ".")
		} else {
			err = rows.Scan(&g.Disco, &g.Descricao, &g.CdNumero, &g.Raiz, &g.Tipo, &g.Emulador, &g.Options)
			if err != nil {
				return nil, err
			}
		}
		games = append(games, g)
	}

	logDebug("SQL sucesso: Busca retornou %d jogos da categoria '%s'", len(games), category)
	return games, nil
}

// InitGameEmulacaoTables cria as tabelas game_emulacao, game_emulador_detalhes e categoria se não existirem
func InitGameEmulacaoTables() error {
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS game_emulacao (
		cdnumero INTEGER,
		disco TEXT,
		execucoes INTEGER DEFAULT 0,
		emulador_escolhido TEXT DEFAULT 'openmsx',
		PRIMARY KEY(cdnumero, disco)
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS game_emulador_detalhes (
		cdnumero INTEGER,
		disco TEXT,
		emulador TEXT,
		chave TEXT,
		valor TEXT,
		PRIMARY KEY(cdnumero, disco, emulador, chave)
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS categoria (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT UNIQUE,
		ativo INTEGER DEFAULT 0
	)`)
	if err != nil {
		return err
	}

	_, _ = db.Exec("DELETE FROM categoria WHERE nome = 'Good MSX'")

	// Popula categorias padrão
	categorias := []struct {
		nome  string
		ativo int
	}{
		{"MSX Stuffs", 1},
		{"MSX Mania", 1},
		{"CAS Collection", 0},
		{"Good MSX 1", 1},
		{"Wave Games", 0},
		{"MSX Tools", 0},
		{"Nemesis Diskpack", 0},
	}

	for _, cat := range categorias {
		_, _ = db.Exec("INSERT INTO categoria (nome, ativo) VALUES (?, ?) ON CONFLICT(nome) DO UPDATE SET ativo = excluded.ativo", cat.nome, cat.ativo)
	}

	return nil
}

// GetOrCreateGameEmulacao retorna o contador de execuções e o emulador configurado do jogo. 
// Cria um registro padrão caso não exista.
func GetOrCreateGameEmulacao(cdNumero int, disco string) (int, string, error) {
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return 0, "", err
	}
	defer db.Close()

	var execucoes int
	var emuladorEscolhido string
	err = db.QueryRow("SELECT execucoes, emulador_escolhido FROM game_emulacao WHERE cdnumero = ? AND disco = ?", cdNumero, disco).Scan(&execucoes, &emuladorEscolhido)
	if err == sql.ErrNoRows {
		emuladorEscolhido = "openmsx"
		_, err = db.Exec("INSERT INTO game_emulacao (cdnumero, disco, execucoes, emulador_escolhido) VALUES (?, ?, 0, ?)", cdNumero, disco, emuladorEscolhido)
		if err != nil {
			return 0, "", err
		}
		return 0, emuladorEscolhido, nil
	} else if err != nil {
		return 0, "", err
	}

	return execucoes, emuladorEscolhido, nil
}

// SaveGameEmulacao salva as preferências principais de emulação para um jogo específico
func SaveGameEmulacao(cdNumero int, disco string, emuladorEscolhido string) error {
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO game_emulacao (cdnumero, disco, execucoes, emulador_escolhido) 
		VALUES (?, ?, 0, ?) 
		ON CONFLICT(cdnumero, disco) 
		DO UPDATE SET emulador_escolhido = excluded.emulador_escolhido`, 
		cdNumero, disco, emuladorEscolhido)
	return err
}

// IncrementGameExecution incrementa o número de execuções do jogo
func IncrementGameExecution(cdNumero int, disco string) error {
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Assegura que o registro de emulação exista primeiro
	_, _, err = GetOrCreateGameEmulacao(cdNumero, disco)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE game_emulacao SET execucoes = execucoes + 1 WHERE cdnumero = ? AND disco = ?", cdNumero, disco)
	return err
}

// GetGameEmuladorDetalhe busca um campo de detalhe customizado do emulador para o jogo
func GetGameEmuladorDetalhe(cdNumero int, disco string, emulador string, chave string) (string, error) {
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var val string
	err = db.QueryRow("SELECT valor FROM game_emulador_detalhes WHERE cdnumero = ? AND disco = ? AND emulador = ? AND chave = ?", 
		cdNumero, disco, emulador, chave).Scan(&val)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return val, err
}

// SetGameEmuladorDetalhe grava/atualiza um detalhe customizado do emulador para o jogo
func SetGameEmuladorDetalhe(cdNumero int, disco string, emulador string, chave string, valor string) error {
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO game_emulador_detalhes (cdnumero, disco, emulador, chave, valor) 
		VALUES (?, ?, ?, ?, ?) 
		ON CONFLICT(cdnumero, disco, emulador, chave) 
		DO UPDATE SET valor = excluded.valor`, 
		cdNumero, disco, emulador, chave, valor)
	return err
}

// HasGameEmulacao verifica se existe um registro configurado na tabela game_emulacao.
// Retorna true se encontrar o registro, o emulador escolhido, e erro se houver.
func HasGameEmulacao(cdNumero int, disco string) (bool, string, error) {
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return false, "", err
	}
	defer db.Close()

	var emuladorEscolhido string
	err = db.QueryRow("SELECT emulador_escolhido FROM game_emulacao WHERE cdnumero = ? AND disco = ?", cdNumero, disco).Scan(&emuladorEscolhido)
	if err == sql.ErrNoRows {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	}
	return true, emuladorEscolhido, nil
}

// Category representa uma linha da tabela categoria
type Category struct {
	ID    int
	Nome  string
	Ativo int
}

// GetCategories retorna todas as categorias ordenadas por id
func GetCategories() ([]Category, error) {
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nome, ativo FROM categoria ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Nome, &c.Ativo); err == nil {
			list = append(list, c)
		}
	}
	return list, nil
}

// ImportMSXMania le games.txt do diretorio do MSX Mania e importa para a tabela msxmania
func ImportMSXMania(msxmaniaDir string) error {
	logDebug("Banco: Iniciando ImportMSXMania do diretório '%s'", msxmaniaDir)
	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Recria a tabela msxmania
	_, err = db.Exec("DROP TABLE IF EXISTS msxmania")
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE msxmania (
		disco TEXT,
		descricao TEXT,
		cdnumero INTEGER,
		emulador TEXT,
		options TEXT
	)`)
	if err != nil {
		return err
	}

	filePath := filepath.Join(msxmaniaDir, "games.txt")
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO msxmania (disco, descricao, cdnumero, emulador, options) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}

		diskField := strings.TrimSpace(parts[1])
		descricao := strings.TrimSpace(parts[2])

		var disco string
		var cdNum int
		diskFieldUpper := strings.ToUpper(diskField)
		if strings.HasPrefix(diskFieldUpper, "SP") {
			disco = "mania" + strings.ToLower(diskField) + ".dsk"
			if diskFieldUpper == "SP1" {
				cdNum = 131
			} else if diskFieldUpper == "SP2" {
				cdNum = 132
			} else {
				cdNum = 131
			}
		} else {
			disco = "mania" + strings.ToLower(diskField) + ".dsk"
			fmt.Sscanf(diskField, "%d", &cdNum)
		}

		_, err = stmt.Exec(disco, descricao, cdNum, "openmsx", "-machine Gradiente_Expert_GPC-1 -ext DDX_3.0")
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ImportGoodMSX1 le as ROMs do diretorio Common/Good_MSX1_Roms e importa para a tabela goodmsx1
func ImportGoodMSX1(rootDir string) error {
	logDebug("Banco: Iniciando ImportGoodMSX1 do diretório raiz '%s'", rootDir)
	
	romsDir := filepath.Join(rootDir, "Common", "Good_MSX1_Roms")
	files, err := os.ReadDir(romsDir)
	if err != nil {
		return fmt.Errorf("não foi possível abrir o diretório %s: %w", romsDir, err)
	}

	dbPath := filepath.Join("data", "msxstuff.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Recria a tabela goodmsx1
	_, err = db.Exec("DROP TABLE IF EXISTS goodmsx1")
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE goodmsx1 (
		disco TEXT,
		descricao TEXT,
		cdnumero INTEGER,
		emulador TEXT,
		options TEXT
	)`)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO goodmsx1 (disco, descricao, cdnumero, emulador, options) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Controla duplicatas para adicionar sufixos %A, %B, %C etc.
	seenCount := make(map[string]int)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		ext := strings.ToLower(filepath.Ext(filename))
		if ext != ".rom" {
			continue
		}

		cleanedName := cleanROMName(filename)
		if cleanedName == "" {
			cleanedName = strings.TrimSuffix(filename, filepath.Ext(filename))
		}

		var finalName string
		count := seenCount[cleanedName]
		if count == 0 {
			finalName = cleanedName
			seenCount[cleanedName] = 1
		} else {
			// count = 1 -> %A, count = 2 -> %B, etc.
			suffixRune := 'A' + rune(count-1)
			finalName = fmt.Sprintf("%s %%%c", cleanedName, suffixRune)
			seenCount[cleanedName] = count + 1
		}

		logDebug("ImportGoodMSX1: Importando '%s' como '%s'", filename, finalName)

		_, err = stmt.Exec(filename, finalName, 0, "openmsx", "-machine Gradiente_Expert_GPC-1")
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func cleanROMName(filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	var sb strings.Builder
	inParen := 0
	inBracket := 0
	for _, r := range name {
		if r == '(' {
			inParen++
		} else if r == '[' {
			inBracket++
		} else if r == ')' {
			if inParen > 0 {
				inParen--
			}
		} else if r == ']' {
			if inBracket > 0 {
				inBracket--
			}
		} else {
			if inParen == 0 && inBracket == 0 {
				sb.WriteRune(r)
			}
		}
	}

	cleaned := sb.String()
	words := strings.Fields(cleaned)
	return strings.Join(words, " ")
}


