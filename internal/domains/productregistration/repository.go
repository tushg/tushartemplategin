package productregistration

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"tushartemplategin/pkg/interfaces"
)

// ProductRepository implements the Repository interface for product data access
type ProductRepository struct {
	db     interfaces.Database
	logger interfaces.Logger
}

// NewProductRepository creates a new product repository
func NewProductRepository(db interfaces.Database, log interfaces.Logger) Repository {
	return &ProductRepository{
		db:     db,
		logger: log,
	}
}

// Create creates a new product in the database
func (r *ProductRepository) Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, error) {
	query := `
		INSERT INTO products (name, description, category, price, sku, stock, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, query,
			product.Name, product.Description, product.Category, product.Price,
			product.SKU, product.Stock, product.IsActive, product.CreatedAt, product.UpdatedAt,
		).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	}); err != nil {
		r.logger.Error(ctx, "Failed to create product", interfaces.Fields{
			"error": err.Error(),
			"sku":   product.SKU,
		})
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	r.logger.Info(ctx, "Product created successfully", interfaces.Fields{
		"id":  product.ID,
		"sku": product.SKU,
	})

	return product, nil
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*ProductRegistration, error) {
	query := `
		SELECT id, name, description, category, price, sku, stock, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &ProductRegistration{}
	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, query, id).Scan(
			&product.ID, &product.Name, &product.Description, &product.Category,
			&product.Price, &product.SKU, &product.Stock, &product.IsActive,
			&product.CreatedAt, &product.UpdatedAt,
		)
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		r.logger.Error(ctx, "Failed to get product by ID", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, id int64, product *ProductRegistration) (*ProductRegistration, error) {
	query := `
		UPDATE products 
		SET name = $1, description = $2, category = $3, price = $4, sku = $5, 
		    stock = $6, is_active = $7, updated_at = $8
		WHERE id = $9
		RETURNING created_at, updated_at
	`

	product.UpdatedAt = time.Now()

	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, query,
			product.Name, product.Description, product.Category, product.Price,
			product.SKU, product.Stock, product.IsActive, product.UpdatedAt, id,
		).Scan(&product.CreatedAt, &product.UpdatedAt)
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		r.logger.Error(ctx, "Failed to update product", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	product.ID = id
	r.logger.Info(ctx, "Product updated successfully", interfaces.Fields{
		"id":  product.ID,
		"sku": product.SKU,
	})

	return product, nil
}

// Delete removes a product from the database
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	var rowsAffected int64
	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		result, execErr := tx.ExecContext(ctx, query, id)
		if execErr != nil {
			return execErr
		}
		var raErr error
		rowsAffected, raErr = result.RowsAffected()
		return raErr
	}); err != nil {
		r.logger.Error(ctx, "Failed to delete product", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	r.logger.Info(ctx, "Product deleted successfully", interfaces.Fields{
		"id": id,
	})

	return nil
}

// List retrieves a list of products with pagination and filtering
func (r *ProductRepository) List(ctx context.Context, req *ProductListRequest) ([]*ProductRegistration, int64, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	// Build WHERE clause
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Category != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, req.Category)
		argIndex++
	}

	if req.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d OR sku ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total records and fetch rows within a single transaction
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	query := fmt.Sprintf(`
		SELECT id, name, description, category, price, sku, stock, is_active, created_at, updated_at
		FROM products
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, req.Limit, offset)

	var (
		total    int64
		products []*ProductRegistration
	)

	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil { // exclude limit/offset for count
			return err
		}
		rows, qerr := tx.QueryContext(ctx, query, args...)
		if qerr != nil {
			return qerr
		}
		defer rows.Close()

		for rows.Next() {
			product := &ProductRegistration{}
			if err := rows.Scan(
				&product.ID, &product.Name, &product.Description, &product.Category,
				&product.Price, &product.SKU, &product.Stock, &product.IsActive,
				&product.CreatedAt, &product.UpdatedAt,
			); err != nil {
				return err
			}
			products = append(products, product)
		}
		return rows.Err()
	}); err != nil {
		r.logger.Error(ctx, "Failed to list products", interfaces.Fields{
			"error": err.Error(),
		})
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	return products, total, nil
}

// GetBySKU retrieves a product by its SKU
func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*ProductRegistration, error) {
	query := `
		SELECT id, name, description, category, price, sku, stock, is_active, created_at, updated_at
		FROM products
		WHERE sku = $1
	`

	product := &ProductRegistration{}
	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, query, sku).Scan(
			&product.ID, &product.Name, &product.Description, &product.Category,
			&product.Price, &product.SKU, &product.Stock, &product.IsActive,
			&product.CreatedAt, &product.UpdatedAt,
		)
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with sku %s not found", sku)
		}
		r.logger.Error(ctx, "Failed to get product by SKU", interfaces.Fields{
			"error": err.Error(),
			"sku":   sku,
		})
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return product, nil
}

// UpdateStock updates the stock quantity of a product
func (r *ProductRepository) UpdateStock(ctx context.Context, id int64, stock int) error {
	query := `
		UPDATE products 
		SET stock = $1, updated_at = $2
		WHERE id = $3
	`

	var rowsAffected int64
	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		result, execErr := tx.ExecContext(ctx, query, stock, time.Now(), id)
		if execErr != nil {
			return execErr
		}
		var raErr error
		rowsAffected, raErr = result.RowsAffected()
		return raErr
	}); err != nil {
		r.logger.Error(ctx, "Failed to update product stock", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
			"stock": stock,
		})
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	r.logger.Info(ctx, "Product stock updated successfully", interfaces.Fields{
		"id":    id,
		"stock": stock,
	})

	return nil
}

// Exists checks if a product exists by ID
func (r *ProductRepository) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)`

	var exists bool
	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, query, id).Scan(&exists)
	}); err != nil {
		r.logger.Error(ctx, "Failed to check product existence", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}

	return exists, nil
}

// SKUExists checks if a SKU already exists (optionally excluding a specific product ID)
func (r *ProductRepository) SKUExists(ctx context.Context, sku string, excludeID *int64) (bool, error) {
	var query string
	var args []interface{}

	if excludeID != nil {
		query = `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 AND id != $2)`
		args = []interface{}{sku, *excludeID}
	} else {
		query = `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1)`
		args = []interface{}{sku}
	}

	var exists bool
	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, query, args...).Scan(&exists)
	}); err != nil {
		r.logger.Error(ctx, "Failed to check SKU existence", interfaces.Fields{
			"error": err.Error(),
			"sku":   sku,
		})
		return false, fmt.Errorf("failed to check SKU existence: %w", err)
	}

	return exists, nil
}
