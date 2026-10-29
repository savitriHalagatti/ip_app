# IP App — Python Flask

Returns the IP address of the machine (or container) it runs on.
Language: **Python 3.12** · Framework: **Flask** · Build tool: **pip + venv** · Port: **8000**

| URL | Returns |
|-----|---------|
| `http://localhost:8000/` | `The IP address is: 172.19.0.2` (plain text) |
| `http://localhost:8000/greet` | `Hi how are you` |

**Build tool in simple points**
- Python is not compiled to a file — the "build" is installing the dependencies.
- `requirements.txt` = list of libraries (like `pom.xml` in Maven, `package.json` in Node.js).
- `pip` = downloads and installs the libraries.
- `venv` = a private folder (`.venv/`) with Python + libraries only for this project, so projects do not clash.
- Flask vs FastAPI: same idea; Flask returns plain text here, FastAPI returns JSON and has `/docs`.

---

## 1. Install (one time)

### Ubuntu
```bash
sudo apt update
sudo apt install -y python3 python3-venv python3-pip git curl
python3 --version
```
Any Python 3.12 or newer works (Ubuntu 24.04 has 3.12, Ubuntu 26.04 has 3.14). The Docker image uses 3.12.

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
Install Python and Git, allow venv scripts, reload PATH:
```powershell
winget install -e --id Python.Python.3.12 --accept-source-agreements --accept-package-agreements
winget install -e --id Git.Git --accept-source-agreements --accept-package-agreements
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned -Force
$env:Path = [Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [Environment]::GetEnvironmentVariable('Path','User')
py -3.12 --version
```

## 2. Get the code
```bash
git clone https://github.com/itdefined-org-apps/ip_app.git
cd ip_app/flask
```
(Already cloned: `cd ip_app && git pull && cd flask`)

## 3. Build and run

### Ubuntu
```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
python app.py
```

### Windows (PowerShell)
```powershell
py -3.12 -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
python app.py
```
You should see: `Running on http://127.0.0.1:8000`. Leave it running.
The red `WARNING: This is a development server` is normal for a demo.

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
You should see `The IP address is: ...` and `Hi how are you`.
On a server (EC2), use `http://<server-public-ip>:8000/` and open port 8000 in the security group.

## 5. Stop and clean up
Press **Control-C** in the first terminal, then:
```bash
deactivate          # leave the venv
rm -rf .venv        # Ubuntu: delete the venv (Windows: Remove-Item -Recurse -Force .venv)
```

---

## 6. Docker (later)
Same commands on Ubuntu and Windows (Docker installation: see the main [README](../README.md)).
```bash
docker build -t flask:latest .
docker run -d --name flask_app -p 8081:8000 flask:latest
curl http://localhost:8081/          # Windows: curl.exe http://localhost:8081/
docker logs flask_app
docker stop flask_app && docker rm flask_app
```
- `-p 8081:8000` = port 8081 on your machine -> port 8000 in the container.
- The IP shown is the **container's** IP, not your machine's.
- To run all four apps together: `docker compose up -d --build` in the `ip_app` folder.

## Tests (as in a real project)

| Check | Tool | Where | Needs |
|-------|------|-------|-------|
| Lint | ruff | all code | nothing |
| Unit tests | pytest + Flask test client | `tests/unit` / `*_test.go` | nothing — the IP lookup is replaced with `monkeypatch`, so the result is always `10.1.2.3` |
| Coverage gate | built into the test tool | — | fails below **80%** (now 92%) |
| Security scan | pip-audit | dependencies | internet |
| Integration tests | pytest + requests | `tests/integration` / `integration_test.go` | the app running (`BASE_URL`) |

- `requirements.txt` = what the **app** needs (goes into the Docker image).
- `requirements-dev.txt` = app + **test tools** (pytest, ruff, pip-audit, requests) — only for developers and Jenkins.
- `pyproject.toml` = pytest and ruff settings.
- Coverage is 92%: the only red line is `app.run(...)`, which runs only when you start the app yourself.

### Run lint, unit tests and security scan
### Ubuntu
```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements-dev.txt
ruff check . && ruff format --check .                        # Lint
pytest --cov=app --cov-report=term-missing --cov-fail-under=80   # Unit tests + coverage
pip-audit -r requirements.txt                                # Security scan
```

### Windows (PowerShell)
```powershell
py -3.12 -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements-dev.txt
ruff check .; ruff format --check .
pytest --cov=app --cov-report=term-missing --cov-fail-under=80
pip-audit -r requirements.txt
```

Lint fails? Fix most problems automatically: `ruff check --fix . && ruff format .`.

### Run the integration tests
Start the app in terminal 1 (`python app.py`), then in terminal 2:
```bash
BASE_URL=http://localhost:8000 pytest tests/integration
```
Windows (PowerShell): `$env:BASE_URL='http://localhost:8000'; pytest tests/integration`

Unit tests: **5**, integration tests: **2**.

## CI/CD with Jenkins
`Jenkinsfile` in this folder runs the same checks, then builds, tests and deploys the Docker image:

| # | Stage | What happens | Fails when |
|---|-------|--------------|------------|
| 1 | Checkout | Gets the code | - |
| 2 | Install | Installs app + test tools | Download error |
| 3 | Lint | ruff | Style or code problem |
| 4 | Unit Test | Unit tests + coverage report | A test fails or coverage < 80% |
| 5 | Security Scan | pip-audit | A dependency has a known vulnerability |
| 6 | Build Image | `docker build -t flask:<build-number>` | Dockerfile error |
| 7 | Container Test | Runs the new image on a free port, runs the integration tests against it | An integration test fails |
| 8 | Push Image | Pushes to Docker Hub (only if `PUSH_IMAGE` is ticked) | Login or push error |
| 9 | Deploy (dev) | Runs the image on the Jenkins server on port **9081** (only if `DEPLOY` is ticked) | App does not answer |

Set up:
1. Jenkins server needs python3 + python3-venv, Docker and curl; add the Jenkins user to the docker group:
   `sudo usermod -aG docker jenkins && sudo systemctl restart jenkins`
2. Plugins: **Pipeline, Git, JUnit, Credentials Binding** (Pipeline Stage View is nice to have).
3. New Item -> **Pipeline** -> *Pipeline script from SCM* -> Git repo URL, branch `*/main`, Script Path **`flask/Jenkinsfile`**.
4. Only for Push Image: Manage Jenkins -> Credentials -> *Username with password*, ID **`dockerhub-creds`**
   (Docker Hub username + an access token from Docker Hub -> Account settings -> Personal access tokens).

Results: build page -> **Tests** (unit + integration), **Build Artifacts** (coverage report), and the app at
`http://<jenkins-server-ip>:9081/`. Open port 9081 in the security group.
If Jenkins itself runs inside Docker, change `TEST_HOST` in the `Jenkinsfile` to `host.docker.internal`.

## If something fails

| Problem | Fix |
|---------|-----|
| Ubuntu: `error: externally-managed-environment` | You ran `pip` outside the venv — run `source .venv/bin/activate` first |
| Windows: `Activate.ps1 cannot be loaded` | Run `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned -Force` |
| Windows: `py` or `python` not found | Close and reopen PowerShell after installing Python |
| `ModuleNotFoundError: No module named 'flask'` | Activate the venv, then `pip install -r requirements.txt` |
| `Address already in use` | Another app is on port 8000 — stop it (all four apps use 8000) |
