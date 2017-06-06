#!/bin/bash

cat > /pipeline/source/bcpayslip-service.json << EOF
{
   "apiVersion": "v1",
   "kind": "Service",
   "metadata": {
      "name": "${MS_NAME}",
      "labels": {
         "name": "${MS_NAME}"
      }
   },
   "spec":{
      "type": "LoadBalancer",
      "ports": [
         {
           "port": 80,
           "targetPort": "${MS_NAME}-web",
           "protocol": "TCP"
         }
      ],
      "selector":{
         "name":"${MS_NAME}"
      }
   }
}
EOF
