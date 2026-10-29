# IP App — Go

Returns the IP address of the machine (or container) it runs on.
Language: **Go** · Framework: none (Go standard library `net/http`) · Build tool: **go** · Port: **8000**

| URL | Returns |
|-----|---------|
| `http://localhost:8000/` | `{"ip":"172.19.0.4","message":"The IP address is: 172.19.0.4"}` |
| `http://localhost:8000/greet` | `{"message":"Hi how are you"}` |

**Build tool in simple points**
- Go is **compiled**: `go build` turns the code into one program file (`ip-app` / `ip-app.exe`).
- That file runs on its own — no Go, no libraries needed on the server. Copy it and run it.
- `go.mod` = module name + Go version (like `pom.xml` / `package.json`). This app uses only the standard library, so no downloads.
- `go run .` = build + run in one step (for development). `go build` = create the file (for release).
- One machine can build for another OS: set `GOOS` / `GOARCH` (see step 5).

---

## 1. Install (one time)

### Ubuntu
Ubuntu's own `golang-go` package is old, so install the latest Go from go.dev (works on Intel and ARM servers):
```bash
sudo apt update
sudo apt install -y curl git
GO_VERSION=$(curl -s 'https://go.dev/VERSION?m=text' | head -1)
GO_FILE="${GO_VERSION}.linux-$(dpkg --print-architecture).tar.gz"
curl -fsSLO "https://go.dev/dl/${GO_FILE}"
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf "${GO_FILE}" && rm "${GO_FILE}"
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
go version
```

### Windows (PowerShell)
**Package manager — winget** (built into Windows 10/11). Check it:
```powershell
winget --version
```
Only if `winget` is not found — open PowerShell **as Administrator** and run:
```powershell
Install-PackageProvider -Name NuGet -Force | Out-Null
Install-Module -Name Microsoft.WinGet.Client -Force -Repository PSGallery | Out-Null
Repair-WinGetPackageManager -AllUsers
```
Install Go and Git, reload PATH:
```powershell
winget install -e --id GoLang.Go --accept-source-agreements --accept-package-agreements
winget install -e --id Git.Git --accept-source-agreements --accept-package-agreements
$env:Path = [Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [Environment]::GetEnvironmentVariable('Path','User')
go version
```

## 2. Get the code
```bash
git clone https://github.com/itdefined-org-apps/ip_app.git
cd ip_app/golang
```
(Already cloned: `cd ip_app && git pull && cd golang`)

## 3. Build and run

### Ubuntu
```bash
go fmt ./...          # format the code
go vet ./...          # check for common mistakes
go build -o ip-app .
ls -lh ip-app         # the whole app is this one file
./ip-app
```

### Windows (PowerShell)
```powershell
go fmt ./...
go vet ./...
go build -o ip-app.exe .
Get-Item ip-app.exe   # the whole app is this one file
.\ip-app.exe
```
You should see: `Listening on port 8000`. Leave it running.
Windows may ask to allow network access the first time — click **Allow**.
Quick development run without creating a file: `go run .`

## 4. Test (in a second terminal)

### Ubuntu
```bash
curl http://localhost:8000/
curl http://localhost:8000/greet
```

### Windows (PowerShell)
```powershell
curl.exe http://localhost:8000/
curl.exe http://localhost:8000/greet
```
You should see the JSON from the table at the top.
On a server (EC2), use `http://<server-public-ip>:8000/` and open port 8000 in the security group.

## 5. Build for another OS (cross-compile)
Build a Windows program on Ubuntu:
```bash
GOOS=windows GOARCH=amd64 go build -o ip-app.exe .
```
Build a Linux program on Windows (PowerShell) — copy `ip-app` to an Ubuntu server and run `./ip-app`:
```powershell
$env:GOOS='linux'; $env:GOARCH='amd64'; go build -o ip-app .; Remove-Item Env:GOOS, Env:GOARCH
```
Use `GOARCH=arm64` for ARM servers (e.g. AWS Graviton `t4g` instances).

## 6. Stop and clean up
Press **Control-C** in the first terminal, then delete the built files:
```bash
rm -f ip-app ip-app.exe        # Windows: Remove-Item ip-app, ip-app.exe -ErrorAction SilentlyContinue
```

---

## 7. Docker (later)
Same commands on Ubuntu and Windows (Docker installation: see the main [README](../README.md)).
```bash
docker build -t golang-ip-app:latest .
docker run -d --name golang_app -p 8083:8000 golang-ip-app:latest
curl http://localhost:8083/          # Windows: curl.exe http://localhost:8083/
docker logs golang_app
docker stop golang_app && docker rm golang_app
```
- The `Dockerfile` has two stages: stage 1 compiles with the Go image, stage 2 copies **only the program** into a small Alpine image.
- Result: `docker images` shows about **22 MB**, compared with 210-320 MB for the Python and Node.js images.
- The image is named `golang-ip-app`, not `golang`, so it does not replace Docker's official `golang` image.
- To run all four apps together: `docker compose up -d --build` in the `ip_app` folder.

