package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// CurrentVersion — простая версия схемы файла настроек.
const CurrentVersion = 1

// ConnectionMode описывает режим установления VPN‑подключения.
type ConnectionMode string

const (
	ConnectionModeAuto             ConnectionMode = "auto"
	ConnectionModeVLESSRealityOnly ConnectionMode = "vless_reality_only"
)

// Language — код языка интерфейса.
type Language string

const (
	LanguageRU Language = "ru"
	LanguageEN Language = "en"
)

type Settings struct {
	Version int `json:"version"`

	Connection ConnectionSettings `json:"connection"`
	Privacy    PrivacySettings    `json:"privacy"`
	App        AppSettings        `json:"app"`
}

type ConnectionSettings struct {
	AutoConnectOnLaunch   bool           `json:"auto_connect_on_launch"`
	AutoReconnect         bool           `json:"auto_reconnect"`
	ReconnectIntervalSecs int            `json:"reconnect_interval_secs"`
	Mode                  ConnectionMode `json:"mode"`
}

type PrivacySettings struct {
	// RememberDevice — только флаг для UI, никакие секреты пока не хранятся.
	RememberDevice bool `json:"remember_device"`
}

type AppSettings struct {
	StartWithWindows bool     `json:"start_with_windows"`
	Language         Language `json:"language"`
}

// Default возвращает настройки по умолчанию.
func Default() Settings {
	return Settings{
		Version: CurrentVersion,
		Connection: ConnectionSettings{
			AutoConnectOnLaunch:   false,
			AutoReconnect:         true,
			ReconnectIntervalSecs: 10,
			Mode:                  ConnectionModeAuto,
		},
		Privacy: PrivacySettings{
			RememberDevice: true,
		},
		App: AppSettings{
			StartWithWindows: false,
			Language:         LanguageRU,
		},
	}
}

// ConfigFilePath возвращает путь к JSON‑файлу настроек в каталоге конфигурации пользователя.
func ConfigFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// Use vendor-like folder, avoid secrets, only simple JSON config.
	return filepath.Join(dir, "VoltaVPN", "settings.json"), nil
}

// Load читает настройки с диска.
// Если файла нет или он некорректен, возвращаются Default и ошибка.
func Load() (Settings, error) {
	path, err := ConfigFilePath()
	if err != nil {
		return Default(), err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Default(), err
	}
	if err != nil {
		return Default(), err
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return Default(), err
	}

	// Простая проверка версии: если версия не совпадает — откатываемся к значениям по умолчанию.
	if s.Version != CurrentVersion {
		return Default(), errors.New("settings version mismatch")
	}
	if err := validate(s); err != nil {
		return Default(), err
	}

	return s, nil
}

// LoadOrDefault возвращает настройки, никогда не падая: при ошибке Load вернёт Default.
func LoadOrDefault() Settings {
	s, err := Load()
	if err != nil {
		return Default()
	}
	return s
}

// Save сохраняет настройки на диск, создавая каталог при необходимости.
func Save(s Settings) error {
	s.Version = CurrentVersion
	if err := validate(s); err != nil {
		return err
	}

	path, err := ConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	// Перезаписываем атомарно: пишем во временный файл и переименовываем.
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}

// Clear удаляет файл настроек и возвращает значения по умолчанию.
func Clear() (Settings, error) {
	path, err := ConfigFilePath()
	if err != nil {
		return Default(), err
	}

	// Remove file if exists. Ignore "not exists" errors.
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Default(), err
	}
	if err := os.Remove(path + ".tmp"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Default(), err
	}

	return Default(), nil
}

func validate(s Settings) error {
	if !isValidReconnectInterval(s.Connection.ReconnectIntervalSecs) {
		return errors.New("invalid reconnect interval")
	}
	if !isValidConnectionMode(s.Connection.Mode) {
		return errors.New("invalid connection mode")
	}
	if !isValidLanguage(s.App.Language) {
		return errors.New("invalid app language")
	}
	return nil
}

func isValidReconnectInterval(v int) bool {
	switch v {
	case 5, 10, 30:
		return true
	default:
		return false
	}
}

func isValidConnectionMode(v ConnectionMode) bool {
	switch v {
	case ConnectionModeAuto, ConnectionModeVLESSRealityOnly:
		return true
	default:
		return false
	}
}

func isValidLanguage(v Language) bool {
	switch v {
	case LanguageRU, LanguageEN:
		return true
	default:
		return false
	}
}
