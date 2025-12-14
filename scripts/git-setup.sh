#!/bin/bash

# =============================================================================
# Git Setup Script - AYO Football Backend
# =============================================================================
# Script ini akan membuat struktur branch dan commit yang professional
# Jalankan: chmod +x scripts/git-setup.sh && ./scripts/git-setup.sh
# =============================================================================

set -e

echo "=========================================="
echo "  AYO Football Backend - Git Setup"
echo "=========================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_step() {
    echo -e "${GREEN}[STEP]${NC} $1"
}

print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# =============================================================================
# STEP 1: Setup develop branch
# =============================================================================
print_step "Creating develop branch..."
git checkout master 2>/dev/null || git checkout main
git checkout -b develop 2>/dev/null || git checkout develop

# =============================================================================
# STEP 2: Feature - Core Infrastructure
# =============================================================================
print_step "Creating feature/core-infrastructure..."
git checkout -b feature/core-infrastructure 2>/dev/null || git checkout feature/core-infrastructure

git add go.mod go.sum .gitignore .env.example 2>/dev/null || true
git commit -m "chore(config): initialize Go module and project configuration

- Add go.mod with required dependencies (gin, gorm, jwt, uuid)
- Add .gitignore for Go projects
- Add .env.example for environment template

Tech Stack:
- Go 1.23
- Gin v1.9.1
- GORM v1.25.5
- JWT v5.2.0" 2>/dev/null || print_info "Already committed"

git add internal/config/
git commit -m "feat(config): add application configuration module

- Implement environment-based configuration loading
- Support PostgreSQL and MySQL database drivers
- Add JWT secret and expiration configuration
- Add server port and mode configuration" 2>/dev/null || print_info "Already committed"

git add internal/infrastructure/database/postgres.go
git commit -m "feat(db): implement database connection with auto-migration

- Add PostgreSQL/MySQL connection with GORM
- Implement auto-migration for all entities
- Add connection logging and error handling
- Support SSL mode configuration" 2>/dev/null || print_info "Already committed"

git checkout develop
git merge --no-ff feature/core-infrastructure -m "merge: feature/core-infrastructure into develop

Includes:
- Go module initialization
- Application configuration
- Database connection setup"

# =============================================================================
# STEP 3: Feature - Auth Module
# =============================================================================
print_step "Creating feature/auth-module..."
git checkout -b feature/auth-module

git add internal/domain/entity/base.go internal/domain/entity/user.go
git commit -m "feat(auth): add user entity with role-based access control

- Implement BaseEntity with UUID primary key
- Add soft delete support (deleted_at)
- Define User entity with email, password, name, role
- Support admin and user roles" 2>/dev/null || print_info "Already committed"

git add internal/domain/repository/user_repository.go internal/infrastructure/database/user_repository.go
git commit -m "feat(auth): implement user repository layer

- Add UserRepository interface
- Implement GORM-based user repository
- Add FindByEmail for authentication
- Add Exists method for validation" 2>/dev/null || print_info "Already committed"

git add internal/infrastructure/security/
git commit -m "feat(auth): implement JWT authentication service

- Add JWTService interface
- Implement HS256 token generation
- Add token validation and claims extraction
- Include user_id, email, role in claims" 2>/dev/null || print_info "Already committed"

git add internal/domain/usecase/auth_usecase.go
git commit -m "feat(auth): implement authentication use case

- Add Login with bcrypt password verification
- Add Register with password hashing
- Implement CreateDefaultAdmin for initial setup
- Add GetUserByID for profile retrieval" 2>/dev/null || print_info "Already committed"

git add internal/delivery/http/handler/auth_handler.go internal/delivery/http/dto/auth_dto.go internal/delivery/http/middleware/
git commit -m "feat(auth): add HTTP handler and middleware

- Implement login endpoint (POST /auth/login)
- Implement register endpoint (POST /auth/register)
- Implement profile endpoint (GET /auth/profile)
- Add JWT authentication middleware
- Add admin authorization middleware
- Add CORS and recovery middleware" 2>/dev/null || print_info "Already committed"

