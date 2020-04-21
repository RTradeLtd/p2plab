package p2plab


items: object & {
    golang: {
        type: "oci"
        source: "docker.io/library/golang:latest"
    }
    mysql: {
        type: "oci"
        source: "docker.io/library/mysql:latest"
    }
}

experiment: Experiment & {
    cluster: Cluster & {
        groups: [ 
            Nodes & {
                size: 10
                instanceType: "t3.micro"
                region: "us-west-1"
            }, 
            Nodes & {
                size: 2
                instanceType: "t3.medium"
                region: "us-east-1"
                labels: [ "neighbors" ]
            } 
        ]
    }
    scenario: Scenario & {
        objects:  [ items ]
        seed: {
            "neighbors": "golang"
        }
        benchmark: {
            "(not neighbors)": "golang"
        }
    }
}