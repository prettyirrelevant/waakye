set dotenv-load

# Runs database migrations
@apply-migrations:
    echo "Attempting to run migrations for {{justfile_directory()}}/$DATABASE_URI..."
    cd api/database/migrations && goose sqlite3 {{justfile_directory()}}/$DATABASE_URI up -v

# Creates a new SQL migration file with the NAME provided
@create-migration NAME:
    echo "Creating migration..."
    cd api/database/migrations && goose create {{NAME}} sql -v


run-dev NAME:
    #!/usr/bin/env python3
    import sys
    import subprocess

    if "{{NAME}}" in {"waakye", "masa"}:
        print("Starting {{NAME}} server...")
        subprocess.run("air -c .air.{{NAME}}.toml", shell=True)
    else:
        print("Unknown server {{NAME}}. Exiting...")
        sys.exit(1)
