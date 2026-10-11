# monolitoScanner — Vulnerability Scanner distribuido

Escáner de vulnerabilidades con arquitectura distribuida: un frontend, un
reverse proxy / load balancer (nginx), dos réplicas de un middleware
(service discovery + load balancing) y 3 microservicios de backend, cada
uno con su propio patrón de diseño y su propia conexión a base de datos.
Orquestado completamente con Docker Compose.

> **Nota:** los 3 módulos de escaneo (`port_scan`, `http_scan`, `tls_scan`)
> devuelven resultados **simulados** por ahora — el propósito actual del
> proyecto es la arquitectura distribuida (discovery, balanceo, patrones
> de diseño), no el escaneo real todavía.

## Arquitectura

```
┌────────────────────────────────────────────────────────────┐
│  FRONTEND (nginx:alpine, puerto host 3000)                  │
│  index.html + style.css + script.js                         │
│  fetch() → POST http://localhost:8080/scan                  │
└──────────────────────────┬─────────────────────────────────┘
                           ▼
┌────────────────────────────────────────────────────────────┐
│  NGINX GATEWAY (reverse proxy + load balancer, host 8080)   │
│  upstream middleware_pool → round robin, max_fails=3        │
│  Único punto de entrada: los middlewares ya no son públicos │
└───────────────┬─────────────────────────┬──────────────────┘
                ▼                         ▼
        ┌───────────────┐         ┌───────────────┐
        │ middleware-1  │         │ middleware-2  │
        │ discovery +   │         │ discovery +   │
        │ round robin   │         │ round robin   │
        └───────┬───────┘         └───────┬───────┘
                └────────────┬────────────┘
        ┌────────────────────┼────────────────────┐
        ▼                    ▼                    ▼
 ┌─────────────┐      ┌─────────────┐      ┌─────────────┐
 │ port-scan-1 │      │ http-scan-1 │      │ tls-scan-1  │
 │ port-scan-2 │      │ (Strategy)  │      │ (Factory)   │
 │ (Repository)│      └──────┬──────┘      └──────┬──────┘
 └──────┬──────┘             │                    │
        └────────────────────┴────────────────────┘
                             │ INSERT INTO scan_results
                             ▼
                  ┌─────────────────────┐
                  │  Postgres (db)       │
                  └─────────────────────┘
```

9 contenedores: `frontend`, `nginx-gateway`, `middleware-1`, `middleware-2`,
`port-scan-1`, `port-scan-2`, `http-scan-1`, `tls-scan-1`, `db`.

Hay **dos niveles de balanceo**:

| Nivel | Quién | Reparte entre | Cómo |
|---|---|---|---|
| 1 | `nginx-gateway` | `middleware-1` y `middleware-2` | `upstream` con round robin |
| 2 | cada `middleware` | las réplicas sanas de cada servicio (ej. `port-scan-1/2`) | round robin en Go + health checks cada 5s |

## Reverse proxy / Load balancer (nginx)

Corre en su propio contenedor (`nginx-gateway/`), entre el frontend y los
middlewares. El frontend ya no llama al middleware directamente: llama a
nginx, y nginx reenvía.

`nginx-gateway/nginx.conf`:

```nginx
upstream middleware_pool {
    server middleware-1:8080 max_fails=3 fail_timeout=30s;
    server middleware-2:8080 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    location /scan {
        proxy_pass http://middleware_pool;
        ...
    }
    location /health { return 200 "nginx gateway alive"; }
}
```

- `upstream middleware_pool`: el grupo de servidores detrás de un solo
  nombre. Sin otra directiva, nginx usa **round robin**.
- `max_fails=3 fail_timeout=30s`: health check **pasivo**. Si una réplica
  falla 3 veces, nginx no le manda tráfico durante 30s. (El health check
  activo, como el de `registry.go`, no existe en nginx open source.)
