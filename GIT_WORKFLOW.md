# Git Workflow - AYO Football Backend

## Branch Strategy

```
main (production)
│
├── develop (development/staging)
│   │
│   ├── feature/auth-module
│   ├── feature/team-management
│   ├── feature/player-management
│   ├── feature/match-management
│   ├── feature/report-module
│   └── feature/api-documentation
│
├── release/v1.0.0
│
└── hotfix/fix-critical-bug
```

## Branch Naming Convention

| Type | Format | Example |
|------|--------|---------|
| Feature | `feature/<feature-name>` | `feature/team-management` |
| Bugfix | `bugfix/<bug-description>` | `bugfix/fix-login-error` |
| Hotfix | `hotfix/<fix-description>` | `hotfix/fix-security-issue` |
| Release | `release/v<version>` | `release/v1.0.0` |
| Develop | `develop` | `develop` |
| Main | `main` | `main` |

## Semantic Commit Messages

### Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

| Type | Description | Example |
|------|-------------|---------|
| `feat` | New feature | `feat(auth): add JWT authentication` |
| `fix` | Bug fix | `fix(player): resolve jersey number validation` |
| `docs` | Documentation | `docs(api): add API documentation` |
| `style` | Formatting, no code change | `style(handler): format code with gofmt` |
| `refactor` | Code refactoring | `refactor(usecase): optimize query performance` |
| `test` | Adding tests | `test(auth): add unit tests for login` |
| `chore` | Maintenance | `chore(deps): update dependencies` |
| `perf` | Performance improvement | `perf(db): add database indexing` |
| `ci` | CI/CD changes | `ci(github): add GitHub Actions workflow` |
| `build` | Build system changes | `build(docker): add Dockerfile` |

### Scopes untuk Project Ini

| Scope | Description |
|-------|-------------|
| `auth` | Authentication module |
| `team` | Team management |
| `player` | Player management |
| `match` | Match management |
| `report` | Report module |
| `api` | API/Handler layer |
| `db` | Database layer |
| `config` | Configuration |
| `docs` | Documentation |
| `deps` | Dependencies |

---

## Alur Commit untuk Project Ini

### Step 1: Setup Branch Structure

```bash
# Pastikan di branch master/main
git checkout master

# Buat branch develop
git checkout -b develop

# Push develop ke remote
git push -u origin develop
```

### Step 2: Buat Feature Branches dan Commit

#### A. Core Infrastructure
```bash
git checkout develop
git checkout -b feature/core-infrastructure

# Commit 1: Project setup
git add go.mod go.sum .gitignore .env.example
git commit -m "chore(config): initialize Go module and project configuration

- Add go.mod with required dependencies
- Add .gitignore for Go projects
- Add .env.example for environment template"

# Commit 2: Configuration
git add internal/config/
git commit -m "feat(config): add application configuration module

- Implement environment-based configuration
- Support PostgreSQL and MySQL drivers
- Add JWT and server configuration"

# Commit 3: Database setup
git add internal/infrastructure/database/postgres.go
git commit -m "feat(db): implement database connection and migration

- Add PostgreSQL connection with GORM
- Implement auto-migration for all entities
- Support multiple database drivers"

# Merge ke develop
git checkout develop
git merge --no-ff feature/core-infrastructure -m "merge: feature/core-infrastructure into develop"
```

#### B. Auth Module
```bash
git checkout develop
git checkout -b feature/auth-module

# Commit 1: User entity
git add internal/domain/entity/base.go internal/domain/entity/user.go
git commit -m "feat(auth): add user entity with role-based access

- Implement BaseEntity with UUID and soft delete
- Add User entity with email, password, role fields
- Define admin and user roles"

# Commit 2: Auth repository
git add internal/domain/repository/user_repository.go
git add internal/infrastructure/database/user_repository.go
git commit -m "feat(auth): implement user repository

- Add UserRepository interface
- Implement GORM-based user repository
- Add FindByEmail and Exists methods"

# Commit 3: JWT service
git add internal/infrastructure/security/
git commit -m "feat(auth): implement JWT authentication service

- Add JWTService interface
- Implement token generation and validation
- Include user claims in token"

# Commit 4: Auth usecase
git add internal/domain/usecase/auth_usecase.go
git commit -m "feat(auth): implement authentication use case

- Add login functionality with password verification
- Add user registration with password hashing
- Implement default admin user creation"

# Commit 5: Auth handler & middleware
git add internal/delivery/http/handler/auth_handler.go
git add internal/delivery/http/middleware/
git add internal/delivery/http/dto/auth_dto.go
git commit -m "feat(auth): add auth handler and middleware

- Implement login and register endpoints
- Add JWT authentication middleware
- Add admin role authorization middleware
- Add CORS and recovery middleware"

# Merge ke develop
git checkout develop
git merge --no-ff feature/auth-module -m "merge: feature/auth-module into develop"
```

