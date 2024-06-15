package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// menyimpan data rsvp
type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

// menginisialisasi `slice` dari pointer ke struct `Rsvp`
var responses = make([]*Rsvp, 0, 10)

// sebuah map kosong untuk menyimpan HTML atau template
var templates = make(map[string]*template.Template, 3)

// memuat beberapa template HTML dari file-file terpisah
// dan menyimpannya ke dalam map `templates`
func loadTemplates() {
	// 5 string array dengan nama file HTML yang akan diload
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	// mengambil nilai `index` dan `name` dari `templateNames`
	// `index` adalah indeks numerik dari elemen saat ini dalam array, dimulai dari 0.
	// `name` adalah nilai string dari elemen saat ini dalam array,
	// yang merupakan nama template.
	for index, name := range templateNames {
		/*
			"layout.html": File HTML yang berisi layout atau struktur dasar yang akan digunakan oleh semua template.
			"name+".html": File HTML yang berisi konten spesifik untuk template dengan nama tertentu.
			  Nama file ini dibentuk dengan menggabungkan name (nama template) dengan ekstensi .html.
		*/
		// `template.ParseFiles` akan membaca dan compile HTML files menjadi objek
		t, err := template.ParseFiles("layout.html", name+".html")

		// jika error saat compile template, fungsi `panic` akan menghentikan program
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

// merender template `welcome` & mengirimkan hasilnya sebagai response HTTP ke client
func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}

// merender template `list` dgn data `responses`
func listHandler(writer http.ResponseWriter, request *http.Request) {
	// `Execute` untuk merender template ke `writer` dan `responses`
	templates["list"].Execute(writer, responses)
}

// mengelola data `form` & `errors` dalam satu struct yg terorganisir
type formData struct {
	*Rsvp
	Errors []string
}

// merender template `form` & mengirimkan hasilnya sebagai response HTTP ke client
func formHandler(writer http.ResponseWriter, request *http.Request) {
	// jika metode HTTP adalah `GET`, render template `form`,
	if request.Method == http.MethodGet {
		// dengan menggunakan objek `formData`.
		templates["form"].Execute(writer, formData{
			// pointer ke struct `Rsvp` yang nil.
			Rsvp: &Rsvp{}, Errors: []string{},
		})
		// jika metode HTTP adalah `POST`,
	} else if request.Method == http.MethodPost {
		// parse form data dari request
		request.ParseForm()

		// variable `responseData` dari type `Rsvp`,
		// diinisialisasi dengan data yg diterima dari form.
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}

		// slice string kosong `errors` untuk menyimpan semua error
		errors := []string{}
		// jika nama, email, phone, dan willattend kosong,
		if responseData.Name == "" {
			errors = append(errors, "Please enter your name")
		}
		if responseData.Email == "" {
			errors = append(errors, "Please enter your email address")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Please enter your phone number")
		}
		// jika panjang slice `errors` lebih dari 0,
		if len(errors) > 0 {
			// `form` akan di render lagi dengan data `responseData` dan `errors`
			templates["form"].Execute(writer, formData{
				Rsvp: &responseData, Errors: errors,
			})
			// jika tidak ada error,
		} else {
			// `responseData` akan ditambahkan ke `responses`
			responses = append(responses, &responseData)
			// dan `thanks` atau `sorry` akan di render,
			// sesuai dengan `willattend`.
			if responseData.WillAttend {
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				templates["sorry"].Execute(writer, responseData.Name)
			}
		}
	}
}

func main() {
	loadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	// menjalankan server HTTP di port 5000 & mendengarkan request yg masuk
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
