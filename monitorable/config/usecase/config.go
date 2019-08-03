package usecase

import (
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/config/models"
	. "github.com/monitoror/monitoror/pkg/monitoror/validator"
)

// Versions
const (
	CurrentVersion = Version1

	Version1 = 2
)

var SupportedVersions = map[int]bool{
	Version1: true,
}

// Tile keys
const (
	TypeKey       = "type"
	LabelKey      = "label"
	ColumnSpanKey = "columnSpan"
	RowSpanKey    = "rowSpan"
	ParamsKey     = "params"
	TilesKey      = "tiles"
	UrlKey        = "url" // Injected by hydrate function
)

var AuthorizedTileKey = map[string]bool{
	TypeKey:       true,
	LabelKey:      true,
	ColumnSpanKey: true,
	RowSpanKey:    true,
	ParamsKey:     true,
	TilesKey:      true,
}

const (
	// Custom Tile Type
	EmptyTileType tiles.TileType = "EMPTY"
	GroupTileType tiles.TileType = "GROUP"
)

type (
	configUsecase struct {
		repository  config.Repository
		tileConfigs map[tiles.TileType]*TileConfig
	}

	// TileConfig struct is used by GetConfig endpoint to check / hydrate config
	TileConfig struct {
		Path      string
		Validator Validator
	}
)

func NewConfigUsecase(repository config.Repository) config.Usecase {
	return &configUsecase{
		repository:  repository,
		tileConfigs: make(map[tiles.TileType]*TileConfig),
	}
}

func (cu *configUsecase) RegisterTile(tileType tiles.TileType, path string, validator Validator) {
	cu.tileConfigs[tileType] = &TileConfig{
		Path:      path,
		Validator: validator,
	}
}

// GetConfig load and parse GetConfig
func (cu *configUsecase) GetConfig(params *models.ConfigParams) (config *models.Config, err error) {
	if params.Url != "" {
		config, err = cu.repository.GetConfigFromUrl(params.Url)
	} else if params.Path != "" {
		config, err = cu.repository.GetConfigFromPath(params.Path)
	}

	if err != nil {
		return
	}

	// Set config to CurrentVersion if config isn't set
	if config.Version == 0 {
		config.Version = CurrentVersion
	}

	return
}