#### C. Team Management
```bash
git checkout develop
git checkout -b feature/team-management

# Commit 1: Team entity & repository
git add internal/domain/entity/team.go
git add internal/domain/repository/team_repository.go
git add internal/infrastructure/database/team_repository.go
git commit -m "feat(team): add team entity and repository

- Implement Team entity with name, logo, founded_year, address, city
- Add TeamRepository interface with CRUD operations
- Implement search and pagination"

# Commit 2: Team usecase
git add internal/domain/usecase/team_usecase.go
git commit -m "feat(team): implement team use case with business logic

- Add CRUD operations for teams
- Implement search functionality
- Add team existence validation"

# Commit 3: Team handler
git add internal/delivery/http/handler/team_handler.go
git add internal/delivery/http/dto/team_dto.go
git commit -m "feat(team): add team HTTP handler and DTOs

- Implement REST endpoints for team CRUD
- Add request validation
- Add pagination support"

# Merge ke develop
git checkout develop
git merge --no-ff feature/team-management -m "merge: feature/team-management into develop"
```

#### D. Player Management
```bash
git checkout develop
git checkout -b feature/player-management

# Commit 1: Player entity & repository
git add internal/domain/entity/player.go
git add internal/domain/repository/player_repository.go
git add internal/infrastructure/database/player_repository.go
git commit -m "feat(player): add player entity and repository

- Implement Player entity with position enum
- Add jersey number uniqueness check per team
- Implement player-team relationship"

# Commit 2: Player usecase
git add internal/domain/usecase/player_usecase.go
git commit -m "feat(player): implement player use case with validations

- Add CRUD operations for players
- Validate jersey number uniqueness within team
- Validate player position enum"

# Commit 3: Player handler
git add internal/delivery/http/handler/player_handler.go
git add internal/delivery/http/dto/player_dto.go
git commit -m "feat(player): add player HTTP handler and DTOs

- Implement REST endpoints for player CRUD
- Add team filter support
- Add Indonesian position names in response"

# Merge ke develop
git checkout develop
git merge --no-ff feature/player-management -m "merge: feature/player-management into develop"
```

#### E. Match Management
```bash
git checkout develop
git checkout -b feature/match-management

# Commit 1: Match & Goal entities
git add internal/domain/entity/match.go
git add internal/domain/entity/goal.go
git commit -m "feat(match): add match and goal entities

- Implement Match entity with status enum
- Add Goal entity for tracking scorers
- Implement match result calculation"

# Commit 2: Match & Goal repositories
git add internal/domain/repository/match_repository.go
git add internal/domain/repository/goal_repository.go
git add internal/infrastructure/database/match_repository.go
git add internal/infrastructure/database/goal_repository.go
git commit -m "feat(match): implement match and goal repositories

- Add match CRUD with team preloading
- Implement goal tracking per match
- Add team win count calculation"

# Commit 3: Match usecase
git add internal/domain/usecase/match_usecase.go
git commit -m "feat(match): implement match use case with result recording

- Add match scheduling functionality
- Implement match result recording with goals
- Validate team existence and match status"

# Commit 4: Match handler
git add internal/delivery/http/handler/match_handler.go
git add internal/delivery/http/dto/match_dto.go
git commit -m "feat(match): add match HTTP handler and DTOs

- Implement REST endpoints for match CRUD
- Add result recording endpoint
- Add status and date filtering"

# Merge ke develop
git checkout develop
git merge --no-ff feature/match-management -m "merge: feature/match-management into develop"
```

#### F. Report Module
```bash
git checkout develop
git checkout -b feature/report-module

# Commit 1: Report usecase
git add internal/domain/usecase/report_usecase.go
git commit -m "feat(report): implement report use case

- Add match report generation
- Implement top scorers calculation
- Add team win accumulation"

# Commit 2: Report handler
git add internal/delivery/http/handler/report_handler.go
git add internal/delivery/http/dto/report_dto.go
git commit -m "feat(report): add report HTTP handler and DTOs

- Implement match reports endpoint
- Add top scorers endpoint
- Include accumulated wins in response"

# Merge ke develop
git checkout develop
git merge --no-ff feature/report-module -m "merge: feature/report-module into develop"
```

#### G. API Router & Entry Point
```bash
git checkout develop
git checkout -b feature/api-setup

# Commit 1: Router setup
git add internal/delivery/http/router.go
git commit -m "feat(api): implement HTTP router with all endpoints

- Setup versioned API routes (/api/v1)
- Configure public and protected routes
- Add admin-only middleware for mutations"

# Commit 2: Main entry point
git add cmd/api/main.go
git commit -m "feat(api): add application entry point

- Initialize all dependencies
- Setup graceful shutdown
- Create default admin user on startup"

# Merge ke develop
git checkout develop
git merge --no-ff feature/api-setup -m "merge: feature/api-setup into develop"
```

