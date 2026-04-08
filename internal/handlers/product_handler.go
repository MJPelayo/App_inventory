package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "inventory-api/internal/db"
    "inventory-api/internal/models"
)

// GetProducts returns all products from the catalog (GET /api/products)
func GetProducts(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT id, name, description, sku, brand, price, cost, 
               category_id, supplier_id, image_url, is_active, created_at, updated_at
        FROM products WHERE is_active = true ORDER BY id
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var p models.Product
        var imageURL sql.NullString
        err := rows.Scan(
            &p.ID, &p.Name, &p.Description, &p.SKU, &p.Brand,
            &p.Price, &p.Cost, &p.CategoryID, &p.SupplierID,
            &imageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
        )
        if imageURL.Valid {
            p.ImageURL = &imageURL.String
        }
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        products = append(products, p)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

// GetProduct returns a single product by ID with its inventory levels (GET /api/products/{id})
func GetProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    // Get product details
    var product models.Product
    var imageURL sql.NullString
    err = db.DB.QueryRow(`
        SELECT id, name, description, sku, brand, price, cost,
               category_id, supplier_id, image_url, is_active, created_at, updated_at
        FROM products WHERE id = $1
    `, id).Scan(
        &product.ID, &product.Name, &product.Description, &product.SKU, &product.Brand,
        &product.Price, &product.Cost, &product.CategoryID, &product.SupplierID,
        &imageURL, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
    )
    if imageURL.Valid {
        product.ImageURL = &imageURL.String
    }
    
    if err != nil {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    // Get inventory levels for this product
    invRows, err := db.DB.Query(`
        SELECT i.id, i.product_id, i.warehouse_id, i.quantity, i.reorder_point, i.max_stock,
               w.name as warehouse_name
        FROM inventory i
        JOIN warehouses w ON i.warehouse_id = w.id
        WHERE i.product_id = $1
    `, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer invRows.Close()

    type InventoryInfo struct {
        ID            int    `json:"id"`
        WarehouseID   int    `json:"warehouse_id"`
        WarehouseName string `json:"warehouse_name"`
        Quantity      int    `json:"quantity"`
        ReorderPoint  int    `json:"reorder_point"`
        MaxStock      int    `json:"max_stock"`
    }

    var inventory []InventoryInfo
    for invRows.Next() {
        var inv InventoryInfo
        err := invRows.Scan(&inv.ID, &inv.WarehouseID, &inv.WarehouseName, &inv.Quantity, &inv.ReorderPoint, &inv.MaxStock)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        inventory = append(inventory, inv)
    }

    result := struct {
        Product   models.Product   `json:"product"`
        Inventory []InventoryInfo  `json:"inventory"`
    }{
        Product:   product,
        Inventory: inventory,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

// CreateProduct adds a new product to the catalog (POST /api/products)
func CreateProduct(w http.ResponseWriter, r *http.Request) {
    var req models.CreateProductRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    var productID int
    err := db.DB.QueryRow(`
        INSERT INTO products (name, description, sku, brand, price, cost, category_id, supplier_id, image_url)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `, req.Name, req.Description, req.SKU, req.Brand, req.Price, req.Cost, req.CategoryID, req.SupplierID, req.ImageURL).Scan(&productID)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return the created product
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id": productID,
        "message": "Product created successfully",
    })
}