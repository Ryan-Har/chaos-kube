# Tiltfile

# Define the image name and path to the Dockerfile
docker_build('pandects/chaos-executor:latest', 'src', 
             dockerfile='src/cmd/chaos-executor/Dockerfile', 
             ignore=['.git'],
            #  live_update=[
            #      sync('src', '/src'),  # Sync changes from local source to container's /src
            #      run('go install ./...', '/src')  # Rebuild Go code in container
            #  ]
             )

docker_build('pandects/chaos-controller:latest', 'src', 
             dockerfile='src/cmd/chaos-controller/Dockerfile', 
             ignore=['.git'],
            #  live_update=[
            #      sync('src', '/src'),  # Sync changes from local source to container's /src
            #      run('go install ./...', '/src')  # Rebuild Go code in container
            #  ]
             )

# Deploy the Helm chart
k8s_yaml(helm(
    'chaos-kube-chart',  # Path to the Helm chart
    namespace='chaos-kube',
    values=['chaos-kube-chart/values/values-dev.yaml']
))


