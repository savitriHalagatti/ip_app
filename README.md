# ip_app — the same app in 4 languages

Each app returns the IP address of the machine (or container) it runs on, plus a `/greet` message.
Use it to compare build tools across languages, then run all of them in Docker.

| Folder | Language | Build tool | What "build" means | Manual run | Docker port |
|--------|----------|-----------|--------------------|-----------|-------------|
| [`fast_api/`](fast_api/README.md) | Python 3.12 (FastAPI) | pip + venv | Install libraries from `requirements.txt` | `localhost:8000` | `8080` |
| [`flask/`](flask/README.md) | Python 3.12 (Flask) | pip + venv | Install libraries from `requirements.txt` | `localhost:8000` | `8081` |
| [`nodejs/`](nodejs/README.md) | Node.js 24 (Express) | npm | Install libraries from `package.json` into `node_modules/` | `localhost:8000` | `8082` |
| [`golang/`](golang/README.md) | Go | go | **Compile** into one program file (`ip-app`) | `localhost:8000` | `8083` |

Every folder README has: install commands (Ubuntu + Windows), build, run, test, clean up, and Docker.
All apps listen on port **8000** when run by hand — run one at a time. Docker maps them to 8080-8083 so all four run together.

**Switch to the next app: stop whatever is running on port 8000**

Press **Control-C** in the terminal where the app runs. If that terminal is closed or the app is in the background:

Ubuntu
```bash
fuser -k 8000/tcp
```
Mac
```bash
kill $(lsof -ti :8000)
```
Windows (PowerShell)
```powershell
Get-NetTCPConnection -LocalPort 8000 -State Listen | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }
```
See what is using the port first (optional): Ubuntu `sudo ss -ltnp 'sport = :8000'`, Mac `lsof -i :8000`,
Windows `Get-NetTCPConnection -LocalPort 8000`. If you started the app with `sudo`, add `sudo` to the stop command.

**Build tools side by side**

| | Python | Node.js | Go | Java (Maven, for comparison) |
|-|--------|---------|----|------------------------------|
| Dependency file | `requirements.txt` | `package.json` | `go.mod` | `pom.xml` |
| Install dependencies | `pip install -r requirements.txt` | `npm install` | `go mod download` | `mvn dependency:resolve` |
| Where they go | `.venv/` | `node_modules/` | Go module cache | `~/.m2/` |
| Build output | none (runs the `.py` files) | none (runs the `.js` files) | one binary file | `.jar` / `.war` file |
| Run | `python app.py` / `fastapi run` | `npm start` | `./ip-app` | `java -jar` / Tomcat |
| Needs the language on the server? | Yes | Yes | **No** | Yes (Java) |

**Tests and CI/CD side by side** (details in each folder's README)

| | Python (FastAPI / Flask) | Node.js | Go |
|-|--------------------------|---------|----|
| Lint | `ruff` | `ESLint` | `gofmt` + `go vet` |
| Unit tests | `pytest` | `Jest` + `supertest` | `go test` (built in) |
| Fake the IP in tests | `monkeypatch` | `jest.spyOn` | replaceable variable |
| Coverage gate (80%) | `pytest-cov` | Jest `coverageThreshold` | `go tool cover` |
| Security scan | `pip-audit` | `npm audit` | `govulncheck` |
| Integration tests | `pytest` + HTTP client | `Jest` + `fetch` | `go test -tags=integration` |
| Test tools listed in | `requirements-dev.txt` | `devDependencies` | none needed |
| Jenkinsfile | `fast_api/Jenkinsfile`, `flask/Jenkinsfile` | `nodejs/Jenkinsfile` | `golang/Jenkinsfile` |
| Dev deploy port | 9080 / 9081 | 9082 | 9083 |

Every Jenkinsfile has the same stages: **Checkout -> Install -> Lint -> Unit Test -> Security Scan -> Build Image ->
Container Test -> Push Image (optional) -> Deploy (dev)**. Only the commands change per language.

---

## Docker (later)

### Install Docker

**Ubuntu**
```bash
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER
newgrp docker
docker version
docker compose version
```

**Windows (PowerShell)** — check winget first: `winget --version` (to install winget, see any folder README).
```powershell
winget install -e --id Docker.DockerDesktop --accept-source-agreements --accept-package-agreements
```
Restart Windows, start **Docker Desktop** once (it sets up WSL 2), then check in a new PowerShell:
```powershell
docker version
docker compose version
```

### Run all four apps
```bash
git clone https://github.com/itdefined-org-apps/ip_app.git
cd ip_app
docker compose up -d --build
docker compose ps
```
Test (Windows: use `curl.exe` instead of `curl`):
```bash
curl http://localhost:8080/      # FastAPI
curl http://localhost:8081/      # Flask
curl http://localhost:8082/      # Node.js
curl http://localhost:8083/      # Go
```
Each one shows its **container's** IP (e.g. `172.19.0.2`) — every container gets its own IP on the Docker network.

**On a server that also runs Tomcat or Jenkins** (they use port 8080), FastAPI fails with
`failed to bind host port 0.0.0.0:8080/tcp: address already in use`. Either stop Tomcat during the Docker demo
(`sudo systemctl stop tomcat10`), or change FastAPI's port in `docker-compose.yml` from `"8080:8000"` to `"8084:8000"`.

**Disk space:** building all four images needs about **2 GB free** (`df -h /`). On a small EC2 disk, free space afterwards with
`docker builder prune -af` (removes the build cache, keeps the images).

```bash
docker compose logs -f golang    # logs of one app (Control-C to stop following)
docker compose down              # stop and remove all containers
docker images                    # compare image sizes: Go ~22 MB, Python/Node 210-320 MB
```

### What `docker-compose.yml` shows
- `build: context:` — builds each image from its folder's `Dockerfile`.
- `ports: "8080:8000"` — your machine's port -> the container's port.
- `fastapi_backend` also shows `environment` variables, a custom `networks` entry and a named `volumes` entry.
- The named volume `fastapi_data` is mounted on `/app`. Docker copies the image's code into it only the **first** time;
  after you change `app.py` and rebuild, run `docker compose down -v` (deletes the volume) so the new code is used.

## Repository files
- `.gitignore` keeps build and test output out of Git: `.venv/`, `node_modules/`, `reports/`, the Go binary.
