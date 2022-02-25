package film

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
	"log"
	"strings"
)

func GetListOfFilmIDs(db *sql.DB) []string {
	var id string
	var listOfIDs []string
	rows, err := db.Query("select id from films")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)
		if err == nil {
			listOfIDs = append(listOfIDs, id)
		}
	}
	return listOfIDs
}
func GetListOfFilmIDsByCategory(categories string, db *sql.DB) []string {
	var id string
	var listOfIDs []string
	listOfCategories := strings.Split(categories, ",")
	qms, listOfCategoriesInterface := func() (string, []interface{}) {
		out := ""
		outInterface := []interface{}{}

		for i := 1; i <= len(listOfCategories); i++ {
			tmp := listOfCategories[i-1]
			outInterface = append(outInterface, tmp)
			//out += fmt.Sprintf("@p%d,", i)	-> azure sql
			out += "?,"
		}
		return out[:len(out)-1], outInterface
	}()

	q := "SELECT id FROM films WHERE category IN ( " + qms + " );" //TODO mssql: The data types text and nvarchar are incompatible in the equal to operator.
	rows, err := db.Query(q, listOfCategoriesInterface...)
	if err != nil {
		log.Println("getListOfFilmIDsByCategory -> \n", err)
		return listOfIDs
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)
		if err == nil {
			listOfIDs = append(listOfIDs, id)
		}
	}
	return listOfIDs
}

func GetListOfFilmCategories(db *sql.DB) []string {
	var name string
	var listOfCategories []string
	rows, err := db.Query("select name from categories")
	if err != nil {

		log.Println("Error reading film categories, returning empty,\nerror msg -> ", err)
		return listOfCategories
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name)
		if err == nil {
			listOfCategories = append(listOfCategories, name)
		}
	}
	return listOfCategories
}
func PostNewFilmCategory(category string, db *sql.DB) error {
	value := strings.ToLower(category)
	_, err := db.Query("insert into categories values ( ? )", value)
	return err
}
func GenerateUUID(table string, db *sql.DB) string {
	id := uuid.New().String()
	//id = "id123"
	row, err := db.Query("select id from "+table+" where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	row.Next()
	err = row.Scan(&id)
	if err == nil {
		id = GenerateUUID(table, db)
	}
	return id
}

type filmData struct {
	Name, Id, Owner, Description, Category, Quality string
	visible                                         bool
}

func (fd *filmData) Jsonify() []byte {
	var jsonData []byte
	jsonData, err := json.Marshal(fd)
	if err != nil {
		log.Println(err)
	}
	return jsonData
}

func NewFilmData(Name, Owner, Description, Category, Quality string, db *sql.DB) (filmData, error) {

	//var data filmData
	var err error
	id := GenerateUUID("films", db)

	data := filmData{Name: Name, Owner: Owner, Description: Description, Category: strings.ToLower(Category), Id: id, Quality: Quality, visible: true}

	return data, err
}

type RawFilmData struct {
	Name        string `json:"title"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Quality     string `json:"quality"`
}

func (r *RawFilmData) CreateFilmData(db *sql.DB) (filmData, error) {
	data, err := NewFilmData(r.Name, r.Owner, r.Description, r.Category, r.Quality, db)
	return data, err
}
func GetDataOfFilmByID(id string, db *sql.DB) (filmData, error) {
	var data filmData
	row, err := db.Query("select id, name, owner, description, category,quality, visible from films where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	row.Next()
	err = row.Scan(&data.Id, &data.Name, &data.Owner, &data.Description, &data.Category, &data.Quality, &data.visible)

	if err == nil {
		return data, nil
	}
	return data, errors.New("Film with id: " + id + " not found in database" + "		->		" + err.Error())
}

func PostFilmData(data filmData, db *sql.DB) error {
	var name string
	row, err := db.Query("select name from categories where name = ?", data.Category)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	row.Next()
	err = row.Scan(&name)
	if err != nil {
		err = PostNewFilmCategory(data.Category, db)
		if err != nil {
			return err
		}
		//na = generateUUID(table)
	}
	_, err = db.Query("insert into films values (?, ?, ?, ?, ?, ?,?)", data.Name, data.Id, data.Owner, data.Description, data.Category, data.visible, data.Quality)

	return err
}
func DeleteFilm(data filmData, db *sql.DB, serviceClient azblob.ServiceClient, ctx context.Context) []error {
	var errors []error
	filmQualityList := strings.Split(data.Quality, ",")
	for _, q := range filmQualityList {
		q = strings.Trim(q, " ")
		out := strings.Split(q, "-")
		resolution, bitrate := out[0], out[1]
		containerName := strings.ToLower(data.Id + "-" + strings.Replace(resolution, ":", "-", -1) + "-" + bitrate)
		container := serviceClient.NewContainerClient(containerName)
		_, err := container.Delete(ctx, nil)
		if err != nil {
			log.Println("deleteFilm() -> error deleting azr container : ", containerName)
			errors = append(errors, err)
		}

		containerManifests := serviceClient.NewContainerClient("manifests")
		manifestName := data.Id + ".m3u8"
		blobManifest := containerManifests.NewBlobClient(manifestName)
		_, err = blobManifest.Delete(ctx, nil)
		if err != nil {
			log.Println("deleteFilm() -> error deleting azr blob(manifest) : ", containerName)
			errors = append(errors, err)
		}

		containerPosters := serviceClient.NewContainerClient("posters")
		posterName := data.Id
		blobPoster := containerPosters.NewBlobClient(posterName)
		_, err = blobPoster.Delete(ctx, nil)
		if err != nil {
			log.Println("deleteFilm() -> error deleting azr blob(poster) : ", containerName)
			errors = append(errors, err)
		}

		_, err = db.Query("DELETE FROM films WHERE id = ?", data.Id)
		if err != nil {
			log.Println(err)
			errors = append(errors, err)
		}

	}
	return errors
}
