# monolitoScanner — Vulnerability Scanner distribuido

Escáner de vulnerabilidades con arquitectura distribuida: un frontend, un
middleware (service discovery + load balancing), y 3 microservicios de
backend, cada uno con su propio patrón de diseño y su propia conexión a
base de datos. Orquestado completamente con Docker Compose.

> **Nota:** los 3 módulos de escaneo (`port_scan`, `http_scan`, `tls_scan`)
> devuelven resultados **simulados** por ahora — el propósito actual del
> proyecto es la arquitectura distribuida (discovery, balanceo, patrones
> de diseño), no el escaneo real todavía.

## Arquitectura

```
┌──────────────────────────────────────────────────────────────────┐
│  FRONTEND (contenedor: nginx, puerto host 3000)                    │
│  index.html + style.css + script.js                                │
│  fetch() → POST http://localhost:8080/scan                         │
└───────────────────────────┬──────────────────────────────────────┘
                             │  {target, modules[]}
                             ▼
┌──────────────────────────────────────────────────────────────────┐
│  MIDDLEWARE / GATEWAY (contenedor, puerto host 8080)                │
│                                                                      │
│  ┌──────────────┐   ┌────────────────┐   ┌──────────────────┐     │
│  │   REGISTRY    │──▶│  HEALTH CHECKS │──▶│  LOAD BALANCER     │     │
│  │ services.json │   │ goroutine, 5s   │   │  round robin       │     │
│  │ doc. estático │   │ GET /health     │   │  sync.Mutex         │     │
│  └──────────────┘   └────────────────┘   └──────────────────┘     │
│                                                                      │
│  Por cada módulo pedido: discovery → balanceo → forwardToService()   │
│  → junta resultados → responde JSON                                  │
└──────┬─────────────────┬───────────────────────┬────────────────────┘
       │ HTTP POST /scan  │ HTTP POST /scan       │ HTTP POST /scan
       ▼                  ▼                       ▼
┌─────────────┐   ┌─────────────┐        ┌─────────────┐
│ port-scan-1 │   │ http-scan-1 │        │ tls-scan-1  │
│ port-scan-2 │   │ (Strategy)  │        │ (Factory)   │
│ (Repository)│   └──────┬──────┘        └──────┬──────┘
└──────┬──────┘          │                      │
       └──────────────────┴──────────────────────┘
                          │ INSERT INTO scan_results
                          ▼
               ┌─────────────────────┐
               │  Postgres (db)       │
               │  tabla scan_results  │
               └─────────────────────┘
```

7 contenedores en total: `frontend`, `middleware`, `port-scan-1`,
`port-scan-2`, `http-scan-1`, `tls-scan-1`, `db`.

## Estructura del proyecto

```
monolitoScanner/
├── docker-compose.yml
├── schema.sql
├── setup.sh
├── frontend/
│   ├── Dockerfile
│   ├── index.html
│   ├── style.css
│   └── script.js
├── middleware/
│   ├── Dockerfile
│   ├── go.mod
│   ├── main.go
│   ├── registry.go        (service discovery)
│   ├── loadbalancer.go    (round robin)
│   └── services.json      (documento estático de discovery)
└── services/
    ├── port-scan/    (Repository Pattern)
    ├── http-scan/    (Strategy Pattern)
    └── tls-scan/     (Factory Pattern)
```

## Requisitos

**El proyecto en sí NO requiere Linux** — corre en cualquier sistema
operativo compatible con Docker (Windows, macOS o Linux), porque todo
corre dentro de contenedores: no necesitas instalar Go, Postgres, ni
nada más en tu máquina, sin importar el sistema operativo.

Lo que **sí es específico de Linux** es el script `setup.sh` incluido
aquí, porque usa `apt` (gestor de paquetes de Debian/Ubuntu/Kali) para
instalar Docker automáticamente si falta. En Windows o macOS, instala
Docker manualmente (ver abajo) y luego levanta el proyecto con
`docker-compose up --build` directamente, sin usar `setup.sh`.

| Sistema operativo | Cómo obtener Docker + Compose |
|---|---|
| **Linux** (Debian/Ubuntu/Kali) | Ejecuta `./setup.sh` — lo instala todo automáticamente |
| **Windows** | Instala [Docker Desktop](https://www.docker.com/products/docker-desktop/) (incluye Compose) |
| **macOS** | Instala [Docker Desktop](https://www.docker.com/products/docker-desktop/) (incluye Compose) |

Además necesitas los puertos **3000**, **8080** y **5432** libres en tu
máquina.

## Cómo levantarlo

### En Linux (Debian/Ubuntu/Kali)

```bash
chmod +x setup.sh
./setup.sh
```

El script:
1. Verifica si Docker está instalado; si no, lo instala (`apt install docker.io`).
2. Se asegura de que el servicio de Docker esté corriendo.
3. Verifica si tienes Docker Compose (plugin `docker compose` o binario
   standalone `docker-compose`); si no tienes ninguno, instala el binario.
4. Te agrega al grupo `docker` (para no necesitar `sudo` en el futuro —
   requiere cerrar sesión y volver a entrar para tomar efecto).
5. Levanta todo el proyecto con `--build -d` (en segundo plano).

### En Windows o macOS (con Docker Desktop ya instalado y corriendo)

```bash
docker compose up --build
```
o, si tu instalación usa el binario standalone con guion:
```bash
docker-compose up --build
```

## Verificar que todo quedó arriba

```bash
docker compose ps
```
Los 7 servicios deben mostrar `Up` (o `running`).

## Probar

**Desde el navegador:**
```
http://localhost:3000
```

**Desde la terminal:**
```bash
curl -X POST http://localhost:8080/scan \
  -H "Content-Type: application/json" \
  -d '{"target":"192.168.1.20","modules":["port_scan","http_scan","tls_scan"]}'
```

**Ver los resultados guardados en la base de datos:**
```bash
docker compose exec db psql -U postgres -d monolito -c "SELECT * FROM scan_results;"
```

## Apagar

```bash
docker compose down          # detiene y elimina contenedores (conserva los datos)
docker compose down -v       # además borra el volumen de la base de datos
```

## Patrones de diseño usados en el backend

| Servicio | Patrón | Dónde |
|---|---|---|
| `port-scan` | Repository | `repository.go` — interfaz `ResultRepository` + `PostgresRepository` |
| `http-scan` | Strategy | `strategy.go` — interfaz `CheckStrategy` con verificaciones intercambiables |
| `tls-scan` | Factory | `factory.go` — función `NewChecker(kind)` que construye el checker correcto |

## Limitaciones conocidas

- El escaneo es simulado, no real (ver nota al inicio).
- Sin autenticación en ningún endpoint.
- Configuración (DSN de Postgres, URL del middleware) hardcodeada en el
  código, no por variables de entorno.
- Sin reintentos automáticos si una instancia falla a mitad de una
  petición (aunque sí se detecta como caída en el siguiente chequeo de
  salud, hasta 5s después).
