#!/bin/bash

cat > /pipeline/source/bcpayslip-deployment.yaml <<EOF
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: ${MS_NAME}
spec:
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  minReadySeconds: 5
  template:
    metadata:
      labels:
        name: ${MS_NAME}
    spec:
      imagePullSecrets:
        - name: pto-registry-creds
      containers:
        - image: priyankhub/${MS_NAME}:${WERCKER_GIT_COMMIT}
          imagePullPolicy: Always
          name: ${MS_NAME}-webapp
          command: ["/go/src/${MS_NAME}/bcpayslip"]
          ports:
            - containerPort: 3001
              name: ${MS_NAME}-w
              protocol: TCP
          env:
            - name: bc_intranet_client_id
              value: "${BC_CLIENT_ID}"
            - name: bc_intranet_client_secret
              value: "${BC_CLIENT_SECRET}"
            - name: bc_app_key
              value: "${BC_APP_KEY}"
            - name: bc_env
              value: "${BC_ENV}"
            - name: bc_host
              value: "${BC_LOCALHOST}"
            - name: bc_mongo_db
              value: "${MS_NAME}"
            - name: PORT
              value: "${BC_PORT}"
            - name: MONGO_URI
              value: "${BC_MONGO_URI}"
EOF
