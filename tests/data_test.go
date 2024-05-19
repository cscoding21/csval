//go:generate csval

package tests

type FooStruct struct {
	Name        string       `csval:"req"`
	Email       string       `csval:"req,email"`
	Password    string       `csval:"min(3),max(11)"`
	ConfirmPass string       `csval:"equals(Password)"`
	Age         int          `csval:"min(18),max(65)"`
	Sub         BarSubStruct `csval:"obj"`
}

type BarSubStruct struct {
	IP   string `csval:"req,ip"`
	Port int
}
