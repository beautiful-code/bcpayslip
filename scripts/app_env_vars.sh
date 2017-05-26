#!/bin/bash

cat > .env << EOF
  bc_intranet_client_id=${BC_CLIENT_ID}
  bc_intranet_client_secret=${BC_CLIENT_SECRET}
  bc_app_key=${BC_APP_KEY}
  bc_env=${BC_ENV}
  bc_host=${BC_LOCALHOST}
  bc_mongo_db="${MS_NAME}"
  PORT=${BC_PORT}
  MONGO_URI=${BC_MONGO_URI}
EOF
