Создание kind

    kind create cluster --config=./k8s/kind.yaml
    kind delete cluster

Создание образа

    docker build -t lp-app:1.0.0 -f Dockerfile .
    docker build -t lp-migrator:1.0.0 -f Dockerfile.migrator .

Создание временного каталога

    mkdir -p ~/kind-temp
    export TMPDIR=~/kind-temp

Сохранение образов

    kind load docker-image lp-migrator:1.0.0
    kind load docker-image lp-app:1.0.0
    kind load docker-image sso-app:1.0.0
    kind load docker-image queue-manager-app:1.0.0
    kind load docker-image notification-app:1.0.0
    kind load docker-image api-gateway-app:1.0.0

Создание секретов

    kubectl create secret generic sso-keys \
  --from-file=private_key.pem=./k8s/private_key.pem \
  --from-file=public_key.pem=./k8s/public_key.pem

Последовательность применения манифестов

    kubectl apply -f ./k8s/configmaps
    kubectl apply -f ./k8s/pvs
    kubectl apply -f ./k8s/services
    kubectl apply -f ./k8s/jobs
    kubectl apply -f ./k8s/deployments/dependencies
    kubectl apply -f ./k8s/deployments

    kubectl delete deployments mongo-deployment
    kubectl apply -f ./k8s/deployments/dependencies/sso-mongodb-deployment.yaml

Проблемы с монгой

    kubectl exec deployment/mongo-deployment -it -- /bin/bash
    mongosh
    use admin
    db.createUser({user: "test", pwd: "test", roles: [{role: "root", db: "admin"}]})

    kubectl delete deployments sso-app-deployment
    kubectl apply -f ./k8s/deployments/sso-deployment.yaml

Логи

    kubectl logs lp-migrator-6874b94cdf-zh8sh -c lp-container

Рестарт

    kubectl rollout restart deployment/lp-app-deployment

Проброс портов

    kubectl port-forward svc/rabbitmq 15672:15672
    kubectl port-forward svc/mongo 27017:27017

{
  "email": "dima@gmail.com",
  "password": "Qwerty!123456"
}

kubectl exec -it api-gateway-app-deployment-7b6b44787d-xj4z4 -- nslookup sso-app-service


kubectl delete configmap api-gateway-config
kubectl delete configmap sso-config
kubectl delete deployments api-gateway-app-deployment
kubectl delete deployments sso-app-deployment
kubectl delete service sso-service

kubectl apply -f ./k8s/configmaps/api-gateway-configmap.yaml
kubectl apply -f ./k8s/configmaps/sso-configmap.yaml
kubectl apply -f ./k8s/services/sso-service.yaml
kubectl apply -f ./k8s/deployments/api-gateway-deployment.yaml
kubectl apply -f ./k8s/deployments/sso-deployment.yaml

