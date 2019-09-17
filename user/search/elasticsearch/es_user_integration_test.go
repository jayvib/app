// +build integration,elasticsearch

package elasticsearch_test

import (
	"context"
	"encoding/json"
	"github.com/jayvib/app/apperr"
	"github.com/jayvib/app/config"
	"github.com/jayvib/app/internal/elasticsearch/testutil"
	"github.com/jayvib/app/log"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/pkg/elasticsearch"
	"github.com/jayvib/app/pkg/validator"
	"github.com/jayvib/app/user"
	usersearches "github.com/jayvib/app/user/search/elasticsearch"
	"github.com/jayvib/app/utils/generateutil"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var client *elastic.Client

const (
	elasticsearchTestURL = "http://localhost:9200"
	index                = "user"
)

func init() {
	log.SetOutput(ioutil.Discard)
	os.Setenv(config.AppEnvironmentKey, config.DevelopmentEnv)
}

func TestIntegration_Delete(t *testing.T) {
	search := usersearches.New(client)
	filename := filepath.Join("testdata", "users.search.1.input")
	loadSampleDataFromFile(t, filename)
	id := "5d724acc90b3232c5337977c"
	err := search.Delete(context.Background(), id)
	require.NoError(t, err)
	_, err = client.Get().Index("user").Id(id).Do(context.Background())
	assert.True(t, elastic.IsNotFound(err))
}

func TestIntegration_SearchByName(t *testing.T) {
	t.Run("search name with size of 10", func(t *testing.T) {
		filename := filepath.Join("testdata", "users.search.1.input")
		loadSampleDataFromFile(t, filename)
		name := "Luffy"
		search := usersearches.New(client)
		size := 10
		res, next, err := search.SearchByName(context.Background(), name, 0, size)
		require.NoError(t, err)
		expectedResultSize := 5
		expectedNext := 0
		assert.Equal(t, expectedNext, next)
		assert.Equal(t, expectedResultSize, len(res))
	})

	t.Run("search name with size of 2", func(t *testing.T) {
		filename := filepath.Join("testdata", "users.search.1.input")
		loadSampleDataFromFile(t, filename)
		name := "Luffy"
		size := 2

		search := usersearches.New(client)
		var next int

		t.Run("expecting next of 3", func(t *testing.T) {
			res, n, err := search.SearchByName(context.Background(), name, 0, size)
			require.NoError(t, err)

			expectedNext := 2
			assert.Equal(t, expectedNext, n)

			expectedResultSize := 2
			assert.Equal(t, expectedResultSize, len(res))
			next = n
		})

		t.Run("expecting next of 4", func(t *testing.T) {
			// Already reach the last item.
			// Another fetch
			res, n, err := search.SearchByName(context.Background(), name, next, size)
			require.NoError(t, err)

			expectedNext := 4
			assert.Equal(t, expectedNext, n)

			expectedResultSize := 2 // just for the sake of readbility
			assert.Equal(t, expectedResultSize, len(res))
			next = n
		})

		t.Run("expecting next of 0", func(t *testing.T) {
			res, n, err := search.SearchByName(context.Background(), name, next, size)
			require.NoError(t, err)

			expectedNext := 0
			assert.Equal(t, expectedNext, n)

			expectedResultSize := 1 // just for the sake of readbility
			assert.Equal(t, expectedResultSize, len(res))
			next = n
		})
	})
}

func TestIntegration_Store(t *testing.T) {
	search := usersearches.New(client)
	newUser := &model.User{
		ID:        generateutil.GenerateID("user"),
		Firstname: "Usop",
		Lastname:  "god",
		Username:  "usop.god",
		Email:     "usup.god@onepiece.com",
		Password:  "bestsniper",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := search.Store(context.Background(), newUser)
	require.NoError(t, err)
	userGot := getUserByIdHelper(t, newUser.ID)
	assert.Equal(t, newUser.ID, userGot.ID)
	assert.Equal(t, userGot.Firstname, newUser.Firstname)
	assert.Equal(t, userGot.Lastname, newUser.Lastname)
	assert.Equal(t, userGot.Email, newUser.Email)
	assert.Equal(t, userGot.Username, newUser.Username)
}

func TestIntegration_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		search := usersearches.New(client)
		filename := filepath.Join("testdata", "users.search.1.input")
		loadSampleDataFromFile(t, filename)
		id := "5d724acc90b3232c5337977c"
		u, err := search.GetByID(context.Background(), id)
		require.NoError(t, err)
		assert.NotNil(t, u)
	})

	t.Run("Not found", func(t *testing.T) {
		search := usersearches.New(client)
		id := "notfound"
		_, err := search.GetByID(context.Background(), id)
		assert.Error(t, err)
		aerr, ok := err.(apperr.Error)
		assert.True(t, ok)
		assert.Equal(t, apperr.NoItemFound, aerr.Code())
	})
}

