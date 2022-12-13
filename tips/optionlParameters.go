package main

import "fmt"

type Foo struct {
	Name    string
	Age     int
	Address string
}
type Option func(f *Foo)

func WithAddress(s string) Option {
	return func(f *Foo) {
		f.Address = s
	}
}

func WithAge(s int) Option {
	return func(f *Foo) {
		f.Age = s
	}
}

func New(name string, opts ...Option) *Foo {
	f := &Foo{
		Name:    name,
		Age:     10, // default value
		Address: "Egypt",
	}
	for _, applyOpt := range opts {
		applyOpt(f)
	}
	return f
}

func main() {

	s := New("islam", WithAddress("Ismailia"))
	fmt.Println(s)

	s2 := New("ahmed", WithAge(20), WithAddress("Cairo"))
	fmt.Println(s2)

	s3 := New("ahmed")
	fmt.Println(s3)
	// router := httprouter.New()
	// router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("not found"))
	// })

	// router.ServeFiles("/v1/movies/swagger/*filepath", http.Dir("swagger"))

	// router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("hello there"))
	// })

	// http.ListenAndServe(":4000", router)
}
