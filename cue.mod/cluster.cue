package p2plab

// defines a set of nodes of size 1 or higher
// a "node" is simply an EC2 instance provisioned of the given type
Nodes :: {
    // default of 1
    size: int 
    instanceType: string
    region: string
    labels?: [...string]
}

// a group is a collection of 1 or more nodes
Group :: {
    members: [...Nodes]
}

// a cluster is a collection of 1 or more groups
Cluster :: {
    groups: [...Group]
}

us_west_1: Nodes & {
    size: 10
    instanceType: "t3.micro"
    region: "us-west-1"
}

us_east_1: Nodes & {
    size: 2
    instanceType: "t3.medium"
    region: "us-east-1"
}

group1: Group & {
    members: [ us_west_1, us_east_1 ]
}

cluster: Cluster & {
    groups: [ group1 ]
}