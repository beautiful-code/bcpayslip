#!/bin/bash

cat > .env << EOF
  export bc_intranet_client_id=${BC_CLIENT_ID}
  export bc_intranet_client_secret=${BC_CLIENT_SECRET}
  export bc_app_key=${BC_APP_KEY}
  export bc_env=${BC_ENV}
  export bc_host=${BC_LOCALHOST}
  export bc_mongo_db="${MS_NAME}db"
  export PORT=${BC_PORT}
  export MONGO_URI=${BC_MONGO_URI}
EOF
