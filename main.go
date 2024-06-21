package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type Order struct {
	ID           int    `json:"id"`
	BouquetName  string `json:"bouquetName"`
	CustomerName string `json:"customerName"`
	CustomerPhone string `json:"customerPhone"`
}

func main() {
	db, err := sql.Open("sqlite3", "./products.db")
	if err != nil {
		log.Fatal("Ошибка открытия базы данных:", err)
	}
	defer db.Close()


	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        name TEXT,
                        image TEXT,
                        description TEXT,
                        price TEXT
                    )`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы товаров:", err)
	}


	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS orders (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        bouquet_name TEXT,
                        customer_name TEXT,
                        customer_phone TEXT
                    )`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы заказов:", err)
	}

	initialProducts := []Product{
		{Name: "Цветущая милость", Image: "img/V1.png", Description: "Очень нежный букет из микса розовых и кремовых роз", Price: "2 350₽"},
		{Name: "Нежное серце", Image: "img/V2.png", Description: "Букет из розовых махровых цветов для романтичных прогулок", Price: "1800₽"},
		{Name: "Нежный поцелуй", Image: "img/v3.png", Description: "Трогательное сочетание розовых и белых роз для первого свидания", Price: "3000₽"},
		{Name: "Белые ночи", Image: "img/v4.png", Description: "Потрясающий букет для ночных свиданий", Price: "3 500₽"},
		{Name: "Романтика", Image: "img/v5.png", Description: "Изящный букет, созданный как полотно известных натюрмортов", Price: "3200₽"},
		{Name: "Магическая красота", Image: "img/v6.png", Description: "Чарующий букет для создания нежной атмосферы", Price: "4800₽"},
		{Name: "СОБЕРИ СВОЙ БУКЕТ", Image: "img/v7.png", Description: "Создай свою неповторимую композицию цветов! Подбери букет на любой случай", Price: "от 1500₽"},
	}

	for _, p := range initialProducts {
		_, err := db.Exec(`INSERT INTO products (name, image, description, price) VALUES (?, ?, ?, ?)`,
			p.Name, p.Image, p.Description, p.Price)
		if err != nil {
			log.Fatal("Ошибка добавления товара:", err)
		}
	}


	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		rows, err := db.Query("SELECT id, name, image, description, price FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Image, &p.Description, &p.Price); err != nil {
				log.Println(err)
				continue
			}
			products = append(products, p)
		}

		if err := rows.Err(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(products)
	})


	http.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")

			rows, err := db.Query("SELECT id, bouquet_name, customer_name, customer_phone FROM orders")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var orders []Order
			for rows.Next() {
				var o Order
				if err := rows.Scan(&o.ID, &o.BouquetName, &o.CustomerName, &o.CustomerPhone); err != nil {
					log.Println(err)
					continue
				}
				orders = append(orders, o)
			}

			if err := rows.Err(); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(orders)
			return
		}

		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")

			var order Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				log.Println("Ошибка при декодировании JSON:", err)
				return
			}
			defer r.Body.Close()

			log.Printf("Получен заказ: %+v\n", order)

			_, err := db.Exec(`INSERT INTO orders (bouquet_name, customer_name, customer_phone) VALUES (?, ?, ?)`,
				order.BouquetName, order.CustomerName, order.CustomerPhone)
			if err != nil {
				http.Error(w, "Failed to insert order", http.StatusInternalServerError)
				log.Println("Ошибка при выполнении запроса INSERT:", err)
				return
			}
			log.Println("Заказ успешно добавлен в базу данных")

			w.WriteHeader(http.StatusCreated)
			response := map[string]string{"status": "order created"}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Println("Ошибка при отправке JSON-ответа:", err)
			}
			return
		}

		// Если метод не поддерживается
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Получен недопустимый метод запроса:", r.Method)
	})

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}