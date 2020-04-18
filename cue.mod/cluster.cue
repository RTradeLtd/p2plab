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

// a cluster is a collection of 1 or more groups of nodes
Cluster :: {
    groups: [...Nodes]
}