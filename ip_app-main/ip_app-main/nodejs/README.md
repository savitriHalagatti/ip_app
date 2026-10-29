# IP App — Node.js Express

Returns the IP address of the machine (or container) it runs on.
Language: **JavaScript (Node.js 24 LTS)** · Framework: **Express** · Build tool: **npm** · Port: **8000**

| URL | Returns |
|-----|---------|
| `http://localhost:8000/` | `{"ip":"172.19.0.3","message":"The IP address is: 172.19.0.3"}` |
| `http://localhost:8000/greet` | `{"message":"Hi how are you"}` |

**Build tool in simple points**
- JavaScript is not compiled to a file — the "build" is installing the dependencies.
- `package.json` = project name, libraries and scripts (like `pom.xml` in Maven, `requirements.txt` in Python).
- `npm install` = downloads the libraries into `node_modules/` and writes `package-lock.json` (exact versions).
- `npm start` = runs the `start` script from `package.json` (`node app.js`).
- `package-lock.json` is committed, and CI/CD uses `npm ci` — it installs exactly the locked versions.

---

## 1. Install (one time)

### Ubuntu
Ubuntu's own `nodejs` package is old, so install Node.js 24 LTS from NodeSource:
```bash
sudo apt update
sudo apt install -y curl git
curl -fsSL https://deb.nodesource.com/setup_24.x | sudo -E bash -
sudo apt install -y nodejs
node -v && npm -v
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
Install Node.js LTS and Git, allow npm scripts, reload PATH:
```powershell
winget install -e --id OpenJS.NodeJS.LTS --accept-source-agreements --accept-package-agreements
winget install -e --id Git.Git --accept-source-agreements --accept-package-agreements
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned -Force
$env:Path = [Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [Environment]::GetEnvironmentVariable('Path','User')
node -v; npm -v
```

## 2. Get the code
```bash
git clone https://github.com/itdefined-org-apps/ip_app.git
cd ip_app/nodejs
```
(Already cloned: `cd ip_app && git pull && cd nodejs`)

## 3. Build and run
Same commands on Ubuntu and Windows:
```bash
npm install
npm start
```
You should see: `added .. packages` and then `Listening on port 8000`. Leave it running.
Look at what the build created: the `node_modules/` folder (versions come from `package-lock.json`).

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

## 5. Stop and clean up
Press **Control-C** in the first terminal, then delete what `npm install` created:
```bash
rm -rf node_modules     # Windows: Remove-Item -Recurse -Force node_modules
```

---

## 6. Docker (later)
Same commands on Ubuntu and Windows (Docker installation: see the main [README](../README.md)).
```bash
docker build -t nodejs:latest .
docker run -d --name nodejs_app -p 8082:8000 nodejs:latest
curl http://localhost:8082/          # Windows: curl.exe http://localhost:8082/
docker logs nodejs_app
docker stop nodejs_app && docker rm nodejs_app
```
- `-p 8082:8000` = port 8082 on your machine -> port 8000 in the container.
- `.dockerignore` keeps your local `node_modules/` out of the image; the image runs its own `npm install`.
- To run all four apps together: `docker compose up -d --build` in the `ip_app` folder.

## Tests (as in a real project)

| Check | Tool | Where | Needs |
|-------|------|-------|-------|
| Lint | ESLint | all code | nothing |
| Unit tests | Jest + supertest | `tests/unit` / `*_test.go` | nothing — the network cards are replaced with `jest.spyOn(os, "networkInterfaces")` |
| Coverage gate | built into the test tool | — | fails below **80%** (now 89%) |
| Security scan | npm audit | dependencies | internet |
| Integration tests | Jest + `fetch` | `tests/integration` / `integration_test.go` | the app running (`BASE_URL`) |

- `package.json` `devDependencies` = test tools (Jest, supertest, ESLint) — not installed in the Docker image (`npm ci --omit=dev`).
- `package-lock.json` = exact versions; it is committed so `npm ci` installs the same versions everywhere.
- `eslint.config.js` = lint rules. `app.js` exports `app` so tests can use it without starting a server.
- Coverage is 89%: the only red lines are `app.listen(...)`, which run only when you start the app yourself.

### Run lint, unit tests and security scan
### Ubuntu
```bash
npm ci                      # install exact versions from package-lock.json
npm run lint                # Lint (ESLint)
npm run test:ci             # Unit tests + coverage (fails below 80%)
npm audit --omit=dev --audit-level=high      # Security scan
```

Windows: same commands in PowerShell.

Lint fails? Fix most problems automatically: `npx eslint --fix .`.

### Run the integration tests
Start the app in terminal 1 (`npm start`), then in terminal 2:
```bash
BASE_URL=http://localhost:8000 npm run test:integration
```
Windows (PowerShell): `$env:BASE_URL='http://localhost:8000'; npm run test:integration`

Unit tests: **5**, integration tests: **2**.

## CI/CD with Jenkins
`Jenkinsfile` in this folder runs the same checks, then builds, tests and deploys the Docker image:

| # | Stage | What happens | Fails when |
|---|-------|--------------|------------|
| 1 | Checkout | Gets the code | - |
| 2 | Install | Installs app + test tools | Download error |
| 3 | Lint | ESLint | Style or code problem |
| 4 | Unit Test | Unit tests + coverage report | A test fails or coverage < 80% |
| 5 | Security Scan | npm audit | A dependency has a known vulnerability |
| 6 | Build Image | `docker build -t nodejs:<build-number>` | Dockerfile error |
| 7 | Container Test | Runs the new image on a free port, runs the integration tests against it | An integration test fails |
| 8 | Push Image | Pushes to Docker Hub (only if `PUSH_IMAGE` is ticked) | Login or push error |
| 9 | Deploy (dev) | Runs the image on the Jenkins server on port **9082** (only if `DEPLOY` is ticked) | App does not answer |

Set up:
1. Jenkins server needs Node.js 24 + npm, Docker and curl; add the Jenkins user to the docker group:
   `sudo usermod -aG docker jenkins && sudo systemctl restart jenkins`
2. Plugins: **Pipeline, Git, JUnit, Credentials Binding** (Pipeline Stage View is nice to have).
3. New Item -> **Pipeline** -> *Pipeline script from SCM* -> Git repo URL, branch `*/main`, Script Path **`nodejs/Jenkinsfile`**.
4. Only for Push Image: Manage Jenkins -> Credentials -> *Username with password*, ID **`dockerhub-creds`**
   (Docker Hub username + an access token from Docker Hub -> Account settings -> Personal access tokens).

Results: build page -> **Tests** (unit + integration), **Build Artifacts** (coverage report), and the app at
`http://<jenkins-server-ip>:9082/`. Open port 9082 in the security group.
If Jenkins itself runs inside Docker, change `TEST_HOST` in the `Jenkinsfile` to `host.docker.internal`.

## If something fails

| Problem | Fix |
|---------|-----|
| Windows: `npm.ps1 cannot be loaded` | Run `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned -Force` |
| Windows: `node` or `npm` not found | Close and reopen PowerShell after installing Node.js |
| `Error: Cannot find module 'express'` | Run `npm install` first |
| Ubuntu: `node -v` shows v12 / v18 | Old Ubuntu Node.js — run the NodeSource commands above |
| `EADDRINUSE` | Another app is on port 8000 — stop it (all four apps use 8000) |
