package dao_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/shopspring/decimal"

	"brightloop/internal/dao"
	"brightloop/internal/model/do"
	"brightloop/internal/model/entity"
	"brightloop/internal/testdb"
)

func TestMain(m *testing.M) {
	if err := testdb.Setup("bl_test_dao"); err != nil {
		fmt.Println("testdb:", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// gf gen dao maps DECIMAL to float64 unless hack/config.yaml's typeMapping says otherwise. This proves money
// survives a write and read through the generated DAO exactly, including values float64 cannot hold.
func TestDecimalRoundTrip(t *testing.T) {
	ctx := context.Background()
	for _, s := range []string{"1180.00", "0.10", "9999999999.99", "1234567890.07"} {
		price := decimal.RequireFromString(s)
		id, err := dao.Products.Ctx(ctx).Data(do.Products{Name: "P " + s, ProductType: "software", SalesPrice: price, CostPrice: price}).InsertAndGetId()
		if err != nil {
			t.Fatal(err)
		}
		var p entity.Products
		if err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, id).Scan(&p); err != nil {
			t.Fatal(err)
		}
		if got := p.SalesPrice.StringFixed(2); got != s {
			t.Errorf("wrote %s, read back %s", s, got)
		}
	}
}
