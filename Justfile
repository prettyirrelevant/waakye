set dotenv-load

# Runs database migrations
@apply-migrations:
    echo "Attempting to run migrations for $DATABASE_NAME..."
    cd api/database/migrations && goose sqlite3 $DATABASE_NAME up -v

# Creates a new SQL migration file with the NAME provided
@create-migration NAME:
    echo "Creating migration..."
    cd api/database/migrations && goose create {{NAME}} sql -v
