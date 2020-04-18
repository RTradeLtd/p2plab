package p2plab

us_west_1: Nodes & {
    size: 10
    instanceType: "t3.micro"
    region: "us-west-1"
}

us_east_1: Nodes & {
    size: 2
    instanceType: "t3.medium"
    region: "us-east-1"
    labels: [ "neighbors" ]
}

clust1: Cluster & {
    groups: [ us_west_1, us_east_1 ]
}

object: Object & {
    type: "oci" 
    source: "docker.io/library/golang:latest"
}

scen1: Scenario & {
    objects: [ object ]
}

experiment: Experiment & {
    cluster: clust1
    scenario: scen1
}