git checkout develop
git merge --no-ff feature/auth-module -m "merge: feature/auth-module into develop

Includes:
- User entity with roles
- JWT authentication
- Auth endpoints and middleware"

# =============================================================================
# STEP 4: Feature - Team Management
# =============================================================================
print_step "Creating feature/team-management..."
git checkout -b feature/team-management

git add internal/domain/entity/team.go internal/domain/repository/team_repository.go internal/infrastructure/database/team_repository.go
git commit -m "feat(team): add team entity and repository

Fields: nama tim, logo tim, tahun berdiri, alamat markas, kota markas

- Implement Team entity with all required fields
- Add TeamRepository interface with CRUD
- Implement search by name and city
- Add pagination support" 2>/dev/null || print_info "Already committed"

git add internal/domain/usecase/team_usecase.go
git commit -m "feat(team): implement team use case with business logic

- Add Create, Update, Delete operations
- Implement GetByID with optional players
- Add GetAll with search and pagination
- Add team existence validation" 2>/dev/null || print_info "Already committed"

git add internal/delivery/http/handler/team_handler.go internal/delivery/http/dto/team_dto.go
git commit -m "feat(team): add team HTTP handler and DTOs

Endpoints:
- GET /teams - List all teams
- GET /teams/:id - Get team by ID
- POST /teams - Create team (admin)
- PUT /teams/:id - Update team (admin)
- DELETE /teams/:id - Soft delete team (admin)

- Add request validation
- Add pagination in response" 2>/dev/null || print_info "Already committed"

git checkout develop
git merge --no-ff feature/team-management -m "merge: feature/team-management into develop

Implements Requirement #1: Pengelolaan informasi tim sepak bola"

# =============================================================================
# STEP 5: Feature - Player Management
# =============================================================================
print_step "Creating feature/player-management..."
git checkout -b feature/player-management

git add internal/domain/entity/player.go internal/domain/repository/player_repository.go internal/infrastructure/database/player_repository.go
git commit -m "feat(player): add player entity and repository

Fields: nama, tinggi, berat, posisi, nomor punggung

- Implement Player entity with position enum
- Positions: forward, midfielder, defender, goalkeeper
- Add jersey number uniqueness validation per team
- Implement player-team relationship (1 player : 1 team)" 2>/dev/null || print_info "Already committed"

git add internal/domain/usecase/player_usecase.go
git commit -m "feat(player): implement player use case with validations

Business Rules:
- Jersey number must be unique within team
- Jersey number range: 1-99
- Position must be valid enum value
- Player must belong to existing team

- Add CRUD operations
- Add GetByTeamID filter" 2>/dev/null || print_info "Already committed"

git add internal/delivery/http/handler/player_handler.go internal/delivery/http/dto/player_dto.go
git commit -m "feat(player): add player HTTP handler and DTOs

Endpoints:
- GET /players - List all players
- GET /players/:id - Get player by ID
- POST /players - Create player (admin)
- PUT /players/:id - Update player (admin)
- DELETE /players/:id - Soft delete player (admin)

- Add team_id filter parameter
- Include Indonesian position names" 2>/dev/null || print_info "Already committed"

git checkout develop
git merge --no-ff feature/player-management -m "merge: feature/player-management into develop

Implements Requirement #2: Pengelolaan informasi pemain
- Nomor punggung unik per tim ✓"

# =============================================================================
# STEP 6: Feature - Match Management
# =============================================================================
print_step "Creating feature/match-management..."
git checkout -b feature/match-management

git add internal/domain/entity/match.go internal/domain/entity/goal.go
git commit -m "feat(match): add match and goal entities

Match fields: tanggal, waktu, tim home, tim away, skor
Goal fields: player_id, team_id, minute, is_own_goal

- Implement Match entity with status enum
- Status: scheduled, ongoing, completed, cancelled
- Add Goal entity for tracking scorers
- Implement match result calculation (home_win/away_win/draw)" 2>/dev/null || print_info "Already committed"

