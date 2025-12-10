package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Handler untuk endpoint GET
func GetProductHandler(w http.ResponseWriter, r *http.Request) error {
	queryIDs := r.URL.Query().Get("ids")
	if queryIDs == "" {
		return fmt.Errorf("ids is required")
	}

	strIDs := strings.Split(queryIDs, ",")
	var ids []int
	for _, s := range strIDs {
		id, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return ValidationError{Msg: fmt.Sprintf("ID harus angka: %s", s)}
		}
		ids = append(ids, id)
	}

	products, err := FetchProductsConcurrently(ids)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    products,
	})
}

func main() {
	http.HandleFunc("/products", WithErrorHandling(GetProductHandler))

	fmt.Println("Server berjalan di http://localhost:8080")
	fmt.Println("Coba endpoint berikut:")
	fmt.Println("1. Sukses:  http://localhost:8080/products?ids=1,2,3,4,5")
	fmt.Println("2. Validasi: http://localhost:8080/products?ids=1,a,3")
	fmt.Println("3. Validasi: http://localhost:8080/products?ids=1,-5")
	fmt.Println("4. NotFound: http://localhost:8080/products?ids=1,999")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
