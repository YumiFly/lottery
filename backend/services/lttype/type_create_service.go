package lttype

import (
	"backend/models"
	"backend/utils"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// LotteryService encapsulates lottery-related business logic
type TypeCreateService struct {
	db *gorm.DB
}

// NewLotteryService creates a new LotteryService instance
func NewTypeCreateService(db *gorm.DB) *TypeCreateService {
	return &TypeCreateService{db: db}
}

// CreateLotteryTypeParams defines the parameters for creating a lottery type
type CreateLotteryTypeParams struct {
	TypeName    string
	Description string
}

// validateCreateLotteryTypeParams validates the parameters for creating a lottery type
func (s *TypeCreateService) validateCreateLotteryTypeParams(ctx context.Context, params CreateLotteryTypeParams) error {
	// Validate name
	if len(params.TypeName) == 0 || len(params.TypeName) > 100 {
		return utils.NewBadRequestError("Name must be between 1 and 100 characters", nil)
	}

	// Check for duplicate name
	var existingType models.LotteryType
	if err := s.db.WithContext(ctx).
		Where("type_name = ?", params.TypeName).
		First(&existingType).Error; err == nil {
		utils.Logger.Warn("Lottery type name already exists", "name", params.TypeName)
		return utils.NewBadRequestError("Lottery type name already exists", nil)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.NewInternalError("Failed to check lottery type name uniqueness", errors.Wrap(err, "database error"))
	}

	return nil
}

// CreateLotteryType creates a new lottery type
//
// Parameters:
//   - ctx: Request context
//   - params: Creation parameters, including name
//
// Returns:
//   - *LotteryType: The created lottery type
//   - error: Creation error or invalid parameters
func (s *TypeCreateService) CreateLotteryType(ctx context.Context, params CreateLotteryTypeParams) (*models.LotteryType, error) {
	// Validate parameters
	if err := s.validateCreateLotteryTypeParams(ctx, params); err != nil {
		return nil, err
	}

	// Construct lottery type
	lotteryType := models.LotteryType{
		TypeID:      uuid.NewString(),
		TypeName:    params.TypeName,
		Description: params.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Log creation attempt
	utils.Logger.Info("Creating lottery type",
		"type_id", lotteryType.TypeID,
		"name", lotteryType.TypeName)

	// Save to database
	if err := s.db.WithContext(ctx).Create(&lotteryType).Error; err != nil {
		utils.Logger.Error("Failed to create lottery type", "error", err)
		return nil, utils.NewInternalError("Failed to create lottery type", err)
	}

	// Log success
	utils.Logger.Info("Lottery type created successfully",
		"type_id", lotteryType.TypeID)

	return &lotteryType, nil
}
