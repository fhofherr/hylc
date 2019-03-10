#!/usr/bin/env sh

if ! command -v pwgen > /dev/null 2>&1
then
    echo "pwgen required"
    exit 1
fi

ENV_FILE="$1"

if [ -z "$ENV_FILE" ]
then
    echo "ENV_FILE not set"
    exit 1
fi

cat << EOF > .env
POSTGRES_PASSWORD=`pwgen -n1 30`

HYDRA_DB_NAME=hydra
HYDRA_DB_USER=hydra
HYDRA_DB_PASS=`pwgen -n1 30`
HYDRA_SYSTEM_SECRET=`pwgen -n1 30`

HYLC_DB_NAME=hylc
HYLC_DB_USER=hylc
HYLC_DB_PASS=`pwgen -n1 30`
EOF
