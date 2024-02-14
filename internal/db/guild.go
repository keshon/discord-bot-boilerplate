package db

import (
	"gorm.io/gorm"
)

type Guild struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

// CreateGuild creates a new guild.
//
// guild: the guild to be created.
// error: an error if the creation fails.
func CreateGuild(guild Guild) error {
	return DB.Create(&guild).Error
}

// GetGuildByID retrieves a Guild by its ID.
//
// guildID string
// *Guild, error
func GetGuildByID(guildID string) (*Guild, error) {
	var guild Guild
	err := DB.Where("id = ?", guildID).First(&guild).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &guild, err
}

// GetAllGuildIDs retrieves all guild IDs from the database.
//
// No parameters.
// Returns a slice of strings and an error.
func GetAllGuildIDs() ([]string, error) {
	var guilds []Guild
	var guildIDs []string

	if err := DB.Find(&guilds).Error; err != nil {
		return nil, err
	}

	for _, guild := range guilds {
		guildIDs = append(guildIDs, guild.ID)
	}

	return guildIDs, nil
}

// DoesGuildExist checks if the guild with the given ID exists in the database.
//
// guildID string
// bool, error
func DoesGuildExist(guildID string) (bool, error) {
	var count int64
	err := DB.Model(&Guild{}).Where("id = ?", guildID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteGuild deletes a guild by its ID.
//
// Parameter: guildID string
// Return type: error
func DeleteGuild(guildID string) error {
	return DB.Where("id = ?", guildID).Delete(&Guild{}).Error
}
