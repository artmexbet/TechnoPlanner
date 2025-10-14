package configer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	apisvc "technoBro/internal/app/configer/svc/api"
)

type Configer struct {
	// Путь к каталогу с описаниями сервисов (internal/app/configer/svc по умолчанию)
	pathToConfigDir string
	outPath         string
}

func New(in, out string) *Configer {
	var _pathToConfigDir = in
	if in == "" {
		_pathToConfigDir = filepath.ToSlash(filepath.Join("internal", "app", "configer", "svc"))
	}
	return &Configer{pathToConfigDir: _pathToConfigDir, outPath: out}
}

// Gen обходит каталоги сервисов и для каждого доступного сервиса создаёт YAML-конфиг
// в cmd/<mapped_service>/config/cfg.yaml
func (c *Configer) Gen() error {
	if c.pathToConfigDir == "" {
		return errors.New("pathToConfigDir is empty")
	}

	entries, err := os.ReadDir(c.pathToConfigDir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", c.pathToConfigDir, err)
	}

	// регистрация известных сервисов и их сборщиков конфигов
	// key: имя каталога сервиса под c.pathToConfigDir
	// value: (конфиг, имя каталога под cmd для записи)
	builders := map[string]func() (any, string){
		"api": func() (any, string) {
			// Значения по умолчанию — из фабрики сервиса
			cfg := apisvc.Default()
			// Мэппим сервис api на бинарь server (cmd/server)
			return cfg, "server"
		},
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		build, ok := builders[name]
		if !ok {
			// неизвестный сервис — пропускаем без ошибки
			continue
		}
		cfg, outName := build()

		// сериализация в YAML
		outBytes, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("yaml marshal for service %s: %w", name, err)
		}

		// путь назначения: cmd/<outName>/config/cfg.yaml
		outDir := filepath.ToSlash(filepath.Join(c.outPath, outName, "config"))
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", outDir, err)
		}
		outPath := filepath.ToSlash(filepath.Join(outDir, "cfg.yaml"))
		if err := os.WriteFile(outPath, outBytes, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", outPath, err)
		}
	}

	_ = time.Second // удерживаем импорт time, если появятся дополнительные билдеры
	return nil
}
