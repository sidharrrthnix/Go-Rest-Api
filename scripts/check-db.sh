#!/bin/bash

# scripts/check-db.sh - Verify MariaDB setup is working

set -e

echo "🔍 Checking MariaDB Docker Setup..."
echo "=================================="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker Desktop.${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Docker is running${NC}"

# Check if containers are up
if ! docker compose ps | grep -q "http-api-mariadb"; then
    echo -e "${YELLOW}⚠ MariaDB container not found. Starting containers...${NC}"
    docker compose up -d
    echo "Waiting 30 seconds for MariaDB to initialize..."
    sleep 30
fi

# Check MariaDB health
HEALTH=$(docker inspect --format='{{.State.Health.Status}}' http-api-mariadb 2>/dev/null || echo "unknown")
if [ "$HEALTH" == "healthy" ]; then
    echo -e "${GREEN}✓ MariaDB container is healthy${NC}"
else
    echo -e "${YELLOW}⚠ MariaDB health status: $HEALTH (may still be starting)${NC}"
fi

# Test database connection
echo ""
echo "Testing database connection..."
if docker exec http-api-mariadb mariadb -u api_user -papi_password -e "SELECT 1" http_api_db > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Database connection successful${NC}"
else
    echo -e "${RED}❌ Database connection failed${NC}"
    exit 1
fi

# Check tables
echo ""
echo "Checking database tables..."
TABLES=$(docker exec http-api-mariadb mariadb -u api_user -papi_password -N -e "SHOW TABLES" http_api_db 2>/dev/null)
if [ -n "$TABLES" ]; then
    echo -e "${GREEN}✓ Tables found:${NC}"
    echo "$TABLES" | while read table; do
        echo "  - $table"
    done
else
    echo -e "${YELLOW}⚠ No tables found (init script may not have run)${NC}"
fi

# Show record counts
echo ""
echo "Record counts:"
docker exec http-api-mariadb mariadb -u api_user -papi_password -e "
    SELECT 'teachers' as table_name, COUNT(*) as count FROM teachers
    UNION ALL
    SELECT 'students', COUNT(*) FROM students;
" http_api_db 2>/dev/null || echo "Could not query tables"

# Check Adminer
echo ""
if docker compose ps | grep -q "http-api-adminer.*running"; then
    echo -e "${GREEN}✓ Adminer GUI is running at: http://localhost:8081${NC}"
else
    echo -e "${YELLOW}⚠ Adminer container is not running${NC}"
fi

echo ""
echo "=================================="
echo -e "${GREEN}🎉 MariaDB Setup Verification Complete!${NC}"
echo ""
echo "Connection Details:"
echo "  Host: localhost"
echo "  Port: 3306"
echo "  Database: http_api_db"
echo "  Username: api_user"
echo "  Password: api_password"
echo ""
echo "Adminer GUI: http://localhost:8081"
echo "  System: MySQL"
echo "  Server: mariadb"
echo "  Username: api_user"
echo "  Password: api_password"
echo "  Database: http_api_db"
