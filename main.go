package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/models/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// ubah Config
type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// Produk represents a product in the cashier system
type Produk struct {
	ID           int     `json:"id"`
	Nama         string  `json:"nama"`
	Harga        int     `json:"harga"`
	Stok         int     `json:"stok"`
	CategoryID   *int    `json:"category_id,omitempty"`
	CategoryName *string `json:"category_name,omitempty"`
}

type Category struct {
	ID          int      `json:"id"`
	Name        string   `json:"nama"`
	Description string   `json:"deskripsi"`
	Products    []Produk `json:"produk,omitempty"`
}

var categories = []Category{
	{
		ID:          1,
		Name:        "Makanan",
		Description: "Kategori untuk semua jenis makanan",
	},
	{
		ID:          2,
		Name:        "Minuman",
		Description: "Kategori untuk semua jenis minuman",
	},
	{
		ID:          3,
		Name:        "Obat & Vitamin",
		Description: "Kategori untuk semua jenis obat dan vitamin",
	},
	{
		ID:          4,
		Name:        "Perawatan Tubuh",
		Description: "Kategori untuk produk perawatan tubuh",
	},
}

// In-memory storage (sementara, nanti ganti database)
var produk = []Produk{
	{
		ID:         1,
		Nama:       "Indomie Godog",
		Harga:      3500,
		Stok:       10,
		CategoryID: func(i int) *int { return &i }(1),
	},
	{
		ID:         2,
		Nama:       "Vit 1000ml",
		Harga:      3000,
		Stok:       40,
		CategoryID: func(i int) *int { return &i }(3),
	},
	{
		ID:         3,
		Nama:       "kecap",
		Harga:      12000,
		Stok:       20,
		CategoryID: func(i int) *int { return &i }(1),
	},
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

// PUT localhost:8080/api/categories/{id}
func updateCategoryByID(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop categories, cari id, ganti sesuai data dari request
	for i := range categories {
		if categories[i].ID == id {
			updateCategory.ID = id
			categories[i] = updateCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}

	http.Error(w, "Category belum ada", http.StatusNotFound)
}

func deleteCategoryByID(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// loop categories cari ID, dapet index yang mau dihapus
	for i, c := range categories {
		if c.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			categories = append(categories[:i], categories[i+1:]...)
			for j, p := range produk {
				if p.CategoryID != nil && *p.CategoryID == id {
					p.CategoryID = nil
					produk[j] = p
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func main() {

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategoryByID(w, r)
		} else if r.Method == "DELETE" {
			deleteCategoryByID(w, r)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

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

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	addr := "0.0.0.0:" + config.Port

	fmt.Printf("Server running di %s\n", addr)

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("gagal running server: %v", err)
	}

}
