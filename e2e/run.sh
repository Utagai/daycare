#!/bin/bash
set -e

echo "=== Daycare E2E Tests ==="
echo ""

# Create test directories
export DAYCARE_CONFIG_DIR="$HOME/.config/daycare"
mkdir -p "$DAYCARE_CONFIG_DIR"
export DAYCARE_PROJECTS_DIR="$HOME/projects"
mkdir -p "$DAYCARE_PROJECTS_DIR"

# Write test config
cat > "$DAYCARE_CONFIG_DIR/config.yaml" << 'EOF'
projectsDir: /root/projects
EOF

echo "Config:"
cat "$DAYCARE_CONFIG_DIR/config.yaml"
echo ""

# Test 1: Create project
echo "Test 1: Creating project 'testproject'..."
daycare pcreate testproject
echo ""

# Test 2: Verify project structure
echo "Test 2: Verifying project structure..."
test -d "$DAYCARE_PROJECTS_DIR/testproject/.daycare" && echo "PASS: .daycare/ exists" || echo "FAIL: .daycare/ missing"
test -f "$DAYCARE_PROJECTS_DIR/testproject/.daycare/pwcd" && echo "PASS: pwcd exists" || echo "FAIL: pwcd missing"
test -f "$DAYCARE_PROJECTS_DIR/testproject/.daycare/wcreate" && echo "PASS: wcreate exists" || echo "FAIL: wcreate missing"
echo ""

# Test 3: Create worktree
echo "Test 3: Creating worktree 'feature-branch'..."
daycare wcreate testproject feature-branch
echo ""

# Test 4: List worktrees
echo "Test 4: Listing worktrees..."
daycare wlist testproject
echo ""

# Test 5: List projects
echo "Test 5: Listing projects..."
daycare plist
echo ""

# Test 6: pcd output
echo "Test 6: Testing pcd output..."
daycare pcd testproject
echo ""

# Test 7: wcd output
echo "Test 7: Testing wcd output..."
daycare wcd testproject feature-branch
echo ""

# Test 8: Completion generation
echo "Test 8: Testing completion generation..."
daycare completion zsh | head -20
echo ""

# Test 9: Setup command
echo "Test 9: Testing setup command..."
daycare setup
echo ""

echo "=== All E2E tests completed ==="
