package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/library/conf/paladin.v2"
	"go-common/library/testing/dockertest"
)

var d *dao
var ctx = context.Background()

func TestMain(m *testing.M) {
	flag.Set("conf", "../../test")
	flag.Parse()
	dockertest.Run("../../test/dockertest.yml")
	var err error
	if err = paladin.Init(); err != nil {
		panic(err)
	}
	var cf func()
	if d, cf, err = newTestDao(); err != nil {
		panic(err)
	}
	ret := m.Run()
	cf()
	os.Exit(ret)
}
