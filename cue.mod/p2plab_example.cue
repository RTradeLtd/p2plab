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

clust: Cluster & {
    groups: [ us_west_1, us_east_1 ]
}

object: Object & {
    type: "oci" 
    source: "docker.io/library/golang:latest"
}

scen: Scenario & {
    objects: [ object ]
    benchmark: {
        "(not neighbors)": "golang"
    }
}


experiment: Experiment & {
    cluster: clust
    scenario: scen
}