package p2plab

// defines a set of nodes of size 1 or higher
// a "node" is simply an EC2 instance provisioned of the given type
Nodes :: {
    // must be greater than or equal to 1
    // default value of this field is 1
    size: >=1 | *1
    instanceType: string
    region: string
    labels?: [...string]
}

// a cluster is a collection of 1 or more groups of nodes
Cluster :: {
    groups: [...Nodes]
}