- `proxy_pass`: reenvía la petición al pool.
- `proxy_set_header ...`: conserva la IP original del cliente
  (`X-Real-IP`, `X-Forwarded-For`).
- `add_header X-Upstream $upstream_addr`: devuelve en la respuesta qué
  middleware atendió la petición (útil para demostrar el balanceo).

El puerto del host no cambió (`8080`), por eso `script.js` no se tocó.

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
├── nginx-gateway/
│   ├── Dockerfile
│   └── nginx.conf         (reverse proxy + load balancer)
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
docker-compose up --build
```
o, si tu instalación usa el binario standalone con guion:
```bash
docker-compose up --build
```

## Verificar que todo quedó arriba

```bash
docker-compose ps
```
Los 9 servicios deben mostrar `Up` (o `running`).

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

**Ver el balanceo de nginx entre los middlewares** (revisa el header
`X-Upstream`; la IP debe alternar entre dos valores):
```bash
for i in 1 2 3 4; do
  curl -si -X POST http://localhost:8080/scan \
    -d '{"target":"192.168.1.20","modules":["port_scan"]}' | grep X-Upstream
done
```

**Ver qué middleware atendió cada petición:**
```bash
docker-compose logs middleware-1 middleware-2 | grep dispatch
```

**Probar la tolerancia a fallos de nginx** (apaga una réplica y sigue
respondiendo con la otra). Con una réplica caída, todas las respuestas
deben ser `200`:
```bash
docker-compose stop middleware-2
for i in 1 2 3 4 5 6; do
  curl -s -o /dev/null -w "%{http_code} " -X POST http://localhost:8080/scan \
    -d '{"target":"192.168.1.20","modules":["port_scan"]}'
done; echo
docker-compose start middleware-2
```
(Si tu instalación usa el binario standalone como en este proyecto, escribe `docker-compose` con
guion en lugar de `docker compose`.)

**Ver los resultados guardados en la base de datos:**
```bash
docker-compose exec db psql -U postgres -d monolito -c "SELECT * FROM scan_results;"
```

## Resultados de las pruebas del gateway

**1. Reparto round robin** — el header `X-Upstream` alterna entre las IPs
internas de `middleware-1` y `middleware-2`:
```
X-Upstream: 172.18.0.7:8080
X-Upstream: 172.18.0.8:8080
X-Upstream: 172.18.0.7:8080
X-Upstream: 172.18.0.8:8080
```

**2. Una réplica caída (`middleware-2` detenido)** — cuando a nginx le toca
la réplica caída, falla la conexión y reintenta con la otra. El header
lista ambos intentos y la **última IP** es la que respondió:
```
X-Upstream: 172.18.0.7:8080
X-Upstream: 172.18.0.8:8080, 172.18.0.7:8080
X-Upstream: 172.18.0.7:8080
X-Upstream: 172.18.0.8:8080, 172.18.0.7:8080
```
Las IPs dobles siguen apareciendo hasta que se acumulan `max_fails=3`
fallos; entonces nginx marca la réplica como caída durante `fail_timeout`
(30s) y deja de intentarla.

**3. Visto desde el cliente** — con una réplica caída, el servicio no se
interrumpe:
```
200 200 200 200 200 200
```

## Apagar

```bash
docker-compose down          # detiene y elimina contenedores (conserva los datos)
docker-compose down -v       # además borra el volumen de la base de datos
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
- El health check de nginx es pasivo (reacciona a fallos reales de
  peticiones); el activo solo lo hace el middleware con sus servicios.
- Las URLs de los middlewares están fijas en `nginx.conf` (no hay
  discovery dinámico a nivel nginx).
- Configuración (DSN de Postgres, URL del middleware) hardcodeada en el
  código, no por variables de entorno.
- Sin reintentos automáticos si una instancia falla a mitad de una
  petición (aunque sí se detecta como caída en el siguiente chequeo de
  salud, hasta 5s después).
