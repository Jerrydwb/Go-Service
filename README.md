# Kardex PDF Service v2.0

Servicio en Go para generación optimizada de PDFs de Kardex Valorizado con soporte para datasets grandes.

## Estructura del Proyecto

```
go-service/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
│
├── internal/                       # Código compartido (no exportable)
│   ├── config/
│   │   └── config.go               # Config global + env vars
│   ├── handler/
│   │   ├── health.go               # GET /health
│   │   └── config.go               # PUT /api/config
│   ├── server/
│   │   ├── server.go               # Setup Gin, middlewares, routes
│   │   └── middleware/
│   │       └── cors.go             # CORS middleware
│   ├── shared/
│   │   ├── types.go                # FlexibleString, tipos genéricos
│   │   └── utils.go                # FormatNumber, GC, ZIP, copy
│   └── tracker/
│       └── tracker.go              # Process tracker + auto-cleanup
│
├── services/                       # Servicios de negocio (escalable)
│   └── kardex/
│       ├── handler.go              # Handlers HTTP (generate, progress)
│       ├── types.go                # Domain types + DTOs
│       ├── generator.go            # Estrategias de generación
│       └── tables.go               # Render de tablas PDF
│
├── scripts/                        # Utilidades CLI
│   └── delete-folder/
│       └── main.go
│
├── go.mod
├── go.sum
└── README.md
```

### Agregar un nuevo servicio

1. Crear carpeta en `services/nombre/` con `handler.go`, `types.go`, `generator.go`
2. Registrar rutas en `internal/server/server.go`:
   ```go
   nombreGroup := api.Group("/nombre")
   nombre.RegisterRoutes(nombreGroup, cfg, trk)
   ```

## Uso

### Compilar

```bash
go build -o kardex-pdf-service ./cmd/server/
```

### Ejecutar

```bash
./kardex-pdf-service
# O con puerto personalizado
PORT=3000 ./kardex-pdf-service
```

### Probar

```bash
./test-go-service.sh
```

## API Endpoints

| Método | Ruta                     | Descripción                 |
| ------ | ------------------------ | --------------------------- |
| POST   | `/api/pdf/generate`      | Genera PDF de Kardex        |
| GET    | `/api/pdf/progress?key=` | Consulta progreso           |
| GET    | `/health`                | Health check con métricas   |
| PUT    | `/api/config`            | Actualiza config en runtime |

### POST `/api/pdf/generate`

**Request:**

```json
{
  "jsonFilePath": "/path/to/kardex-data.json",
  "outputDirectory": "/path/to/output",
  "outputFilename": "kardex-test.pdf"
}
```

**Response:**

```json
{
  "success": true,
  "filename": "kardex-test.pdf",
  "totalPages": 5,
  "totalInsumos": 3,
  "totalMovements": 7,
  "strategy": "single",
  "duration": "5.2ms",
  "memoryUsedMb": 2
}
```

## Estrategias de Generación

| Estrategia | Condición             | Comportamiento                      |
| ---------- | --------------------- | ----------------------------------- |
| Single     | ≤1000 movimientos     | Un solo PDF                         |
| Batch      | >1000 movimientos     | Un PDF con procesamiento en batches |
| ZIP        | >5000 + `UseZip=true` | Múltiples PDFs en ZIP               |

## Configuración

```go
ItemsPerBatch:        50,    // Insumos por batch
MaxMovementsPerBatch: 1500,  // Movimientos por batch
SingleFileThreshold:  1000,  // Límite para archivo único
MergeThreshold:       5000,  // Límite para ZIP
ProcessCleanupDelay:  5000,  // Auto-cleanup de procesos (ms)
UseZip:               false, // Activar estrategia ZIP
```

## Fixes incluidos vs versión anterior

- **Totales corregidos**: columna Saldo Final muestra costo neto (entradas - salidas) en vez de repetir salidas
- **Headers en page break**: se reimprimen labels completos (no vacíos)
- **ProcessTracker auto-cleanup**: procesos completados se limpian automáticamente
- **Código sin duplicación**: headers de tablas en funciones DRY
