#!/usr/bin/env sh

set -e

die() {
    [ -n "$1" ] && echo "$1"
    [ -n "$2" ] && exit $2
    exit 1
}

[ -n "$HYDRA_DB_NAME" ] || die "HYDRA_DB_NAME not set"
[ -n "$HYDRA_DB_USER" ] || die "HYDRA_DB_USER not set"
[ -n "$HYDRA_DB_PASS" ] || die "HYDRA_DB_PASS not set"

[ -n "$HYLC_DB_NAME" ] || die "HYLC_DB_NAME not set"
[ -n "$HYLC_DB_USER" ] || die "HYLC_DB_USER not set"
[ -n "$HYLC_DB_PASS" ] || die "HYLC_DB_PASS not set"


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOF
    CREATE USER $HYDRA_DB_USER ENCRYPTED PASSWORD '$HYDRA_DB_PASS';
    CREATE DATABASE $HYDRA_DB_NAME OWNER $HYDRA_DB_USER ENCODING UTF8;

    CREATE USER $HYLC_DB_USER ENCRYPTED PASSWORD '$HYLC_DB_PASS';
    CREATE DATABASE $HYLC_DB_NAME OWNER $HYLC_DB_USER ENCODING UTF8;
EOF