func TestIntegration_Update(t *testing.T) {
	filename := filepath.Join("testdata", "users.search.1.input")
	search := usersearches.New(client)
	loadSampleDataFromFile(t, filename)
	updatedUser := &model.User{
		ID:        "5d724acc90b3232c5337977c",
		Username:  "updated.username",
		UpdatedAt: time.Now(),
	}

	err := search.Update(context.Background(), updatedUser)
	require.NoError(t, err)

	resUser := getUserByIdHelper(t, updatedUser.ID)
	assert.Equal(t, updatedUser.Username, resUser.Username)

	t.Run("Empty ID", func(t *testing.T) {
		updatedUser := &model.User{
			Username:  "updated.username",
			UpdatedAt: time.Now(),
		}
		err := search.Update(context.Background(), updatedUser)
		assert.Error(t, err)
		if _, ok := err.(apperr.Error); ok {
			aerr := err.(apperr.Error)
			assert.Equal(t, aerr.Code(), apperr.ValidationErr)
			errs := aerr.OrigErrors()
			if assert.Len(t, errs, 1) {
				verr := errs[0]
				assert.IsType(t, validator.ValidationErr{}, verr)
				assert.Len(t, verr, 1)
			}
		} else {
			assert.True(t, ok, "expecting an err that implements an apperr.Error interface")
		}
	})
}

func TestIntegration_Search(t *testing.T) {
	filename := filepath.Join("testdata", "users.search.1.input")
	search := newSearchAndLoadSampleData(t, filename)

	t.Run("Search by firstname with a size of 2", func(t *testing.T) {
		query := "firstname=luffy"
		input := user.SearchInput{
			Query: query,
			Size:  2,
		}

		result, err := search.Search(context.Background(), input)
		require.NoError(t, err)
		require.Len(t, result.Users, 2)
		require.Equal(t, result.TotalHits, 5)
		require.Equal(t, result.Next, 2)

		t.Run("fetch another 2 results", func(t *testing.T) {
			input.From = result.Next
			result, err = search.Search(context.Background(), input)
			require.NoError(t, err)
			require.Len(t, result.Users, 2)
			require.Equal(t, result.TotalHits, 5)
			require.Equal(t, 4, result.Next)
			t.Run("fetch last 1 result", func(t *testing.T) {
				input.From = result.Next
				result, err = search.Search(context.Background(), input)
				require.NoError(t, err)
				require.Len(t, result.Users, 1)
				require.Equal(t, result.TotalHits, 5)
				require.Equal(t, 0, result.Next)
			})
		})
	})

	t.Run("Search by firstname with a size of 5", func(t *testing.T) {
		query := "firstname=luffy"
		input := user.SearchInput{
			Query: query,
			Size:  5,
		}

		result, err := search.Search(context.Background(), input)
		require.NoError(t, err)
		require.Len(t, result.Users, 5)
		assert.Equal(t, result.TotalHits, 5)
		assert.Equal(t, 0, result.Next)
	})

	t.Run("Missing query string", func(t *testing.T) {
		_, err := search.Search(context.Background(), user.SearchInput{})
		assert.Error(t, err)
		aerr, ok := err.(apperr.Error)
		if assert.True(t, ok, "expecting an error implements the apperr.Error interface") {
			// expecting that the code has an ValidationError
			assert.Equal(t, apperr.ValidationErr, aerr.Code())
		}

	})
}

