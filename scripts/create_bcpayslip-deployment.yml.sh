#!/bin/bash

cat > bcpayslip-deployment.yaml <<EOF
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: ${MS_NAME}
spec:
  replicas: 2
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
          command: ["./go/src/${MS_NAME}/app"]
          ports:
            - containerPort: 3001
              name: ${MS_NAME}-web-server
              protocol: TCP
EOF