#### H. Documentation & DevOps
```bash
git checkout develop
git checkout -b feature/documentation

# Commit 1: API Documentation
git add docs/API_DOCUMENTATION.md
git commit -m "docs(api): add comprehensive API documentation

- Document all endpoints with examples
- Add request/response formats
- Include validation rules and error codes"

# Commit 2: Postman Collection
git add docs/postman_collection.json
git commit -m "docs(api): add Postman collection for API testing

- Add all 29 API requests
- Include auto-save variables for tokens
- Add descriptions for each endpoint"

# Commit 3: DevOps files
git add Dockerfile docker-compose.yml Makefile
git commit -m "build(docker): add Docker and build configuration

- Add multi-stage Dockerfile
- Add docker-compose for local development
- Add Makefile with common commands"

# Commit 4: README update
git add README.md
git commit -m "docs(readme): update README with project information

- Add project description and features
- Include setup instructions
- Add API endpoint summary"

# Merge ke develop
git checkout develop
git merge --no-ff feature/documentation -m "merge: feature/documentation into develop"
```

### Step 3: Create Release

```bash
# Buat release branch
git checkout develop
git checkout -b release/v1.0.0

# Final checks dan version bump jika diperlukan
git commit --allow-empty -m "chore(release): prepare v1.0.0 release

Release Notes:
- Initial release of AYO Football API
- Features: Team, Player, Match, Report management
- JWT Authentication with role-based access
- PostgreSQL database with soft delete
- Full API documentation and Postman collection"

# Merge ke main
git checkout main
git merge --no-ff release/v1.0.0 -m "merge: release/v1.0.0 into main"

# Tag the release
git tag -a v1.0.0 -m "v1.0.0 - Initial Release

Features:
- Team Management (CRUD)
- Player Management with jersey number validation
- Match Scheduling and Result Recording
- Reports with top scorers and team statistics
- JWT Authentication
- Soft Delete mechanism
- API Documentation"

# Merge release back to develop
git checkout develop
git merge --no-ff release/v1.0.0 -m "merge: release/v1.0.0 back into develop"

# Delete release branch
git branch -d release/v1.0.0
```

### Step 4: Push All Branches

```bash
# Push main dengan tags
git push origin main --tags

# Push develop
git push origin develop

# Push feature branches (optional, untuk history)
git push origin feature/core-infrastructure
git push origin feature/auth-module
git push origin feature/team-management
git push origin feature/player-management
git push origin feature/match-management
git push origin feature/report-module
git push origin feature/api-setup
git push origin feature/documentation
```

---

## Alternatif: Single Branch dengan Atomic Commits

Jika ingin lebih simple tapi tetap professional:

```bash
# Dari branch master
git checkout -b develop

# Commit berurutan
git add .gitignore .env.example go.mod go.sum
git commit -m "chore: initialize project with Go modules and configuration"

git add internal/config/ internal/infrastructure/database/postgres.go
git commit -m "feat(db): add database configuration and connection"

git add internal/domain/entity/
git commit -m "feat(entity): add domain entities (User, Team, Player, Match, Goal)"

git add internal/domain/repository/ internal/infrastructure/database/*_repository.go
git commit -m "feat(repository): implement repository layer with GORM"

git add internal/infrastructure/security/
git commit -m "feat(auth): add JWT authentication service"

git add internal/domain/usecase/
git commit -m "feat(usecase): implement business logic layer"

git add internal/delivery/http/
git commit -m "feat(api): add HTTP handlers, DTOs, and router"

git add cmd/api/main.go
git commit -m "feat(app): add application entry point with graceful shutdown"

git add docs/
git commit -m "docs: add API documentation and Postman collection"

git add Dockerfile docker-compose.yml Makefile
git commit -m "build: add Docker and build configuration"

git add README.md
git commit -m "docs: update README with project information"

# Merge ke main dan tag
git checkout main
git merge --no-ff develop -m "merge: develop into main - v1.0.0 release"
git tag -a v1.0.0 -m "v1.0.0 - Initial Release"

# Push
git push origin main develop --tags
```

---

## Git Flow Diagram

```
main     ─────●─────────────────────────────────●──────── (v1.0.0)
              │                                 │
              │                                 │
develop  ─────●──●──●──●──●──●──●──●──●──●──●──●────────
              │  │  │  │  │  │  │  │  │  │  │  │
              │  │  │  │  │  │  │  │  │  │  │  │
features      │  └──┘  └──┘  └──┘  └──┘  └──┘  │
              │   A     B     C     D     E    │
              │                                │
release       └────────────────────────────────┘
                        v1.0.0

A = feature/core-infrastructure
B = feature/auth-module
C = feature/team-management
D = feature/player-management
E = feature/match-management + report + docs
```

---

## Tips

1. **Selalu pull sebelum push**: `git pull origin <branch> --rebase`
2. **Review sebelum commit**: `git diff --staged`
3. **Gunakan interactive rebase untuk cleanup**: `git rebase -i HEAD~n`
4. **Protect main branch** di GitHub settings
5. **Require PR reviews** untuk merge ke main/develop