## Tests (as in a real project)

| Check | Tool | Where | Needs |
|-------|------|-------|-------|
| Lint | gofmt + go vet | all code | nothing |
| Unit tests | `go test` + `httptest` | `tests/unit` / `*_test.go` | nothing — the network cards are replaced through the `interfaceAddrs` variable |
| Coverage gate | built into the test tool | — | fails below **80%** (now 89%) |
| Security scan | govulncheck `-scan package` | dependencies + Go standard library | internet |
| Integration tests | `go test -tags=integration` | `tests/integration` / `integration_test.go` | the app running (`BASE_URL`) |

- `main_test.go` = unit tests (`go test` finds every `*_test.go` file).
- `integration_test.go` starts with `//go:build integration`, so it runs **only** with `-tags=integration`.
- No test libraries to install — testing is built into Go.
- Coverage is 89%: the only red lines are in `main()`, which runs only when you start the app yourself.

### Run lint, unit tests and security scan
### Ubuntu
```bash
test -z "$(gofmt -l .)" && go vet ./...        # Lint (no output = OK)
go test -v -cover ./...                           # Unit tests + coverage
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 -scan package ./...   # Security scan
```

### Windows (PowerShell)
```powershell
gofmt -l .; go vet ./...
go test -v -cover ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 -scan package ./...
```

Lint fails? Fix most problems automatically: `gofmt -w .`.

On a small server (1 GB RAM, e.g. t3.micro) govulncheck's default mode is killed for lack of memory, so this project uses
`-scan package` (checks every imported package, about 20 MB of RAM). On a bigger server you can drop `-scan package`
for the default mode, which also checks whether the vulnerable function is actually called.

### Run the integration tests
Start the app in terminal 1 (`./ip-app`), then in terminal 2:
```bash
BASE_URL=http://localhost:8000 go test -tags=integration -run Live -v ./...
```
Windows (PowerShell): `$env:BASE_URL='http://localhost:8000'; go test -tags=integration -run Live -v ./...`

Unit tests: **7**, integration tests: **2**.

## CI/CD with Jenkins
`Jenkinsfile` in this folder runs the same checks, then builds, tests and deploys the Docker image:

| # | Stage | What happens | Fails when |
|---|-------|--------------|------------|
| 1 | Checkout | Gets the code | - |
| 2 | Install | Installs app + test tools | Download error |
| 3 | Lint | gofmt + go vet | Style or code problem |
| 4 | Unit Test | Unit tests + coverage report | A test fails or coverage < 80% |
| 5 | Security Scan | govulncheck | A dependency has a known vulnerability |
| 6 | Build Image | `docker build -t golang-ip-app:<build-number>` | Dockerfile error |
| 7 | Container Test | Runs the new image on a free port, runs the integration tests against it | An integration test fails |
| 8 | Push Image | Pushes to Docker Hub (only if `PUSH_IMAGE` is ticked) | Login or push error |
| 9 | Deploy (dev) | Runs the image on the Jenkins server on port **9083** (only if `DEPLOY` is ticked) | App does not answer |

Set up:
1. Jenkins server needs Go in /usr/local/go/bin, Docker and curl; add the Jenkins user to the docker group:
   `sudo usermod -aG docker jenkins && sudo systemctl restart jenkins`
2. Plugins: **Pipeline, Git, JUnit, Credentials Binding** (Pipeline Stage View is nice to have).
3. New Item -> **Pipeline** -> *Pipeline script from SCM* -> Git repo URL, branch `*/main`, Script Path **`golang/Jenkinsfile`**.
4. Only for Push Image: Manage Jenkins -> Credentials -> *Username with password*, ID **`dockerhub-creds`**
   (Docker Hub username + an access token from Docker Hub -> Account settings -> Personal access tokens).

Results: build page -> **Tests** (unit + integration), **Build Artifacts** (coverage report), and the app at
`http://<jenkins-server-ip>:9083/`. Open port 9083 in the security group.
If Jenkins itself runs inside Docker, change `TEST_HOST` in the `Jenkinsfile` to `host.docker.internal`.

## If something fails

| Problem | Fix |
|---------|-----|
| `go: command not found` (Ubuntu) | Run `export PATH=$PATH:/usr/local/go/bin`, or open a new terminal |
| `go` not found (Windows) | Close and reopen PowerShell after installing Go |
| `go: go.mod requires go >= 1.26` | Your Go is too old — install the latest with the commands in step 1 |
| `./ip-app: Permission denied` | You copied the file from Windows — run `chmod +x ip-app` |
| `address already in use` | Another app is on port 8000 — stop it (all four apps use 8000) |