git add internal/domain/repository/match_repository.go internal/domain/repository/goal_repository.go
git add internal/infrastructure/database/match_repository.go internal/infrastructure/database/goal_repository.go
git commit -m "feat(match): implement match and goal repositories

- Add match CRUD with team preloading
- Implement GetCompletedMatches for reports
- Add GetTeamWinCount for statistics
- Implement goal tracking with player info
- Add GetTopScorers aggregation" 2>/dev/null || print_info "Already committed"

git add internal/domain/usecase/match_usecase.go
git commit -m "feat(match): implement match use case with result recording

Features:
- Match scheduling with team validation
- Result recording with goal details
- Auto status change to 'completed'
- Validation: home_team != away_team

- Add CRUD operations
- Add date range filtering" 2>/dev/null || print_info "Already committed"

git add internal/delivery/http/handler/match_handler.go internal/delivery/http/dto/match_dto.go
git commit -m "feat(match): add match HTTP handler and DTOs

Endpoints:
- GET /matches - List all matches
- GET /matches/:id - Get match with goals
- POST /matches - Create match schedule (admin)
- PUT /matches/:id - Update match (admin)
- DELETE /matches/:id - Soft delete match (admin)
- POST /matches/:id/result - Record result (admin)

- Add status and date filtering
- Include goal scorer details in response" 2>/dev/null || print_info "Already committed"

git checkout develop
git merge --no-ff feature/match-management -m "merge: feature/match-management into develop

Implements:
- Requirement #3: Pengelolaan jadwal pertandingan
- Requirement #4: Pencatatan hasil pertandingan"

# =============================================================================
# STEP 7: Feature - Report Module
# =============================================================================
print_step "Creating feature/report-module..."
git checkout -b feature/report-module

git add internal/domain/usecase/report_usecase.go
git commit -m "feat(report): implement report use case

Report includes:
- Match info (jadwal, tim home/away, skor)
- Match result status (Home Win/Away Win/Draw)
- Top scorer in match
- Accumulated wins for home team
- Accumulated wins for away team

- Add GetMatchReport for single match
- Add GetAllMatchReports with pagination
- Add GetTopScorers leaderboard" 2>/dev/null || print_info "Already committed"

git add internal/delivery/http/handler/report_handler.go internal/delivery/http/dto/report_dto.go
git commit -m "feat(report): add report HTTP handler and DTOs

Endpoints:
- GET /reports/matches - All match reports
- GET /reports/matches/:id - Single match report
- GET /reports/top-scorers - Top scorers leaderboard

Response includes all required report fields" 2>/dev/null || print_info "Already committed"

git checkout develop
git merge --no-ff feature/report-module -m "merge: feature/report-module into develop

Implements Requirement #5: Data report pertandingan"

# =============================================================================
# STEP 8: Feature - API Setup
# =============================================================================
print_step "Creating feature/api-setup..."
git checkout -b feature/api-setup

git add internal/delivery/http/router.go
git commit -m "feat(api): implement HTTP router with versioned endpoints

