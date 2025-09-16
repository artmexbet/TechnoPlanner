# Techno Planner
Приложение для планирования и распределения оборудования между местами и ответственными

## Линтинг
- Конфигурация: .golangci.yml
- Make-таргеты:
  - lint — запуск через Docker (при наличии)
  - lint-fix — автопочинка через Docker
  - lint-local — локальный запуск (нужен установленный golangci-lint)
  - lint-local-fix — локальная автопочинка

Установка локально (Windows/PowerShell):
```powershell
choco install golangci-lint -y
# либо через установщик от проекта
# iwr -useb https://raw.githubusercontent.com/golangci/golangci-lint/master/install.ps1 | iex
```

Запуск:
```powershell
make lint-local
# или через Docker
make lint
```

Примечание по версиям:
- Если видите ошибку вида: "the Go language version (go1.23) used to build golangci-lint is lower than the targeted Go version (1.24)",
  она решена настройкой run.go=1.23 в .golangci.yml. Можно также обновить образ/версию golangci-lint, собранный с более новой Go.
