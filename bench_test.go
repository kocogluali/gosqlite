package sqlite_test

import (
	. "github.com/gwenn/gosqlite"
	"testing"
)

func fill(db *Conn, n int) {
	db.Exec("DROP TABLE IF EXISTS test")
	db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY NOT NULL, float_num REAL, int_num INTEGER, a_string TEXT)")
	s, _ := db.Prepare("INSERT INTO test (float_num, int_num, a_string) VALUES (?, ?, ?)")

	db.Begin()
	for i := 0; i < n; i++ {
		s.Exec(float64(i)*float64(3.14), i, "hello")
	}
	s.Finalize()
	db.Commit()
}

func BenchmarkValuesScan(b *testing.B) {
	b.StopTimer()
	db, _ := Open("")
	defer db.Close()
	fill(db, 1000)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cs, _ := db.Prepare("SELECT float_num, int_num, a_string FROM test")

		values := make([]interface{}, 3)
		for Must(cs.Next()) {
			cs.ScanValues(values)
		}
		cs.Finalize()
	}
}

func BenchmarkScan(b *testing.B) {
	b.StopTimer()
	db, _ := Open("")
	defer db.Close()
	fill(db, 1000)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cs, _ := db.Prepare("SELECT float_num, int_num, a_string FROM test")

		var fnum float64
		var inum int64
		var sstr string

		for Must(cs.Next()) {
			cs.Scan(&fnum, &inum, &sstr)
		}
		cs.Finalize()
	}
}

func BenchmarkNamedScan(b *testing.B) {
	b.StopTimer()
	db, _ := Open("")
	defer db.Close()
	fill(db, 1000)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cs, _ := db.Prepare("SELECT float_num, int_num, a_string FROM test")

		var fnum float64
		var inum int64
		var sstr string

		for Must(cs.Next()) {
			cs.NamedScan("float_num", &fnum, "int_num", &inum, "a_string", &sstr)
		}
		cs.Finalize()
	}
}

func BenchmarkInsert(b *testing.B) {
	db, _ := Open("")
	defer db.Close()
	fill(db, b.N)
}

func BenchmarkNamedInsert(b *testing.B) {
	db, _ := Open("")
	defer db.Close()
	db.Exec("DROP TABLE IF EXISTS test")
	db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY NOT NULL," +
		" float_num REAL, int_num INTEGER, a_string TEXT)")
	s, _ := db.Prepare("INSERT INTO test (float_num, int_num, a_string)" +
		" VALUES (:f, :i, :s)")
	defer s.Finalize()

	db.Begin()
	for i := 0; i < b.N; i++ {
		s.NamedBind("f", float64(i)*float64(3.14), "i", i, "s", "hello")
		s.Next()
	}
	db.Commit()
}
