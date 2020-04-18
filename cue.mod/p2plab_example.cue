package p2plab

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
        objects: [ 
            Object & {
                type: "oci" 
                source: "docker.io/library/golang:latest"
            } 
        ]
        benchmark: {
            "(not neighbors)": "golang"
        }
    }
}