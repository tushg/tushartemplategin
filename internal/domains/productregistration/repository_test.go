package productregistration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of interfaces.DBInterface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) BeginTx(ctx context.Context, opts interface{}) (interface{}, error) {
	args := m.Called(ctx, opts)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) Stats() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockDB) SetMaxOpenConns(n int) {
	m.Called(n)
}

func (m *MockDB) SetMaxIdleConns(n int) {
	m.Called(n)
}

func (m *MockDB) SetConnMaxLifetime(d time.Duration) {
	m.Called(d)
}

func (m *MockDB) SetConnMaxIdleTime(d time.Duration) {
	m.Called(d)
}

func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) interface{} {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0)
}

func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0), mockArgs.Error(1)
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0), mockArgs.Error(1)
}

// MockLogger is a mock implementation of interfaces.Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(ctx context.Context, msg string, fields interface{}) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Info(ctx context.Context, msg string, fields interface{}) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Warn(ctx context.Context, msg string, fields interface{}) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Error(ctx context.Context, msg string, fields interface{}) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Fatal(ctx context.Context, msg string, err error, fields interface{}) {
	m.Called(ctx, msg, err, fields)
}

func TestNewProductRepository(t *testing.T) {
	mockDB := &MockDB{}
	mockLogger := &MockLogger{}

	repo := NewProductRepository(mockDB, mockLogger)

	assert.NotNil(t, repo)
	assert.IsType(t, &ProductRepository{}, repo)
}

func TestProductRepository_Create(t *testing.T) {
	mockDB := &MockDB{}
	mockLogger := &MockLogger{}
	repo := NewProductRepository(mockDB, mockLogger)

	ctx := context.Background()
	product := &ProductRegistration{
		Name:        "Test Product",
		Description: "Test Description",
		Category:    "Test Category",
		Price:       99.99,
		SKU:         "TEST-001",
		Stock:       100,
		IsActive:    true,
	}

	// Mock the database call
	mockRow := &MockRow{}
	mockRow.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).
		Return(nil).
		Run(func(args mock.Arguments) {
			// Set the returned values
			*args[0].(*int64) = 1
			*args[1].(*time.Time) = time.Now()
			*args[2].(*time.Time) = time.Now()
		})

	mockDB.On("QueryRowContext", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]interface {}")).
		Return(mockRow)

	mockLogger.On("Info", ctx, "Product created successfully", mock.AnythingOfType("interfaces.Fields"))

	result, err := repo.Create(ctx, product)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Test Product", result.Name)
	assert.Equal(t, "TEST-001", result.SKU)

	mockDB.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// MockRow is a mock implementation of sql.Row
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func TestProductService_CreateProduct(t *testing.T) {
	mockRepo := &MockRepository{}
	mockLogger := &MockLogger{}
	service := NewProductService(mockRepo, mockLogger)

	ctx := context.Background()
	req := &CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Category:    "Test Category",
		Price:       99.99,
		SKU:         "TEST-001",
		Stock:       100,
		IsActive:    true,
	}

	// Mock repository calls
	mockRepo.On("SKUExists", ctx, "TEST-001", (*int64)(nil)).Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*productregistration.ProductRegistration")).
		Return(&ProductRegistration{
			ID:          1,
			Name:        "Test Product",
			Description: "Test Description",
			Category:    "Test Category",
			Price:       99.99,
			SKU:         "TEST-001",
			Stock:       100,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)

	mockLogger.On("Info", ctx, "Creating new product", mock.AnythingOfType("interfaces.Fields"))
	mockLogger.On("Info", ctx, "Product created successfully", mock.AnythingOfType("interfaces.Fields"))

	result, err := service.CreateProduct(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Test Product", result.Name)
	assert.Equal(t, "TEST-001", result.SKU)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// MockRepository is a mock implementation of Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(*ProductRegistration), args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*ProductRegistration, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*ProductRegistration), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, id int64, product *ProductRegistration) (*ProductRegistration, error) {
	args := m.Called(ctx, id, product)
	return args.Get(0).(*ProductRegistration), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, req *ProductListRequest) ([]*ProductRegistration, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*ProductRegistration), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetBySKU(ctx context.Context, sku string) (*ProductRegistration, error) {
	args := m.Called(ctx, sku)
	return args.Get(0).(*ProductRegistration), args.Error(1)
}

func (m *MockRepository) UpdateStock(ctx context.Context, id int64, stock int) error {
	args := m.Called(ctx, id, stock)
	return args.Error(0)
}

func (m *MockRepository) Exists(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) SKUExists(ctx context.Context, sku string, excludeID *int64) (bool, error) {
	args := m.Called(ctx, sku, excludeID)
	return args.Bool(0), args.Error(1)
}
