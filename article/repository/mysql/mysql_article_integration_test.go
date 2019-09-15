// +build integration,mysql

package mysql_test

import (
	"context"
	"github.com/jayvib/clean-architecture/article"
	mysqlrepo "github.com/jayvib/clean-architecture/article/repository/mysql"
	"github.com/jayvib/clean-architecture/config"
	"github.com/jayvib/clean-architecture/pkg/database/mysql"
	apptime "github.com/jayvib/clean-architecture/time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func setupRepo() (article.Repository, error) {
	conf, err := config.New()
	if err != nil {
		return nil, err
	}
	db, err := mysql.New(conf)
	if err != nil {
		return nil, err
	}
	repo := mysqlrepo.New(db)
	return repo, nil
}

func TestMySQL_GetByTitle(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)
	db, err := mysql.New(conf)
	require.NoError(t, err)
	repo := mysqlrepo.New(db)
	res, err := repo.GetByTitle(context.Background(), "Pirate King")
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestMySQL_GetByID(t *testing.T) {
	//t.SkipNow()
	repo, err := setupRepo()
	require.NoError(t, err)
	res, err := repo.GetByID(context.Background(), "uniquefjeid1")
	assert.Nil(t, res)
	require.Error(t, err)
}

func TestEncodeCursor(t *testing.T) {
	t.SkipNow()
	hourago := time.Now().Add((-1) * time.Hour)
	cursor := apptime.EncodeCursor(hourago)
	_ = cursor
}
