package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Produk represents a product in the cashier system
type Produk struct {
	ID    int    	`json:"id"`
	Nama  string 	`json:"nama"`
	Harga int    	`json:"harga"`
	Stok  int    	`json:"stok"`
	CategoryID *int  `json:"category_id,omitempty"`
	CategoryName *string `json:"category_name,omitempty"`
}

type Category struct {
	ID    		int    `json:"id"`
	Name  		string `json:"nama"`
	Description string `json:"deskripsi"`
	Products   	[]Produk `json:"produk,omitempty"`
}

var categories = []Category{
	{
		ID: 1,
		Name: "Makanan",
		Description: "Kategori untuk semua jenis makanan",
	},
	{
		ID: 2,
		Name: "Minuman",
		Description: "Kategori untuk semua jenis minuman",
	},
	{
		ID: 3,
		Name: "Obat & Vitamin",
		Description: "Kategori untuk semua jenis obat dan vitamin",
	},
	{
		ID: 4,
		Name: "Perawatan Tubuh",
		Description: "Kategori untuk produk perawatan tubuh",
	},
}
// In-memory storage (sementara, nanti ganti database)
var produk = []Produk{
	{
		ID: 1, 
		Nama: "Indomie Godog", 
		Harga: 3500, 
		Stok: 10,
		CategoryID: func(i int) *int { return &i }(1),
	},
	{
		ID: 2, 
		Nama: "Vit 1000ml", 
		Harga: 3000, 
		Stok: 40,
		CategoryID: func(i int) *int { return &i }(3),
	},
	{
		ID: 3, 
		Nama: "kecap", 
		Harga: 12000, 
		Stok: 20,
		CategoryID: func(i int) *int { return &i }(1),
	},
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /api/produk/123 -> ID = 123
	fmt.Println("URL Path:", r.URL.Path)

	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// Cari produk dengan ID tersebut
	for _, p := range produk {
		if p.ID == id {
			if p.CategoryID != nil {
				for _, cat := range categories {
					if cat.ID == *p.CategoryID {
						p.CategoryName = &cat.Name
						break
					}
				}
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

// PUT localhost:8080/api/produk/{id}
func updateProduk(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop produk, cari id, ganti sesuai data dari request
	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}
	
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	
	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}
	
	// loop produk cari ID, dapet index yang mau dihapus
	for i, p := range produk {
		if p.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			produk = append(produk[:i], produk[i+1:]...)
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}
// PUT localhost:8080/api/produk/{id}
func updateProduk(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop produk, cari id, ganti sesuai data dari request
	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}
	
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /api/categories/123 -> ID = 123
	fmt.Println("URL Path:", r.URL.Path)
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	// fmt.Println("Category ID:", err)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Cari kategori dengan ID tersebut
	for _, p := range categories {
		if p.ID == id {
			var productsInCategory []Produk
			for _, prod := range produk {
				if prod.CategoryID != nil && *prod.CategoryID == id {
					productsInCategory = append(productsInCategory, prod)
				}
			}
			p.Products = productsInCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func main() {
	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		}
		// } else if r.Method == "PUT" {
		// 	updateProduk(w, r)
		// } else if r.Method == "DELETE" {
		// 	deleteProduk(w, r)
		// }
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
			
		} else if r.Method == "POST" {
			// baca data dari request
			var categoryBaru Category
			err := json.NewDecoder(r.Body).Decode(&categoryBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// masukkin data ke dalam variable produk
			categoryBaru.ID = len(categories) + 1
			categories = append(categories, categoryBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(categoryBaru)
		}
	})

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			// baca data dari request
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// cek categoryID jika disediakan
			if produkBaru.CategoryID != nil {
				// cari category dengan ID tersebut
				categoryFound := false
				for _, cat := range categories {
					if cat.ID == *produkBaru.CategoryID {
						categoryFound = true
						break
					}
				}
				
				// jika category tidak ditemukan, return error
				if !categoryFound {
					http.Error(w, "Category not found", http.StatusBadRequest)
					return
				}
			}

			// masukkin data ke dalam variable produk
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(produkBaru)
		}
	})

	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}