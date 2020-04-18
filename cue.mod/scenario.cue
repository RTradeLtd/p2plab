package p2plab

Object :: {
    type: string
    source: string
}

Scenario :: {
    objects: [...Object]
    // enable any fields for benchmark
    benchmark:  { ... }
}

Seed :: {
    neighbors: string
}