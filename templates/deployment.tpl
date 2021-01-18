apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.UserName}}
  labels:
    app: spider
    email: {{.Email}}
  namespace: {{.Namespace}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: spider
      email: {{.Email}}
  template:
    metadata:
      labels:
        app: spider
        email: {{.Email}}
    spec:
      securityContext:
        runAsUser: 1001
        runAsGroup: 1001
        fsGroup: 2000
     containers:
      - name: {{.UserName}}
        image: {{.Image}}
        ports:
        - containerPort: 3000
          name: "ide"
        env:
        - name: USER
          value: {{.UserName}}