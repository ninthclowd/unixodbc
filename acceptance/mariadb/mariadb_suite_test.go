package mariadb_test

import (
	"context"
	"database/sql"
	_ "github.com/ninthclowd/unixodbc"
	"github.com/onsi/ginkgo/v2/dsl/core"
	"github.com/tommy351/goldga"
	"runtime/trace"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	db   *sql.DB
	conn *sql.Conn
)

func TestMariadb(t *testing.T) {
	goldga.DefaultSerializer = &goldga.JSONSerializer{
		EscapeHTML:   false,
		IndentPrefix: "",
		Indent:       "	",
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mariadb Suite")
}

var _ = BeforeSuite(func() {
	trace.Start(core.GinkgoWriter)
	var err error
	db, err = sql.Open("unixodbc", "DSN=MariaDB")
	Expect(err).To(BeNil())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err = db.Conn(ctx)
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	//conn.Close() //TODO freezes
	//db.Close()
})
