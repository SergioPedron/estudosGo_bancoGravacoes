package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

// Uma coleção de campos igual ao retorno do Select
type Album struct {
	ID      int64
	Titulo  string
	Artista string
	Preco   float32
}

// Por enquanto uma variável global de identificar do banco para ser utilizada em outas funções aqui neste mesmo package
var db *sql.DB

func main() {
	// Propriedades de conexão ao banco
	cfg := mysql.Config{
		User:   "sergio",
		Passwd: "sergio",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "gravacoes",
		//AllowNativePasswords: true,                          Conforme configuração do banco ...
		//User:   os.Getenv("DBUSER"),                         Quando utilizado Environment Variables:  $ export DBUSER=username
		//Passwd: os.Getenv("DBPASS"),                         O tutorial utilizava, mas não é necessário neste nosso exemplo
	}

	// Pega o identificar do banco
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	// Verifica a conexão
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Conectado ao banco de gravações!")

	// Exibe uma struct que foi retornada de um SQL pela função abaixo
	albuns, err := albunsPorArtista("Ira")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albuns encontrados: %v\n", albuns)

	// Pesquisa UM album por ID
	alb, err := albumPorID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album encotrado: %v\n", alb)

	novoAlbum := Album{Titulo: "Improvável Certeza", Artista: "Marcelo Bonfa", Preco: 30.12}
	id, err := adicionaAlbum(novoAlbum)
	if err != nil {
		log.Fatal("Erro ao incluir novo album: %w", err)

	}
	fmt.Printf("ID do novo Album cadastrado: %v\n", id)
}

//--------------------------------------------------------------------------------------------------------------------------------------------------//

// Função que recebe um nome de artista e retorna um slice de Album (struct de registro retornardos do banco) com os albuns encontrados.
func albunsPorArtista(nome string) ([]Album, error) {
	var albuns []Album

	// Para consultar várias linhas deve-se utilizar 'Query' que retorna rows de um select que são percorridas logo abaixo
	// Sempre utiliar parâmetros no SQL para evitar riscos de injeção de SQL
	rows, err := db.Query("Select * From album Where Artista = ?", nome)
	if err != nil {
		return nil, fmt.Errorf("albunsPorArtista %q: %v", nome, err)
	}
	defer rows.Close() // Garante o fechamento antes do retorno da função.  Seria um equivalente ao Try/Finally

	// Percorre as rows retornadas do select
	for rows.Next() {
		var alb Album
		// scan atribui os campos do select aos ponteiros que apontam para os campos da struct Album
		if err := rows.Scan(&alb.ID, &alb.Titulo, &alb.Artista, &alb.Preco); err != nil {
			return nil, fmt.Errorf("albunsPorArtista %q: %v", nome, err)
		}
		//
		albuns = append(albuns, alb) // adiciona a struct do select a slice dos albuns
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", nome, err)
	}
	return albuns, nil
}

//--------------------------------------------------------------------------------------------------------------------------------------------------//

// Função que retorna UM ALBUM a partir de um ID recebido como parâmetro.  Utiliza uma query de uma única linha.
func albumPorID(id int64) (Album, error) {
	var alb Album
	// QueryRow retorna uma única row.  Não retorna erro pois eles são postergados para o scan
	row := db.QueryRow("Select * From album Where Id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Titulo, &alb.Artista, &alb.Preco); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumPorID %d: não encontrou este album", id)
		}
		return alb, fmt.Errorf("albumsPorId %d: %v", id, err)
	}
	return alb, nil
}

//--------------------------------------------------------------------------------------------------------------------------------------------------//

// Adicionar um album recebido por parâmetro e retornar seu novo ID
func adicionaAlbum(alb Album) (int64, error) {
	result, err := db.Exec("Insert into album (Titulo, Artista, Preco) Values (?, ?, ?)", alb.Titulo, alb.Artista, alb.Preco)
	if err != nil {
		return 0, fmt.Errorf("adicionaAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("adicionaAlbum: %v", err)
	}
	return id, nil
}

//--------------------------------------------------------------------------------------------------------------------------------------------------//