func TestMain(m *testing.M) {
	//log.SetOutput(ioutil.Discard)
	var err error
	client, err = testutil.Setup(elasticsearchTestURL, index, "user.mapping.es.json")
	if err != nil {
		panic(err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func newSearchAndLoadSampleData(t *testing.T, filename string) *usersearches.Search {
	loadSampleDataFromFile(t, filename)
	return usersearches.New(client)
}

func getUserByIdHelper(t *testing.T, id string) *model.User {
	t.Helper()
	result, err := client.Get().Index("user").Id(id).Do(context.Background())
	require.NoError(t, err)

	if !result.Found {
		return nil
	}

	var u model.User
	err = json.Unmarshal(result.Source, &u)
	require.NoError(t, err)
	return &u
}

func loadSampleData() (err error) {
	client, err := elasticsearch.NewClient()
	if err != nil {
		return err
	}
	defer func() {
		log.Infof("Flushing")
		_, err = client.Flush("user").Do(context.Background())
	}()

	for _, u := range users {
		_, err = client.Index().Index("user").Id(u.ID).BodyJson(u).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func loadSampleDataFromFile(t *testing.T, filename string) {
	t.Helper()
	file, err := os.Open(filename)
	require.NoError(t, err)
	defer func() {
		err = file.Close()
		require.NoError(t, err)
		_, err = client.Flush("user").Do(context.Background())
		require.NoError(t, err)
	}()
	var users []*model.User
	err = json.NewDecoder(file).Decode(&users)
	require.NoError(t, err)
	for _, u := range users {
		_, err = client.Index().Index("user").Id(u.ID).BodyJson(u).Do(context.Background())
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second) // give ES a time to process
}

func writeUserIntoFile() (err error) {
	filename := filepath.Join("testdata", "users.json")
	bite, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bite, 0644)
}

var users = []*model.User{
	{
		ID:        "id1",
		Firstname: "Luffy",
		Lastname:  "Monkey",
		Email:     "luffy.monkey@onepiece.com",
		Username:  "luffy.monkey",
		Password:  "pirateking",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "id2",
		Firstname: "Luffy",
		Lastname:  "Monkeys",
		Email:     "luffy.monkeys@onepiece.com",
		Username:  "luffy.monkeys",
		Password:  "pirateking",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "id3",
		Firstname: "Luffy",
		Lastname:  "Monkeyz",
		Email:     "luffy.monkeyz@onepiece.com",
		Username:  "luffy.monkeyz",
		Password:  "pirateking",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "id4",
		Firstname: "Luffy",
		Lastname:  "Monkoyz",
		Email:     "luffy.monkoyz@onepiece.com",
		Username:  "luffy.monkoyz",
		Password:  "pirateking",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "id5",
		Firstname: "Luffy",
		Lastname:  "Monkayz",
		Email:     "luffy.monkayz@onepiece.com",
		Username:  "luffy.monkayz",
		Password:  "pirateking",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "id6",
		Firstname: "Sanji",
		Lastname:  "Vinsmoke",
		Email:     "luffy.monkey@onepiece.com",
		Username:  "sanji.vinsmoke",
		Password:  "bestcook",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "id7",
		Firstname: "Roronoa",
		Lastname:  "Zoro",
		Email:     "roronoa.zoro@onepiece.com",
		Username:  "roronoa.zoro",
		Password:  "bestswordsman",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
