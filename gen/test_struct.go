//go:generate go run github.com/cscoding21/csgen/gen/run.go

package gen

type TestStruct struct {
	Name     string        `csval:"req"`
	Email    string        `csval:"req,email"`
	Password string        `csval:"min(3),max(11)"`
	Age      int           `csval:"min(18),max(65)"`
	Sub      TestSubStruct `csval:"obj"`
}

type TestSubStruct struct {
	IP   string `csval:"req,ip"`
	Port int
}