Routes:
- /health - Health check
- /api/v1/auth/* - Authentication
- /api/v1/teams/* - Team management
- /api/v1/players/* - Player management
- /api/v1/matches/* - Match management
- /api/v1/reports/* - Reports

- Configure public and protected routes
- Add admin-only middleware for mutations" 2>/dev/null || print_info "Already committed"

git add cmd/api/main.go
git commit -m "feat(api): add application entry point

- Initialize all dependencies with DI
- Setup graceful shutdown handler
- Create default admin on startup
- Configure HTTP server with timeouts" 2>/dev/null || print_info "Already committed"

git checkout develop
git merge --no-ff feature/api-setup -m "merge: feature/api-setup into develop

API server ready to run"

# =============================================================================
# STEP 9: Feature - Documentation
# =============================================================================
print_step "Creating feature/documentation..."
git checkout -b feature/documentation

git add docs/API_DOCUMENTATION.md
git commit -m "docs(api): add comprehensive API documentation

Contents:
- System description and features
- Authentication guide
- All endpoints with examples
- Request/response formats
- Validation rules
- Error codes and handling
- Setup instructions" 2>/dev/null || print_info "Already committed"

git add docs/postman_collection.json
git commit -m "docs(api): add Postman collection for testing

Collection includes:
- 29 API requests organized in folders
- Auto-save variables for tokens
- Detailed descriptions
- Example request bodies" 2>/dev/null || print_info "Already committed"

git add Dockerfile docker-compose.yml Makefile 2>/dev/null || true
git commit -m "build(docker): add containerization support

- Multi-stage Dockerfile for optimized image
- docker-compose for local development
- Makefile with common commands" 2>/dev/null || print_info "Already committed"

git add README.md
git commit -m "docs(readme): update project documentation

- Project description and features
- Tech stack overview
- Setup and running instructions
- API endpoint summary
- Default credentials" 2>/dev/null || print_info "Already committed"

git checkout develop
git merge --no-ff feature/documentation -m "merge: feature/documentation into develop

Documentation complete"

# =============================================================================
# STEP 10: Create Release
# =============================================================================
print_step "Creating release/v1.0.0..."
git checkout -b release/v1.0.0

git commit --allow-empty -m "chore(release): prepare v1.0.0

Release Notes - AYO Football API v1.0.0
=======================================

Features:
- Team Management (CRUD with soft delete)
- Player Management with jersey number validation
- Match Scheduling and Result Recording
- Reports with statistics

Technical:
- Golang + GIN Framework
- PostgreSQL with GORM
- JWT Authentication
- Clean Architecture
- Soft Delete mechanism

Documentation:
- API Documentation (Markdown)
- Postman Collection"

# Merge to main
print_step "Merging to main..."
git checkout master 2>/dev/null || git checkout main
git merge --no-ff release/v1.0.0 -m "merge: release/v1.0.0 - Initial Release

AYO Football API v1.0.0

All requirements implemented:
1. ✓ Pengelolaan informasi tim sepak bola
2. ✓ Pengelolaan informasi pemain (jersey unik per tim)
3. ✓ Pengelolaan jadwal pertandingan
4. ✓ Pencatatan hasil pertandingan
5. ✓ Data report pertandingan

Technical requirements:
a. ✓ Database storage (PostgreSQL)
b. ✓ Soft delete mechanism
c. ✓ Assumptions documented
d. ✓ Security (JWT Authentication)"

# Create tag
print_step "Creating tag v1.0.0..."
git tag -a v1.0.0 -m "v1.0.0 - Initial Release

AYO Football Backend API

Features:
- Team Management (CRUD)
- Player Management with jersey number validation
- Match Scheduling and Result Recording
- Reports with top scorers and team statistics

Technical:
- Golang 1.23 + GIN Framework
- PostgreSQL + GORM
- JWT Authentication
- Clean Architecture
- Soft Delete

Documentation:
- Full API documentation
- Postman collection included"

# Merge back to develop
git checkout develop
git merge --no-ff release/v1.0.0 -m "merge: release/v1.0.0 back into develop"

# Cleanup
git branch -d release/v1.0.0

print_success "Git setup complete!"
echo ""
echo "=========================================="
echo "  Summary"
echo "=========================================="
echo ""
echo "Branches created:"
echo "  - main (production)"
echo "  - develop (development)"
echo "  - feature/core-infrastructure"
echo "  - feature/auth-module"
echo "  - feature/team-management"
echo "  - feature/player-management"
echo "  - feature/match-management"
echo "  - feature/report-module"
echo "  - feature/api-setup"
echo "  - feature/documentation"
echo ""
echo "Tag created:"
echo "  - v1.0.0"
echo ""
echo "To push to remote:"
echo "  git push origin main develop --tags"
echo "  git push origin --all  # (to push all branches)"
echo ""
echo "=========================================="
