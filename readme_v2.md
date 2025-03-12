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

Удаление всех подов

    kubectl delete pod -l app=lp

Применить новый манифест

    kubectl apply -f ./k8s/deployments/lp-deployment.yaml
    kubectl apply -f ./k8s/configmaps/lp-configmap.yaml

    kubectl delete deployments lp-app-deployment
    kubectl delete configmaps lp-config

    kubectl apply -f ./k8s/configmaps
    kubectl apply -f ./k8s/pvs
    kubectl apply -f ./k8s/services
    kubectl apply -f ./k8s/jobs
    kubectl apply -f ./k8s/deployments

    kubectl logs lp-migrator-6874b94cdf-zh8sh -c lp-container

    kubectl delete job db-migration
    kubectl apply -f ./k8s/deployments/lp-migrator-job.yaml

Рестарт

    kubectl rollout restart deployment/lp-app-deployment

Проброс портов

    kubectl port-forward svc/rabbitmq 15672:15672
    kubectl port-forward svc/mongo 27017:27017


Создание секретов

    kubectl create secret generic sso-keys \
  --from-file=private_key.pem=./k8s/private_key.pem \
  --from-file=public_key.pem=./k8s/public_key.pem

kubectl apply -f ./k8s/configmaps/sso-configmap.yaml
